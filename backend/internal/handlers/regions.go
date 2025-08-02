package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ys-1052/crewee/internal/services"
)

const (
	trueBoolStr = "true"
)

// RegionsHandler handles regions-related HTTP requests
type RegionsHandler struct {
	regionsService *services.RegionsService
}

// NewRegionsHandler creates a new RegionsHandler
func NewRegionsHandler(regionsService *services.RegionsService) *RegionsHandler {
	return &RegionsHandler{
		regionsService: regionsService,
	}
}

// GetRegions handles GET /api/v1/regions
func (h *RegionsHandler) GetRegions(c echo.Context) error {
	ctx := c.Request().Context()

	// Check query parameters for filtering
	regionType := c.QueryParam("type")
	prefectureCode := c.QueryParam("prefecture")
	hierarchy := c.QueryParam("hierarchy") == trueBoolStr

	var regions interface{}
	var err error

	switch {
	case hierarchy:
		regions, err = h.regionsService.GetRegionHierarchy(ctx)
	case regionType == "prefecture":
		regions, err = h.regionsService.GetPrefectures(ctx)
	case prefectureCode != "":
		regions, err = h.regionsService.GetMunicipalitiesByPrefecture(ctx, prefectureCode)
	default:
		regions, err = h.regionsService.GetAllRegions(ctx)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch regions")
	}

	response := NewSuccessResponse(regions)
	return c.JSON(http.StatusOK, response)
}

// GetRegionByCode handles GET /api/v1/regions/:code
func (h *RegionsHandler) GetRegionByCode(c echo.Context) error {
	ctx := c.Request().Context()
	code := c.Param("code")

	region, err := h.regionsService.GetRegionByCode(ctx, code)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Region not found")
	}

	response := NewSuccessResponse(region)
	return c.JSON(http.StatusOK, response)
}
