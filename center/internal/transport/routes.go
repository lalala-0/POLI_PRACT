package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// API группа
	api := router.Group("/api")
	{
		// Хосты
		hosts := api.Group("/hosts")
		{
			hosts.GET("", GetHosts)
			hosts.GET("/:id", GetHostByID)
			hosts.POST("", CreateHost)
			hosts.PUT("/:id", UpdateHost)
			hosts.DELETE("/:id", DeleteHost)
			hosts.GET("/master", GetMasterHost)
			hosts.PUT("/:id/master", SetMasterHost)

			// Процессы хоста
			hosts.GET("/:id/processes", GetProcessesByHostID)
			hosts.POST("/:id/processes", CreateProcess)
			hosts.DELETE("/:id/processes/:process_id", DeleteProcess)

			// Контейнеры хоста
			hosts.GET("/:id/containers", GetContainersByHostID)
			hosts.POST("/:id/containers", CreateContainer)
			hosts.DELETE("/:id/containers/:container_id", DeleteContainer)

			// Правила оповещений хоста
			hosts.GET("/:id/alerts", GetAlertsByHostID)
			hosts.POST("/:id/alerts", CreateAlert)
			hosts.PUT("/:id/alerts/:alert_id", UpdateAlert)
			hosts.DELETE("/:id/alerts/:alert_id", DeleteAlert)
			hosts.PATCH("/:id/alerts/:alert_id/status", EnableDisableAlert)
		}

		// Метрики
		metrics := api.Group("/metrics")
		{
			metrics.POST("", ReceiveMetrics)                         // Приём метрик от агентов
			metrics.GET("/:host_id/system", GetSystemMetrics)        // Системные метрики
			metrics.GET("/:host_id/processes", GetProcessMetrics)    // Метрики процессов
			metrics.GET("/:host_id/containers", GetContainerMetrics) // Метрики контейнеров
			metrics.GET("/:host_id/network", GetNetworkMetrics)      // Сетевые метрики
		}

		// Проверка состояния системы
		api.GET("/health", GetHealth)
	}
}
