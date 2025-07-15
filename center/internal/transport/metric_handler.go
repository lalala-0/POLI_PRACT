package api

import (
	"center/internal/models"
	"center/internal/services"
	"context"
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

	// Параметры диапазона времени
	from, to, err := parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Параметры диапазона времени
	from, to, err := parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Параметры диапазона времени
	from, to, err := parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Параметры диапазона времени
	from, to, err := parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	metrics, err := h.service.MetricRepo.GetNetworkMetricsInRange(ctx, hostID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
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
