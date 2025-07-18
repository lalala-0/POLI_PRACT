package models

import "time"

type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type PostgresConfig struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	Driver          string        `yaml:"driver"`
	MaxOpenConns    uint64        `yaml:"maxOpenConns"`
	MaxIdleConns    uint64        `yaml:"maxIdleConns"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

type MongoDBConfig struct {
	URI                    string        `yaml:"uri"`
	DBName                 string        `yaml:"dbname"`
	ConnectTimeout         time.Duration `yaml:"connectTimeout"`
	MaxPoolSize            uint64        `yaml:"maxPoolSize"`
	MinPoolSize            uint64        `yaml:"minPoolSize"`
	ServerSelectionTimeout time.Duration `yaml:"serverSelectionTimeout"`
}

type MetricsConfig struct {
	PollInterval      time.Duration        `yaml:"poll_interval"`
	MetricsTTLDays    int                  `yaml:"metrics_ttl_days"`
	SelfCheckInterval time.Duration        `yaml:"self_check_interval"`
	System            SystemMetricsConfig  `yaml:"system"`
	Process           FeatureToggle        `yaml:"process"`
	Network           NetworkMetricsConfig `yaml:"network"`
	Container         FeatureToggle        `yaml:"container"`
}

type SystemMetricsConfig struct {
	Enabled      bool `yaml:"enabled"`
	CollectCPU   bool `yaml:"collect_cpu"`
	CollectRAM   bool `yaml:"collect_ram"`
	CollectDisks bool `yaml:"collect_disks"`
}

type NetworkMetricsConfig struct {
	Enabled    bool `yaml:"enabled"`
	MonitorTCP bool `yaml:"monitor_tcp"`
	MonitorUDP bool `yaml:"monitor_udp"`
}

type FeatureToggle struct {
	Enabled bool `yaml:"enabled"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`    // в мегабайтах
	MaxBackups int    `yaml:"max_backups"` // количество файлов
	MaxAge     int    `yaml:"max_age"`     // в днях
	Compress   bool   `yaml:"compress"`
}

type InitialDataConfig struct {
	Hosts []HostConfig `yaml:"hosts"`
}

type HostConfig struct {
	Hostname   string        `yaml:"hostname"`
	IPAddress  string        `yaml:"ip_address"`
	AgentPort  int           `mapstructure:"agent_port"`
	Priority   int           `yaml:"priority"`
	IsMaster   bool          `yaml:"is_master"`
	Processes  []string      `yaml:"processes"`
	Containers []string      `yaml:"containers"`
	Alerts     []AlertConfig `yaml:"alerts"`
}

type AlertConfig struct {
	Telegram                TelegramConfig `yaml:"telegram"`
	Email                   EmailConfig    `yaml:"email"`
	FailureThresholdPercent float64        `yaml:"failure_threshold_percent"`
	IntervalSeconds         int            `yaml:"interval_seconds"`
}

type TelegramConfig struct {
	Token  string `yaml:"token"`
	ChatID string `yaml:"chat_id"`
}

type EmailConfig struct {
	To       []string `yaml:"to"`
	SMTPHost string   `yaml:"smtp_host"`
	SMTPPort int      `yaml:"smtp_port"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}
