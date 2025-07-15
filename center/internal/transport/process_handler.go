package api

import (
	"center/internal/models"
	"net/http"
	"strconv"
	"context"
	"time"


	"github.com/gin-gonic/gin"
)

type Handler struct {
	hostService *services.HostService
}

func NewHandler(hostService *services.HostService) *Handler {
	return &Handler{hostService: hostService}
}

// Host Handlers

// GetHosts возвращает список всех хостов
func (h *Handler) GetHosts(c *gin.Context) {
	ctx := c.Request.Context()
	hosts, err := h.hostService.hostRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hosts)
}

// GetHost возвращает информацию о конкретном хосте
func (h *Handler) GetHostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	host, err := h.hostService.GetHost(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}
	c.JSON(http.StatusOK, host)
}

// CreateHost создает новый хост
func (h *Handler) CreateHost(c *gin.Context) {
	var hostInput models.HostInput
	if err := c.ShouldBindJSON(&hostInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	id, err := h.hostService.CreateHost(ctx, hostInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) UpdateHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var hostInput models.HostInput
	if err := c.ShouldBindJSON(&hostInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	if err := h.hostService.UpdateHost(ctx, id, hostInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.hostService.DeleteHost(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetMasterHost(c *gin.Context) {
	ctx := c.Request.Context()
	host, err := h.hostService.hostRepo.GetMaster(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, host)
}

func (h *Handler) SetMasterHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.hostService.SetMasterHost(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Process Handlers
func (h *Handler) GetProcessesByHostID(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	ctx := c.Request.Context()
	processes, err := h.hostService.processRepo.GetByHostID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, processes)
}

func (h *Handler) CreateProcess(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var processInput models.ProcessInput
	if err := c.ShouldBindJSON(&processInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	id, err := h.hostService.AddProcess(ctx, hostID, processInput.ProcessName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) DeleteProcess(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	processID, err := strconv.Atoi(c.Param("process_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid process ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.hostService.processRepo.Delete(ctx, processID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}


// Container Handlers
func (h *Handler) GetContainerByHostID(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	ctx := c.Request.Context()
	container, err := h.hostService.containerRepo.GetByHostID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, container)
}

func (h *Handler) CreateContainer(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var containerInput models.ContainerInput
	if err := c.ShouldBindJSON(&containerInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	id, err := h.hostService.AddContainer(ctx, hostID, containerInput.ProcessName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) DeleteContainer(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	containerID, err := strconv.Atoi(c.Param("container_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid container ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.hostService.containerRepo.Delete(ctx, containerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}


// Alert Handlers
func (h *Handler) GetAlertsByHostID(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	ctx := c.Request.Context()
	alerts, err := h.hostService.alertRepo.GetByHostID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

func (h *Handler) CreateAlert(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var alertInput models.AlertInput
	if err := c.ShouldBindJSON(&alertInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	id, err := h.hostService.CreateAlertRule(ctx, hostID, alertInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ReceiveMetrics принимает метрики от агента
func (h *Handler) ReceiveMetrics(c *gin.Context) {
	var metricsData struct {
		System    *models.SystemMetrics    `json:"system"`
		Processes *models.ProcessMetrics   `json:"processes"`
		Containers *models.ContainerMetrics `json:"containers"`
		Network   *models.NetworkMetrics   `json:"network"`
	}

	if err := c.ShouldBindJSON(&metricsData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	if metricsData.System != nil {
		if err := h.hostService.SaveSystemMetrics(ctx, metricsData.System); err != nil {
			log.Printf("Error saving system metrics: %v", err)
		}
	}

	if metricsData.Processes != nil {
		if err := h.hostService.SaveProcessMetrics(ctx, metricsData.Processes); err != nil {
			log.Printf("Error saving process metrics: %v", err)
		}
	}

	if metricsData.Containers != nil {
		if err := h.hostService.SaveContainerMetrics(ctx, metricsData.Containers); err != nil {
			log.Printf("Error saving container metrics: %v", err)
		}
	}

	if metricsData.Network != nil {
		if err := h.hostService.SaveNetworkMetrics(ctx, metricsData.Network); err != nil {
			log.Printf("Error saving network metrics: %v", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "metrics received"})
}

