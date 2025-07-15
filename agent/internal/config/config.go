package config

import (
	"os"
	"time"
	
	"gopkg.in/yaml.v3"
)

type AgentConfig struct {
	HostID        string        `yaml:"host_id"`
	ServerAddress string        `yaml:"server_address"`
	PollInterval  time.Duration `yaml:"poll_interval"`
	Port          string        `yaml:"port"`
	Processes     []string      `yaml:"processes"`
	Containers    []string      `yaml:"containers"`
}

func LoadAgentConfig(path string) (*AgentConfig, error) {
	// Читаем файл конфигурации
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Декодируем YAML
	var cfg AgentConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Переопределение из переменных окружения
	if hostID := os.Getenv("HOST_ID"); hostID != "" {
		cfg.HostID = hostID
	}
	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		cfg.ServerAddress = serverAddr
	}
	if pollInterval := os.Getenv("POLL_INTERVAL"); pollInterval != "" {
		if dur, err := time.ParseDuration(pollInterval); err == nil {
			cfg.PollInterval = dur
		}
	}
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}
	
	// Установка значений по умолчанию, если не заданы
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 60 * time.Second
	}
	if cfg.Port == "" {
		cfg.Port = "8081"
	}

	return &cfg, nil
}