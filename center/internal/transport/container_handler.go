package api

import (
	"center/internal/models"
	"center/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContainerHandler struct {
	service *services.HostService
}

func NewContainerHandler(service *services.HostService) *ContainerHandler {
	return &ContainerHandler{service: service}
}

// GetContainersByHostID
// @Summary Получить контейнеры для хоста
// @Description Возвращает все контейнеры для указанного хоста
// @Tags Containers
// @Produce json
// @Param id path int true "ID хоста"
// @Success 200 {array} models.Container
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/containers [get]
func (h *ContainerHandler) GetContainersByHostID(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	ctx := c.Request.Context()
	container, err := h.service.ContainerRepo.GetByHostID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, container)
}

// CreateContainer
// @Summary Добавить контейнер для мониторинга
// @Description Добавляет новый контейнер для мониторинга на указанном хосте
// @Tags Containers
// @Accept json
// @Produce json
// @Param id path int true "ID хоста"
// @Param container body models.ContainerInput true "Данные контейнера"
// @Success 201 {object} map[string]int "ID созданного контейнера"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/containers [post]
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
	//log.Println("-----------------------------", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//log.Println("-----------------------------", hostID)
	host, err := h.service.GetHost(ctx, hostID)
	//log.Println("-----------------------------", err)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Println("-----------------------------", host)
	log.Println("-----------------------------", *host)
	err = h.service.SendContainerConfigurationToAgent(ctx, *host)
	log.Println("-----------------------------", err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// DeleteContainer
// @Summary Удалить контейнер из мониторинга
// @Description Удаляет контейнер из списка мониторинга
// @Tags Containers
// @Param id path int true "ID хоста"
// @Param container_id path int true "ID контейнера"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/containers/{container_id} [delete]
func (h *ContainerHandler) DeleteContainer(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("id"))
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
	if err := h.service.ContainerRepo.Delete(ctx, containerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
