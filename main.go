package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/patrickmn/go-cache"

	"github.com/cifer76/oversee/collectors"
	"github.com/cifer76/oversee/entity"
)

var (
	interestedWords []string
	tgbot           *tgbotapi.BotAPI
	sentLinks       = cache.New(24*time.Hour, 60*time.Second)
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

func sendNews(news entity.PieceOfNews) {
	// dont send duplicate articles
	if _, found := sentLinks.Get(news.Link); found {
		return
	}

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
	sentLinks.Set(news.Link, true, 0)
}

func main() {
	tgbot, _ = tgbotapi.NewBotAPI("407954143:AAGDxLmxcr5DGVE3GY_Ih9pe8GIh-P0EhDI")
	//bot.Debug = true
	collectors.Init()

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
