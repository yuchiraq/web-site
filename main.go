package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация Gin
	r := gin.Default()

	gin.SetMode(gin.ReleaseMode) // Устанавливаем режим release

	// Настройка логирования в файл
	logFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Генерация sitemap.xml перед запуском сервера
	if err := GenerateSitemap(); err != nil {
		log.Fatalf("Ошибка генерации sitemap.xml: %v", err)
	}
	log.Println("sitemap.xml успешно сгенерирован")

	// Загрузка шаблонов
	if err := loadTemplates(r); err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}

	// Обработчик для главной страницы
	r.GET("/", indexHandler)

	// Автоматическая регистрация остальных маршрутов
	if err := registerRoutes(r); err != nil {
		log.Fatalf("Ошибка регистрации маршрутов: %v", err)
	}

	// Статические файлы
	r.Static("/static", "./static")
	r.StaticFile("/sitemap.xml", "./static/data/sitemap.xml")
	r.StaticFile("/robots.txt", "./static/data/robots.txt")
	r.StaticFile("/favicon.ico", "./static/images/favicon.ico")
	r.StaticFile("/favicon.svg", "./static/images/logo_n.svg")
	r.StaticFile("/favicon", "./static/images/logo_n.svg")

	// Обработчик формы
	r.POST("/submit", submitHandler)
	// RSS-фид
	r.GET("/rss.xml", generateRSSFeed)

	// Обработка 404
	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{
			"Title": "404 - Страница не найдена",
		})
	})

	log.Println("Запуск HTTP-сервера на порту 8088")
	if err := r.Run(":8088"); err != nil {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	}
}
