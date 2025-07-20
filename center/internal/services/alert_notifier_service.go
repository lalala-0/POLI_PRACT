package services

import (
	"bytes"
	"center/internal/config"
	"center/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type checkResult struct {
	success int
	fail    int
}

type AlertNotifierService struct {
	cfg           config.AlertsConfig // конфигурация из YAML
	hostService   *HostService
	checksCounter *checkResult
	alertRules    map[int][]models.AlertRule // Кэш правил алертов
	logMu         sync.Mutex
	ruleMu        sync.RWMutex
}

// Конструктор
func NewAlertNotifierService(cfg config.AlertsConfig, hostService *HostService) *AlertNotifierService {
	checksCounter := &checkResult{0, 0}
	service := &AlertNotifierService{
		cfg:           cfg,
		hostService:   hostService,
		checksCounter: checksCounter,
		alertRules:    make(map[int][]models.AlertRule),
	}
	service.refreshAlertRules(context.Background())
	// Загрузка правил при инициализации
	//go service.loadAlertRules(context.Background())
	return service
}

// loadAlertRules каждые 5 минут обновляет данные
func (s *AlertNotifierService) loadAlertRules(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		s.refreshAlertRules(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.refreshAlertRules(ctx)
		}
	}
}

func (s *AlertNotifierService) refreshAlertRules(ctx context.Context) {
	hosts, err := s.hostService.GetAllHosts(ctx)
	if err != nil {
		log.Printf("Failed to load hosts for alert rules: %v", err)
		return
	}

	newRules := make(map[int][]models.AlertRule)
	for _, host := range hosts {
		rules, err := s.hostService.GetAlertsByHostID(ctx, host.ID)
		if err != nil {
			log.Printf("Failed to get alerts for host %d: %v", host.ID, err)
			continue
		}
		newRules[host.ID] = rules
	}

	s.ruleMu.Lock()
	s.alertRules = newRules
	s.ruleMu.Unlock()
}

// InvalidateAlertRules обновляет алерты в кэше
func (s *AlertNotifierService) InvalidateAlertRules() {
	s.ruleMu.Lock()
	defer s.ruleMu.Unlock()

	// Очищаем кэш, чтобы при следующем обращении загрузились свежие данные
	s.alertRules = make(map[int][]models.AlertRule)

	// Инициируем немедленную перезагрузку
	go s.refreshAlertRules(context.Background())
}

// CheckHostAlerts проверяет метрики на соответствие правилам алертов
func (s *AlertNotifierService) CheckHostAlerts(ctx context.Context, host *models.Host, metrics *models.Metrics) {
	s.ruleMu.RLock()
	rules, ok := s.alertRules[host.ID]
	s.ruleMu.RUnlock()

	//log.Println("------------1-------------------", host.ID, ":   ", rules)
	if !ok {
		return
	}

	//log.Println("------------2-------------------")
	host, err := s.hostService.GetHost(ctx, host.ID)
	if err != nil {
		log.Printf("Failed to get host for alert: %v", err)
		return
	}

	//log.Println("------------3-------------------")
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		triggered, current := s.evaluateRule(metrics, rule)
		if !triggered {
			continue
		}
		message := fmt.Sprintf("🔔 ALERT: Host %s (%s): %s %s %.2f (current: %s)",
			host.Hostname, host.IPAddress, rule.MetricName, rule.Condition, rule.ThresholdValue, current)
		log.Println(message)
		s.sendAlert(message)
	}
}

func (s *AlertNotifierService) evaluateRule(metrics *models.Metrics, rule models.AlertRule) (bool, string) {
	// Парсим имя метрики: тип.имя.поле
	parts := strings.Split(rule.MetricName, ".")
	l := len(parts)
	if l < 2 || l > 3 {
		return false, "invalid metric name"
	}

	metricType := parts[0]
	var objectName, fieldName string
	if l == 2 {
		fieldName = parts[1]
	} else {
		objectName = parts[1]
		fieldName = parts[2]
	}

	switch metricType {
	case "system":
		return s.evaluateSystemMetric(metrics.SystemMetrics, rule, fieldName)
	case "process":
		return s.evaluateProcessMetric(metrics.ProcessesInfo, rule, objectName, fieldName)
	case "container":
		return s.evaluateContainerMetric(metrics.ContainersInfo, rule, objectName, fieldName)
	case "network":
		return s.evaluateNetworkMetric(metrics.PortsInfo, rule, objectName, fieldName)
	default:
		return false, "unknown metric type"
	}
}

func (s *AlertNotifierService) evaluateSystemMetric(system models.SystemDetails, rule models.AlertRule, fieldName string) (bool, string) {
	var value float64
	var current string

	switch fieldName {
	case "cpu_usage_percent":
		value = system.CPU.UsagePercent
		current = fmt.Sprintf("%.2f%%", value)
	case "memory_usage_percent":
		value = system.RAM.UsagePercent
		current = fmt.Sprintf("%.2f%%", value)
	case "disk_usage_percent":
		value = system.Disk.UsagePercent
		current = fmt.Sprintf("%.2f%%", value)
	default:
		return false, "unknown system metric"
	}

	return s.compare(value, rule), current
}

