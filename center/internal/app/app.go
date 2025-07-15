package app

import (
	"center/internal/transport"
	"center/internal/config"
	"center/internal/detebase/mongo"
	"center/internal/detebase/mongo/repositories"
	"center/internal/detebase/postgres"
	"center/internal/detebase/postgres/repositories"
	"center/internal/models"
	"center/internal/services"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


// App представляет основное приложение центра мониторинга
type App struct {
	cfg             *config.Config
	//logger          *log.Logger
	pgDB            *sql.DB
	mongoDB         *mongo.MongoClient
	hostService     *services.HostService
	pollerService   *services.PollerService
	maintenanceSvc  *services.MaintenanceService
	handler         *api.Handler
	server          *http.Server
	router          *gin.Engine
}



func NewApp(cfg *config.Config) *App 
{
	// Инициализация логгера
	//logger := log.New(os.Stdout, "[MONITORING] ", log.LstdFlags|log.Lshortfile)


	// Инициализация подключения к PostgreSQL
	err := db.InitPostgres(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Валидация структуры БД 
    if err := db.ensurePostgresStructure(); err != nil {
        log.Fatalf("Database structure verification failed: %v", err)
    }
	
	// // Автомиграции
	// if err := runMigrations(pgDB); err != nil {
	// 	log.Fatalf("Migrations failed: %v", err)
	// }

	// Инициализация подключения к MongoDB
	mongoDB, err := db.InitMongo(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	//defer mongoDB.Disconnect()
	defer func() {
		if err := mongoDB.Client().Disconnect(context.Background()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()
	
	// Создание TTL индексов для автоматического удаления старых данных
	if err := createTTLIndexes(mongoDB, cfg.Metrics.MetricsTTLDays); err != nil {
		log.Printf("Failed to create TTL indexes: %v", err)
	}

	// Загрузка начальных данных
	if err := a.loadInitialData(); err != nil {
		log.Printf("Initial data loading error: %v", err)
	}

	// Инициализация репозиториев
	hostRepo := repositories.NewPostgresHostRepository(db.DB)
	processRepo := repositories.NewPostgresProcessRepository(db.DB)
	containerRepo := repositories.NewPostgresContainerRepository(db.DB)
	alertRepo := repositories.NewPostgresAlertRepository(db.DB)
	metricRepo := repositories.NewMongoMetricRepository(mongoDB)

	// Инициализация сервисов
	hostService := services.NewHostService(
		hostRepo,
		processRepo,
		containerRepo,
		alertRepo,
		metricRepo,
	)
	
	pollerService := services.NewPollerService(
		hostService,
		cfg.Metrics.PollInterval,
	)
	
	maintenanceService := services.NewMaintenanceService(
		metricRepo,
		hostRepo,
		cfg.Metrics.MetricsTTLDays,
		cfg.Metrics.SelfCheckInterval,
	)

	// Инициализация обработчиков API
	hostHandler := api.NewHostHandler(hostService)
	processHandler := api.NewProcessHandler(hostService)
	containerHandler := api.NewContainerHandler(hostService)
	alertHandler := api.NewAlertHandler(hostService)
	metricHandler := api.NewMetricHandler(hostService)

	// Создаем общий обработчик
	handler = &api.Handler{
		HostHandler:     hostHandler,
		ProcessHandler:  processHandler,
		ContainerHandler: containerHandler,
		AlertHandler:    alertHandler,
		MetricHandler:   metricHandler,
	}

	// Создание Gin роутера
	router := gin.Default()

	// Настройка маршрутов
	api.SetupRoutes(router, a.handler)

	// Создание HTTP-сервера с таймаутами
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &App{
		cfg:				cfg,
		//logger: 			logger,
		pgDB:				db.DB
		mongoDB:			mongoDB
		hostService:		hostService,
		pollerService:		pollerService,
		maintenanceService:	maintenanceService,
		handler:			handler,
		router:				router,
	}

}

func (a *App) Run(ctx context.Context, wg *sync.WaitGroup) {
	// Запуск HTTP сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.log.Printf("Starting server on port %s", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Fatalf("ListenAndServe: %v", err)
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
			a.log.Println("Sending initial configuration to agents...")
			a.sendInitialConfigToAgents(ctx)
		case <-ctx.Done():
			return
		}
	}()
}



// func runMigrations(db *gorm.DB) error {
// 	return db.AutoMigrate(
// 		&models.Host{},
// 		&models.Process{},
// 		&models.Container{},
// 		&models.AlertRule{},
// 	)
// }


// createTTLIndexes создает TTL индексы в MongoDB
func createTTLIndexes(db *mongo.Database, ttlDays int) error {
	collections := []string{
		"system_metrics",
		"process_metrics",
		"container_metrics",
		"network_metrics",
	}
	
	ttlSeconds := int32(ttlDays * 24 * 60 * 60)
	
	for _, collection := range collections {
		model := mongo.IndexModel{
			Keys: bson.M{"timestamp": 1},
			Options: options.Index().SetExpireAfterSeconds(ttlSeconds),
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		_, err := db.Collection(collection).Indexes().CreateOne(ctx, model)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// Close освобождает ресурсы приложения
func (a *App) Close() {
	// Создаем контекст с таймаутом для graceful shutdown сервера
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	if a.server != nil {
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			a.log.Printf("HTTP server shutdown error: %v", err)
		}
	}
	
	// Закрываем соединения с БД
	if a.pgDB != nil {
		a.pgDB.Close()
	}
	
	if a.mongoDB != nil {
		if err := a.mongoDB.Client().Disconnect(context.Background()); err != nil {
			a.log.Printf("MongoDB disconnect error: %v", err)
		}
	}
}


// sendInitialConfigToAgents отправляет начальную конфигурацию на агентов
func (a *App) sendInitialConfigToAgents(ctx context.Context) {
	hosts, err := a.hostService.GetAllHosts(ctx)
	if err != nil {
		a.log.Printf("Failed to get hosts for initial config: %v", err)
		return
	}

	for _, host := range hosts {
		if err := a.hostService.SendConfigurationToAgent(ctx, host); err != nil {
			a.log.Printf("Failed to send config to host %s: %v", host.Hostname, err)
		} else {
			a.log.Printf("Configuration sent to host %s", host.Hostname)
		}
	}
}


func (a *App) loadInitialData() error {
	// Проверка, есть ли уже данные
	var count int
	if err := a.pgDB.QueryRow("SELECT COUNT(*) FROM hosts").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		return nil // Данные уже есть, пропускаем
	}

	ctx := context.Background()

	// Загрузка данных из конфига
	for _, hostCfg := range a.cfg.InitialData.Hosts {
		// Создание хоста через сервис
		hostID, err := a.hostService.CreateHost(ctx, models.HostInput{
			Hostname:  hostCfg.Hostname,
			IPAddress: hostCfg.IPAddress,
			AgentPort: hostCfg.AgentPort,
			Priority:  hostCfg.Priority,
		})
		if err != nil {
			log.Printf("Failed to create host %s: %v", hostCfg.Hostname, err)
			continue
		}

		// Добавление процессов
		for _, process := range hostCfg.Processes {
			if _, err := a.hostService.AddProcess(ctx, hostID, process); err != nil {
				log.Printf("Failed to add process %s to host %s: %v", process, hostCfg.Hostname, err)
			}
		}

		// Добавление контейнеров
		for _, container := range hostCfg.Containers {
			if _, err := a.hostService.AddContainer(ctx, hostID, container); err != nil {
				log.Printf("Failed to add container %s to host %s: %v", container, hostCfg.Hostname, err)
			}
		}

		// Добавление правил оповещений
		for _, alert := range hostCfg.Alerts {
			if _, err := a.hostService.CreateAlertRule(ctx, hostID, models.AlertInput{
				MetricName:     alert.MetricName,
				ThresholdValue: alert.ThresholdValue,
				Condition:      alert.Condition,
				Enabled:        alert.Enabled,
			}); err != nil {
				log.Printf("Failed to add alert for %s to host %s: %v", alert.MetricName, hostCfg.Hostname, err)
			}
		}
	}
	
	log.Println("Initial data loaded from config using services")
	return nil
}