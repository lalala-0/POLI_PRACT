package collectors

import (
	"agent/internal/models"
	"context"
	"encoding/json"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerCollector собирает метрики Docker-контейнеров
type DockerCollector struct {
	client     *client.Client
	containers []string // список отслеживаемых контейнеров
}

func NewDockerCollector(containers []string) (*DockerCollector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerCollector{
		client:     cli,
		containers: containers,
	}, nil
}

func (c *DockerCollector) Collect(metrics *models.AgentMetrics) error {
	ctx := context.Background()

	containers, err := c.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	var containerInfos []models.ContainerInfo

	for _, container := range containers {
		// Если указан список контейнеров и текущего в нем нет - пропускаем
		if len(c.containers) > 0 {
			found := false
			cleanName := ""
			for _, name := range container.Names {
				cleanName = strings.TrimPrefix(name, "/")
				for _, targetName := range c.containers {
					if targetName == cleanName {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}

		// Получаем статистику контейнера
		stats, err := c.client.ContainerStats(ctx, container.ID, false)
		if err != nil {
			continue
		}
		defer stats.Body.Close()

		var statsJSON types.StatsJSON
		if err := json.NewDecoder(stats.Body).Decode(&statsJSON); err != nil {
			continue
		}

		// Рассчитываем проценты использования CPU и памяти
		cpuPercent := calculateCPUPercent(&statsJSON)
		memoryUsage := statsJSON.MemoryStats.Usage
		memoryLimit := statsJSON.MemoryStats.Limit

		memPercent := 0.0
		if memoryLimit > 0 {
			memPercent = float64(memoryUsage) / float64(memoryLimit) * 100.0
		}

		// Получаем статус контейнера
		status := "unknown"
		if container.State != "" {
			status = container.State
		}

		containerInfos = append(containerInfos, models.ContainerInfo{
			ID:         container.ID[:12],
			//Name:       strings.TrimPrefix(container.Names[0], "/"),
			Name:		cleanName,
			Image:      container.Image,
			Status:     status,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
		})
	}

	metrics.Containers = containerInfos
	return nil
}

func calculateCPUPercent(stats types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		numCPUs := len(stats.CPUStats.CPUUsage.PercpuUsage)
		if numCPUs == 0 {
			numCPUs = 1
		}
		return (cpuDelta / systemDelta) * float64(numCPUs) * 100.0
	}

	return 0.0
}
