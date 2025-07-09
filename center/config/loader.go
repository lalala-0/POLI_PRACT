package config

import (
    "os"
    "gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server      ServerConfig      `yaml:"server"`
	Postgres    PostgresConfig    `yaml:"postgres"`
	MongoDB     MongoDBConfig     `yaml:"mongodb"`
	Metrics     MetricsConfig     `yaml:"metrics"`
	Logging     LoggingConfig     `yaml:"logging"`
	InitialData InitialDataConfig `yaml:"initial_data"`
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
    
    return cfg, nil
}