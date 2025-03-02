package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация Gin
	r := gin.Default()

	// Статические файлы
	r.Static("/static", "./static")

	// Маршруты для страниц
	r.GET("/", func(c *gin.Context) {
		c.File("templates/index.html")
	})

	r.GET("/services", func(c *gin.Context) {
		c.File("templates/services.html")
	})

	r.GET("/projects", func(c *gin.Context) {
		c.File("templates/projects.html")
	})

	r.GET("/contacts", func(c *gin.Context) {
		c.File("templates/contacts.html")
	})

	// Запуск сервера
	r.Run(":8088")
}
