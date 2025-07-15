package app

import (
	"center/internal/config"
	mgdb "center/internal/database/mongodb"
	"center/internal/database/mongodb/repositories"
	pgdb "center/internal/database/postgres"
	pg_repo "center/internal/database/postgres/repositories"
	"center/internal/services"
	api "center/internal/transport"
	"sync"

	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// App представляет основное приложение центра мониторинга
type App struct {
	cfg *config.AppConfig
	//logger          *log.Logger
	pgDB           *sql.DB
	mongoDB        *mgdb.MongoDatabase
	hostService    *services.HostService
	pollerService  *services.PollerService
	maintenanceSvc *services.MaintenanceService
	handler        *api.Handler
	server         *http.Server
	router         *gin.Engine
}

func NewApp(cfg *config.AppConfig) *App {
	// Инициализация логгера
	//logger := log.New(os.Stdout, "[MONITORING] ", log.LstdFlags|log.Lshortfile)

	// Инициализация подключения к PostgreSQL
	err := pgdb.InitPostgres(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Валидация структуры БД
	err = pgdb.EnsurePostgresStructure()
	if err != nil {
		log.Fatalf("Database structure verification failed: %v", err)
	}

	// Инициализация подключения к MongoDB
	mongoDB, err := mgdb.InitMongo(cfg.MongoDB, cfg.Metrics.MetricsTTLDays)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	//defer mongoDB.Disconnect()
	defer func() {
		if err := mongoDB.Client.Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()

	// Инициализация репозиториев
	hostRepo := pg_repo.NewPostgresHostRepository(pgdb.DB)
	processRepo := pg_repo.NewPostgresProcessRepository(pgdb.DB)
	containerRepo := pg_repo.NewPostgresContainerRepository(pgdb.DB)
	alertRepo := pg_repo.NewPostgresAlertRepository(pgdb.DB)
	metricRepo := repositories.NewMongoMetricRepository(mongoDB.Database)

	// Инициализация сервисов
	hostService := services.NewHostService(
		*hostRepo,
		*processRepo,
		*containerRepo,
		*alertRepo,
		*metricRepo,
	)

	pollerService := services.NewPollerService(
		hostService,
		cfg.Metrics.PollInterval,
	)

	maintenanceService := services.NewMaintenanceService(
		*metricRepo,
		*hostRepo,
		cfg.Metrics,
	)

	// Инициализация обработчиков API
	hostHandler := api.NewHostHandler(hostService)
	processHandler := api.NewProcessHandler(hostService)
	containerHandler := api.NewContainerHandler(hostService)
	alertHandler := api.NewAlertHandler(hostService)
	metricHandler := api.NewMetricHandler(hostService)

	// Создаем общий обработчик
	handler := &api.Handler{
		HostHandler:      hostHandler,
		ProcessHandler:   processHandler,
		ContainerHandler: containerHandler,
		AlertHandler:     alertHandler,
		MetricHandler:    metricHandler,
	}

	// Создание Gin роутера
	router := gin.Default()

	// Настройка маршрутов
	api.SetupRoutes(router, handler)

	// Создание HTTP-сервера с таймаутами
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &App{
		cfg: cfg,
		//logger: 			logger,
		pgDB:           pgdb.DB,
		mongoDB:        mongoDB,
		hostService:    hostService,
		pollerService:  pollerService,
		maintenanceSvc: maintenanceService,
		server:         server,
		handler:        handler,
		router:         router,
	}

}

func (a *App) Run(ctx context.Context, wg *sync.WaitGroup) {

	// Загрузка начальных данных
	if err := a.hostService.LoadInitialData(ctx, a.cfg); err != nil {
		log.Printf("Initial data loading error: %v", err)
	}

	// Запуск HTTP сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Запуск фоновых задач
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.pollerService.Start(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.maintenanceSvc.StartCleanupRoutine(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.maintenanceSvc.StartSelfCheckRoutine(ctx)
	}()

	// Отправка начальной конфигурации на агентов
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-time.After(5 * time.Second): // Даем серверу время запуститься
			log.Printf("Sending initial configuration to agents...")
			a.sendInitialConfigToAgents(ctx)
		case <-ctx.Done():
			return
		}
	}()
}

// Close освобождает ресурсы приложения
func (a *App) Close() {
	// Создаем контекст с таймаутом для graceful shutdown сервера
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if a.server != nil {
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}

	// Закрываем соединения с БД
	if a.pgDB != nil {
		a.pgDB.Close()
	}

	if a.mongoDB != nil {
		if err := a.mongoDB.Client.Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}
}

// sendInitialConfigToAgents отправляет начальную конфигурацию на агентов
func (a *App) sendInitialConfigToAgents(ctx context.Context) {
	hosts, err := a.hostService.GetAllHosts(ctx)
	if err != nil {
		log.Printf("Failed to get hosts for initial config: %v", err)
		return
	}

	for _, host := range hosts {
		if err := a.hostService.SendConfigurationToAgent(ctx, host); err != nil {
			log.Printf("Failed to send config to host %s: %v", host.Hostname, err)
		} else {
			log.Printf("Configuration sent to host %s", host.Hostname)
		}
	}
}
