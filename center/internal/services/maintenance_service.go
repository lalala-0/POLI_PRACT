package services

import (
	"center/internal/config"
	"center/internal/database/mongodb/repositories"
	pg_repo "center/internal/database/postgres/repositories"
	"context"
	"log"

	"time"
)

// MaintenanceService отвечает за фоновое обслуживание системы
type MaintenanceService struct {
	metricRepo repositories.MongoMetricRepository
	hostRepo   pg_repo.PostgresHostRepository
	config     config.MetricsConfig
}

func NewMaintenanceService(
	metricRepo repositories.MongoMetricRepository,
	hostRepo pg_repo.PostgresHostRepository,
	config config.MetricsConfig,
) *MaintenanceService {
	return &MaintenanceService{
		metricRepo: metricRepo,
		hostRepo:   hostRepo,
		config:     config,
	}
}

// StartCleanupRoutine запускает фоновую очистку старых данных
func (s *MaintenanceService) StartCleanupRoutine(ctx context.Context) {
	// Очистка при старте
	s.cleanupOldMetrics(ctx)

	ticker := time.NewTicker(24 * time.Hour) // Очистка раз в сутки
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupOldMetrics(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// cleanupOldMetrics удаляет метрики старше заданного срока
func (s *MaintenanceService) cleanupOldMetrics(ctx context.Context) {
	threshold := time.Now().AddDate(0, 0, -s.config.MetricsTTLDays)

	collections := []string{
		"system_metrics",
		"process_metrics",
		"container_metrics",
		"network_metrics",
	}

	for _, collection := range collections {
		if err := s.metricRepo.CleanupOldMetrics(ctx, collection, threshold); err != nil {
			log.Printf("Failed to cleanup %s: %v", collection, err)
		} else {
			log.Printf("Cleaned up old %s metrics", collection)
		}
	}
}

// StartSelfCheckRoutine запускает самодиагностику системы
func (s *MaintenanceService) StartSelfCheckRoutine(ctx context.Context) {
	// Первая проверка сразу при запуске
	s.selfCheck(ctx)

	ticker := time.NewTicker(s.config.SelfCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.selfCheck(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// selfCheck выполняет проверки работоспособности системы
func (s *MaintenanceService) selfCheck(ctx context.Context) {
	// Проверка подключения к БД
	//if err := s.hostRepo.Ping(ctx); err != nil {
	//	log.Printf("SELF-CHECK FAILED: PostgreSQL connection error: %v", err)
	//}

	if err := s.metricRepo.Ping(ctx); err != nil {
		log.Printf("SELF-CHECK FAILED: MongoDB connection error: %v", err)
	}

	// Проверка количества активных хостов
	hosts, err := s.hostRepo.GetAll(ctx)
	if err != nil {
		log.Printf("SELF-CHECK FAILED: Could not retrieve hosts: %v", err)
		return
	}

	activeCount := 0
	for _, host := range hosts {
		if host.Status == "active" {
			activeCount++
		}
	}

	if activeCount == 0 {
		log.Println("SELF-CHECK WARNING: No active hosts detected")
	}

	// Проверка мастер-хоста
	master, err := s.hostRepo.GetMaster(ctx)
	if err != nil || master == nil {
		log.Println("SELF-CHECK FAILED: Master host not found")
	} else if master.Status != "active" {
		log.Printf("SELF-CHECK FAILED: Master host %s is not active", master.Hostname)
	}

	log.Println("SELF-CHECK: System health verified")
}
