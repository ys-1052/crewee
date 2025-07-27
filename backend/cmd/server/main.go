package main

import (
	"fmt"
	"log"

	"github.com/ys-1052/crewee/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Starting Crewee server on port %s in %s mode\n", cfg.Port, cfg.Env)

	// TODO: Initialize database connection
	// TODO: Initialize Echo server
	// TODO: Setup middleware
	// TODO: Setup routes
	// TODO: Start server

	log.Println("Server setup completed")
}
