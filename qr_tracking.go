package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultQRStatsPath = "data/qr_stats.json"

var qrCampaignSanitizer = regexp.MustCompile(`[^a-z0-9_-]+`)
var minskLocation = time.FixedZone("Europe/Minsk", 3*60*60)

type QRStats struct {
	Total      uint64            `json:"total"`
	LastScanAt string            `json:"lastScanAt,omitempty"`
	ByCampaign map[string]uint64 `json:"byCampaign"`
	ByDate     map[string]uint64 `json:"byDate"`
}

type QRScanTracker struct {
	mu    sync.Mutex
	path  string
	stats QRStats
	ready bool
	now   func() time.Time
}

func NewQRScanTracker(path string) *QRScanTracker {
	return &QRScanTracker{
		path: path,
		now:  time.Now,
		stats: QRStats{
			ByCampaign: make(map[string]uint64),
			ByDate:     make(map[string]uint64),
		},
	}
}

func (tracker *QRScanTracker) Load() error {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	file, err := os.ReadFile(tracker.path)
	if errors.Is(err, os.ErrNotExist) {
		tracker.ready = true
		return nil
	}
	if err != nil {
		return fmt.Errorf("read QR statistics: %w", err)
	}

	var stats QRStats
	if err := json.Unmarshal(file, &stats); err != nil {
		return fmt.Errorf("decode QR statistics: %w", err)
	}
	if stats.ByCampaign == nil {
		stats.ByCampaign = make(map[string]uint64)
	}
	if stats.ByDate == nil {
		stats.ByDate = make(map[string]uint64)
	}

	tracker.stats = stats
	tracker.ready = true
	return nil
}

func (tracker *QRScanTracker) Record(campaign string) (QRStats, error) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	if !tracker.ready {
		return QRStats{}, errors.New("QR statistics storage is not ready")
	}

	campaign = normalizeQRCampaign(campaign)
	now := tracker.now().In(minskLocation)
	dateKey := now.Format("2006-01-02")

	nextStats := cloneQRStats(tracker.stats)
	nextStats.Total++
	nextStats.LastScanAt = now.Format(time.RFC3339)
	nextStats.ByCampaign[campaign]++
	nextStats.ByDate[dateKey]++

	if err := writeJSONFile(tracker.path, nextStats); err != nil {
		return cloneQRStats(tracker.stats), err
	}

	tracker.stats = nextStats
	return cloneQRStats(tracker.stats), nil
}

func cloneQRStats(stats QRStats) QRStats {
	copyStats := QRStats{
		Total:      stats.Total,
		LastScanAt: stats.LastScanAt,
		ByCampaign: make(map[string]uint64, len(stats.ByCampaign)),
		ByDate:     make(map[string]uint64, len(stats.ByDate)),
	}
	for campaign, count := range stats.ByCampaign {
		copyStats.ByCampaign[campaign] = count
	}
	for date, count := range stats.ByDate {
		copyStats.ByDate[date] = count
	}
	return copyStats
}

func writeJSONFile(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create QR statistics directory: %w", err)
	}

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode QR statistics: %w", err)
	}
	data = append(data, '\n')

	temporaryFile, err := os.CreateTemp(filepath.Dir(path), ".qr_stats_*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary QR statistics file: %w", err)
	}
	temporaryPath := temporaryFile.Name()
	defer os.Remove(temporaryPath)

	if _, err := temporaryFile.Write(data); err != nil {
		temporaryFile.Close()
		return fmt.Errorf("write QR statistics: %w", err)
	}
	if err := temporaryFile.Sync(); err != nil {
		temporaryFile.Close()
		return fmt.Errorf("sync QR statistics: %w", err)
	}
	if err := temporaryFile.Close(); err != nil {
		return fmt.Errorf("close QR statistics: %w", err)
	}

	renameErr := os.Rename(temporaryPath, path)
	if renameErr == nil {
		return nil
	}

	// Windows cannot replace an existing file with os.Rename. The fallback is
	// only used there; Unix deployments keep the atomic rename above.
	if runtime.GOOS != "windows" {
		return fmt.Errorf("replace QR statistics: %w", renameErr)
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("replace QR statistics: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace QR statistics: %w", err)
	}
	return nil
}

func normalizeQRCampaign(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = qrCampaignSanitizer.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-_")
	if value == "" {
		return "main"
	}
	if len(value) > 64 {
		value = strings.Trim(value[:64], "-_")
	}
	if value == "" {
		return "main"
	}
	return value
}

func qrRedirectHandler(tracker *QRScanTracker) gin.HandlerFunc {
	return func(c *gin.Context) {
		campaign := normalizeQRCampaign(c.Query("campaign"))
		stats, err := tracker.Record(campaign)
		if err != nil {
			log.Printf("failed to count QR scan for campaign %s: %v", campaign, err)
		} else {
			log.Printf("QR scan counted: campaign=%s total=%d", campaign, stats.Total)
		}

		query := url.Values{}
		query.Set("utm_source", "qr")
		query.Set("utm_medium", "offline")
		query.Set("utm_campaign", campaign)

		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("X-Robots-Tag", "noindex, nofollow")
		c.Redirect(http.StatusFound, "/?"+query.Encode())
	}
}
