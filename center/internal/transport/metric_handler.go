package api

import (
	"center/internal/models"
	"center/internal/services"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	service *services.HostService
}

func NewMetricHandler(service *services.HostService) *MetricHandler {
	return &MetricHandler{service: service}
}

// ReceiveMetrics принимает метрики от агента
func (h *MetricHandler) ReceiveMetrics(c *gin.Context) {
	var metricsData struct {
		System     *models.SystemMetrics    `json:"system"`
		Processes  *models.ProcessMetrics   `json:"processes"`
		Containers *models.ContainerMetrics `json:"containers"`
		Network    *models.NetworkMetrics   `json:"network"`
	}

	if err := c.ShouldBindJSON(&metricsData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	if metricsData.System != nil {
		if err := h.service.SaveSystemMetrics(ctx, metricsData.System); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save system metrics"})
			return
		}
	}

	if metricsData.Processes != nil {
		if err := h.service.SaveProcessMetrics(ctx, metricsData.Processes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save process metrics"})
			return
		}
	}

	if metricsData.Containers != nil {
		if err := h.service.SaveContainerMetrics(ctx, metricsData.Containers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save container metrics"})
			return
		}
	}

	if metricsData.Network != nil {
		if err := h.service.SaveNetworkMetrics(ctx, metricsData.Network); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save network metrics"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "metrics received"})
}

func (h *MetricHandler) GetSystemMetrics(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("host_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	from, to := time.Now().Add(time.Duration(-14*24)*time.Hour), time.Now()
	//// Параметры диапазона времени
	//from, to, err := parseTimeRange(c)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	ctx := c.Request.Context()
	metrics, err := h.service.MetricRepo.GetSystemMetricsInRange(ctx, hostID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *MetricHandler) GetProcessMetrics(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("host_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	from, to := time.Now().Add(time.Duration(-14*24)*time.Hour), time.Now()
	//// Параметры диапазона времени
	//from, to, err := parseTimeRange(c)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	ctx := c.Request.Context()
	metrics, err := h.service.MetricRepo.GetProcessMetricsInRange(ctx, hostID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *MetricHandler) GetContainerMetrics(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("host_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	from, to := time.Now().Add(time.Duration(-14*24)*time.Hour), time.Now()
	//// Параметры диапазона времени
	//from, to, err := parseTimeRange(c)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	ctx := c.Request.Context()
	metrics, err := h.service.MetricRepo.GetContainerMetricsInRange(ctx, hostID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *MetricHandler) GetNetworkMetrics(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("host_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	from, to := time.Now().Add(time.Duration(-5*24)*time.Hour), time.Now()
	//// Параметры диапазона времени
	//from, to, err := parseTimeRange(c)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	ctx := c.Request.Context()
	metrics, err := h.service.MetricRepo.GetNetworkMetricsInRange(ctx, hostID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetMetrics возвращает агрегированные метрики по всем хостам
func (h *MetricHandler) GetMetrics(c *gin.Context) {
	//// Параметры диапазона времени
	//from, to, err := parseTimeRange(c)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	from, to := time.Now().Add(time.Duration(-14*24)*time.Hour), time.Now()

	ctx := c.Request.Context()

	// Получаем список всех хостов
	hosts, err := h.service.GetAllHosts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get hosts"})
		return
	}

	// Собираем последние метрики для каждого хоста
	response := make(map[string]interface{})
	for _, host := range hosts {
		hostMetrics := gin.H{}

		// Системные метрики
		if systemMetrics, err := h.service.MetricRepo.GetSystemMetricsInRange(ctx, host.ID, from, to); err == nil {
			hostMetrics["system"] = systemMetrics
		}

		// Метрики процессов
		if processMetrics, err := h.service.MetricRepo.GetProcessMetricsInRange(ctx, host.ID, from, to); err == nil {
			hostMetrics["processes"] = processMetrics
		}

		// Метрики контейнеров
		if containerMetrics, err := h.service.MetricRepo.GetContainerMetricsInRange(ctx, host.ID, from, to); err == nil {
			hostMetrics["containers"] = containerMetrics
		}

		// Сетевые метрики
		if networkMetrics, err := h.service.MetricRepo.GetNetworkMetricsInRange(ctx, host.ID, from, to); err == nil {
			hostMetrics["network"] = networkMetrics
		}

		response[host.Hostname] = hostMetrics
	}

	c.JSON(http.StatusOK, response)
}

// GetHostMetrics возвращает все метрики для конкретного хоста
func (h *MetricHandler) GetHostMetrics(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("host_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	//// Параметры диапазона времени
	//from, to, err := parseTimeRange(c)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}

	from, to := time.Now().Add(time.Duration(-14*24)*time.Hour), time.Now()
	//log.Printf("-----------------------------Getting host metrics from %s to %s", from, to)

	ctx := c.Request.Context()
	response := gin.H{}

	// Системные метрики
	if systemMetrics, err := h.service.MetricRepo.GetSystemMetricsInRange(ctx, hostID, from, to); err == nil {
		response["system"] = systemMetrics
	} else {
		log.Printf("Error getting system metrics: %v", err)
	}

	// Метрики процессов
	if processMetrics, err := h.service.MetricRepo.GetProcessMetricsInRange(ctx, hostID, from, to); err == nil {
		response["processes"] = processMetrics
	} else {
		log.Printf("Error getting process metrics: %v", err)
	}

	// Метрики контейнеров
	if containerMetrics, err := h.service.MetricRepo.GetContainerMetricsInRange(ctx, hostID, from, to); err == nil {
		response["containers"] = containerMetrics
	} else {
		log.Printf("Error getting container metrics: %v", err)
	}

	// Сетевые метрики
	if networkMetrics, err := h.service.MetricRepo.GetNetworkMetricsInRange(ctx, hostID, from, to); err == nil {
		response["network"] = networkMetrics
	} else {
		log.Printf("Error getting network metrics: %v", err)
	}
	//log.Printf("-------------response: %+v", response)

	c.JSON(http.StatusOK, response)
}

func (h *MetricHandler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// parseTimeRange парсит параметры времени из запроса
func parseTimeRange(c *gin.Context) (from, to time.Time, err error) {
	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")
	if fromStr == "" || toStr == "" {
		return from, to, errors.New("both from and to parameters are required")
	}

	from, err = time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return from, to, errors.New("invalid from parameter format, use RFC3339")
	}

	to, err = time.Parse(time.RFC3339, toStr)
	if err != nil {
		return from, to, errors.New("invalid to parameter format, use RFC3339")
	}

	return from, to, nil
}
