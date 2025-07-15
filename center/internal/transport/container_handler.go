package api

import (
	"center/internal/models"
	"net/http"
	"strconv"
	"context"
	"time"


	"github.com/gin-gonic/gin"
)
type ContainerHandler struct {
	service *services.HostService
}

func NewContainerHandler(service *services.HostService) *ContainerHandler {
	return &ContainerHandler{service: service}
}

// Container Handlers
func (h *ContainerHandler) GetContainerByHostID(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	ctx := c.Request.Context()
	container, err := h.service.containerRepo.GetByHostID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, container)
}

func (h *ContainerHandler) CreateContainer(c *gin.Context) {
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
	id, err := h.service.AddContainer(ctx, hostID, containerInput.ContainerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

//id - host
//container_id - container
func (h *ContainerHandler) DeleteContainer(c *gin.Context) {
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
	if err := h.service.containerRepo.Delete(ctx, containerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
