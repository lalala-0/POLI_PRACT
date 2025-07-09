package service

import (
	"agent/internal/models"
	"sync"
	"time"
)

// MetricsServiceInterface определяет методы для работы с метриками
type MetricsServiceInterface interface {
	UpdateProcessConfig(processes []string) error
	GetProcessConfig() []string
	IsProcessConfigSet() bool
	UpdateContainerConfig(containers []string) error
	GetContainerConfig() []string
	IsContainerConfigSet() bool
	ProcessMetrics(metrics models.AgentMetrics)
}

// MetricsService предоставляет методы для работы с метриками
type MetricsService struct {
	processConfig      []string
	containerConfig    []string
	collectionInterval time.Duration
	mu                 sync.RWMutex
	processConfigSet   bool
	containerConfigSet bool
}

// NewMetricsService создает новый сервис метрик
func NewMetricsService() *MetricsService {
	return &MetricsService{
		processConfig:      []string{},
		containerConfig:    []string{},
		processConfigSet:   false,
		containerConfigSet: false,
	}
}

// UpdateProcessConfig обновляет список отслеживаемых процессов
func (s *MetricsService) UpdateProcessConfig(processes []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processConfig = processes
	s.processConfigSet = true
	return nil
}

// GetProcessConfig возвращает текущий список отслеживаемых процессов
func (s *MetricsService) GetProcessConfig() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.processConfig
}

// IsProcessConfigSet проверяет, установлен ли список отслеживаемых процессов
func (s *MetricsService) IsProcessConfigSet() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.processConfigSet
}

// UpdateContainerConfig обновляет список отслеживаемых контейнеров
func (s *MetricsService) UpdateContainerConfig(containers []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.containerConfig = containers
	s.containerConfigSet = true
	return nil
}

// GetContainerConfig возвращает текущий список отслеживаемых контейнеров
func (s *MetricsService) GetContainerConfig() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.containerConfig
}

// IsContainerConfigSet проверяет, установлен ли список отслеживаемых контейнеров
func (s *MetricsService) IsContainerConfigSet() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.containerConfigSet
}

// ProcessMetrics обрабатывает собранные метрики
// TODO
func (s *MetricsService) ProcessMetrics(metrics models.AgentMetrics) {
	// Здесь может быть логика анализа или фильтрации метрик
}

// UpdateCollectionInterval обновляет интервал сбора метрик
func (s *MetricsService) UpdateCollectionInterval(interval time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.collectionInterval = interval
	return nil
}
