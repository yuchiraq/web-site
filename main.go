package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.Use(gin.Logger(), gin.Recovery(), canonicalHostMiddleware(), cacheControlMiddleware())

	logFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if err := loadSettings(); err != nil {
		log.Printf("warning: settings.json was not loaded: %v", err)
	}

	qrScanTracker := NewQRScanTracker(defaultQRStatsPath)
	if err := qrScanTracker.Load(); err != nil {
		log.Printf("warning: QR statistics were not loaded and counting is disabled: %v", err)
	}

	if err := GenerateSitemap(); err != nil {
		log.Fatalf("failed to generate sitemap.xml: %v", err)
	}
	log.Println("sitemap.xml generated")

	if err := loadTemplates(r); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	r.GET("/", indexHandler)
	r.HEAD("/", indexHandler)
	r.GET("/qr", qrRedirectHandler(qrScanTracker))

	if err := registerRoutes(r); err != nil {
		log.Fatalf("failed to register routes: %v", err)
	}

	r.Static("/static", "./static")
	r.StaticFile("/sitemap.xml", "./static/data/sitemap.xml")
	r.StaticFile("/robots.txt", "./static/data/robots.txt")
	r.StaticFile("/favicon.ico", "./static/images/favicon.ico")
	r.StaticFile("/favicon.svg", "./static/images/logo_n.svg")
	r.StaticFile("/favicon", "./static/images/logo_n.svg")

	r.POST("/submit", submitHandler)
	r.GET("/rss.xml", generateRSSFeed)

	r.NoRoute(func(c *gin.Context) {
		c.Header("X-Robots-Tag", "noindex, follow")
		c.HTML(http.StatusNotFound, "404.html", pageData("/404"))
	})

	server := &http.Server{
		Addr:              ":8088",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("starting HTTP server on :8088")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

func canonicalHostMiddleware() gin.HandlerFunc {
	publicPaths := make(map[string]bool)
	for _, path := range publicPagePaths() {
		publicPaths[path] = true
	}
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Next()
			return
		}
		host := strings.ToLower(c.Request.Host)
		if colonIndex := strings.Index(host, ":"); colonIndex != -1 {
			host = host[:colonIndex]
		}

		path := c.Request.URL.Path
		normalizedPath := strings.TrimSuffix(strings.ToLower(path), "/")
		normalizedPath = strings.TrimSuffix(normalizedPath, ".html")
		if normalizedPath == "" || normalizedPath == "/index" {
			normalizedPath = "/"
		}
		if !publicPaths[normalizedPath] && normalizedPath != "/qr" {
			normalizedPath = path
		}

		// TLS is normally terminated by the production reverse proxy. Only
		// redirect the apex scheme when the proxy explicitly reports HTTP;
		// an absent forwarding header must not create a redirect loop.
		forceHTTPS := host == "avayusstroi.by" && strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "http")
		if host == "www.avayusstroi.by" || forceHTTPS || normalizedPath != path {
			target := normalizedPath
			if c.Request.URL.RawQuery != "" {
				target += "?" + c.Request.URL.RawQuery
			}
			if host == "avayusstroi.by" || host == "www.avayusstroi.by" {
				target = siteURL + target
			}
			c.Redirect(http.StatusMovedPermanently, target)
			c.Abort()
			return
		}

		c.Next()
	}
}

func cacheControlMiddleware() gin.HandlerFunc {
	const oneYearInSeconds = 31536000

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		switch {
		case path == "/sitemap.xml" || path == "/robots.txt" || path == "/rss.xml":
			c.Header("Cache-Control", "public, max-age=3600")
		case len(path) > 7 && path[:7] == "/static":
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
			c.Header("Expires", time.Now().Add(oneYearInSeconds*time.Second).UTC().Format(http.TimeFormat))
		default:
			c.Header("Cache-Control", "no-cache")
		}

		c.Next()
	}
}
