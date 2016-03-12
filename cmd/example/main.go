package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jeremiahyan/gocrawl"
	"github.com/jeremiahyan/goquery"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	html, _ := doc.Html()
	fmt.Printf("Visit: %s\n %s \n", ctx.URL(), html)
	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	if ctx.URL().Host == "chicklitclub.com" {
		return true
	}
	return false
}

func main() {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = false
	opts.MaxVisits = 100

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("http://chicklitclub.com")
}
