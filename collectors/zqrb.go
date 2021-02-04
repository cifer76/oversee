package collectors

import (
	"fmt"

	"github.com/gocolly/colly/v2"

	"github.com/cifer76/oversee/entity"
)

func NewZqrbCollector(news chan entity.PieceOfNews) Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	// handler for the following sites:
	// http://www.zqrb.cn
	c.OnHTML("p a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		news <- entity.PieceOfNews{
			Title:  e.Text,
			Link:   link,
			Source: "Zqrb",
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	return Collector{
		colly: c,
		site:  "http://www.zqrb.cn",
	}
}
