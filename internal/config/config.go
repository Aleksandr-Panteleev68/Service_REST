package config

import (
	"errors"
	"fmt"
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
	Port        int        `yaml:"port"`
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
		return nil, fmt.Errorf("failed to read file: %s: %w", path, err)
	}

	if len(file) == 0{
		return nil, errors.New("config file is empty")
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	var errs []error

	if cfg.Env == "" {
		errs = append(errs, errors.New("enviroment cannot be empty"))
	}

	if cfg.Database.Host == "" {
		errs = append(errs, errors.New("database host cannot be empty"))
	}

	if cfg.Database.Port <= 0 {
		errs = append(errs, errors.New("database port cannot be < 0"))
	}

	if cfg.Database.User == "" {
		errs = append(errs, errors.New("database user cannot be empty"))
	}

	if cfg.Database.Password == ""{
		errs = append(errs, errors.New("database password cannot be empty"))
	}

	if cfg.Database.DBName == ""{
		errs = append(errs, errors.New("database name cannot be empty"))
	}

	if cfg.Database.MaxConns <= 0 {
		errs = append(errs, errors.New("database max connections must be > 0"))
	}

	if cfg.HTTPServer.Address == "" {
		errs = append(errs, errors.New("http server address cannot be empty"))
	}

	if cfg.HTTPServer.Port <= 0 {
		errs = append(errs, errors.New("http server port cannot be empty"))
	}

	if cfg.HTTPServer.Timeout <= 0 {
		errs = append(errs, errors.New("http server timeout must be > 0"))
	}

	if cfg.HTTPServer.IdleTimeout <= 0 {
		errs = append(errs, errors.New("http server idle timeout must be > 0"))
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("validation errors: %v", errs)
	}

	return &cfg, nil
}
