package services

import (
	"bytes"
	"center/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

type AlertNotifierService struct {
	cfg         models.AlertConfig // конфигурация из YAML
	hostService HostService
	failLog     map[int][]pollResult
	logMu       sync.Mutex
}

// Конструктор
func NewAlertNotifierService(cfg models.AlertConfig, hostService HostService) *AlertNotifierService {
	return &AlertNotifierService{
		cfg:         cfg,
		hostService: hostService,
		failLog:     make(map[int][]pollResult),
	}
}

func (s *AlertNotifierService) alertMonitor(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second) // можно взять из конфига
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkAlerts(ctx)
		}
	}
}

func (s *AlertNotifierService) checkAlerts(ctx context.Context) {
	s.logMu.Lock()
	defer s.logMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-60 * time.Second) // интервал

	for hostID, logs := range s.failLog {
		var recent []pollResult
		var failed int

		// фильтрация по времени
		for _, l := range logs {
			if l.Timestamp.After(cutoff) {
				recent = append(recent, l)
				if !l.Success {
					failed++
				}
			}
		}

		// сохраняем только свежие
		s.failLog[hostID] = recent

		total := len(recent)
		if total == 0 {
			continue
		}

		failureRate := float64(failed) / float64(total) * 100

		if failureRate >= 90.0 {
			host, _ := s.hostService.GetHost(ctx, hostID)
			message := fmt.Sprintf("🚨 ALERT for host %s (%s): %.0f%% failures in last minute", host.Hostname, host.IPAddress, failureRate)

			log.Println(message)

			// отправка в Telegram и Email
			go s.sendAlert(message)
		}
	}
}

func (s *AlertNotifierService) sendAlert(message string) {
	cfg := s.cfg // или передай конфиг как поле

	if err := SendTelegramMessage(cfg.Telegram.Token, cfg.Telegram.ChatID, message); err != nil {
		log.Printf("Telegram alert failed: %v", err)
	}
	if err := SendEmailAlert(cfg.Email, "🔔 Alert triggered", message); err != nil {
		log.Printf("Email alert failed: %v", err)
	}
}

func (s *AlertNotifierService) recordPollResult(hostID int, success bool) {
	s.logMu.Lock()
	defer s.logMu.Unlock()

	s.failLog[hostID] = append(s.failLog[hostID], pollResult{
		Timestamp: time.Now(),
		Success:   success,
	})
}

func SendTelegramMessage(token, chatID, text string) error {
	url := "https://api.telegram.org/bot" + token + "/sendMessage"
	data := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}
	body, _ := json.Marshal(data)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	return err
}

// SendEmailAlert отправляет email-уведомление по заданной конфигурации.
func SendEmailAlert(cfg models.EmailConfig, subject, message string) error {
	// Формируем заголовки и тело письма
	from := cfg.Username
	to := strings.Join(cfg.To, ",")
	body := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, message)

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)

	if err := smtp.SendMail(addr, auth, from, cfg.To, []byte(body)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
