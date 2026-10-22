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
	for _, route := range publicPagePaths() {
		if route == "/" {
			continue
		}
		templatePath := templateForPath(route)
		if _, err := os.Stat(templatePath); err != nil {
			return fmt.Errorf("public page %s: %w", route, err)
		}
		templateName := filepath.Base(templatePath)
		handler := func(c *gin.Context) {
			data := pageData(route)

			if strings.HasPrefix(route, "/services") {
				photos := getProjectPhotos()
				primaryPhotos, secondaryPhotos := splitProjectPhotos(photos)
				data["Photos"] = primaryPhotos
				data["PhotosSecondary"] = secondaryPhotos
			}

			c.HTML(http.StatusOK, templateName, data)
		}
		router.GET(route, handler)
		router.HEAD(route, handler)

		fmt.Printf("Registered route: %s\n", route)
	}
	return nil
}