func (s *AlertNotifierService) evaluateProcessMetric(processes []models.ProcessInfo, rule models.AlertRule, processName, fieldName string) (bool, string) {
	for _, proc := range processes {
		if proc.Name == processName {
			var value float64
			var current string

			switch fieldName {
			case "cpu_percent":
				value = proc.CPUPercent
				current = fmt.Sprintf("%.2f%%", value)
			case "memory_mb":
				value = proc.MemoryMB
				current = fmt.Sprintf("%.2fMB", value)
			default:
				return false, "unknown process metric"
			}

			return s.compare(value, rule), current
		}
	}
	return false, "process not found"
}

func (s *AlertNotifierService) evaluateContainerMetric(containers []models.ContainerInfo, rule models.AlertRule, containerName, fieldName string) (bool, string) {
	for _, cont := range containers {
		if cont.Name == containerName {
			var value float64
			var current string

			switch fieldName {
			case "cpu_percent":
				value = cont.CPUPercent
				current = fmt.Sprintf("%.2f%%", value)
			case "memory_percent":
				value = cont.MemoryPercent
				current = fmt.Sprintf("%.2f%%", value)
			case "status":
				// Преобразуем статус в числовое значение
				if cont.Status == "running" {
					value = 1
				} else {
					value = 0
				}
				current = cont.Status
			default:
				return false, "unknown container metric"
			}

			return s.compare(value, rule), current
		}
	}
	return false, "container not found"
}

func (s *AlertNotifierService) evaluateNetworkMetric(ports []models.PortInfo, rule models.AlertRule, portStr, fieldName string) (bool, string) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false, "invalid port"
	}

	for _, p := range ports {
		if p.LocalPort == port {
			var value float64
			var current string

			switch fieldName {
			case "status":
				// Преобразуем статус в числовое значение
				if p.State == "LISTEN" {
					value = 1
				} else {
					value = 0
				}
				current = p.State
			default:
				return false, "unknown network metric"
			}

			return s.compare(value, rule), current
		}
	}
	return false, "port not found"
}

func (s *AlertNotifierService) compare(value float64, rule models.AlertRule) bool {
	switch rule.Condition {
	case ">":
		return value > rule.ThresholdValue
	case "<":
		return value < rule.ThresholdValue
	case "=":
		return value == rule.ThresholdValue
	case ">=":
		return value >= rule.ThresholdValue
	case "<=":
		return value <= rule.ThresholdValue
	case "!=":
		return value != rule.ThresholdValue
	default:
		return false
	}
}

// sendAlert отправляет уведомления во все каналы
func (s *AlertNotifierService) sendAlert(message string) {
	// Отправка в Telegram для всех чатов
	for _, chatID := range s.cfg.Telegram.ChatIDs {
		if err := SendTelegramMessage(s.cfg.Telegram.Token, chatID, message); err != nil {
			log.Printf("Telegram alert to %s failed: %v", chatID, err)
		}
	}

	//// Отправка email всем получателям
	//if err := SendEmailAlert(s.cfg.Email, "🔔 Monitoring Alert", message); err != nil {
	//	log.Printf("Email alert failed: %v", err)
	//}
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
func SendEmailAlert(cfg config.EmailConfig, subject, message string) error {
	// Формируем заголовки и тело письма
	from := cfg.Username
	to := strings.Join(cfg.To, ",")
	body := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, message)

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	//log.Println("              auth:", auth, ", addr:", addr, ", body:", body, "to:", to)
	if err := smtp.SendMail(addr, auth, from, cfg.To, []byte(body)); err != nil {
		log.Printf("failed to send email: %w", err)
		return err
	}
	return nil
}

func (s *AlertNotifierService) recordCheckResult(success bool) {
	s.logMu.Lock()
	defer s.logMu.Unlock()
	if success {
		s.checksCounter.success++
	} else {
		s.checksCounter.fail++
	}
}

//
//// Start запускает мониторинг
//func (s *AlertNotifierService) Start(ctx context.Context) {
//	go s.alertMonitor(ctx)
//}

//
//// AfterPoll вызывается после каждого цикла опроса хостов
//func (s *AlertNotifierService) AfterPoll(results map[int]bool) {
//	for hostID, success := range results {
//		s.recordPollResult(hostID, success)
//	}
//}

func (s *AlertNotifierService) AlertMonitor(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.cfg.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checksCounterAlerts()
		}
	}
}

func (s *AlertNotifierService) checksCounterAlerts() {
	s.logMu.Lock()
	defer s.logMu.Unlock()
	failureRate := float64(s.checksCounter.fail) / float64(s.checksCounter.fail+s.checksCounter.success) * 100.
	s.checksCounter.success = 0
	s.checksCounter.fail = 0
	// Проверяем порог срабатывания
	if failureRate >= s.cfg.FailureThresholdPercent {
		message := fmt.Sprintf("🚨 ALERT Monitoring center failed: %.0f%% failures in last %d seconds",
			failureRate, s.cfg.IntervalSeconds)

		log.Println(message)

		// отправка в Telegram и Email
		go s.sendAlert(message)
	}
}
