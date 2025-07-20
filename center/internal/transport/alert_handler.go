package api

import (
	"center/internal/models"
	"center/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AlertHandler обработчик оповещений
type AlertHandler struct {
	hostService  *services.HostService
	alertService *services.AlertNotifierService
}

func NewAlertHandler(hostService *services.HostService, alertService *services.AlertNotifierService) *AlertHandler {
	return &AlertHandler{
		hostService:  hostService,
		alertService: alertService,
	}
}

// GetAlertsByHostID
// @Summary Получить правила оповещений для хоста
// @Description Возвращает все правила оповещений для указанного хоста
// @Tags Alerts
// @Produce json
// @Param id path int true "ID хоста"
// @Success 200 {array} models.AlertRule
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/alerts [get]
func (h *AlertHandler) GetAlertsByHostID(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	ctx := c.Request.Context()
	alerts, err := h.hostService.AlertRepo.GetByHostID(ctx, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}

// CreateAlert
// @Summary Создать правило оповещения
// @Description Создает новое правило оповещения для указанного хоста
// @Tags Alerts
// @Accept json
// @Produce json
// @Param id path int true "ID хоста"
// @Param alert body models.AlertInput true "Данные правила оповещения"
// @Success 201 {object} map[string]int "ID созданного правила"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/alerts [post]
func (h *AlertHandler) CreateAlert(c *gin.Context) {
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
	log.Println("                               Create alert")
	// Принудительное обновление кэша
	h.alertService.InvalidateAlertRules()

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateAlert
// @Summary Обновить правило оповещения
// @Description Обновляет существующее правило оповещения
// @Tags Alerts
// @Accept json
// @Param id path int true "ID хоста"
// @Param alert_id path int true "ID правила оповещения"
// @Param alert body models.AlertInput true "Обновленные данные правила"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/alerts/{alert_id} [put]
func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	alertID, err := strconv.Atoi(c.Param("alert_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	ctx := c.Request.Context()
	alert, err := h.hostService.AlertRepo.GetByID(ctx, alertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if alert == nil || alert.HostID != hostID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	var alertInput models.AlertInput
	if err := c.ShouldBindJSON(&alertInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert.MetricName = alertInput.MetricName
	alert.ThresholdValue = alertInput.ThresholdValue
	alert.Condition = alertInput.Condition
	alert.Enabled = alertInput.Enabled

	if err := h.hostService.AlertRepo.Update(ctx, alert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("                               Update alert")
	// Принудительное обновление кэша
	h.alertService.InvalidateAlertRules()
	c.Status(http.StatusNoContent)
}

// DeleteAlert
// @Summary Удалить правило оповещения
// @Description Удаляет правило оповещения по ID
// @Tags Alerts
// @Param id path int true "ID хоста"
// @Param alert_id path int true "ID правила оповещения"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/alerts/{alert_id} [delete]
func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	alertID, err := strconv.Atoi(c.Param("alert_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	ctx := c.Request.Context()
	alert, err := h.hostService.AlertRepo.GetByID(ctx, alertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if alert == nil || alert.HostID != hostID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	if err := h.hostService.AlertRepo.Delete(ctx, alertID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("                               Delete alert")
	// Принудительное обновление кэша
	h.alertService.InvalidateAlertRules()

	c.Status(http.StatusNoContent)
}

// EnableDisableAlert
// @Summary Включить/выключить правило оповещения
// @Description Изменяет статус активности правила оповещения
// @Tags Alerts
// @Accept json
// @Param id path int true "ID хоста"
// @Param alert_id path int true "ID правила оповещения"
// @Param status body object{enabled=bool} true "Статус активности"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /hosts/{id}/alerts/{alert_id}/status [patch]
func (h *AlertHandler) EnableDisableAlert(c *gin.Context) {
	hostID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	alertID, err := strconv.Atoi(c.Param("alert_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	var status struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	alert, err := h.hostService.AlertRepo.GetByID(ctx, alertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if alert == nil || alert.HostID != hostID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	if err := h.hostService.AlertRepo.SetEnabled(ctx, alertID, status.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("                               En/Dis alert")
	// Принудительное обновление кэша
	h.alertService.InvalidateAlertRules()
	c.Status(http.StatusNoContent)
}
