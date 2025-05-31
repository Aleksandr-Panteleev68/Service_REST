package main

import (
	"log"

	"pet-project/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("C:/Users...")
	if err != nil {
		log.Fatalf("Failed to loading config: %v", err)
	}

	log.Printf("Config loaded successfully! Environment: %s", cfg.Env)
}
