package main

import (
	"log"
	"os"

	"pet-project/internal/config"
)

func main() {
	configPath := os.Getenv("PET_CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
		log.Printf("PET_CONFIG_PATH not set, using default: %s", configPath)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to loading config: %v", err)
	}

	log.Printf("Config loaded successfully! Environment: %s", cfg.Env)
}
