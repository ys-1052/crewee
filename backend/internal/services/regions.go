package services

import (
	"context"
	"database/sql"

	generated "github.com/ys-1052/crewee/internal/database/generated"
)

// RegionsService handles regions-related business logic
type RegionsService struct {
	queries *generated.Queries
}

// NewRegionsService creates a new RegionsService
func NewRegionsService(db *sql.DB) *RegionsService {
	return &RegionsService{
		queries: generated.New(db),
	}
}

// GetAllRegions returns all regions
func (s *RegionsService) GetAllRegions(ctx context.Context) ([]generated.Region, error) {
	return s.queries.GetAllRegions(ctx)
}

// GetPrefectures returns all prefectures
func (s *RegionsService) GetPrefectures(ctx context.Context) ([]generated.Region, error) {
	return s.queries.GetPrefectures(ctx)
}

// GetMunicipalitiesByPrefecture returns municipalities for a specific prefecture
func (s *RegionsService) GetMunicipalitiesByPrefecture(
	ctx context.Context, prefectureCode string,
) ([]generated.Region, error) {
	return s.queries.GetMunicipalitiesByPrefecture(ctx, sql.NullString{String: prefectureCode, Valid: true})
}

// GetRegionByCode returns a region by JIS code
func (s *RegionsService) GetRegionByCode(ctx context.Context, code string) (generated.Region, error) {
	return s.queries.GetRegionByCode(ctx, code)
}

// GetRegionHierarchy returns regions with parent information
func (s *RegionsService) GetRegionHierarchy(ctx context.Context) ([]generated.GetRegionHierarchyRow, error) {
	return s.queries.GetRegionHierarchy(ctx)
}
