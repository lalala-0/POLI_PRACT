package models

import "time"

// Container представляет контейнер, который нужно мониторить на хосте
type Container struct {
	ID            int    `json:"id" db:"id"`
	HostID        int    `json:"host_id" db:"host_id"`
	ContainerName string `json:"container_name" binding:"required" db:"container_name"`
}

// ContainerInput представляет данные для добавления контейнера
type ContainerInput struct {
	ContainerName string `json:"container_name" binding:"required"`
}

// ContainerMetrics представляет метрики контейнеров
type ContainerMetrics struct {
	HostID     int             `json:"host_id" bson:"host_id"`
	Timestamp  time.Time       `json:"timestamp" bson:"timestamp"`
	Containers []ContainerInfo `json:"containers" bson:"containers"`
}

// ContainerInfo представляет информацию о контейнере
type ContainerInfo struct {
	Name          string  `json:"name" bson:"name"`
	ID            string  `json:"id" bson:"id"`
	Image         string  `json:"image" bson:"image"`
	Status        string  `json:"status" bson:"status"`
	CPUPercent    float64 `json:"cpu_percent" bson:"cpu_percent"`
	MemoryPercent float64 `json:"mem_percent" bson:"mem_percent"`
}
