package mgdb

import (
	"center/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//type MongoDBConfig struct {
//	URI                    string
//	DBName                 string
//	ConnectTimeout         time.Duration
//	MaxPoolSize            uint64
//	MinPoolSize            uint64
//	ServerSelectionTimeout time.Duration
//}

type MongoDatabase struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func InitMongo(cfg config.MongoDBConfig, tllDays int) (*MongoDatabase, error) {
	// Создаем контекст с таймаутом подключения
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	// Создаем клиент с опциями
	var clientOptions = options.Client().ApplyURI(cfg.URI)

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
	cl, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Проверяем подключение
	//err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB")
	log.Printf("Connection pool: min=%d, max=%d", cfg.MinPoolSize, cfg.MaxPoolSize)

	var db = cl.Database(cfg.DBName)

	// Создание TTL индексов для автоматического удаления старых данных
	if err := createTTLIndexes(db, tllDays); err != nil {
		log.Printf("Failed to create TTL indexes: %v", err)
	}

	return &MongoDatabase{
		Client:   cl,
		Database: db,
	}, nil
}

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
			Keys:    bson.M{"timestamp": 1},
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
