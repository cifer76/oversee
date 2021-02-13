package collectors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/cifer76/oversee/entity"
)

type FutuArticle struct {
	Title string `json:"title"`
	Link  string `json:"url"`
}

func NewFutuCollector(news chan entity.PieceOfNews) Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	// handler for futu news
	c.OnResponse(func(resp *colly.Response) {
		if !strings.Contains(strings.ToLower(resp.Headers.Get("Content-Type")), "json") {
			return
		}

		var rsp map[string]json.RawMessage
		json.Unmarshal(resp.Body, &rsp)
		json.Unmarshal(rsp["data"], &rsp)

		as := []FutuArticle{}
		json.Unmarshal(rsp["list"], &as)

		for _, a := range as {
			link := "https://t.me/iv?url=" + a.Link + "&rhash=66c8e82d1af8a5"
			news <- entity.PieceOfNews{
				Title:  a.Title,
				Link:   link,
				Source: "Futu",
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	return Collector{
		colly: c,
		site:  "https://news.futunn.com/client/market-list?lang=zh-cn",
	}
}
