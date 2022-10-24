package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
)

var (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	headers   = map[string]string{"user-agent": userAgent, "Cookie": "SUB=_2AkMWIuNSf8NxqwJRmP8dy2rhaoV2ygrEieKgfhKJJRMxHRl-yT9jqk86tRB6PaLNvQZR6zYUcYVT1zSjoSreQHidcUq7"}
)

type News struct {
	Title  string `json:"title,omitempty"`
	Url    string `json:"url,omitempty"`
	Hotnum string `json:"hotnum,omitempty"`
}

func RestyFetch(url string) {
	client := resty.New()
	r, err := client.R().SetHeaders(headers).Get(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(r.String()))
	dom.Find("td.td-02").Each(func(i int, selection *goquery.Selection) {
		parse(selection)
	})
}

func CollyFetch(url string) {
	c := colly.NewCollector(colly.UserAgent(userAgent))
	file, _ := os.Create("weibo.csv")
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	_ = writer.Write([]string{"title", "url", "hotnum"})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", headers["Cookie"])
		r.Headers.Set("user-agent", userAgent)
		color.Red("Visiting：%s", r.URL.String())
	})

	c.OnHTML("td.td-02", func(e *colly.HTMLElement) {
		e.DOM.Each(func(i int, selection *goquery.Selection) {
			news := parse(selection)
			if news != nil {
				_ = writer.Write([]string{news.Title, news.Url, news.Hotnum})
			}
		})
	})

	_ = c.Visit(url)
}

func parse(selection *goquery.Selection) (news *News) {
	title := selection.Find("a").Text()
	href := selection.Find("a").AttrOr("href", "")
	hotnum := strings.TrimSpace(selection.Find("span").Text())
	if !strings.HasPrefix(href, "javascript:void") && hotnum != "" {
		news := News{
			Title:  title,
			Url:    href,
			Hotnum: hotnum,
		}
		jsonStr, _ := json.MarshalIndent(&news, "", " ")
		color.Green(string(jsonStr))
		return &news
	}
	return nil
}

func main() {
	url := "https://s.weibo.com/top/summary?cate=realtimehot"
	CollyFetch(url)
	RestyFetch(url)
}
