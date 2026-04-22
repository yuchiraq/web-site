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

func registerRoutes(router *gin.Engine) error {
	return filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		normalizedPath := filepath.ToSlash(path)
		if strings.Contains(normalizedPath, "templates/static/") ||
			strings.HasSuffix(normalizedPath, "header.html") ||
			strings.HasSuffix(normalizedPath, "footer.html") ||
			strings.HasSuffix(normalizedPath, "head.html") ||
			strings.HasSuffix(normalizedPath, "phone_button.html") ||
			strings.HasSuffix(normalizedPath, "index.html") {
			return nil
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		route := strings.TrimPrefix(normalizedPath, "templates/")
		route = strings.TrimSuffix(route, ".html")
		route = "/" + route
		templateName := filepath.Base(path)

		router.GET(route, func(c *gin.Context) {
			data := pageData(route)

			if strings.HasPrefix(route, "/services") {
				data["Photos"] = getProjectPhotos()
			}

			c.HTML(http.StatusOK, templateName, data)
		})

		fmt.Printf("Registered route: %s\n", route)
		return nil
	})
}
