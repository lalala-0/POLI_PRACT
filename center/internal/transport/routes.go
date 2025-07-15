package api

import (
	"github.com/gin-gonic/gin"
)


// Handler структура для группировки обработчиков
type Handler struct {
	HostHandler     *HostHandler
	MetricHandler   *MetricHandler
	ProcessHandler  *ProcessHandler
	ContainerHandler *ContainerHandler
	AlertHandler    *AlertHandler
}

func SetupRoutes(router *gin.Engine, handler *Handler) {
	// API группа
	api := router.Group("/api")
	{
		// Хосты
		hosts := api.Group("/hosts")
		{
			hosts.GET("", handler.HostHandler.GetHosts)
			hosts.GET("/:id", handler.HostHandler.GetHostByID)
			hosts.POST("", handler.HostHandler.CreateHost)
			hosts.PUT("/:id", handler.HostHandler.UpdateHost)
			hosts.DELETE("/:id", handler.HostHandler.DeleteHost)
			hosts.GET("/master", handler.HostHandler.GetMasterHost)
			hosts.PUT("/:id/master", handler.HostHandler.SetMasterHost)

			// Процессы хоста
			hosts.GET("/:id/processes", handler.ProcessHandler.GetProcessesByHostID)
			hosts.POST("/:id/processes", handler.ProcessHandler.CreateProcess)
			hosts.DELETE("/:id/processes/:process_id", handler.ProcessHandler.DeleteProcess)

			// Контейнеры хоста
			hosts.GET("/:id/containers", handler.ContainerHandler.GetContainersByHostID)
			hosts.POST("/:id/containers", handler.ContainerHandler.CreateContainer)
			hosts.DELETE("/:id/containers/:container_id", handler.ContainerHandler.DeleteContainer)

			// Правила оповещений хоста
			hosts.GET("/:id/alerts", handler.AlertHandler.GetAlertsByHostID)
			hosts.POST("/:id/alerts", handler.AlertHandler.CreateAlert)
			hosts.PUT("/:id/alerts/:alert_id", handler.AlertHandler.UpdateAlert)
			hosts.DELETE("/:id/alerts/:alert_id", handler.AlertHandler.DeleteAlert)
			hosts.PATCH("/:id/alerts/:alert_id/status", handler.AlertHandler.EnableDisableAlert)
		}

		// Метрики
		metrics := api.Group("/metrics")
		{
			metrics.POST("", handler.MetricHandler.ReceiveMetrics)
			metrics.GET("/:host_id/system", handler.MetricHandler.GetSystemMetrics)
			metrics.GET("/:host_id/processes", handler.MetricHandler.GetProcessMetrics)
			metrics.GET("/:host_id/containers", handler.MetricHandler.GetContainerMetrics)
			metrics.GET("/:host_id/network", handler.MetricHandler.GetNetworkMetrics)
		}

		// Проверка состояния системы
		api.GET("/health", handler.MetricHandler.GetHealth)
	}
}
