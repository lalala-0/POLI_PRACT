package models

import "time"

// AgentMetrics - корневая структура всех метрик, собираемых агентом
type AgentMetrics struct {
	HostID     string          `json:"host_id"`
	Timestamp  time.Time       `json:"timestamp"`
	System     SystemMetrics   `json:"system,omitempty"`
	Processes  []ProcessInfo   `json:"processes,omitempty"`
	Ports      []PortInfo      `json:"ports,omitempty"`
	Containers []ContainerInfo `json:"containers,omitempty"`
}

// NewAgentMetrics создает новую структуру метрик с заполненным ID хоста и временной меткой
func NewAgentMetrics(hostID string) AgentMetrics {
	return AgentMetrics{
		HostID:    hostID,
		Timestamp: time.Now(),
	}
}

// SystemMetrics содержит информацию о системных ресурсах
type SystemMetrics struct {
	CPU  CPUMetrics  `json:"cpu"`
	RAM  RAMMetrics  `json:"ram"`
	Disk DiskMetrics `json:"disk"`
}

// CPUMetrics содержит информацию о загрузке процессора
type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"` // Процент использования CPU
}

// RAMMetrics содержит информацию об использовании памяти
type RAMMetrics struct {
	Total        uint64  `json:"total"`         // Общий объем в байтах
	Used         uint64  `json:"used"`          // Используемый объем в байтах
	Free         uint64  `json:"free"`          // Свободный объем в байтах
	UsagePercent float64 `json:"usage_percent"` // Процент использования
}

// DiskMetrics содержит информацию об использовании диска
type DiskMetrics struct {
	Total        uint64  `json:"total"`         // Общий объем в байтах
	Used         uint64  `json:"used"`          // Используемый объем в байтах
	Free         uint64  `json:"free"`          // Свободный объем в байтах
	UsagePercent float64 `json:"usage_percent"` // Процент использования
}

// ProcessInfo содержит информацию о процессе
type ProcessInfo struct {
	PID        int32   `json:"pid"`         // ID процесса
	Name       string  `json:"name"`        // Имя процесса
	CPUPercent float64 `json:"cpu_percent"` // Процент использования CPU
	MemPercent float64 `json:"mem_percent"` // Процент использования памяти
}

// PortInfo содержит информацию об открытом сетевом порте
type PortInfo struct {
	Port     uint16 `json:"port"`     // Номер порта
	Protocol string `json:"protocol"` // Протокол (TCP/UDP)
	State    string `json:"state"`    // Состояние (LISTEN, etc.)
}

// ContainerInfo содержит информацию о Docker-контейнере
type ContainerInfo struct {
	ID         string  `json:"id"`          // Короткий ID контейнера
	Name       string  `json:"name"`        // Имя контейнера
	Image      string  `json:"image"`       // Образ контейнера
	Status     string  `json:"status"`      // Статус (running, stopped, etc.)
	CPUPercent float64 `json:"cpu_percent"` // Процент использования CPU
	MemPercent float64 `json:"mem_percent"` // Процент использования памяти
}
