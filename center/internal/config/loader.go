package config

import (
	"center/internal/models"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type AppConfig struct {
	Server      models.ServerConfig      `yaml:"server"`
	Postgres    models.PostgresConfig    `yaml:"postgres"`
	MongoDB     models.MongoDBConfig     `yaml:"mongodb"`
	Metrics     models.MetricsConfig     `yaml:"metrics"`
	Logging     models.LoggingConfig     `yaml:"logging"`
	Alerts      models.AlertConfig       `yaml:"alerts"`
	InitialData models.InitialDataConfig `yaml:"initial_data"`
}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &AppConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	// Установка значений по умолчанию для mongoDB
	if cfg.MongoDB.ConnectTimeout <= 0 {
		cfg.MongoDB.ConnectTimeout = 5 * time.Second
	}
	if cfg.MongoDB.MaxPoolSize <= 0 {
		cfg.MongoDB.MaxPoolSize = 100
	}
	if cfg.MongoDB.ServerSelectionTimeout <= 0 {
		cfg.MongoDB.ServerSelectionTimeout = 30 * time.Second
	}
	// Таймауты сервера по умолчанию
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30 * time.Second
	}
	return cfg, nil
}
