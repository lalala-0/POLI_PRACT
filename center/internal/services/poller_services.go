package services

import (
	"bytes"
	"center/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// PollerService отвечает за периодический опрос агентов
type PollerService struct {
	hostService  *HostService
	alertService *AlertNotifierService
	interval     time.Duration
	httpClient   *http.Client
	//logger      *log.Logger
}

func NewPollerService(hostService *HostService, alertService *AlertNotifierService, pollInterval time.Duration) *PollerService {
	return &PollerService{
		hostService:  hostService,
		alertService: alertService,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		interval:     pollInterval,
	}
}

// Start запускает процесс опроса хостов
func (s *PollerService) Start(ctx context.Context) {
	log.Println("Starting poller service...")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Первый опрос сразу при запуске
	s.pollHosts(ctx)

	for {
		select {
		case <-ticker.C:
			s.pollHosts(ctx)
		case <-ctx.Done():
			log.Println("Stopping poller service...")
			return
		}
	}
}

// pollHosts опрашивает все активные хосты
func (s *PollerService) pollHosts(ctx context.Context) {
	log.Println("Polling hosts...")
	hosts, err := s.hostService.GetAllHosts(ctx)
	if err != nil {
		log.Printf("Error getting hosts: %v", err)
		return
	}

	for _, host := range hosts {
		go s.pollHost(ctx, host)
	}
}

// pollHost опрашивает конкретный хост
func (s *PollerService) pollHost(ctx context.Context, host models.Host) {
	log.Printf("Polling host %s at %s:%d", host.Hostname, host.IPAddress, host.AgentPort)
	// Формируем URL для запроса метрик
	url := fmt.Sprintf("http://%s:%d/metrics", host.IPAddress, host.AgentPort)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("[%s] Error creating request: %v", host.Hostname, err)
		s.updateHostStatus(ctx, host.ID, "error")
		s.alertService.recordPollResult(host.ID, false)

		return
	}

	start := time.Now()
	resp, err := s.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		log.Printf("[%s] Polling error: %v", host.Hostname, err)
		s.updateHostStatus(ctx, host.ID, "down")
		s.alertService.recordPollResult(host.ID, false)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[%s] Unexpected status: %d", host.Hostname, resp.StatusCode)
		s.updateHostStatus(ctx, host.ID, "unstable")
		s.alertService.recordPollResult(host.ID, false)
		return
	}

	var metrics models.Metrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		log.Printf("[%s] Error decoding metrics: %v", host.Hostname, err)
		s.alertService.recordPollResult(host.ID, false)
		return
	}

	//log.Printf("[%s] Metrics: %+v", host.Hostname, metrics)

	// Обновляем статус хоста
	s.updateHostStatus(ctx, host.ID, "active")

	// Сохраняем метрики
	s.hostService.ProcessHostMetrics(ctx, host.ID, metrics)

	log.Printf("[%s] Metrics collected in %v", host.Hostname, duration)
	s.alertService.recordPollResult(host.ID, true)
	// Вызов проверки алертов после успешного получения метрик
	s.alertService.CheckHostAlerts(ctx, &host, &metrics)

}

// updateHostStatus обновляет статус хоста в БД
func (s *PollerService) updateHostStatus(ctx context.Context, hostID int, status string) {
	if err := s.hostService.UpdateHostStatus(ctx, hostID, status); err != nil {
		log.Printf("Error updating host status: %v", err)
	}
}

// ProcessHostMetrics обрабатывает и сохраняет метрики хоста
func (s *HostService) ProcessHostMetrics(ctx context.Context, hostID int, metrics models.Metrics) {
	// Сохраняем системные метрики
	systemMetrics := models.SystemMetrics{
		HostID:    hostID,
		Timestamp: metrics.Timestamp,
		System:    metrics.SystemMetrics,
	}

	if err := s.SaveSystemMetrics(ctx, &systemMetrics); err != nil {
		log.Printf("Error saving system metrics: %v", err)
	}

	// Сохраняем метрики процессов
	if len(metrics.ProcessesInfo) > 0 {
		processMetrics := models.ProcessMetrics{
			HostID:    hostID,
			Timestamp: metrics.Timestamp,
			Processes: metrics.ProcessesInfo,
		}
		if err := s.SaveProcessMetrics(ctx, &processMetrics); err != nil {
			log.Printf("Error saving process metrics: %v", err)
		}
	}

	// Сохраняем сетевые метрики
	if len(metrics.PortsInfo) > 0 {
		networkMetrics := models.NetworkMetrics{
			HostID:    hostID,
			Timestamp: metrics.Timestamp,
			Ports:     metrics.PortsInfo,
		}
		if err := s.SaveNetworkMetrics(ctx, &networkMetrics); err != nil {
			log.Printf("Error saving network metrics: %v", err)
		}
	}

	// Сохраняем метрики контейнеров
	if len(metrics.ContainersInfo) > 0 {
		containerMetrics := models.ContainerMetrics{
			HostID:     hostID,
			Timestamp:  metrics.Timestamp,
			Containers: metrics.ContainersInfo,
		}
		if err := s.SaveContainerMetrics(ctx, &containerMetrics); err != nil {
			log.Printf("Error saving container metrics: %v", err)
		}
	}
}

// SendConfigurationToAgent отправляет конфигурацию на агент
func (s *HostService) SendConfigurationToAgent(ctx context.Context, host models.Host) error {
	if err := s.SendProcessConfigurationToAgent(ctx, host); err != nil {
		return err
	}
	return s.SendContainerConfigurationToAgent(ctx, host)
}

// SendProcessConfigurationToAgent отправляет конфигурацию process на агент
func (s *HostService) SendProcessConfigurationToAgent(ctx context.Context, host models.Host) error {
	// Отправка конфигурации процессов
	processes, err := s.ProcessRepo.GetByHostID(ctx, host.ID)
	if err != nil {
		return err
	}

	processNames := make([]string, 0, len(processes))
	for _, p := range processes {
		processNames = append(processNames, p.ProcessName)
	}

	return s.sendToAgent(ctx, host, "/config/processes", map[string]interface{}{
		"processes": processNames,
	})
}

// SendContainerConfigurationToAgent отправляет конфигурацию container на агент
func (s *HostService) SendContainerConfigurationToAgent(ctx context.Context, host models.Host) error {
	// Отправка конфигурации контейнеров
	containers, err := s.ContainerRepo.GetByHostID(ctx, host.ID)
	if err != nil {
		return err
	}

	containerNames := make([]string, 0, len(containers))
	for _, c := range containers {
		containerNames = append(containerNames, c.ContainerName)
	}

	return s.sendToAgent(ctx, host, "/config/containers", map[string]interface{}{
		"containers": containerNames,
	})
}

// sendToAgent отправляет данные на агент
func (s *HostService) sendToAgent(ctx context.Context, host models.Host, endpoint string, data interface{}) error {
	url := fmt.Sprintf("http://%s:%d%s", host.IPAddress, host.AgentPort, endpoint)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	return nil
}
