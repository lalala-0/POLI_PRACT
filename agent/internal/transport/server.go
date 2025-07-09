package transport

import (
	"agent/internal/models"
	"agent/internal/services"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server представляет HTTP-сервер агента мониторинга
type Server struct {
	port           string
	router         *gin.Engine
	metricsCh      <-chan models.AgentMetrics
	metricsService *services.MetricsService
	lastMetrics    models.AgentMetrics
}

// NewServer создает новый экземпляр HTTP-сервера
func NewServer(port string, metricsCh <-chan models.AgentMetrics, metricsService *services.MetricsService) *Server {
	router := gin.Default()
	server := &Server{
		port:           port,
		router:         router,
		metricsCh:      metricsCh,
		metricsService: metricsService,
	}

	server.setupRoutes()
	return server
}

// setupRoutes настраивает маршруты для HTTP-сервера
func (s *Server) setupRoutes() {
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/metrics", s.getMetrics)
	s.router.GET("/metrics/system", s.getSystemMetrics)
	s.router.GET("/metrics/processes", s.getProcessMetrics)
	s.router.GET("/metrics/network", s.getNetworkMetrics)
	s.router.GET("/metrics/containers", s.getContainerMetrics)

	// API для обновления конфигурации
	s.router.POST("/config/processes", s.updateProcessConfig)
	s.router.POST("/config/containers", s.updateContainerConfig)
}

// Start запускает HTTP-сервер и слушает обновления метрик
func (s *Server) Start(ctx context.Context) {
	// Запуск HTTP-сервера
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	go func() {
		log.Printf("HTTP-сервер запущен на порту %s", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска HTTP-сервера: %v", err)
		}
	}()

	// Получение метрик из канала
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case metrics := <-s.metricsCh:
				s.lastMetrics = metrics
			}
		}
	}()

	// Ждем сигнал завершения
	<-ctx.Done()
	log.Println("Останавливаем HTTP-сервер...")

	// Graceful shutdown с таймаутом
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке HTTP-сервера: %v", err)
	}

	log.Println("HTTP-сервер остановлен")
}
