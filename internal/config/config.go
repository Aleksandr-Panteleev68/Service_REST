package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env        string     `yaml:"env"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Database   Database   `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	MaxConns int    `yaml:"max_conns"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("File not found", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed parsing YAML: %w", err)
	}

	if cfg.Database.Port <= 0 {
		return nil, fmt.Errorf("port cannot be empty: %w", err)
	}

	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("cannot be empty: %w", err)
	}

	if cfg.HTTPServer.Address == "" {
		return nil, fmt.Errorf("cannot be empty: %w", err)
	}

	if cfg.Database.User == "" {
		return nil, fmt.Errorf("no User data: %w", err)
	}

	return &cfg, nil
}
