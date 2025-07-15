package models

import (
	"time"
)

// Host представляет информацию о хосте для мониторинга
type Host struct {
	ID        int       `json:"id" db:"id"`
	Hostname  string    `json:"hostname" binding:"required" db:"hostname"`
	IPAddress string    `json:"ip_address" binding:"required" db:"ip_address"`
	AgentPort int       `json:"agent_port" db:"agent_port"` 
	Priority  int       `json:"priority" db:"priority"`
	IsMaster  bool      `json:"is_master" db:"is_master"`
	Status    string    `json:"status" db:"status"`
	//LastCheck time.Time `json:"last_check" db:"last_check"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// HostInput представляет данные для создания/обновления хоста
type HostInput struct {
	Hostname  string `json:"hostname" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	Priority  int    `json:"priority"`
}

// SystemMetrics представляет основные метрики хоста
type SystemMetrics struct {
	HostID    int           `json:"host_id" bson:"host_id"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
	CPU       CPUMetrics    `json:"cpu" bson:"cpu"`
	Memory    MemoryMetrics `json:"memory" bson:"memory"`
	Disk      []DiskMetrics `json:"disk" bson:"disk"`
}

// CPUMetrics представляет метрики CPU
type CPUMetrics struct {
	UsagePercent float64   `json:"usage_percent" bson:"usage_percent"`
	LoadAvg      []float64 `json:"load_avg" bson:"load_avg"`
}

// MemoryMetrics представляет метрики памяти
type MemoryMetrics struct {
	TotalMB      float64 `json:"total_mb" bson:"total_mb"`
	UsedMB       float64 `json:"used_mb" bson:"used_mb"`
	UsagePercent float64 `json:"usage_percent" bson:"usage_percent"`
}

// DiskMetrics представляет метрики диска
type DiskMetrics struct {
	MountPoint   string  `json:"mount_point" bson:"mount_point"`
	TotalGB      float64 `json:"total_gb" bson:"total_gb"`
	UsedGB       float64 `json:"used_gb" bson:"used_gb"`
	UsagePercent float64 `json:"usage_percent" bson:"usage_percent"`
}

// NetworkMetrics представляет метрики сетевых портов
type NetworkMetrics struct {
	HostID    int        `json:"host_id" bson:"host_id"`
	Timestamp time.Time  `json:"timestamp" bson:"timestamp"`
	Ports     []PortInfo `json:"ports" bson:"ports"`
}

// PortInfo представляет информацию о сетевом порте
type PortInfo struct {
	Protocol  string `json:"protocol" bson:"protocol"`
	LocalPort int    `json:"local_port" bson:"local_port"`
	State     string `json:"state" bson:"state"`
	Process   string `json:"process" bson:"process"`
}
