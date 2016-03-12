package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jeremiahyan/gocrawl"
	"github.com/jeremiahyan/goquery"
	"strings"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	html, _ := doc.Html()
	urlString := ctx.URL().String()
	fmt.Printf("Try: %s ", urlString)
	if !strings.Contains(urlString, "reviews") {
		fmt.Printf(" Passed...\n")
		return nil, true;
	}
	fmt.Printf("\nVisit: %s\n\n", ctx.URL())
	arrayLv0 := strings.Split(html, "<a name=\"")
	for i0, contentLv0 := range arrayLv0 {
		if i0 > 0 && strings.Contains(contentLv0, "</p>") {
			arrayLv1 := strings.Split(contentLv0, "</p>")
			paragraph := arrayLv1[0]
			fmt.Printf("HTML: %d: \n%s\n", i0, paragraph)

			var miniTitle, title string
			author := ""
			miniTitle = strings.Split(paragraph, "\"></a>")[0]

			if strings.Contains(paragraph, "<h4 class=\"text-color-left\">") {
				author = strings.Split(paragraph, "<h4 class=\"text-color-left\">")[1]
				author = strings.Split(author, "</h4>")[0]
			}

			if strings.Contains(paragraph, "<h4>") {
				title = strings.Split(paragraph, "<h4>")[1]
				title = strings.Split(title, "</h4>")[0]
			} else if strings.Contains(paragraph, "<h3>") {
				title = strings.Split(paragraph, "<h3>")[1]
				title = strings.Split(title, "</h3>")[0]
			}

			fmt.Printf("\n======\nMiniTitle: %s\nAuthor: %s\nTitle: %s\n", miniTitle, author, title)
		}
	}

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
