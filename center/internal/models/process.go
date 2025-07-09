package models

import "time"

// Process представляет процесс, который нужно мониторить на хосте
type Process struct {
	ID          int    `json:"id" db:"id"`
	HostID      int    `json:"host_id" db:"host_id"`
	ProcessName string `json:"process_name" binding:"required" db:"process_name"`
}

// ProcessInput представляет данные для добавления процесса
type ProcessInput struct {
	ProcessName string `json:"process_name" binding:"required"`
}

// ProcessMetrics представляет метрики процессов
type ProcessMetrics struct {
	HostID    int           `json:"host_id" bson:"host_id"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
	Processes []ProcessInfo `json:"processes" bson:"processes"`
}

// ProcessInfo представляет информацию о процессе
type ProcessInfo struct {
	Name       string  `json:"name" bson:"name"`
	PID        int     `json:"pid" bson:"pid"`
	Status     string  `json:"status" bson:"status"`
	CPUPercent float64 `json:"cpu_percent" bson:"cpu_percent"`
	MemoryMB   float64 `json:"memory_mb" bson:"memory_mb"`
}
