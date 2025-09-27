package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContactRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
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

func init() {
	// Загрузка настроек из settings.json
	file, err := os.ReadFile("settings.json")
	if err != nil {
		log.Fatalf("Ошибка чтения settings.json: %v", err)
	}

	if err := json.Unmarshal(file, &settings); err != nil {
		log.Fatalf("Ошибка разбора settings.json: %v", err)
	}
}

func submitHandler(c *gin.Context) {
	var req ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Валидация данных
	if req.Name == "" || req.Phone == "" {
		log.Printf("Ошибка: пустое имя или телефон")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and phone are required"})
		return
	}

	log.Printf("Получена заявка: Имя: %s, Телефон: %s", req.Name, req.Phone)

	// Отправляем уведомление на email
	/*if err := sendEmail(req.Name, req.Phone); err != nil {
		log.Printf("Ошибка при отправке email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email notification"})
		return
	}*/

	// Отправляем уведомление в Telegram всем указанным чатам
	for _, chatID := range settings.Telegram.ChatIDs {
		if err := sendTelegramMessage(chatID, req.Name, req.Phone); err != nil {
			log.Printf("Ошибка при отправке Telegram-уведомления в чат %s: %v", chatID, err)
		}
	}

	// Логирование успешной заявки
	log.Printf("Успешная заявка: Имя: %s, Телефон: %s", req.Name, req.Phone)

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func sendEmail(name, phone string) error {
	from := settings.SMTP.From
	password := settings.SMTP.Password
	to := settings.SMTP.To
	smtpHost := settings.SMTP.SMTPHost
	smtpPort := settings.SMTP.SMTPPort

	subject := "Новая заявка с сайта"
	body := fmt.Sprintf("Имя: %s\nТелефон: %s", name, phone)
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
}

func sendTelegramMessage(chatID, name, phone string) error {
	message := fmt.Sprintf("Новая заявка с сайта:\nИмя: %s\nТелефон: %s", name, phone)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", settings.Telegram.BotToken)

	formData := url.Values{}
	formData.Set("chat_id", chatID)
	formData.Set("text", message)

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка при отправке сообщения в Telegram: %s", resp.Status)
	}

	return nil
}
