package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация Gin
	r := gin.Default()

	// Загрузка шаблонов (включая templates/static)
	if err := loadTemplates(r); err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}

	// Автоматическая регистрация маршрутов (исключая templates/static)
	if err := registerRoutes(r); err != nil {
		log.Fatalf("Ошибка регистрации маршрутов: %v", err)
	}

	// Статические файлы
	r.Static("/static", "./static")
	r.StaticFile("/sitemap.xml", "./static/data/sitemap.xml")
	r.StaticFile("/robots.txt", "./static/data/robots.txt")
	r.StaticFile("/favicon.ico", "./static/images/favicon.ico")

	// Запуск сервера
	log.Println("Сервер запущен на http://localhost:80")
	if err := r.Run(":80"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
