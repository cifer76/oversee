package collectors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/cifer76/oversee/entity"
)

func sendZqrbNews(news chan entity.PieceOfNews, title, link string) {
	title = strings.TrimSpace(title)
	link = "https://t.me/iv?url=" + link + "&rhash=c8ed79f38e16c2"
	news <- entity.PieceOfNews{
		Title:  title,
		Link:   link,
		Source: "Zqrb",
	}
}

func NewZqrbCollector(news chan entity.PieceOfNews) Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.OnHTML("div.first-nleft1 ul li a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})
	c.OnHTML("div.first-nleft2 ul li a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})

	c.OnHTML("div.first-left1 div.focusx ul li a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})

	c.OnHTML("div.first-left2 div.focusx ul li a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})

	c.OnHTML("div.first-left1 p a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})
	c.OnHTML("div.first-left2 p a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})

	c.OnHTML("div.third-left1 p a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})
	c.OnHTML("div.third-left2 p a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})
	c.OnHTML("div.third-left3 p a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendZqrbNews(news, e.Text, link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	return Collector{
		colly: c,
		site:  "http://www.zqrb.cn",
	}
}
