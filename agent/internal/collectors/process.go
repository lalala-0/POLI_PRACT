package collectors

import (
	"agent/internal/models"
	"github.com/shirou/gopsutil/process"
)

// ProcessCollector собирает информацию о процессах
type ProcessCollector struct {
	processes []string // список отслеживаемых процессов
}

func NewProcessCollector(processes []string) *ProcessCollector {
	return &ProcessCollector{
		processes: processes,
	}
}

func (c *ProcessCollector) ChangeConfig(collType CollectorType, newconfig []string) {
	if collType == Process {
		c.processes = newconfig
	}
}

func (c *ProcessCollector) Collect(metrics *models.AgentMetrics) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	var processInfos []models.ProcessInfo

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		// Если указан список процессов и текущего в нем нет - пропускаем
		if len(c.processes) > 0 {
			found := false
			for _, targetName := range c.processes {
				if targetName == name {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		pid := p.Pid
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()

		processInfo := models.ProcessInfo{
			PID:        pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemPercent: float64(memPercent),
		}

		processInfos = append(processInfos, processInfo)
	}

	metrics.Processes = processInfos
	return nil
}
