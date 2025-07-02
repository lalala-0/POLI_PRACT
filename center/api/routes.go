package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Маршруты для хостов
	hosts := api.Group("/hosts")
	{
		hosts.GET("", GetHosts)
		hosts.GET("/:id", GetHost)
		hosts.POST("", CreateHost)
		hosts.PUT("/:id", UpdateHost)
		hosts.DELETE("/:id", DeleteHost)
		hosts.GET("/master", GetMasterHost)
		hosts.PUT("/:id/master", SetMasterHost)

		// Маршруты для процессов
		hosts.GET("/:id/processes", GetHostProcesses)
		hosts.POST("/:id/processes", AddHostProcess)
		hosts.DELETE("/:id/processes/:process_id", DeleteHostProcess)

		// Маршруты для контейнеров
		hosts.GET("/:id/containers", GetHostContainers)
		hosts.POST("/:id/containers", AddHostContainer)
		hosts.DELETE("/:id/containers/:container_id", DeleteHostContainer)

		// Маршруты для правил оповещений
		hosts.GET("/:id/alerts", GetHostAlerts)
		hosts.POST("/:id/alerts", AddHostAlert)
		hosts.PUT("/:id/alerts/:alert_id", UpdateHostAlert)
		hosts.DELETE("/:id/alerts/:alert_id", DeleteHostAlert)
	}

	// Маршруты для метрик
	metrics := api.Group("/metrics")
	{
		metrics.POST("", ReceiveMetrics)
		metrics.GET("/:host_id", GetHostMetrics)
		metrics.GET("/:host_id/processes", GetHostProcessMetrics)
		metrics.GET("/:host_id/containers", GetHostContainerMetrics)
		metrics.GET("/:host_id/network", GetHostNetworkMetrics)
	}
}
