package main

import (
	"encoding/xml"
	"os"
	"time"
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
	// Список URL вашего сайта
	urls := []URL{
		{
			Loc:        "https://avayusstroi.by/",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "daily",
			Priority:   "1.0",
		},
		{
			Loc:        "https://avayusstroi.by/services",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.8",
		},
		{
			Loc:        "https://avayusstroi.by/rent",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.8",
		},
		{
			Loc:        "https://avayusstroi.by/contacts",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   "0.5",
		},
		{
			Loc:        "https://avayusstroi.by/services/drainage",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		},
		{
			Loc:        "https://avayusstroi.by/services/plumbing",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		},
		{
			Loc:        "https://avayusstroi.by/services/sewerage",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		},
		{
			Loc:        "https://avayusstroi.by/services/storm_sewer",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		},
		{
			Loc:        "https://avayusstroi.by/services/topas",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		},
		{
			Loc:        "https://avayusstroi.by/services/water_lowering",
			LastMod:    time.Now().Format("2006-01-02"),
			ChangeFreq: "weekly",
			Priority:   "0.7",
		},
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
