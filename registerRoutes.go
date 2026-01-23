package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	engine.LoadHTMLFiles(files...)
	return nil
}

// Функция для автоматической регистрации маршрутов
func registerRoutes(router *gin.Engine) error {
	// Карта заголовков для SEO
	titleMap := map[string]string{
		"/":                        "ЧСУП АВАЮССТРОЙ | Строительство в Бресте",
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
		// Пропускаем вложенные шаблоны (header, footer и т.д.)
		if strings.Contains(path, "templates/static") || strings.Contains(path, "header.html") ||
			strings.Contains(path, "footer.html") || strings.Contains(path, "head.html") ||
			strings.Contains(path, "phone_button.html") {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			route := strings.TrimPrefix(path, "templates/")
			route = strings.TrimSuffix(route, ".html")
			if route == "index" {
				route = "/"
			} else {
				route = "/" + route
			}
			// Регистрируем маршрут
			router.GET(route, func(c *gin.Context) {
				title := titleMap[route]
				if title == "" {
					title = "ЧСУП АВАЮССТРОЙ"
				}

				if route == "/" {
					// Get photos for the carousel
					photoPath := "static/images/photos"
					files, err := os.ReadDir(photoPath)
					if err != nil {
						// Log the error and continue without photos
						fmt.Println("Error reading photos directory:", err)
						c.HTML(http.StatusOK, filepath.Base(path), gin.H{
							"Title": title,
						})
						return
					}

					var photos []string
					for _, file := range files {
						if !file.IsDir() {
							// Use filepath.ToSlash for cross-platform compatibility
							photos = append(photos, filepath.ToSlash(filepath.Join("/", "static", "images", "photos", file.Name())))
						}
					}

					// Shuffle photos
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					r.Shuffle(len(photos), func(i, j int) {
						photos[i], photos[j] = photos[j], photos[i]
					})

					c.HTML(http.StatusOK, filepath.Base(path), gin.H{
						"Title":  title,
						"Photos": photos,
					})
				} else {
					c.HTML(http.StatusOK, filepath.Base(path), gin.H{
						"Title": title,
					})
				}
			})
			fmt.Printf("Зарегистрирован маршрут: %s
", route)
		}
		return nil
	})
	return err
}