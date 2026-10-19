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
	data := pageData("/")
	photos := getProjectPhotos()
	primaryPhotos, secondaryPhotos := splitProjectPhotos(photos)

	data["Photos"] = primaryPhotos
	data["PhotosSecondary"] = secondaryPhotos

	c.HTML(http.StatusOK, "index.html", data)
}

func splitProjectPhotos(photos []string) ([]string, []string) {
	if len(photos) < 2 {
		return photos, photos
	}

	primary := make([]string, 0, (len(photos)+1)/2)
	secondary := make([]string, 0, len(photos)/2)
	for index, photo := range photos {
		if index%2 == 0 {
			primary = append(primary, photo)
			continue
		}
		secondary = append(secondary, photo)
	}

	return primary, secondary
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
