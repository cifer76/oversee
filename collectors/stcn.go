package collectors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/cifer76/oversee/entity"
)

func sendNews(news chan entity.PieceOfNews, title, link string) {
	title = strings.TrimSpace(title)
	news <- entity.PieceOfNews{
		Title:  title,
		Link:   link,
		Source: "Stcn",
	}
}

func NewStcnCollector(news chan entity.PieceOfNews) Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		//colly.Debugger(&debug.LogDebugger{}),
	)

	// hot news
	c.OnHTML("div.maj_left a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// caijing focused news
	c.OnHTML("div.maj_right div.caijing p a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// friend media left
	c.OnHTML("div.box div.left ul li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// friend media right
	c.OnHTML("div.box div.maj ul li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// clude sections
	c.OnHTML("div.box dl.clude dd a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - stock
	c.OnHTML("div.box div.gushi ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - data
	c.OnHTML("div.box div.shuju ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - company
	c.OnHTML("div.box div.gongsi ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - deep
	c.OnHTML("div.box div.shendu ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - orgnazition
	c.OnHTML("div.box div.jigou ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - fund
	c.OnHTML("div.box div.jijin ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - chuangtou
	c.OnHTML("div.box div.chuangtou ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - shengufaxing
	c.OnHTML("div.box div.shengufaxing ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})
	// section - zhuanti
	c.OnHTML("div.box div.zhuanti ul.list li a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		sendNews(news, e.Text, link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	return Collector{
		colly: c,
		site:  "https://www.stcn.com",
	}
}
