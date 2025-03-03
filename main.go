package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Загрузка шаблонов
	r.LoadHTMLGlob("templates/*.html")

	// Статические файлы
	r.Static("/static", "./static")

	// Маршруты
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"Title": "АВАЮССТРОЙ | Главная",
		})
	})

	r.GET("/services", func(c *gin.Context) {
		c.HTML(200, "services.html", gin.H{
			"Title": "АВАЮССТРОЙ | Услуги",
		})
	})

	r.GET("/rent", func(c *gin.Context) {
		c.HTML(200, "rent.html", gin.H{
			"Title": "АВАЮССТРОЙ | Аренда",
		})
	})

	r.GET("/contacts", func(c *gin.Context) {
		c.HTML(200, "contacts.html", gin.H{
			"Title": "АВАЮССТРОЙ | Контакты",
		})
	})

	r.GET("/water_lowering", func(c *gin.Context) {
		c.HTML(200, "service_water_lowering.html", gin.H{
			"Title": "АВАЮССТРОЙ | Водопонижение",
		})
	})

	r.GET("/drainage", func(c *gin.Context) {
		c.HTML(200, "service_drainage.html", gin.H{
			"Title": "АВАЮССТРОЙ | Дренажные системы",
		})
	})

	r.GET("/storm_sewer", func(c *gin.Context) {
		c.HTML(200, "service_storm_sewer.html", gin.H{
			"Title": "АВАЮССТРОЙ | Ливневая канализация",
		})
	})

	r.GET("/plumbing", func(c *gin.Context) {
		c.HTML(200, "service_plumbing.html", gin.H{
			"Title": "АВАЮССТРОЙ | Водопровод",
		})
	})

	r.GET("/sewerage", func(c *gin.Context) {
		c.HTML(200, "service_sewerage.html", gin.H{
			"Title": "АВАЮССТРОЙ | Канализация",
		})
	})

	// Запуск сервера
	r.Run(":8088")
}
