package service

import (
	coll "agent/internal/collectors"
	"agent/internal/config"
	"agent/internal/models"
	//"log"
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
	Collectors         []coll.Collector
	collectionInterval time.Duration
	mu                 sync.RWMutex
	processConfigSet   bool
	containerConfigSet bool
}

// NewMetricsService создает новый сервис метрик
func NewMetricsService(cfg *config.AgentConfig) *MetricsService {
	// Инициализация коллекторов
	Collectors := []coll.Collector{
		coll.NewSystemCollector(),
		coll.NewProcessCollector(cfg.Processes),
		coll.NewNetworkCollector(),
	}

	// Docker коллектор добавляем, если он доступен
	if dockerCollector, err := coll.NewDockerCollector(cfg.Containers); err == nil {
		Collectors = append(Collectors, dockerCollector)
	}
	return &MetricsService{
		processConfig:      []string{},
		containerConfig:    []string{},
		Collectors:         Collectors,
		processConfigSet:   false,
		containerConfigSet: false,
	}
}

// UpdateProcessConfig обновляет список отслеживаемых процессов
func (s *MetricsService) UpdateProcessConfig(processes []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	//for _, process := range processes {
	//	s.processConfig = append(s.processConfig, process)
	//}
	//s.Collectors = append(s.Collectors, coll.NewProcessCollector(processes))
	s.processConfig = processes
	for _, c := range s.Collectors {
		c.ChangeConfig(coll.Process, processes)
	}
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
	//for _, container := range containers {
	//	s.containerConfig = append(s.containerConfig, container)
	//}
	//dcoll, err := coll.NewDockerCollector(containers)
	//if err != nil {
	//	return err
	//}
	//s.Collectors = append(s.Collectors, dcoll)
	s.containerConfig = containers
	for _, c := range s.Collectors {
		c.ChangeConfig(coll.Docker, containers)
	}
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
