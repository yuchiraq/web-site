package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// GenerateSitemap создаёт или обновляет sitemap.xml
func GenerateSitemap() error {
	urls := []URL{}

	// Сканируем папку templates
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Пропускаем вложенные шаблоны
		if strings.Contains(path, "header.html") || strings.Contains(path, "footer.html") ||
			strings.Contains(path, "head.html") || strings.Contains(path, "phone_button.html") {
			return nil
		}
		if strings.Contains(path, "templates/static") || strings.Contains(path, "404.html") {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			route := strings.TrimPrefix(path, "templates/")
			route = strings.TrimSuffix(route, ".html")
			if route == "index" {
				route = "/"
			} else {
				route = "/" + route
			}
			priority := "0.5"
			changeFreq := "monthly"
			if route == "/" {
				priority = "1.0"
				changeFreq = "daily"
			} else if strings.HasPrefix(route, "/services") {
				priority = "0.7"
				changeFreq = "weekly"
			}
			// Получаем дату последнего изменения файла
			lastMod := info.ModTime().Format("2006-01-02")
			urls = append(urls, URL{
				Loc:        "https://avayusstroi.by" + route,
				LastMod:    lastMod,
				ChangeFreq: changeFreq,
				Priority:   priority,
			})
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Создаём структуру для sitemap
	urlSet := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	// Создаём файл sitemap.xml
	file, err := os.Create("static/data/sitemap.xml")
	if err != nil {
		return err
	}
	defer file.Close()

	// Кодируем данные в XML и записываем в файл
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlSet); err != nil {
		return err
	}

	return nil
}
