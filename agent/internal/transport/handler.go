package transport

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// healthCheck проверяет работоспособность агента
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Агент работает",
	})
}

// getMetrics возвращает все собранные метрики
func (s *Server) getMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, s.lastMetrics)
}

// getSystemMetrics возвращает только системные метрики
func (s *Server) getSystemMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"host_id":   s.lastMetrics.HostID,
		"timestamp": s.lastMetrics.Timestamp,
		"system":    s.lastMetrics.System,
	})
}

// getProcessMetrics возвращает только метрики процессов
func (s *Server) getProcessMetrics(c *gin.Context) {
	if !s.metricsService.IsProcessConfigSet() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список отслеживаемых процессов не настроен",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host_id":   s.lastMetrics.HostID,
		"timestamp": s.lastMetrics.Timestamp,
		"processes": s.lastMetrics.Processes,
	})
}

// getNetworkMetrics возвращает только сетевые метрики
func (s *Server) getNetworkMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"host_id":   s.lastMetrics.HostID,
		"timestamp": s.lastMetrics.Timestamp,
		"ports":     s.lastMetrics.Ports,
	})
}

// getContainerMetrics возвращает только метрики контейнеров
func (s *Server) getContainerMetrics(c *gin.Context) {
	if !s.metricsService.IsContainerConfigSet() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список отслеживаемых контейнеров не настроен",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host_id":    s.lastMetrics.HostID,
		"timestamp":  s.lastMetrics.Timestamp,
		"containers": s.lastMetrics.Containers,
	})
}

// updateProcessConfig обновляет список отслеживаемых процессов
func (s *Server) updateProcessConfig(c *gin.Context) {
	var config struct {
		Processes []string `json:"processes"`
	}

	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	if len(config.Processes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список процессов не может быть пустым",
		})
		return
	}

	// Обновляем конфигурацию через сервис
	if err := s.metricsService.UpdateProcessConfig(config.Processes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить конфигурацию",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Конфигурация процессов обновлена",
	})
}

// updateContainerConfig обновляет список отслеживаемых контейнеров
func (s *Server) updateContainerConfig(c *gin.Context) {
	var config struct {
		Containers []string `json:"containers"`
	}

	if err := c.BindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Некорректный формат данных",
		})
		return
	}

	if len(config.Containers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Список контейнеров не может быть пустым",
		})
		return
	}

	if err := s.metricsService.UpdateContainerConfig(config.Containers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Не удалось обновить конфигурацию",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Конфигурация контейнеров обновлена",
	})
}
