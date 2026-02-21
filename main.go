package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // Устанавливаем режим release
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), cacheControlMiddleware())

	// Настройка логирования в файл
	logFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка открытия файла логов: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if err := loadSettings(); err != nil {
		log.Printf("Предупреждение: settings.json не загружен: %v", err)
	}

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

	server := &http.Server{
		Addr:              ":8088",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Запуск HTTP-сервера на порту 8088")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	}
}

func cacheControlMiddleware() gin.HandlerFunc {
	const oneYearInSeconds = 31536000

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		switch {
		case path == "/sitemap.xml" || path == "/robots.txt" || path == "/rss.xml":
			c.Header("Cache-Control", "public, max-age=3600")
		case len(path) > 7 && path[:7] == "/static":
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
			c.Header("Expires", time.Now().Add(oneYearInSeconds*time.Second).UTC().Format(http.TimeFormat))
		default:
			c.Header("Cache-Control", "no-cache")
		}

		c.Next()
	}
}
