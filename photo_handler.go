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

	// Get photos for the carousel
	photoPath := "static/images/photos"
	files, err := os.ReadDir(photoPath)
	if err != nil {
		// Log the error and continue without photos
		fmt.Println("Error reading photos directory:", err)
		c.HTML(http.StatusOK, "index.html", gin.H{
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

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":  title,
		"Photos": photos,
	})
}
