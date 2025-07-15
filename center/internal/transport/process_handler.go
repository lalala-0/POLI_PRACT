package api

import (
	"center/internal/models"
	"center/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProcessHandler struct {
	service *services.HostService
}

func NewProcessHandler(service *services.HostService) *ProcessHandler {
	return &ProcessHandler{service: service}
}

// Process ProcessHandlers
func (h *ProcessHandler) GetProcessesByHostID(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	ctx := c.Request.Context()
	processes, err := h.service.ProcessRepo.GetByHostID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, processes)
}

func (h *ProcessHandler) CreateProcess(c *gin.Context) {
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
	id, err := h.service.AddProcess(ctx, hostID, processInput.ProcessName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ProcessHandler) DeleteProcess(c *gin.Context) {
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
	if err := h.service.ProcessRepo.Delete(ctx, processID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
