package services

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	generated "github.com/ys-1052/crewee/internal/database/generated"
)

// SportsService handles sports-related business logic
type SportsService struct {
	queries *generated.Queries
}

// NewSportsService creates a new SportsService
func NewSportsService(db *sql.DB) *SportsService {
	return &SportsService{
		queries: generated.New(db),
	}
}

// GetAllSports returns all sports
func (s *SportsService) GetAllSports(ctx context.Context) ([]generated.Sport, error) {
	return s.queries.GetAllSports(ctx)
}

// GetActiveSports returns only active sports
func (s *SportsService) GetActiveSports(ctx context.Context) ([]generated.Sport, error) {
	return s.queries.GetActiveSports(ctx)
}

// GetSportByID returns a sport by ID
func (s *SportsService) GetSportByID(ctx context.Context, id string) (generated.Sport, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return generated.Sport{}, err
	}
	return s.queries.GetSportByID(ctx, uuid)
}

// GetSportByCode returns a sport by code
func (s *SportsService) GetSportByCode(ctx context.Context, code string) (generated.Sport, error) {
	return s.queries.GetSportByCode(ctx, code)
}
