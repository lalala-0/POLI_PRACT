package api

import (
	"center/internal/models"
	"center/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @title Monitoring Center API
// @version 1.0
// @description API for monitoring hosts and containers
// @host localhost:8080
// @BasePath /api
// @schemes http

type HostHandler struct {
	service *services.HostService
}

func NewHostHandler(service *services.HostService) *HostHandler {
	return &HostHandler{service: service}
}

// GetHosts возвращает список всех хостов
// @Summary Получить список всех хостов
// @Description Возвращает список всех зарегистрированных хостов
// @Tags Hosts
// @Produce json
// @Success 200 {array} models.Host
// @Failure 500 {object} map[string]string
// @Router /hosts [get]
func (h *HostHandler) GetHosts(c *gin.Context) {
	ctx := c.Request.Context()
	hosts, err := h.service.HostRepo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hosts)
}

// GetHostByID возвращает информацию о конкретном хосте
// @Summary Получить хост по ID
// @Description Возвращает информацию о хосте по его ID
// @Tags Hosts
// @Produce json
// @Param id path int true "ID хоста"
// @Success 200 {object} models.Host
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /hosts/{id} [get]
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
// @Summary Создать новый хост
// @Description Добавляет новый хост в систему
// @Tags Hosts
// @Accept json
// @Produce json
// @Param host body models.HostInput true "Данные хоста"
// @Success 201 {object} map[string]int "ID созданного хоста"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts [post]
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

// UpdateHost
// @Summary Обновить данные хоста
// @Description Обновляет информацию о существующем хосте
// @Tags Hosts
// @Accept json
// @Param id path int true "ID хоста"
// @Param host body models.HostInput true "Обновленные данные хоста"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id} [put]
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

// DeleteHost
// @Summary Удалить хост
// @Description Удаляет хост из системы по ID
// @Tags Hosts
// @Param id path int true "ID хоста"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id} [delete]
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

// GetMasterHost
// @Summary Получить мастер-хост
// @Description Возвращает информацию о текущем мастер-хосте
// @Tags Hosts
// @Produce json
// @Success 200 {object} models.Host
// @Failure 500 {object} map[string]string
// @Router /hosts/master [get]
func (h *HostHandler) GetMasterHost(c *gin.Context) {
	ctx := c.Request.Context()
	host, err := h.service.HostRepo.GetMaster(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, host)
}

// SetMasterHost
// @Summary Установить мастер-хост
// @Description Назначает указанный хост мастер-хостом
// @Tags Hosts
// @Param id path int true "ID хоста"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/master [put]
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
