package main

import (
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Turbo   string     `xml:"xmlns:turbo,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Turbo       string   `xml:"turbo,attr"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description"`
	Content     RSSCDATA `xml:"turbo:content"`
}

type RSSCDATA struct {
	HTML string `xml:",cdata"`
}

func generateRSSFeed(c *gin.Context) {
	feed, err := buildRSSFeed()
	if err != nil {
		c.Header("X-Robots-Tag", "noindex")
		c.String(http.StatusInternalServerError, "Не удалось сформировать RSS")
		_ = c.Error(err)
		return
	}
	encoded, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось сформировать RSS")
		return
	}
	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", append([]byte(xml.Header), encoded...))
}

func buildRSSFeed() (RSSFeed, error) {
	feed := RSSFeed{
		Version: "2.0", Turbo: "http://turbo.yandex.ru",
		Channel: RSSChannel{Title: siteName + " — услуги и спецтехника", Link: siteURL + "/", Description: "Инженерные работы и аренда спецтехники в Бресте и Брестской области", Language: "ru-BY"},
	}
	for _, path := range publicPagePaths() {
		if !strings.HasPrefix(path, "/services/") && !strings.HasPrefix(path, "/rent/") {
			continue
		}
		seo := seoDataForPath(path)
		content := "<header><h1>" + html.EscapeString(seo.Title) + "</h1><figure><img src=\"" + html.EscapeString(seo.Image) + "\"></figure></header>"
		if rental := rentalPageForPath(path); rental != nil {
			content += "<p>" + html.EscapeString(rental.Lead) + "</p><h2>Задачи и условия аренды</h2><ul>"
			for _, text := range append(append([]string{}, rental.Tasks...), rental.Conditions...) {
				content += "<li>" + html.EscapeString(text) + "</li>"
			}
			content += "</ul>"
		} else {
			file, err := os.ReadFile(templateForPath(path))
			if err != nil {
				return feed, fmt.Errorf("read RSS source %s: %w", path, err)
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(file)))
			if err != nil {
				return feed, fmt.Errorf("parse RSS source %s: %w", path, err)
			}
			doc.Find("main h2, main h3, main p, main li").Each(func(_ int, selection *goquery.Selection) {
				text := strings.TrimSpace(selection.Text())
				if text == "" || strings.Contains(text, "{{") {
					return
				}
				tag := "p"
				if selection.Is("h2,h3") {
					tag = "h2"
				}
				content += "<" + tag + ">" + html.EscapeString(text) + "</" + tag + ">"
			})
		}
		if advice, ok := serviceAdviceByPath[path]; ok {
			content += "<h2>" + html.EscapeString(advice.TasksHeading) + "</h2>"
			for _, text := range append(append([]string{}, advice.Tasks...), advice.Steps...) {
				content += "<p>" + html.EscapeString(text) + "</p>"
			}
			content += "<h2>Расчёт стоимости</h2><p>" + html.EscapeString(advice.PriceFactors) + "</p><p>" + html.EscapeString(advice.RequestData) + "</p>"
		}
		content += "<p><a href=\"" + seo.Canonical + "\">Подробнее об услуге и заявка</a></p><p><a href=\"tel:+375297202791\">+375 (29) 720-27-91</a></p>"
		feed.Channel.Items = append(feed.Channel.Items, RSSItem{Turbo: "true", Title: seo.Title, Link: seo.Canonical, GUID: seo.Canonical, Description: seo.Description, Content: RSSCDATA{HTML: content}})
	}
	return feed, nil
}
