package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ys-1052/crewee/internal/services"
)

const (
	activeOnlyParam = "true"
)

// SportsHandler handles sports-related HTTP requests
type SportsHandler struct {
	sportsService *services.SportsService
}

// NewSportsHandler creates a new SportsHandler
func NewSportsHandler(sportsService *services.SportsService) *SportsHandler {
	return &SportsHandler{
		sportsService: sportsService,
	}
}

// GetSports handles GET /api/v1/sports
func (h *SportsHandler) GetSports(c echo.Context) error {
	ctx := c.Request().Context()

	// Check if only active sports are requested
	activeOnly := c.QueryParam("active") == activeOnlyParam

	var sports interface{}
	var err error

	if activeOnly {
		sports, err = h.sportsService.GetActiveSports(ctx)
	} else {
		sports, err = h.sportsService.GetAllSports(ctx)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch sports")
	}

	response := NewSuccessResponse(sports)
	return c.JSON(http.StatusOK, response)
}

// GetSportByID handles GET /api/v1/sports/:id
func (h *SportsHandler) GetSportByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	sport, err := h.sportsService.GetSportByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Sport not found")
	}

	response := NewSuccessResponse(sport)
	return c.JSON(http.StatusOK, response)
}
