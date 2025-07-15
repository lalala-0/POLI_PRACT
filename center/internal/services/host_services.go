package services

import (
	"POLI_PRACT/center/internal/models"
	"POLI_PRACT/center/internal/repositories"
	"net/http"
	"time"
	"context"
	"errors"
	"log"
)


// HostService реализует бизнес-логику работы с хостами
type HostService struct {
	hostRepo     repositories.HostRepository
	processRepo  repositories.ProcessRepository
	containerRepo repositories.ContainerRepository
	alertRepo    repositories.AlertRepository
	metricRepo   repositories.MetricRepository
}

func NewHostService(
	hostRepo repositories.HostRepository,
	processRepo repositories.ProcessRepository,
	containerRepo repositories.ContainerRepository,
	alertRepo repositories.AlertRepository,
	metricRepo repositories.MetricRepository,
) *HostService {
	return &HostService{
		hostRepo:     hostRepo,
		processRepo:  processRepo,
		containerRepo: containerRepo,
		alertRepo:    alertRepo,
		metricRepo:   metricRepo,
	}
}

// Host Operations
func (s *HostService) CreateHost(ctx context.Context, hostInput models.HostInput) (int, error) {
	host := models.Host{
		Hostname:  hostInput.Hostname,
		IPAddress: hostInput.IPAddress,
		Priority:  hostInput.Priority,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.hostRepo.Create(ctx, &host)
}

func (s *HostService) GetHost(ctx context.Context, id int) (*models.Host, error) {
	return s.hostRepo.GetByID(ctx, id)
}

func (s *HostService) UpdateHost(ctx context.Context, id int, hostInput models.HostInput) error {
	host, err := s.hostRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	host.Hostname = hostInput.Hostname
	host.IPAddress = hostInput.IPAddress
	host.Priority = hostInput.Priority
	host.UpdatedAt = time.Now()

	return s.hostRepo.Update(ctx, host)
}

func (s *HostService) DeleteHost(ctx context.Context, id int) error {
	return s.hostRepo.Delete(ctx, id)
}

func (s *HostService) SetMasterHost(ctx context.Context, id int) error {
	// Сброс текущего мастера
	currentMaster, err := s.hostRepo.GetMaster(ctx)
	if err == nil && currentMaster != nil {
		currentMaster.IsMaster = false
		if err := s.hostRepo.Update(ctx, currentMaster); err != nil {
			log.Printf("Error resetting master: %v", err)
		}
	}

	// Установка нового мастера
	host, err := s.hostRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	host.IsMaster = true
	return s.hostRepo.Update(ctx, host)
}

// Process Operations
func (s *HostService) AddProcess(ctx context.Context, hostID int, processName string) (int, error) {
	exists, err := s.processRepo.Exists(ctx, hostID, processName)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("process already monitored")
	}

	process := &models.Process{
		HostID:      hostID,
		ProcessName: processName,
	}
	return s.processRepo.Create(ctx, process)
}

// Container Operations
func (s *HostService) AddContainer(ctx context.Context, hostID int, containerName string) (int, error) {
	exists, err := s.containerRepo.Exists(ctx, hostID, containerName)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("container already monitored")
	}

	container := &models.Container{
		HostID:        hostID,
		ContainerName: containerName,
	}
	return s.containerRepo.Create(ctx, container)
}

// Alert Operations
func (s *HostService) CreateAlertRule(ctx context.Context, hostID int, alertInput models.AlertInput) (int, error) {
	rule := &models.AlertRule{
		HostID:         hostID,
		MetricName:     alertInput.MetricName,
		ThresholdValue: alertInput.ThresholdValue,
		Condition:      alertInput.Condition,
		Enabled:        alertInput.Enabled,
	}
	return s.alertRepo.Create(ctx, rule)
}


// type AlertService struct {
// 	logger *logrus.Logger
// }

// func (a *AlertService) CheckAlerts(host models.Host, metrics models.Metrics) {
// 	for _, alert := range host.Alerts {
// 		if !alert.Enabled {
// 			continue
// 		}
// 		// Логика проверки метрик
// 		if alert.MetricName == "cpu_usage_percent" {
// 			if alert.Condition == ">" && metrics.CPU.UsagePercent > alert.Threshold {
// 				a.logger.Warnf("ALERT for %s: CPU usage %.2f%%", host.Hostname, metrics.CPU.UsagePercent)
// 			}
// 		}
// 		// Аналогично для RAM, Disk и др.
// 	}
// }

// type MaintenanceService struct {
// 	metricsRepo *repositories.MetricsRepository
// 	config      *config.Config
// }

// func (m *MaintenanceService) StartCleanupRoutine() {
// 	ticker := time.NewTicker(24 * time.Hour)
// 	for range ticker.C {
// 		m.metricsRepo.CleanupOldMetrics(m.config.Metrics.TTLDays)
// 	}
// }


// Metrics Operations
func (s *HostService) SaveSystemMetrics(ctx context.Context, metrics *models.SystemMetrics) error {
	return s.metricRepo.SaveSystemMetrics(ctx, metrics)
}

func (s *HostService) SaveProcessMetrics(ctx context.Context, metrics *models.ProcessMetrics) error {
	return s.metricRepo.SaveProcessMetrics(ctx, metrics)
}

func (s *HostService) SaveContainerMetrics(ctx context.Context, metrics *models.ContainerMetrics) error {
	return s.metricRepo.SaveContainerMetrics(ctx, metrics)
}

func (s *HostService) SaveNetworkMetrics(ctx context.Context, metrics *models.NetworkMetrics) error {
	return s.metricRepo.SaveNetworkMetrics(ctx, metrics)
}
