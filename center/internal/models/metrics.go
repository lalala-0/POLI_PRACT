package models

import "time"

type Metrics struct {
	HostID         int             `json:"host_id"`
	Timestamp      time.Time       `json:"timestamp"`
	SystemMetrics  SystemDetails   `json:"system,omitempty"`
	ProcessesInfo  []ProcessInfo   `json:"processes,omitempty"`
	PortsInfo      []PortInfo      `json:"ports,omitempty"`
	ContainersInfo []ContainerInfo `json:"containers,omitempty"`
}

// HostMetricsResponse представляет все метрики хоста за период времени
type HostMetricsResponse struct {
	System     []SystemMetrics    `json:"system"`
	Processes  []ProcessMetrics   `json:"processes"`
	Containers []ContainerMetrics `json:"containers"`
	Network    []NetworkMetrics   `json:"network"`
}
