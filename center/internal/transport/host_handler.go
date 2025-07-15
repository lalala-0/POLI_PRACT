package api

import (
	"center/internal/models"
	"net/http"
	"strconv"
	"context"
	"time"


	"github.com/gin-gonic/gin"
)


type HostHandler struct {
	service *services.HostService
}

func NewHostHandler(service *services.HostService) *HostHandler {
	return &HostHandler{service: service}
}

// GetHosts возвращает список всех хостов
func (h *HostHandler) GetHosts(c *gin.Context) {
	ctx := c.Request.Context()
	hosts, err := h.service.hostRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hosts)
}

// GetHost возвращает информацию о конкретном хосте
func (h *HostHandler) GetHostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	host, err := h.service.GetHost(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Host not found"})
		return
	}
	c.JSON(http.StatusOK, host)
}

// CreateHost создает новый хост
func (h *HostHandler) CreateHost(c *gin.Context) {
	var hostInput models.HostInput
	if err := c.ShouldBindJSON(&hostInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	id, err := h.service.CreateHost(ctx, hostInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *HostHandler) UpdateHost(c *gin.Context) {
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
	if err := h.service.UpdateHost(ctx, id, hostInput); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *HostHandler) DeleteHost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.service.DeleteHost(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *HostHandler) GetMasterHost(c *gin.Context) {
	ctx := c.Request.Context()
	host, err := h.service.hostRepo.GetMaster(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, host)
}

func (h *HostHandler) SetMasterHost(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    ctx := c.Request.Context()
    if err := h.service.SetMasterHost(ctx, id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}