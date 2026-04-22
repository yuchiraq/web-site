package main

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
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

func GenerateSitemap() error {
	urls := make([]URL, 0)

	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		normalizedPath := filepath.ToSlash(path)
		if strings.Contains(normalizedPath, "templates/static/") ||
			strings.HasSuffix(normalizedPath, "header.html") ||
			strings.HasSuffix(normalizedPath, "footer.html") ||
			strings.HasSuffix(normalizedPath, "head.html") ||
			strings.HasSuffix(normalizedPath, "phone_button.html") ||
			strings.HasSuffix(normalizedPath, "404.html") {
			return nil
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		route := strings.TrimPrefix(normalizedPath, "templates/")
		route = strings.TrimSuffix(route, ".html")
		if route == "index" {
			route = "/"
		} else {
			route = "/" + route
		}

		priority := "0.5"
		changeFreq := "monthly"
		switch {
		case route == "/":
			priority = "1.0"
			changeFreq = "daily"
		case route == "/services" || route == "/rent" || route == "/contacts":
			priority = "0.8"
			changeFreq = "weekly"
		case strings.HasPrefix(route, "/services/"):
			priority = "0.7"
			changeFreq = "weekly"
		}

		urls = append(urls, URL{
			Loc:        canonicalURL(route),
			LastMod:    info.ModTime().Format("2006-01-02"),
			ChangeFreq: changeFreq,
			Priority:   priority,
		})

		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(urls, func(i, j int) bool {
		return urls[i].Loc < urls[j].Loc
	})

	urlSet := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	file, err := os.Create("static/data/sitemap.xml")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(xml.Header); err != nil {
		return err
	}

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	return encoder.Encode(urlSet)
}
