package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"
	"github.com/yanyiwu/gojieba"

	"github.com/cifer76/oversee/collectors"
	"github.com/cifer76/oversee/entity"
)

var (
	interestedWords []string
	tgbot           *tgbotapi.BotAPI

	sentNews      = cache.New(24*time.Hour, 60*time.Second)
	jieba         = gojieba.NewJieba()
	cacheFileName = "sentNews.bin"
)

func requestInterestedWords() []string {
	gistURL := "https://gist.githubusercontent.com/cifer76/7f14f5bd02b98abbc11d16662266a572/raw"

	interestedWords := []string{}
	resp, err := http.Get(gistURL)
	if err != nil {
		return interestedWords
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return interestedWords
	}
	interestedWords = strings.Fields(string(body))
	return interestedWords
}

func needCheckInterest(source string) bool {
	return source != "Stcn" && source != "Zqrb"
}

func loadSentNews() {
	sentNews.LoadFile(cacheFileName)
}

func saveSentNews() {
	sentNews.SaveFile(cacheFileName)
}

func checkDuplicates(piece entity.PieceOfNews) bool {
	dup := false
	for _, item := range sentNews.Items() {
		v := item.Object.(entity.PieceOfNews)
		if piece.Link == v.Link {
			dup = true
			fmt.Printf("Found duplicates:\n\t  new: %v\n\texist: %v\n", piece, v)
			break
		}

		// construct a map for easy check
		check := make(map[string]bool, len(v.Tags))
		for _, t := range v.Tags {
			check[t] = true
		}

		// cut off news title
		count := 0
		words := jieba.Cut(piece.Title, true)
		piece.Tags = words
		for _, w := range words {
			if ok, _ := check[w]; ok {
				count++
			}
		}

		// if overlap rate exceeds 60%, take them as duplicated
		if float64(count)/float64(len(check)) >= 0.6 {
			dup = true
			fmt.Printf("Found duplicates:\n\tnew: %v\n\texist: %v\n", piece, v)
			break
		}
	}
	if !dup {
		sentNews.Set(piece.Link, piece, 0)
	}
	return dup
}

func sendNews(news entity.PieceOfNews) {

	if needCheckInterest(news.Source) {
		interested := false
		for _, w := range interestedWords {
			if strings.Contains(news.Title, w) {
				interested = true
				break
			}
		}
		if !interested {
			return
		}
	}

	// dont send duplicate articles
	if checkDuplicates(news) {
		return
	}

	content := fmt.Sprintf("<a href=\"%s\">%s</a>", news.Link, news.Title)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:           406797693,
			ReplyToMessageID: 0,
		},
		Text:      content,
		ParseMode: "HTML",
	}
	if _, err := tgbot.Send(msg); err != nil {
		fmt.Printf("tg send fail, error: %v, message: %v\n", err, content)
	}
}

func main() {
	// initialize tgbot
	tgbot, _ = tgbotapi.NewBotAPI("407954143:AAGDxLmxcr5DGVE3GY_Ih9pe8GIh-P0EhDI")
	//bot.Debug = true

	// initialize collectors
	collectors.Init()

	// load cache
	loadSentNews()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		saveSentNews()
		os.Exit(0)
	}()

	// update interested words
	go func() {
		for {
			interestedWords = requestInterestedWords()
			time.Sleep(5 * time.Minute)
		}
	}()

	for news := range collectors.Visit() {
		fmt.Printf("Source: %s, %q -> %s\n", news.Source, news.Title, news.Link)
		sendNews(news)
	}
}
