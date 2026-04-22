package main

import "github.com/gin-gonic/gin"

func pageData(path string) gin.H {
	seo := seoDataForPath(path)
	return gin.H{
		"Title":    seo.Title,
		"SEO":      seo,
		"SiteMeta": currentSiteMeta(),
	}
}
