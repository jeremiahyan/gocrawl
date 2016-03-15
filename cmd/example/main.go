package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jeremiahyan/gocrawl"
	"github.com/jeremiahyan/goquery"
	"github.com/jeremiahyan/goxlsx"
	"strings"
	"regexp"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	triesCount++

	html, _ := doc.Html()
	urlString := ctx.URL().String()
	fmt.Printf("Try %d: %s ", triesCount, urlString)
	if !strings.Contains(urlString, "reviews") {
		fmt.Printf(" Passed...\n")
		return nil, true;
	}
	fmt.Printf("\nVisit: %s\n\n", ctx.URL())

	prefixH2Author := "<h2 class=\"impact\">"
	prefixH3 := "<h3>"
	prefixH4 := "<h4>"
	prefixH4Author := "<h4 class=\"text-color-left\">"
	prefixP := "<p>"

	suffixH2 := "</h2>"
	suffixH3 := "</h3>"
	suffixH4 := "</h4>"
	suffixP := "</p>"

	var author, title, miniTitle, yearT, summary string
	author = ""

	arrayLv0 := strings.Split(html, "<a name=\"")
	for i0, contentLv0 := range arrayLv0 {
		if i0 > 0 {
			if !strings.Contains(contentLv0, suffixP) {
				if strings.Contains(contentLv0, prefixH2Author) {
					author = strings.Split(contentLv0, prefixH2Author)[1]
					author = strings.Split(author, suffixH2)[0]
				}
				continue
			}

			arrayLv1 := strings.Split(contentLv0, "</p>")
			paragraph := arrayLv1[0]
			fmt.Printf("HTML: %d: \n%s\n", i0, paragraph)

			miniTitle = strings.Split(paragraph, "\"></a>")[0]

			if strings.Contains(paragraph, prefixH4Author) {
				author = strings.Split(paragraph, prefixH4Author)[1]
				author = strings.Split(author, suffixH4)[0]
			}

			if author == "" {
				author = strings.Split(ctx.URL().String(), "http://www.chicklitclub.com/")[1]
				author = strings.Split(author, "-")[0]
			}

			if strings.Contains(paragraph, prefixH4) {
				title = strings.Split(paragraph, prefixH4)[1]
				title = strings.Split(title, suffixH4)[0]
			} else if strings.Contains(paragraph, prefixH3) {
				title = strings.Split(paragraph, prefixH3)[1]
				title = strings.Split(title, suffixH3)[0]
			}

			if strings.Contains(paragraph, prefixP) {
				summary = strings.Split(paragraph, prefixP)[1]
				summary = strings.Split(summary, suffixP)[0]
			}

			if strings.Contains(title, "(") && strings.Contains(title, ")") {
				//aYear := strings.Split(title, "(")[1]
				//aYear = strings.Split(aYear, ")")[0]
				reg := regexp.MustCompile(".*[0-9]{4}.*")
				yearT = reg.FindString(title)
			}

			rowNum++
			row := xlsxSheet.AddRow()
			numCell := row.AddCell()
			authorCell := row.AddCell()
			titleCell := row.AddCell()
			miniTitleCell := row.AddCell()
			yearCell := row.AddCell()
			summaryCell := row.AddCell()
			URLCell := row.AddCell()

			numCell.Value = fmt.Sprintf("%d", rowNum)
			authorCell.Value = author
			titleCell.Value = title
			miniTitleCell.Value = miniTitle
			yearCell.Value = yearT
			summaryCell.Value = summary
			URLCell.Value = ctx.URL().String()

			fmt.Printf("\n======\nMiniTitle: %s\nAuthor: %s\nTitle: %s\nYear: %s\n", miniTitle, author, title, yearT)

			author = ""
			title = ""
			miniTitle = ""
			yearT = ""
			summary = ""

			saveXLSX()
		}
	}

	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	if ctx.URL().Host == "www.chicklitclub.com" {
		return true
	}
	return false
}

func createXLSX()  {
	var headerRow *xlsx.Row

	xlsxFile = xlsx.NewFile()
	xlsxSheet, err = xlsxFile.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	headerRow = xlsxSheet.AddRow()

	numCell := headerRow.AddCell()
	authorCell := headerRow.AddCell()
	titleCell := headerRow.AddCell()
	miniTitleCell := headerRow.AddCell()
	yearCell := headerRow.AddCell()
	summaryCell := headerRow.AddCell()
	URLCell := headerRow.AddCell()

	numCell.Value = "No."
	authorCell.Value = "Author"
	titleCell.Value = "Title"
	yearCell.Value = "Year"
	miniTitleCell.Value = "Mini Title"
	summaryCell.Value = "Summary"
	URLCell.Value = "URL"
}

func saveXLSX()  {
	err = xlsxFile.Save("CrawlResult.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

var xlsxFile *xlsx.File
var xlsxSheet *xlsx.Sheet
var rowNum = 0

var triesCount = 0
var err error

func main() {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 300 * time.Millisecond
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = false
	opts.MaxVisits = 100000

	createXLSX()

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("http://www.chicklitclub.com")

	saveXLSX()
}
