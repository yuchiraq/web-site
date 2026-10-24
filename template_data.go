package main

import "github.com/gin-gonic/gin"

func pageData(path string) gin.H {
	seo := seoDataForPath(path)
	return gin.H{
		"Title":          seo.Title,
		"SEO":            seo,
		"SiteMeta":       currentSiteMeta(),
		"StructuredData": structuredDataForPath(path, seo),
		"Rental":         rentalPageForPath(path),
		"Rentals":        rentalPages,
		"ServiceAdvice":  serviceAdviceByPath[path],
		"Topic":          topicPageForPath(path),
		"ServiceTopic":   serviceTopics[path],
	}
}
