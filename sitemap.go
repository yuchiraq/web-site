package main

import (
	"encoding/xml"
	"os"
)

type URL struct {
	Loc string `xml:"loc"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

func GenerateSitemap() error {
	urls := make([]URL, 0, len(publicPagePaths()))
	for _, route := range publicPagePaths() {
		urls = append(urls, URL{Loc: canonicalURL(route)})
	}

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
