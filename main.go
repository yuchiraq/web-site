package main

import (
	"log"
	"net/http"
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

	// Автоматическая регистрация маршрутов
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

	// // Запуск HTTPS-сервера на порту 443
	// go func() {
	// 	log.Println("Запуск HTTPS-сервера на порту 443")
	// 	if err := http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/avayusstroi.by/fullchain.pem", "/etc/letsencrypt/live/avayusstroi.by/privkey.pem", r); err != nil {
	// 		log.Fatalf("Ошибка запуска HTTPS-сервера: %v", err)
	// 	}
	// }()

	// // Запуск HTTP-сервера на порту 80 для перенаправления на HTTPS
	// log.Println("Запуск HTTP-сервера на порту 80 для перенаправления на HTTPS")
	// if err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// Перенаправляем все запросы на HTTPS
	// 	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
	// })); err != nil {
	// 	log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	// }
	log.Println("Запуск HTTP-сервера на порту 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
	}
}
