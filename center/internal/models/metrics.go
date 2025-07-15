package models

import "time"

type Metrics struct {
	ContainersInfo []ContainerInfo `yaml:"containers"`
	HostID         int             `json:"host_id" bson:"host_id"`
	PortsInfo      []PortInfo      `yaml:"ports"`
	ProcessesInfo  []ProcessInfo   `yaml:"processes"`
	SystemMetrics  SystemDetails   `yaml:"system"`
	Timestamp      time.Time       `bson:"timestamp"`
}
