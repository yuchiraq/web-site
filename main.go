package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "Главная | АВАЮССТРОЙ"})
	})

	r.GET("/services", func(c *gin.Context) {
		c.HTML(http.StatusOK, "services.html", gin.H{"title": "Услуги | АВАЮССТРОЙ"})
	})

	r.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.html", gin.H{"title": "Контакты | АВАЮССТРОЙ"})
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.html", gin.H{"title": "О нас | АВАЮССТРОЙ"})
	})

	r.Run(":8088")
}
