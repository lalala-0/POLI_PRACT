package collectors

import (
	"agent/internal/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

// SystemCollector собирает метрики CPU, RAM и дисков
type SystemCollector struct{}

func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

func (c *SystemCollector) Collect(metrics *models.AgentMetrics) error {
	// Сбор CPU метрик
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return err
	}

	// Сбор метрик памяти
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	// Сбор метрик диска
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return err
	}

	metrics.System = models.SystemMetrics{
		CPU: models.CPUMetrics{
			UsagePercent: cpuPercent[0],
		},
		RAM: models.RAMMetrics{
			Total:        memory.Total,
			Used:         memory.Used,
			Free:         memory.Free,
			UsagePercent: memory.UsedPercent,
		},
		Disk: models.DiskMetrics{
			Total:        diskUsage.Total,
			Used:         diskUsage.Used,
			Free:         diskUsage.Free,
			UsagePercent: diskUsage.UsedPercent,
		},
	}

	return nil
}
