package main

import (
	"fmt"
	"log"
	"os"

	"pet-project/internal/config"
	"pet-project/internal/storage"
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

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.Close()

	fmt.Println("Succesfully connected to the database!")
}
