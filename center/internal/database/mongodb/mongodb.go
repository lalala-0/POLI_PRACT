package db
import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type MongoDBConfig struct {
    URI     			  string
    DBName   			  string
	ConnectTimeout        time.Duration 
	MaxPoolSize           uint64        
	MinPoolSize           uint64        
	ServerSelectionTimeout time.Duration 
}

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func InitMongo(cfg MongoDBConfig) (*MongoClient, error)  {
	// Создаем контекст с таймаутом подключения
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	// Создаем клиент с опциями
	clientOptions := options.Client().ApplyURI(cfg.URI)
	
	// Настройки пула подключений
	if cfg.MaxPoolSize > 0 {
		clientOptions.SetMaxPoolSize(cfg.MaxPoolSize)
	}
	if cfg.MinPoolSize > 0 {
		clientOptions.SetMinPoolSize(cfg.MinPoolSize)
	}
	if cfg.ServerSelectionTimeout > 0 {
		clientOptions.SetServerSelectionTimeout(cfg.ServerSelectionTimeout)
	}

	// Устанавливаем подключение
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	log.Printf("Connection pool: min=%d, max=%d", cfg.MinPoolSize, cfg.MaxPoolSize)

	db := client.Database(cfg.DBName)

	return &MongoClient{
		Client:   client,
		Database: db,
	}, nil
}