package config

import "time"

// AppConfig представляет полную конфигурацию центра мониторинга
type AppConfig struct {
	Server   ServerConfig   `yaml:"server" json:"server"`
	Postgres PostgresConfig `yaml:"postgres" json:"postgres"`
	MongoDB  MongoDBConfig  `yaml:"mongodb" json:"mongodb"`
	Metrics  MetricsConfig  `yaml:"metrics" json:"metrics"`
	Logging  LoggingConfig  `yaml:"logging" json:"logging"`
	Alerts   AlertsConfig   `yaml:"alerts" json:"alerts"`

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
	Host            string        `yaml:"host" json:"host" env:"PG_HOST"`
	Port            string        `yaml:"port" json:"port" env:"PG_PORT"`
	User            string        `yaml:"user" json:"user" env:"PG_USER"`
	Password        string        `yaml:"password" json:"password" env:"PG_PASSWORD"`
	DBName          string        `yaml:"dbname" json:"dbname" env:"PG_DBNAME"`
	SSLMode         string        `yaml:"sslmode" json:"sslmode" env:"PG_SSLMODE"`
	Driver          string        `yaml:"driver" json:"driver" env:"PG_DRIVER"`
	MaxOpenConns    uint64        `yaml:"maxOpenConns" json:"maxOpenConns" env:"PG_MAX_OPEN_CONNS"`
	MaxIdleConns    uint64        `yaml:"maxIdleConns" json:"maxIdleConns" env:"PG_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime" json:"connMaxLifetime" env:"PG_CONN_MAX_LIFETIME"`
}

// MongoDBConfig содержит настройки подключения к MongoDB
type MongoDBConfig struct {
	URI                    string        `yaml:"uri" json:"uri" env:"MONGO_URI"`
	DBName                 string        `yaml:"dbname" json:"dbname" env:"MONGO_DBNAME"`
	ConnectTimeout         time.Duration `yaml:"connectTimeout" json:"connectTimeout" env:"MONGO_CONNECT_TIMEOUT"`
	MaxPoolSize            uint64        `yaml:"maxPoolSize" json:"maxPoolSize" env:"MONGO_MAX_POOL_SIZE"`
	MinPoolSize            uint64        `yaml:"minPoolSize" json:"minPoolSize" env:"MONGO_MIN_POOL_SIZE"`
	ServerSelectionTimeout time.Duration `yaml:"serverSelectionTimeout"`
}

// MetricsConfig содержит настройки сбора метрик
type MetricsConfig struct {
	PollInterval      time.Duration `yaml:"poll_interval" json:"poll_interval"`
	MetricsTTLDays    int           `yaml:"metrics_ttl_days" json:"metrics_ttl_days"`
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
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // в мегабайтах
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // количество файлов
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // в днях
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
	AgentPort  int               `yaml:"agent_port" json:"agent_port"`
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

type AlertsConfig struct {
	Telegram                TelegramConfig `yaml:"telegram" json:"telegram"`
	Email                   EmailConfig    `yaml:"email" json:"email"`
	FailureThresholdPercent float64        `yaml:"failure_threshold_percent" json:"failure_threshold_percent"`
	IntervalSeconds         int            `yaml:"interval_seconds" json:"interval_seconds"`
}

type TelegramConfig struct {
	Token   string   `yaml:"token" json:"token"`
	ChatIDs []string `yaml:"chat_ids" json:"chat_ids"`
}

type EmailConfig struct {
	To       []string `yaml:"to" json:"to"`
	SMTPHost string   `yaml:"smtp_host" json:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port" json:"smtp_port"`
	Username string   `yaml:"username" json:"username"`
	Password string   `yaml:"password" json:"password"`
}
