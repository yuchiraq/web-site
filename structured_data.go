package main

import (
	"encoding/json"
	"html/template"
	"strings"
)

func structuredDataForPath(path string, seo SEOData) template.JS {
	if seo.Canonical == "" {
		return ""
	}
	organizationID := siteURL + "/#organization"
	graph := []any{
		map[string]any{
			"@type": "LocalBusiness", "@id": organizationID,
			"name": siteName, "legalName": "ЧСУП «АВАЮССТРОЙ»", "alternateName": "АВАЮССТРОЙ",
			"url": siteURL + "/", "logo": siteURL + "/static/images/logo_n.png",
			"image": siteURL + "/static/images/visitka.jpg", "email": "avajsstroi@mail.ru",
			"telephone":  []string{"+375297202791", "+375336709273"},
			"address":    map[string]any{"@type": "PostalAddress", "streetAddress": "ул. Янки Купалы, 109", "addressLocality": "Брест", "postalCode": "224032", "addressCountry": "BY"},
			"areaServed": []string{"Брест", "Брестский район", "Брестская область"},
			"sameAs":     []string{"https://www.instagram.com/avayusstroi/"},
		},
		map[string]any{
			"@type": "WebSite", "@id": siteURL + "/#website", "url": siteURL + "/",
			"name": siteName, "alternateName": []string{"АВАЮССТРОЙ", "Avayusstroi"},
			"inLanguage": "ru-BY", "publisher": map[string]string{"@id": organizationID},
		},
		map[string]any{
			"@type": seo.PageType, "@id": seo.Canonical + "#webpage", "url": seo.Canonical,
			"name": seo.Title, "description": seo.Description, "inLanguage": "ru-BY",
			"isPartOf":           map[string]string{"@id": siteURL + "/#website"},
			"about":              map[string]string{"@id": organizationID},
			"primaryImageOfPage": map[string]string{"@type": "ImageObject", "url": seo.Image, "caption": seo.ImageAlt},
		},
	}

	if path == "/services" || path == "/rent" {
		items := []any{}
		for _, publicPath := range publicPagePaths() {
			if path == "/services" && !strings.HasPrefix(publicPath, "/services/") {
				continue
			}
			if path == "/rent" && !strings.HasPrefix(publicPath, "/rent/") {
				continue
			}
			itemSEO := seoDataForPath(publicPath)
			items = append(items, map[string]any{"@type": "ListItem", "position": len(items) + 1, "url": itemSEO.Canonical, "name": itemSEO.Title})
		}
		graph = append(graph, map[string]any{"@type": "ItemList", "@id": seo.Canonical + "#catalog", "name": seo.Title, "itemListElement": items})
	}

	if path == "/services" || path == "/contacts" {
		name := "Услуги"
		if path == "/contacts" {
			name = "Контакты"
		}
		graph = append(graph, breadcrumbData([]string{"Главная", name}, []string{"/", path}))
	}

	if rental := rentalPageForPath(path); rental != nil {
		graph = append(graph,
			map[string]any{
				"@type": "Service", "@id": seo.Canonical + "#service", "url": seo.Canonical,
				"name": rental.Heading, "serviceType": rental.Heading, "description": rental.Lead,
				"image": seo.Image, "provider": map[string]string{"@id": organizationID},
				"areaServed": []string{"Брест", "Брестская область"},
			},
			breadcrumbData([]string{"Главная", "Аренда спецтехники", rental.Name}, []string{"/", "/rent", path}),
		)
		questions := []any{}
		for _, faq := range rental.FAQ {
			questions = append(questions, map[string]any{"@type": "Question", "name": faq.Question, "acceptedAnswer": map[string]string{"@type": "Answer", "text": faq.Answer}})
		}
		graph = append(graph, map[string]any{"@type": "FAQPage", "@id": seo.Canonical + "#faq", "mainEntity": questions})
	}
	if advice, ok := serviceAdviceByPath[path]; ok {
		questions := []any{}
		for _, faq := range advice.FAQ {
			questions = append(questions, map[string]any{"@type": "Question", "name": faq.Question, "acceptedAnswer": map[string]string{"@type": "Answer", "text": faq.Answer}})
		}
		graph = append(graph, map[string]any{"@type": "FAQPage", "@id": seo.Canonical + "#faq", "mainEntity": questions})
	}

	encoded, err := json.Marshal(map[string]any{"@context": "https://schema.org", "@graph": graph})
	if err != nil {
		return ""
	}
	// Every value is encoded by encoding/json, including HTML-sensitive characters.
	return template.JS(encoded)
}

func breadcrumbData(names, paths []string) map[string]any {
	items := make([]any, 0, len(names))
	for index, name := range names {
		items = append(items, map[string]any{"@type": "ListItem", "position": index + 1, "name": name, "item": canonicalURL(paths[index])})
	}
	return map[string]any{"@type": "BreadcrumbList", "itemListElement": items}
}
