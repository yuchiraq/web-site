package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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

func TestPublicInternalLinksAndImagesExist(t *testing.T) {
	router := seoTestRouter(t)
	known := map[string]bool{"/qr": true}
	for _, path := range publicPagePaths() {
		known[path] = true
	}
	for _, path := range publicPagePaths() {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body.String()))
		if err != nil {
			t.Fatal(err)
		}
		doc.Find("a[href^='/'], img[src^='/']").Each(func(_ int, node *goquery.Selection) {
			raw := node.AttrOr("href", node.AttrOr("src", ""))
			parsed, err := url.Parse(raw)
			if err != nil {
				t.Errorf("%s: invalid URL %q", path, raw)
				return
			}
			if strings.HasPrefix(parsed.Path, "/static/") {
				if _, err := os.Stat("." + parsed.Path); err != nil {
					t.Errorf("%s: missing asset %s", path, raw)
				}
			} else if !known[parsed.Path] {
				t.Errorf("%s: broken internal link %s", path, raw)
			}
		})
	}
}

func TestCampaignDraftAndSitemap(t *testing.T) {
	data, err := os.ReadFile("docs/marketing/campaign-draft.json")
	if err != nil {
		t.Fatal(err)
	}
	var draft struct {
		Status    string
		Campaigns []struct {
			Groups []struct {
				ID, Landing                       string
				Keywords, Headlines, Descriptions []string
			}
		}
	}
	if err := json.Unmarshal(data, &draft); err != nil {
		t.Fatal(err)
	}
	if draft.Status != "DRAFT_NOT_LAUNCHED" {
		t.Fatal("campaign status must be explicit")
	}
	groups, keywordCount := 0, 0
	keywords := map[string]bool{}
	for _, campaign := range draft.Campaigns {
		for _, group := range campaign.Groups {
			groups++
			if seoDataForPath(group.Landing).Canonical == "" {
				t.Errorf("unknown ad landing %s", group.Landing)
			}
			if len(group.Headlines) < 3 || len(group.Descriptions) < 2 {
				t.Errorf("incomplete ad group %s", group.ID)
			}
			for _, headline := range group.Headlines {
				if len([]rune(headline)) > 30 {
					t.Errorf("headline too long: %s", headline)
				}
			}
			for _, description := range group.Descriptions {
				if len([]rune(description)) > 90 {
					t.Errorf("description too long: %s", description)
				}
			}
			for _, keyword := range group.Keywords {
				if keywords[keyword] {
					t.Errorf("duplicate draft keyword %s", keyword)
				}
				keywords[keyword] = true
				keywordCount++
			}
		}
	}
	if groups != 8 || keywordCount != 80 {
		t.Errorf("unexpected campaign size: %d groups, %d keywords", groups, keywordCount)
	}
	xmlData, err := os.ReadFile("static/data/sitemap.xml")
	if err != nil {
		t.Fatal(err)
	}
	var sitemap URLSet
	if err := xml.Unmarshal(xmlData, &sitemap); err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{}
	for _, path := range publicPagePaths() {
		want[canonicalURL(path)] = true
	}
	if len(sitemap.URLs) != len(want) {
		t.Errorf("sitemap count=%d, public registry=%d", len(sitemap.URLs), len(want))
	}
	for _, item := range sitemap.URLs {
		if !want[item.Loc] {
			t.Errorf("unknown or duplicate sitemap URL %s", item.Loc)
		}
		delete(want, item.Loc)
	}
	if len(want) != 0 {
		t.Errorf("missing sitemap entries: %v", want)
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
	if len(feed.Channel.Items) != 13 {
		t.Fatalf("RSS contains %d items, want 13 service, rental and guide pages", len(feed.Channel.Items))
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

func TestTopicPagesAndFAQMatchVisibleContent(t *testing.T) {
	router := seoTestRouter(t)
	paths := map[string]bool{}
	for _, path := range publicPagePaths() {
		if paths[path] {
			t.Fatalf("duplicate public path: %s", path)
		}
		paths[path] = true
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(response.Body.String()))
		visible := []RentalFAQ{}
		doc.Find("main details").Each(func(_ int, detail *goquery.Selection) {
			visible = append(visible, RentalFAQ{Question: strings.TrimSpace(detail.Find("summary").Text()), Answer: strings.TrimSpace(detail.Find("p").Text())})
		})
		questions := []RentalFAQ{}
		faqCount := 0
		doc.Find("script[type='application/ld+json']").Each(func(_ int, script *goquery.Selection) {
			var root map[string]any
			if err := json.Unmarshal([]byte(script.Text()), &root); err != nil {
				t.Fatal(err)
			}
			nodes := []any{root}
			if graph, ok := root["@graph"].([]any); ok {
				nodes = graph
			}
			for _, raw := range nodes {
				node := raw.(map[string]any)
				if node["@type"] != "FAQPage" {
					continue
				}
				faqCount++
				for _, q := range node["mainEntity"].([]any) {
					question := q.(map[string]any)
					answer := question["acceptedAnswer"].(map[string]any)
					questions = append(questions, RentalFAQ{Question: question["name"].(string), Answer: answer["text"].(string)})
				}
			}
		})
		if faqCount > 1 || len(questions) != len(visible) {
			t.Errorf("%s: FAQ mismatch, %d graph nodes, %d schema / %d visible answers", path, faqCount, len(questions), len(visible))
		}
		for i := range questions {
			if i < len(visible) && questions[i] != visible[i] {
				t.Errorf("%s: FAQ answer differs at %d", path, i)
			}
		}
		if topic := topicPageForPath(path); topic != nil {
			for _, section := range topic.Sections {
				if doc.Find("main #"+section.ID+" h2").Text() != section.Heading {
					t.Errorf("%s: content not rendered server-side: %s", path, section.ID)
				}
			}
			if doc.Find("main .open-modal-button").Length() < 2 {
				t.Errorf("%s: missing lead CTAs", path)
			}
		}
	}
	for _, path := range []string{"/services/sewerage-without-pumping", "/services/wellpoint-dewatering", "/guides/autonomous-sewerage-cost", "/guides/sewerage-high-groundwater", "/guides", "/regions/brest-region"} {
		if !paths[path] {
			t.Errorf("missing marketing page: %s", path)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodHead, path, nil))
		if response.Code != http.StatusOK {
			t.Errorf("HEAD %s = %d", path, response.Code)
		}
	}
}
