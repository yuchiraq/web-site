package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context) {
	title := "ЧСУП АВАЮССТРОЙ | Строительство в Бресте" // Title for the main page
	photos := getProjectPhotos()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":  title,
		"Photos": photos,
	})
}

func getProjectPhotos() []string {
	projectPhotosOnce.Do(func() {
		photoPath := "static/images/photos"
		files, err := os.ReadDir(photoPath)
		if err != nil {
			fmt.Println("Error reading photos directory:", err)
			projectPhotos = []string{}
			return
		}

		photos := make([]string, 0, len(files))
		for _, file := range files {
			if !file.IsDir() {
				photos = append(photos, filepath.ToSlash(filepath.Join("/", "static", "images", "photos", file.Name())))
			}
		}
		projectPhotos = photos
	})

	if len(projectPhotos) == 0 {
		return nil
	}

	photos := make([]string, len(projectPhotos))
	copy(photos, projectPhotos)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(photos), func(i, j int) {
		photos[i], photos[j] = photos[j], photos[i]
	})

	return photos
}
