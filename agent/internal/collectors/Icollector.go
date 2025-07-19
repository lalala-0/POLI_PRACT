package collectors

import (
	"agent/internal/models"
)

type CollectorType int

const (
	Docker  CollectorType = iota // iota = 0
	Process                      // iota = 1
	Network                      // iota = 2
	System                       // iota = 3
)

// Collector определяет интерфейс для всех сборщиков метрик
type Collector interface {
	// Collect собирает метрики и записывает их в переданную структуру
	Collect(metrics *models.AgentMetrics) error
	// ChangeConfig изменяет конфигурацию коллектора
	ChangeConfig(collType CollectorType, newconfig []string)
}
