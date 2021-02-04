package collectors

import (
	"fmt"

	"github.com/gocolly/colly/v2"

	"github.com/cifer76/oversee/entity"
)

func NewSinaCollector(news chan entity.PieceOfNews) Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	// handler for the following sites:
	// https://rsshub.app/sina/finance
	// https://a.jiemian.com/index.php?m=article&a=rss
	c.OnXML("//item", func(e *colly.XMLElement) {
		// Print link
		title := e.ChildText("title")
		link := e.ChildText("link")
		news <- entity.PieceOfNews{
			Title:  title,
			Link:   link,
			Source: "Sina",
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	return Collector{
		colly: c,
		site:  "https://rsshub.app/sina/finance",
	}
}
