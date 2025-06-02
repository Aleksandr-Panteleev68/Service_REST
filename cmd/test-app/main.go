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

	db, err := storage.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.DB.Close()

	user := storage.User{
		FirstName: "Иван",
		LastName: "Иванов",
		Age: 25,
		IsMarried: false,
		Password: "securepassword123",
	}
	id, err := db.CreateUser(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("Create user with ID: %d\n", id)

	fetchedUser, err := db.GetUserByID(id)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Printf("Fetched user: %+v\n",fetchedUser)
}
