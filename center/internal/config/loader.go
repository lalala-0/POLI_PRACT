package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

//type AppConfig struct {
//	Server      ServerConfig      `yaml:"server"`
//	Postgres    PostgresConfig    `yaml:"postgres"`
//	MongoDB     MongoDBConfig     `yaml:"mongodb"`
//	Metrics     MetricsConfig     `yaml:"metrics"`
//	Logging     LoggingConfig     `yaml:"logging"`
//	Alerts      AlertsConfig      `yaml:"alerts"`
//	InitialData InitialDataConfig `yaml:"initial_data"`
//}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := defaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	//// Установка значений по умолчанию для mongoDB
	//if cfg.MongoDB.ConnectTimeout <= 0 {
	//	cfg.MongoDB.ConnectTimeout = 5 * time.Second
	//}
	//if cfg.MongoDB.MaxPoolSize <= 0 {
	//	cfg.MongoDB.MaxPoolSize = 100
	//}
	//if cfg.MongoDB.ServerSelectionTimeout <= 0 {
	//	cfg.MongoDB.ServerSelectionTimeout = 30 * time.Second
	//}
	//// Таймауты сервера по умолчанию
	//if cfg.Server.ReadTimeout == 0 {
	//	cfg.Server.ReadTimeout = 30 * time.Second
	//}
	//if cfg.Server.WriteTimeout == 0 {
	//	cfg.Server.WriteTimeout = 30 * time.Second
	//}
	return cfg, nil
}

// defaultConfig возвращает конфиг с значениями по умолчанию
func defaultConfig() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Postgres: PostgresConfig{
			Host:            "localhost",
			Port:            "5432",
			User:            "postgres",
			Password:        "",
			DBName:          "monitoring",
			SSLMode:         "disable",
			Driver:          "postgres",
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		MongoDB: MongoDBConfig{
			URI:                    "mongodb://localhost:27017",
			DBName:                 "metrics",
			ConnectTimeout:         5 * time.Second,
			MaxPoolSize:            20,
			MinPoolSize:            5,
			ServerSelectionTimeout: 10 * time.Second,
		},
		Metrics: MetricsConfig{
			PollInterval:      60 * time.Second,
			MetricsTTLDays:    14,
			SelfCheckInterval: 5 * time.Minute,
			System: SystemMetricsConfig{
				Enabled:      true,
				CollectCPU:   true,
				CollectRAM:   true,
				CollectDisks: true,
			},
			Process: ProcessMetricsConfig{
				Enabled: true,
			},
			Network: NetworkMetricsConfig{
				Enabled:    true,
				MonitorTCP: true,
				MonitorUDP: true,
			},
			Container: ContainerMetricsConfig{
				Enabled: true,
			},
		},
		Logging: LoggingConfig{
			Level:      "info",
			FilePath:   "./monitoring.log",
			MaxSize:    100, // MB
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   true,
		},
		Alerts: AlertsConfig{
			Telegram: TelegramConfig{
				ChatIDs: []string{},
			},
			Email: EmailConfig{
				To: []string{},
			},
		},
		InitialData: InitialDataConfig{
			Hosts: []HostConfig{},
		},
	}
}
