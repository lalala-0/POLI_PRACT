package api

import (
	"center/internal/models"
	"center/internal/services"
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

// GetProcessesByHostID godoc
// @Summary Получить процессы для хоста
// @Description Возвращает все процессы для указанного хоста
// @Tags Processes
// @Produce json
// @Param id path int true "ID хоста"
// @Success 200 {array} models.Process
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/processes [get]
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

// CreateProcess
// @Summary Добавить процесс для мониторинга
// @Description Добавляет новый процесс для мониторинга на указанном хосте
// @Tags Processes
// @Accept json
// @Produce json
// @Param id path int true "ID хоста"
// @Param process body models.ProcessInput true "Данные процесса"
// @Success 201 {object} map[string]int "ID созданного процесса"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/processes [post]
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
	host, err := h.service.GetHost(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err = h.service.SendProcessConfigurationToAgent(ctx, *host); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeleteProcess
// @Summary Удалить процесс из мониторинга
// @Description Удаляет процесс из списка мониторинга
// @Tags Processes
// @Param id path int true "ID хоста"
// @Param process_id path int true "ID процесса"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/processes/{process_id} [delete]
func (h *ProcessHandler) DeleteProcess(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
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
