package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ContactRequest struct {
	Name        string `json:"name" binding:"max=120"`
	Phone       string `json:"phone" binding:"max=80"`
	City        string `json:"city" binding:"max=160"`
	Task        string `json:"task" binding:"max=1000"`
	Page        string `json:"page"`
	Source      string `json:"source"`
	CTAText     string `json:"ctaText"`
	LandingPage string `json:"landingPage"`
	Referrer    string `json:"referrer"`
	UTMSource   string `json:"utmSource"`
	UTMMedium   string `json:"utmMedium"`
	UTMCampaign string `json:"utmCampaign"`
	UTMContent  string `json:"utmContent"`
	UTMTerm     string `json:"utmTerm"`
	GCLID       string `json:"gclid"`
	FBCLID      string `json:"fbclid"`
	YCLID       string `json:"yclid"`
}

type Settings struct {
	SMTP struct {
		From     string `json:"from"`
		Password string `json:"password"`
		To       string `json:"to"`
		SMTPHost string `json:"smtpHost"`
		SMTPPort string `json:"smtpPort"`
	} `json:"smtp"`
	Telegram struct {
		BotToken string   `json:"botToken"`
		ChatIDs  []string `json:"chatIds"`
	} `json:"telegram"`
}

var settings Settings
var telegramHTTPClient = &http.Client{Timeout: 10 * time.Second}

func loadSettings() error {
	file, err := os.ReadFile("settings.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &settings); err != nil {
		return err
	}

	return nil
}

func submitHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16*1024)
	var req ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.City = strings.TrimSpace(req.City)
	req.Task = strings.TrimSpace(req.Task)
	req.Page = strings.TrimSpace(req.Page)
	req.Source = strings.TrimSpace(req.Source)
	req.CTAText = strings.TrimSpace(req.CTAText)
	req.LandingPage = strings.TrimSpace(req.LandingPage)
	req.Referrer = strings.TrimSpace(req.Referrer)
	req.UTMSource = strings.TrimSpace(req.UTMSource)
	req.UTMMedium = strings.TrimSpace(req.UTMMedium)
	req.UTMCampaign = strings.TrimSpace(req.UTMCampaign)
	req.UTMContent = strings.TrimSpace(req.UTMContent)
	req.UTMTerm = strings.TrimSpace(req.UTMTerm)
	req.GCLID = strings.TrimSpace(req.GCLID)
	req.FBCLID = strings.TrimSpace(req.FBCLID)
	req.YCLID = strings.TrimSpace(req.YCLID)

	if req.Name == "" || req.Phone == "" {
		log.Printf("Ошибка: пустое имя или телефон")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and phone are required"})
		return
	}

	if settings.Telegram.BotToken == "" || len(settings.Telegram.ChatIDs) == 0 {
		log.Printf("Заявка не доставлена: Telegram не настроен")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Delivery unavailable"})
		return
	}
	delivered := false
	deliveryContext, cancelDelivery := context.WithTimeout(c.Request.Context(), 12*time.Second)
	defer cancelDelivery()
	results := make(chan bool, len(settings.Telegram.ChatIDs))
	for _, chatID := range settings.Telegram.ChatIDs {
		go func(id string) {
			err := sendTelegramMessage(deliveryContext, id, req)
			if err != nil {
				// HTTP errors can include the bot token in the request URL.
				log.Printf("Ошибка доставки заявки в Telegram")
			}
			results <- err == nil
		}(chatID)
	}
waitForDelivery:
	for range settings.Telegram.ChatIDs {
		select {
		case success := <-results:
			delivered = delivered || success
		case <-deliveryContext.Done():
			break waitForDelivery
		}
	}
	if !delivered {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Delivery failed"})
		return
	}
	log.Printf("Заявка доставлена, источник=%s", req.Source)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func sendTelegramMessage(ctx context.Context, chatID string, req ContactRequest) error {
	lines := []string{
		"Новая заявка с сайта:",
		fmt.Sprintf("Имя: %s", req.Name),
		fmt.Sprintf("Телефон: %s", req.Phone),
	}
	if req.City != "" {
		lines = append(lines, "Населённый пункт: "+req.City)
	}
	if req.Task != "" {
		lines = append(lines, "Задача: "+req.Task)
	}

	if req.Source != "" {
		lines = append(lines, fmt.Sprintf("Источник CTA: %s", req.Source))
	}
	if req.CTAText != "" {
		lines = append(lines, fmt.Sprintf("Текст CTA: %s", req.CTAText))
	}
	if req.Page != "" {
		lines = append(lines, fmt.Sprintf("Страница: %s", req.Page))
	}
	if req.LandingPage != "" {
		lines = append(lines, fmt.Sprintf("Landing page: %s", req.LandingPage))
	}
	if req.Referrer != "" {
		lines = append(lines, fmt.Sprintf("Referrer: %s", req.Referrer))
	}

	utmParts := make([]string, 0, 5)
	if req.UTMSource != "" {
		utmParts = append(utmParts, "source="+req.UTMSource)
	}
	if req.UTMMedium != "" {
		utmParts = append(utmParts, "medium="+req.UTMMedium)
	}
	if req.UTMCampaign != "" {
		utmParts = append(utmParts, "campaign="+req.UTMCampaign)
	}
	if req.UTMContent != "" {
		utmParts = append(utmParts, "content="+req.UTMContent)
	}
	if req.UTMTerm != "" {
		utmParts = append(utmParts, "term="+req.UTMTerm)
	}
	if len(utmParts) > 0 {
		lines = append(lines, "UTM: "+strings.Join(utmParts, ", "))
	}

	clickIDs := make([]string, 0, 3)
	if req.GCLID != "" {
		clickIDs = append(clickIDs, "gclid="+req.GCLID)
	}
	if req.FBCLID != "" {
		clickIDs = append(clickIDs, "fbclid="+req.FBCLID)
	}
	if req.YCLID != "" {
		clickIDs = append(clickIDs, "yclid="+req.YCLID)
	}
	if len(clickIDs) > 0 {
		lines = append(lines, "Click IDs: "+strings.Join(clickIDs, ", "))
	}

	message := strings.Join(lines, "\n")
	if runes := []rune(message); len(runes) > 3900 {
		message = string(runes[:3900]) + "\n[Контекст сокращён]"
	}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", settings.Telegram.BotToken)

	formData := url.Values{}
	formData.Set("chat_id", chatID)
	formData.Set("text", message)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := telegramHTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка при отправке сообщения в Telegram: %s", resp.Status)
	}
	var confirmation struct {
		OK bool `json:"ok"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 64*1024)).Decode(&confirmation); err != nil || !confirmation.OK {
		return fmt.Errorf("Telegram did not confirm delivery")
	}

	return nil
}
