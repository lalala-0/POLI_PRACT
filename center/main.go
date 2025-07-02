package main

import (
	"POLI_PRACT/api"
	"POLI_PRACT/config"
	"POLI_PRACT/db"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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
