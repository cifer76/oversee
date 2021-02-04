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
