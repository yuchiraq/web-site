package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ContactRequest struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
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
	var req ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
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

	log.Printf("Получена заявка: имя=%s телефон=%s источник=%s страница=%s", req.Name, req.Phone, req.Source, req.Page)

	if settings.Telegram.BotToken == "" || len(settings.Telegram.ChatIDs) == 0 {
		log.Printf("Предупреждение: Telegram не настроен, уведомления не отправлены")
	} else {
		for _, chatID := range settings.Telegram.ChatIDs {
			if err := sendTelegramMessage(chatID, req); err != nil {
				log.Printf("Ошибка при отправке Telegram-уведомления в чат %s: %v", chatID, err)
			}
		}
	}

	log.Printf("Успешная заявка: имя=%s телефон=%s источник=%s", req.Name, req.Phone, req.Source)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func sendTelegramMessage(chatID string, req ContactRequest) error {
	lines := []string{
		"Новая заявка с сайта:",
		fmt.Sprintf("Имя: %s", req.Name),
		fmt.Sprintf("Телефон: %s", req.Phone),
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
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", settings.Telegram.BotToken)

	formData := url.Values{}
	formData.Set("chat_id", chatID)
	formData.Set("text", message)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка при отправке сообщения в Telegram: %s", resp.Status)
	}

	return nil
}
