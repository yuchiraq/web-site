package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func seoTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.Use(canonicalHostMiddleware())
	if err := loadTemplates(router); err != nil {
		t.Fatal(err)
	}
	router.GET("/", indexHandler)
	router.HEAD("/", indexHandler)
	if err := registerRoutes(router); err != nil {
		t.Fatal(err)
	}
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", pageData("/404"))
	})
	return router
}

func TestPublicPagesSEOAndLinks(t *testing.T) {
	router := seoTestRouter(t)
	titles := map[string]string{}
	descriptions := map[string]string{}
	linkedPaths := map[string]bool{"/": true}

	for _, path := range publicPagePaths() {
		t.Run(path, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path+"?utm_source=test", nil))
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d", response.Code)
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body.String()))
			if err != nil {
				t.Fatal(err)
			}
			if doc.Find("main h1").Length() != 1 {
				t.Fatal("expected exactly one main H1")
			}
			canonical := doc.Find("head link[rel='canonical']")
			if canonical.Length() != 1 || canonical.AttrOr("href", "") != canonicalURL(path) {
				t.Fatalf("unexpected canonical: %s", canonical.AttrOr("href", ""))
			}
			if doc.Find("meta[property='og:url']").AttrOr("content", "") != canonicalURL(path) {
				t.Fatal("og:url differs from canonical")
			}
			if strings.Contains(doc.Find("meta[name='robots']").AttrOr("content", ""), "noindex") {
				t.Fatal("public page is noindex")
			}
			if path == "/" {
				if !strings.Contains(doc.Find("main h1").Text(), "Бресте") {
					t.Error("homepage H1 must include the service location")
				}
				preloads := doc.Find("link[rel='preload'][as='image']")
				if preloads.Length() != 2 || preloads.Filter("[media='(max-width: 768px)']").Length() != 1 || preloads.Filter("[media='(min-width: 769px)']").Length() != 1 {
					t.Error("homepage hero preloads must match desktop and mobile backgrounds")
				}
			}
			title := strings.TrimSpace(doc.Find("title").Text())
			description := doc.Find("meta[name='description']").AttrOr("content", "")
			if title == "" || description == "" {
				t.Fatal("missing title or description")
			}
			if other, found := titles[title]; found {
				t.Errorf("title duplicates %s", other)
			}
			if other, found := descriptions[description]; found {
				t.Errorf("description duplicates %s", other)
			}
			titles[title], descriptions[description] = path, path
			doc.Find("script[type='application/ld+json']").Each(func(_ int, item *goquery.Selection) {
				if !json.Valid([]byte(item.Text())) {
					t.Error("invalid JSON-LD")
				}
			})
			doc.Find("a[href]").Each(func(_ int, item *goquery.Selection) {
				href := item.AttrOr("href", "")
				if strings.HasPrefix(href, "/") {
					linkedPaths[strings.Split(href, "#")[0]] = true
				}
			})
			if doc.Find("a.qoncept-credit").AttrOr("href", "") != "https://qoncept.by/" {
				t.Error("missing Qoncept credit")
			}
		})
	}
	for _, path := range publicPagePaths() {
		if !linkedPaths[path] {
			t.Errorf("orphan page without internal links: %s", path)
		}
	}
}

func TestCanonicalRedirectsAndNotFound(t *testing.T) {
	router := seoTestRouter(t)
	tests := []struct {
		URL      string
		Forward  string
		Status   int
		Location string
	}{
		{"http://www.avayusstroi.by/services/topas/?utm_source=qr", "https", 301, "https://avayusstroi.by/services/topas?utm_source=qr"},
		{"http://avayusstroi.by/services", "http", 301, "https://avayusstroi.by/services"},
		{"http://avayusstroi.by/services", "https", 200, ""},
		{"http://localhost/services/topas.html", "", 301, "/services/topas"},
		{"http://localhost/index.html", "", 301, "/"},
		{"http://localhost/services/", "", 301, "/services"},
		{"http://localhost/404", "", 404, ""},
		{"http://localhost/not-a-real-page", "", 404, ""},
		{"http://localhost/static/rental_detail", "", 404, ""},
	}
	for _, test := range tests {
		t.Run(test.URL+test.Forward, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.URL, nil)
			request.Header.Set("X-Forwarded-Proto", test.Forward)
			router.ServeHTTP(response, request)
			if response.Code != test.Status || response.Header().Get("Location") != test.Location {
				t.Fatalf("got %d %q, want %d %q", response.Code, response.Header().Get("Location"), test.Status, test.Location)
			}
			if test.Status == http.StatusNotFound && !strings.Contains(response.Body.String(), "noindex") {
				t.Fatal("404 must be noindex")
			}
		})
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodHead, "/", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("HEAD / status = %d", response.Code)
	}
}

func TestRSSUsesCanonicalMetadata(t *testing.T) {
	feed, err := buildRSSFeed()
	if err != nil {
		t.Fatal(err)
	}
	if len(feed.Channel.Items) != 9 {
		t.Fatalf("RSS contains %d items, want 9", len(feed.Channel.Items))
	}
	for _, item := range feed.Channel.Items {
		path := strings.TrimPrefix(item.Link, siteURL)
		seo := seoDataForPath(path)
		if item.Title != seo.Title || item.Description != seo.Description || item.GUID != item.Link {
			t.Errorf("RSS metadata mismatch for %s", path)
		}
		if strings.Contains(item.Content.HTML, "{{") {
			t.Errorf("unrendered template in RSS for %s", path)
		}
	}
	encoded, err := xml.Marshal(feed)
	if err != nil {
		t.Fatal(err)
	}
	var parsed any
	if err := xml.Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("RSS is not valid XML: %v", err)
	}
}
