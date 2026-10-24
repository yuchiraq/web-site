package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type leadTestTransport func(*http.Request) (*http.Response, error)

func (f leadTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestLeadDeliveryAndOptionalRegionalContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousSettings, previousClient := settings, telegramHTTPClient
	t.Cleanup(func() { settings, telegramHTTPClient = previousSettings, previousClient })
	router := gin.New()
	router.POST("/submit", submitHandler)
	post := func(payload string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(payload))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		return w
	}
	settings = Settings{}
	if got := post(`{"name":"Тест","phone":"000"}`).Code; got != http.StatusServiceUnavailable {
		t.Fatalf("unconfigured delivery = %d", got)
	}
	settings.Telegram.BotToken = "test-token"
	settings.Telegram.ChatIDs = []string{"test-chat"}
	calls := 0
	telegramHTTPClient = &http.Client{Transport: leadTestTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		message := r.Form.Get("text")
		if !strings.Contains(message, "Населённый пункт: Кобрин") || !strings.Contains(message, "Задача: Дом, 4 человека") || !strings.Contains(message, "campaign=sewerage_brest") {
			t.Errorf("missing regional or campaign context: %s", message)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"ok":true}`)), Header: make(http.Header)}, nil
	})}
	response := post(`{"name":"Тест","phone":"000","city":" Кобрин ","task":"Дом, 4 человека","utmCampaign":"sewerage_brest"}`)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"success":true`) || calls != 1 {
		t.Fatalf("unexpected delivery: %d %s calls=%d", response.Code, response.Body.String(), calls)
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Error("lead response must not be cached")
	}
	telegramHTTPClient = &http.Client{Transport: leadTestTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadGateway, Status: "502 Bad Gateway", Body: io.NopCloser(strings.NewReader(`{}`)), Header: make(http.Header)}, nil
	})}
	if got := post(`{"name":"Тест","phone":"000"}`).Code; got != http.StatusBadGateway {
		t.Errorf("failed delivery returned success: %d", got)
	}
	if got := post(`{"name":"Тест","phone":" "}`).Code; got != http.StatusBadRequest {
		t.Errorf("empty phone = %d", got)
	}
	if got := post(`{"name":"Тест","phone":"000","task":"` + strings.Repeat("я", 1001) + `"}`).Code; got != http.StatusBadRequest {
		t.Errorf("overlong task = %d", got)
	}
}

func TestLeadRecipientsAreContactedConcurrently(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previousSettings, previousClient := settings, telegramHTTPClient
	t.Cleanup(func() { settings, telegramHTTPClient = previousSettings, previousClient })
	settings = Settings{}
	settings.Telegram.BotToken = "test-token"
	settings.Telegram.ChatIDs = []string{"first", "second"}
	started := make(chan bool, 2)
	release := make(chan bool)
	telegramHTTPClient = &http.Client{Transport: leadTestTransport(func(r *http.Request) (*http.Response, error) {
		started <- true
		select {
		case <-release:
		case <-r.Context().Done():
			return nil, r.Context().Err()
		}
		_ = r.ParseForm()
		status := http.StatusOK
		if r.Form.Get("chat_id") == "second" {
			status = http.StatusBadGateway
		}
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(`{"ok":true}`)), Header: make(http.Header)}, nil
	})}
	router := gin.New()
	router.POST("/submit", submitHandler)
	response := httptest.NewRecorder()
	done := make(chan bool)
	go func() {
		r := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(`{"name":"Тест","phone":"000"}`))
		r.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(response, r)
		done <- true
	}()
	concurrent := true
	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			concurrent = false
		}
	}
	close(release)
	<-done
	if !concurrent {
		t.Error("recipient requests must run concurrently")
	}
	if response.Code != http.StatusOK {
		t.Errorf("one delivered recipient must be sufficient: %d", response.Code)
	}
}
