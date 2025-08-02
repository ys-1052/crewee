package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/ys-1052/crewee/internal/config"
	"github.com/ys-1052/crewee/internal/database"
	"github.com/ys-1052/crewee/internal/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Starting Crewee server on port %s in %s mode\n", cfg.Port, cfg.Env)

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Database connection established")

	// Initialize Echo server
	e := echo.New()

	// Setup routes and middleware
	handlers.Router(e, db)

	// Start server
	fmt.Printf("Server listening on :%s\n", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Printf("Failed to start server: %v", err)
		return
	}
}
