package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Функция для загрузки шаблонов (включая templates/static)
func loadTemplates(engine *gin.Engine) error {
	var files []string
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Добавляем только HTML-файлы
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	engine.LoadHTMLFiles(files...)
	return nil
}

// Функция для автоматической регистрации маршрутов (исключая templates/static)
func registerRoutes(router *gin.Engine) error {
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Исключаем папку static
		if strings.Contains(path, "templates/static") {
			return nil
		}
		// Регистрируем маршруты только для HTML-файлов
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			// Убираем "templates/" и расширение .html
			route := strings.TrimPrefix(path, "templates/")
			route = strings.TrimSuffix(route, ".html")
			// Делаем index.html корневым маршрутом
			if route == "index" {
				route = "/"
			} else {
				route = "/" + route
			}
			// Регистрируем маршрут
			router.GET(route, func(c *gin.Context) {
				c.HTML(http.StatusOK, filepath.Base(path), gin.H{
					"Title": strings.Title(strings.ReplaceAll(route, "/", " ")),
				})
			})
			fmt.Printf("Зарегистрирован маршрут: %s\n", route)
		}
		return nil
	})
	return err
}
