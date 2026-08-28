package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestQRScanTrackerPersistsCounts(t *testing.T) {
	statsPath := filepath.Join(t.TempDir(), "qr_stats.json")
	tracker := NewQRScanTracker(statsPath)
	tracker.now = func() time.Time {
		return time.Date(2026, time.August, 28, 12, 30, 0, 0, time.Local)
	}

	if err := tracker.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if _, err := tracker.Record("catalog-2026"); err != nil {
		t.Fatalf("Record() error = %v", err)
	}
	if _, err := tracker.Record("catalog-2026"); err != nil {
		t.Fatalf("Record() second error = %v", err)
	}

	reloaded := NewQRScanTracker(statsPath)
	if err := reloaded.Load(); err != nil {
		t.Fatalf("reloaded.Load() error = %v", err)
	}

	if reloaded.stats.Total != 2 {
		t.Fatalf("Total = %d, want 2", reloaded.stats.Total)
	}
	if reloaded.stats.ByCampaign["catalog-2026"] != 2 {
		t.Fatalf("campaign count = %d, want 2", reloaded.stats.ByCampaign["catalog-2026"])
	}
	if reloaded.stats.ByDate["2026-08-28"] != 2 {
		t.Fatalf("date count = %d, want 2", reloaded.stats.ByDate["2026-08-28"])
	}

	data, err := json.Marshal(reloaded.stats)
	if err != nil || len(data) == 0 {
		t.Fatalf("stored statistics are not valid JSON: %v", err)
	}
}

func TestQRRedirectHandlerCountsAndRedirects(t *testing.T) {
	gin.SetMode(gin.TestMode)
	statsPath := filepath.Join(t.TempDir(), "qr_stats.json")
	tracker := NewQRScanTracker(statsPath)
	if err := tracker.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	router := gin.New()
	router.GET("/qr", qrRedirectHandler(tracker))

	request := httptest.NewRequest(http.MethodGet, "/qr?campaign=Street%20Banner", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusFound)
	}
	if location := response.Header().Get("Location"); location != "/?utm_campaign=street-banner&utm_medium=offline&utm_source=qr" {
		t.Fatalf("Location = %q", location)
	}
	if tracker.stats.Total != 1 || tracker.stats.ByCampaign["street-banner"] != 1 {
		t.Fatalf("unexpected statistics: %+v", tracker.stats)
	}
	if cacheControl := response.Header().Get("Cache-Control"); cacheControl != "no-store, no-cache, must-revalidate" {
		t.Fatalf("Cache-Control = %q", cacheControl)
	}
	if robots := response.Header().Get("X-Robots-Tag"); robots != "noindex, nofollow" {
		t.Fatalf("X-Robots-Tag = %q", robots)
	}
}

func TestNormalizeQRCampaign(t *testing.T) {
	tests := map[string]string{
		"":                 "main",
		"  Street Banner ": "street-banner",
		"Каталог 2026":     "2026",
		"summer_email":     "summer_email",
	}

	for input, want := range tests {
		if got := normalizeQRCampaign(input); got != want {
			t.Errorf("normalizeQRCampaign(%q) = %q, want %q", input, got, want)
		}
	}
}
