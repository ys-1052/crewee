package handlers

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ys-1052/crewee/internal/middleware"
	"github.com/ys-1052/crewee/internal/services"
)

// Router sets up all the routes
func Router(e *echo.Echo, db *sql.DB) {
	// Apply global middleware
	e.Use(middleware.LoggerConfig())
	e.Use(middleware.CORSConfig())

	// Custom error handler
	e.HTTPErrorHandler = middleware.ErrorHandler

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// API v1 routes
	api := e.Group("/api/v1")

	// Initialize services
	sportsService := services.NewSportsService(db)
	regionsService := services.NewRegionsService(db)

	// Initialize handlers
	sportsHandler := NewSportsHandler(sportsService)
	regionsHandler := NewRegionsHandler(regionsService)

	// Sports routes
	api.GET("/sports", sportsHandler.GetSports)
	api.GET("/sports/:id", sportsHandler.GetSportByID)

	// Regions routes
	api.GET("/regions", regionsHandler.GetRegions)
	api.GET("/regions/:code", regionsHandler.GetRegionByCode)
}
