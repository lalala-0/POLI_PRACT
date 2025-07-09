package app

import (
	"POLI_PRACT/center/api"
	"POLI_PRACT/center/config"
	"POLI_PRACT/center/internal/database"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func run() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("../../config/config.yml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
	// Валидация конфигурации БД
    if err := cfg.Postgres.Validate(); err != nil {
        log.Fatalf("Invalid DB config: %v", err)
    }

    postgresConfig := db.PostgresConfig{
        Driver:   cfg.Postgres.Driver,
        Host:     cfg.Postgres.Host,
        Port:     cfg.Postgres.Port,
        User:     cfg.Postgres.User,
        Password: cfg.Postgres.Password,
        DBName:   cfg.Postgres.DBName,
        SSLMode:  cfg.Postgres.SSLMode,
    }
	// Инициализация подключения к PostgreSQL
	err = db.InitPostgres(cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Инициализация подключения к MongoDB
	err = db.InitMongo(cfg.Mongo.URI, cfg.Mongo.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Создание Gin роутера
	r := gin.Default()

	// Настройка маршрутов
	api.SetupRoutes(r)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
