package application

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/sergiojaa/soccer-manager-api/internal/teams/infrastructure"
)

var (
	ErrInvalidTeamName    = errors.New("team name is required")
	ErrInvalidTeamCountry = errors.New("team country is required")
)

type UpdateTeamService struct {
	teamRepo TeamUpdaterRepository
}

type TeamUpdaterRepository interface {
	UpdateByUserID(ctx context.Context, userID int64, name string, country string) error
}

func NewUpdateTeamService(db *sql.DB) *UpdateTeamService {
	return NewUpdateTeamServiceWithRepository(infrastructure.NewTeamRepository(db))
}

func NewUpdateTeamServiceWithRepository(teamRepo TeamUpdaterRepository) *UpdateTeamService {
	return &UpdateTeamService{
		teamRepo: teamRepo,
	}
}

func (s *UpdateTeamService) Execute(
	ctx context.Context,
	userID int64,
	name string,
	country string,
) error {
	name = strings.TrimSpace(name)
	country = strings.TrimSpace(country)

	if name == "" {
		return ErrInvalidTeamName
	}

	if country == "" {
		return ErrInvalidTeamCountry
	}

	err := s.teamRepo.UpdateByUserID(ctx, userID, name, country)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTeamNotFound
		}

		return err
	}

	return nil
}
