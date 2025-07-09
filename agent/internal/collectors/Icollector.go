package collectors

import (
	"agent/internal/models"
)

// Collector определяет интерфейс для всех сборщиков метрик
type Collector interface {
	// Collect собирает метрики и записывает их в переданную структуру
	Collect(metrics *models.AgentMetrics) error
}
