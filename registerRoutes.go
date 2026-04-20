package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	projectPhotos     []string
	projectPhotosOnce sync.Once
)

// Функция для загрузки шаблонов
func loadTemplates(engine *gin.Engine) error {
	var files []string
	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	sort.Strings(files)
	engine.LoadHTMLFiles(files...)
	return nil
}

// Функция для автоматической регистрации маршрутов
func registerRoutes(router *gin.Engine) error {
	// Карта заголовков для SEO
	titleMap := map[string]string{
		"/services":                "Услуги | ЧСУП АВАЮССТРОЙ",
		"/rent":                    "Аренда техники | ЧСУП АВАЮССТРОЙ",
		"/contacts":                "Контакты | ЧСУП АВАЮССТРОЙ",
		"/services/drainage":       "Дренаж | ЧСУП АВАЮССТРОЙ",
		"/services/plumbing":       "Водоснабжение | ЧСУП АВАЮССТРОЙ",
		"/services/sewerage":       "Канализация | ЧСУП АВАЮССТРОЙ",
		"/services/storm_sewer":    "Ливневая канализация | ЧСУП АВАЮССТРОЙ",
		"/services/topas":          "Септики TOPAS | ЧСУП АВАЮССТРОЙ",
		"/services/water_lowering": "Водопонижение | ЧСУП АВАЮССТРОЙ",
	}

	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		normalizedPath := filepath.ToSlash(path)
		// Пропускаем вложенные шаблоны (header, footer и т.д.) и index.html
		if strings.Contains(normalizedPath, "templates/static/") || strings.HasSuffix(normalizedPath, "header.html") ||
			strings.HasSuffix(normalizedPath, "footer.html") || strings.HasSuffix(normalizedPath, "head.html") ||
			strings.HasSuffix(normalizedPath, "phone_button.html") || strings.HasSuffix(normalizedPath, "index.html") {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			route := strings.TrimPrefix(normalizedPath, "templates/")
			route = strings.TrimSuffix(route, ".html")
			route = "/" + route
			templateName := filepath.Base(path)

			// Регистрируем маршрут
			router.GET(route, func(c *gin.Context) {
				title := titleMap[route]
				if title == "" {
					title = "ЧСУП АВАЮССТРОЙ"
				}

				data := gin.H{
					"Title": title,
				}
				if strings.HasPrefix(route, "/services") {
					data["Photos"] = getProjectPhotos()
				}

				c.HTML(http.StatusOK, templateName, data)
			})
			fmt.Printf("Зарегистрирован маршрут: %s\n", route)
		}
		return nil
	})
	return err
}
