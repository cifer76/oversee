package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/patrickmn/go-cache"
)

var (
	interestedWords []string
	sentLinks       = cache.New(60*time.Minute, 60*time.Second)
)

func requestInterestedWords() []string {
	gistURL := "https://gist.githubusercontent.com/cifer76/7f14f5bd02b98abbc11d16662266a572/raw/81c942b968c06f3e1c74418d569d2a37275e20ee/interestedWords"

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

func main() {

	// heroku requires us to listen on a port
	go func() {
		port := os.Getenv("PORT")
		http.ListenAndServe("0.0.0.0:"+port, nil)
	}()

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Debugger(&debug.LogDebugger{}),
	)

	bot, err := tgbotapi.NewBotAPI("407954143:AAGDxLmxcr5DGVE3GY_Ih9pe8GIh-P0EhDI")
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true

	// Find and visit all links
	c.OnHTML(".news-li", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		// e.Request.Visit(e.Attr("href"))
	})

	// handler for:
	// https://rsshub.app/sina/finance
	// https://a.jiemian.com/index.php?m=article&a=rss
	c.OnXML("//item", func(e *colly.XMLElement) {
		// Print link
		title := e.ChildText("title")
		link := e.ChildText("link")
		//pubDate := e.ChildText("pubDate")
		desc := e.ChildText("description")

		fmt.Printf("Article found: %s -> %s\n", title, link)
		//fmt.Printf("\t\t Published on: %s\n", pubDate)
		//fmt.Printf("\t\t Brief: %s\n", desc)

		if _, found := sentLinks.Get(link); found {
			return
		}

		for _, w := range interestedWords {
			// avoid sending duplicate articles

			if strings.Contains(title, w) || strings.Contains(desc, w) {
				content := fmt.Sprintf("[%s](%s)", title, link)
				msg := tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:           406797693,
						ReplyToMessageID: 0,
					},
					Text:                  content,
					ParseMode:             "MarkdownV2",
					DisableWebPagePreview: false,
				}
				bot.Send(msg)

				sentLinks.Set(link, true, 0)

				break
			}
		}
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		// e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	for {
		interestedWords = requestInterestedWords()
		if err := c.Visit("https://a.jiemian.com/index.php?m=article&a=rss"); err != nil {
			fmt.Printf("%v\n", err)
		}
		if err := c.Visit("https://rsshub.app/sina/finance"); err != nil {
			fmt.Printf("%v\n", err)
		}
		time.Sleep(10 * time.Minute)
	}
}
