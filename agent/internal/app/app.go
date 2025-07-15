package app

import (
	coll "agent/internal/collectors"
	"agent/internal/config"
	"agent/internal/models"
	service "agent/internal/services"
	"agent/internal/transport"
	"context"
	"log"
	"sync"
	"time"
)

type App struct {
	cfg            *config.AgentConfig
	collectors     []coll.Collector
	metricsService *service.MetricsService
}

func NewApp(cfg *config.AgentConfig) *App {

	// Инициализация коллекторов
	collectors := []coll.Collector{
		coll.NewSystemCollector(),
		coll.NewProcessCollector(cfg.Processes),
		coll.NewNetworkCollector(),
	}

	// Docker коллектор добавляем, если он доступен
	if dockerCollector, err := coll.NewDockerCollector(cfg.Containers); err == nil {
		collectors = append(collectors, dockerCollector)
	}

	metricsService := service.NewMetricsService()

	return &App{
		cfg:            cfg,
		collectors:     collectors,
		metricsService: metricsService,
	}
}

func (a *App) Run(ctx context.Context, wg *sync.WaitGroup) {
	// Канал для передачи метрик между горутинами
	metricsCh := make(chan models.AgentMetrics, 10)

	// Горутина 1: HTTP-сервер для обработки запросов
	wg.Add(1)
	go func() {
		defer wg.Done()

		server := transport.NewServer(a.cfg.Port, metricsCh, a.metricsService)
		server.Start(ctx)
	}()

	// Горутина 2: Сбор метрик
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(a.cfg.PollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Сбор метрик от всех коллекторов
				metrics := models.NewAgentMetrics(a.cfg.HostID)

				for _, c := range a.collectors {
					if err := c.Collect(&metrics); err != nil {
						log.Printf("Collection error: %v", err)
					}
				}

				// Отправка метрик в сервис для обработки
				a.metricsService.ProcessMetrics(metrics)

				// Отправка метрик в канал для HTTP-сервера
				select {
				case metricsCh <- metrics:
					// Успешно отправлено
				default:
					log.Println("Metrics channel full, skipping")
				}
			}
		}
	}()
}
