package services

import (
	"center/internal/models"
	"center/internal/repositories/postgres"
	"center/internal/repositories/mongodb"
	"net/http"
	"time"
	"context"
	"errors"
	"log"
)


type MaintenanceService struct {
	metricsRepo *repositories.MetricsRepository
	config      *config.Config
}

func (m *MaintenanceService) StartCleanupRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		m.metricsRepo.CleanupOldMetrics(m.config.Metrics.TTLDays)
	}
}