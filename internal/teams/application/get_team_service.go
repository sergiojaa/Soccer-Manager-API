package application

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sergiojaa/soccer-manager-api/internal/teams/infrastructure"
)

var ErrTeamNotFound = errors.New("team not found")

type GetTeamService struct {
	teamRepo TeamReaderRepository
}

type TeamReaderRepository interface {
	FindByUserID(ctx context.Context, userID int64) (*infrastructure.TeamView, error)
}

func NewGetTeamService(db *sql.DB) *GetTeamService {
	return NewGetTeamServiceWithRepository(infrastructure.NewTeamRepository(db))
}

func NewGetTeamServiceWithRepository(teamRepo TeamReaderRepository) *GetTeamService {
	return &GetTeamService{
		teamRepo: teamRepo,
	}
}

func (s *GetTeamService) Execute(
	ctx context.Context,
	userID int64,
) (*infrastructure.TeamView, error) {
	team, err := s.teamRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamNotFound
		}

		return nil, err
	}

	return team, nil
}
