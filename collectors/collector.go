package collectors

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"

	"github.com/cifer76/oversee/entity"
)

type Collector struct {
	colly    *colly.Collector
	site     string
	interval int // in seconds
}

var (
	collectors []Collector
	news       chan entity.PieceOfNews
)

func Visit() chan entity.PieceOfNews {
	for _, c := range collectors {
		go func(c Collector) {
			for {
				if err := c.colly.Visit(c.site); err != nil {
					fmt.Printf("%v\n", err)
				}

				intv := 300
				if c.interval != 0 {
					intv = c.interval
				}
				time.Sleep(time.Duration(intv) * time.Second)
			}
		}(c)
	}
	return news
}

func Init() {
	news = make(chan entity.PieceOfNews, 1024)

	//collectors = append(collectors, NewZqrbCollector(news))
	collectors = append(collectors, NewSinaCollector(news))
	// collectors = append(collectors, NewStcnCollector(news))
	// collectors = append(collectors, NewJiemianCollector(news))
	// collectors = append(collectors, NewFutuCollector(news))
}
