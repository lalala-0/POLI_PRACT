package models

import "time"

// Config представляет полную конфигурацию центра мониторинга
type Config struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	Postgres PostgresConfig `yaml:"postgres" json:"postgres"`
	MongoDB  MongoDBConfig  `yaml:"mongodb" json:"mongodb"`
	Metrics  MetricsConfig  `yaml:"metrics" json:"metrics"`
	Logging  LoggingConfig  `yaml:"logging" json:"logging"`

	// Инициальные данные для БД
	InitialData InitialDataConfig `yaml:"initial_data" json:"initial_data"`
}

// ServerConfig содержит настройки HTTP-сервера
type ServerConfig struct {
	Port         string        `yaml:"port" json:"port" env:"SERVER_PORT"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
}

// PostgresConfig содержит настройки подключения к PostgreSQL
type PostgresConfig struct {
	Host     string `yaml:"host" json:"host" env:"PG_HOST"`
	Port     string `yaml:"port" json:"port" env:"PG_PORT"`
	User     string `yaml:"user" json:"user" env:"PG_USER"`
	Password string `yaml:"password" json:"password" env:"PG_PASSWORD"`
	DBName   string `yaml:"dbname" json:"dbname" env:"PG_DBNAME"`
	SSLMode  string `yaml:"sslmode" json:"sslmode" env:"PG_SSLMODE"`
}

// MongoDBConfig содержит настройки подключения к MongoDB
type MongoDBConfig struct {
	URI    string `yaml:"uri" json:"uri" env:"MONGO_URI"`
	DBName string `yaml:"dbname" json:"dbname" env:"MONGO_DBNAME"`
}

// MetricsConfig содержит настройки сбора метрик
type MetricsConfig struct {
	PollInterval      time.Duration `yaml:"poll_interval" json:"poll_interval"`
	MetricsTTLDays    int32         `yaml:"metrics_ttl_days" json:"metrics_ttl_days"`
	SelfCheckInterval time.Duration `yaml:"self_check_interval" json:"self_check_interval"`

	// Настройки сбора метрик
	System    SystemMetricsConfig    `yaml:"system" json:"system"`
	Process   ProcessMetricsConfig   `yaml:"process" json:"process"`
	Network   NetworkMetricsConfig   `yaml:"network" json:"network"`
	Container ContainerMetricsConfig `yaml:"container" json:"container"`
}

// SystemMetricsConfig содержит настройки сбора системных метрик
type SystemMetricsConfig struct {
	Enabled      bool `yaml:"enabled" json:"enabled"`
	CollectCPU   bool `yaml:"collect_cpu" json:"collect_cpu"`
	CollectRAM   bool `yaml:"collect_ram" json:"collect_ram"`
	CollectDisks bool `yaml:"collect_disks" json:"collect_disks"`
}

// ProcessMetricsConfig содержит настройки сбора метрик процессов
type ProcessMetricsConfig struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

// NetworkMetricsConfig содержит настройки сбора сетевых метрик
type NetworkMetricsConfig struct {
	Enabled    bool `yaml:"enabled" json:"enabled"`
	MonitorTCP bool `yaml:"monitor_tcp" json:"monitor_tcp"`
	MonitorUDP bool `yaml:"monitor_udp" json:"monitor_udp"`
}

// ContainerMetricsConfig содержит настройки сбора метрик контейнеров
type ContainerMetricsConfig struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

// LoggingConfig содержит настройки логирования
type LoggingConfig struct {
	Level      string `yaml:"level" json:"level" env:"LOG_LEVEL"`
	FilePath   string `yaml:"file_path" json:"file_path"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
	Compress   bool   `yaml:"compress" json:"compress"`
}

// InitialDataConfig содержит начальные данные для БД
type InitialDataConfig struct {
	Hosts []HostConfig `yaml:"hosts" json:"hosts"`
}

// HostConfig представляет конфигурацию хоста для мониторинга
type HostConfig struct {
	Hostname   string            `yaml:"hostname" json:"hostname"`
	IPAddress  string            `yaml:"ip_address" json:"ip_address"`
	Priority   int               `yaml:"priority" json:"priority"`
	IsMaster   bool              `yaml:"is_master" json:"is_master"`
	Processes  []string          `yaml:"processes" json:"processes"`
	Containers []string          `yaml:"containers" json:"containers"`
	Alerts     []AlertRuleConfig `yaml:"alerts" json:"alerts"`
}

// AlertRuleConfig представляет конфигурацию правила оповещения
type AlertRuleConfig struct {
	MetricName     string  `yaml:"metric_name" json:"metric_name"`
	ThresholdValue float64 `yaml:"threshold_value" json:"threshold_value"`
	Condition      string  `yaml:"condition" json:"condition"`
	Enabled        bool    `yaml:"enabled" json:"enabled"`
}
