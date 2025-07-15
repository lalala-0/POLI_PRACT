package models

import (
	"time"
)

// Host представляет информацию о хосте для мониторинга
type Host struct {
	ID        int    `json:"id" db:"id"`
	Hostname  string `json:"hostname" binding:"required" db:"hostname"`
	IPAddress string `json:"ip_address" binding:"required" db:"ip_address"`
	AgentPort int    `json:"agent_port" db:"agent_port"`
	Priority  int    `json:"priority" db:"priority"`
	IsMaster  bool   `json:"is_master" db:"is_master"`
	Status    string `json:"status" db:"status"`
	//LastCheck time.Time `json:"last_check" db:"last_check"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// HostInput представляет данные для создания/обновления хоста
type HostInput struct {
	Hostname  string `json:"hostname" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	AgentPort int    `json:"agent_port" binding:"required"`
	Priority  int    `json:"priority"`
}

// SystemMetrics представляет основные метрики хоста
type SystemMetrics struct {
	HostID    int           `json:"host_id" bson:"host_id"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
	System    SystemDetails `json:"system" bson:"system"`
}

type SystemDetails struct {
	CPU  CPUInfo  `json:"cpu" bson:"cpu"`
	RAM  RAMInfo  `json:"ram" bson:"ram"`
	Disk DiskInfo `json:"disk" bson:"disk"`
}

type CPUInfo struct {
	UsagePercent float64 `json:"usage_percent" bson:"usage_percent"`
}

type RAMInfo struct {
	Total        uint64  `json:"total" bson:"total"`
	Used         uint64  `json:"used" bson:"used"`
	Free         uint64  `json:"free" bson:"free"`
	UsagePercent float64 `json:"usage_percent" bson:"usage_percent"`
}

type DiskInfo struct {
	Total        uint64  `json:"total" bson:"total"`
	Used         uint64  `json:"used" bson:"used"`
	Free         uint64  `json:"free" bson:"free"`
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
