package collectors

import (
	"fmt"

	"github.com/gocolly/colly/v2"

	"github.com/cifer76/oversee/entity"
)

func NewJiemianCollector(news chan entity.PieceOfNews) Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnXML("//item", func(e *colly.XMLElement) {
		// Print link
		title := e.ChildText("title")
		link := e.ChildText("link")
		link = "https://t.me/iv?url=" + link + "&rhash=c8ed79f38e16c2"
		news <- entity.PieceOfNews{
			Title:  title,
			Link:   link,
			Source: "Jiemian",
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	return Collector{
		colly: c,
		site:  "https://a.jiemian.com/index.php?m=article&a=rss",
	}
}
