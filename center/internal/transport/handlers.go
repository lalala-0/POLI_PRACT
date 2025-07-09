package api

import (
	"POLI_PRACT/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetHosts возвращает список всех хостов
func GetHosts(c *gin.Context) {
	hosts, err := services.GetAllHosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hosts)
}

// CreateHost создает новый хост
func CreateHost(c *gin.Context) {
	var host models.Host
	if err := c.ShouldBindJSON(&host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := services.CreateHost(host)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	host.ID = id
	c.JSON(http.StatusCreated, host)
}

// GetHost возвращает информацию о конкретном хосте
func GetHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	host, err := services.GetHostByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if host == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}

	c.JSON(http.StatusOK, host)
}

// ReceiveMetrics принимает метрики от агента
func ReceiveMetrics(c *gin.Context) {
	var metrics models.MetricData
	if err := c.ShouldBindJSON(&metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SaveMetrics(metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "metrics received"})
}
