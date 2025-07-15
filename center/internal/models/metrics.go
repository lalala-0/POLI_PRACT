package models

import (
	"time"
)

// AgentMetrics - корневая структура всех метрик, собираемых агентом
type AgentMetrics struct {
	HostID     string          `json:"host_id"`    // Уникальный идентификатор хоста
	Timestamp  time.Time       `json:"timestamp"`  // Временная метка сбора метрик
	System     SystemMetrics   `json:"system"`     // Системные метрики (CPU, RAM, Disk)
	Processes  []ProcessInfo   `json:"processes"`  // Метрики процессов
	Ports      []PortInfo      `json:"ports"`      // Информация о сетевых портах
	Containers []ContainerInfo `json:"containers"` // Метрики Docker контейнеров
}

// NewAgentMetrics создает новую структуру метрик с заполненным ID хоста и временной меткой
func NewAgentMetrics(hostID string) AgentMetrics {
	return AgentMetrics{
		HostID:    hostID,
		Timestamp: time.Now().UTC(),
		System: SystemMetrics{
			CPU:  CPUMetrics{},
			RAM:  RAMMetrics{},
			Disk: DiskMetrics{},
		},
		Processes:  make([]ProcessInfo, 0),
		Ports:      make([]PortInfo, 0),
		Containers: make([]ContainerInfo, 0),
	}
}
