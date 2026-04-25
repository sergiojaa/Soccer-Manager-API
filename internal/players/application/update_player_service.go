package application

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/sergiojaa/soccer-manager-api/internal/players/infrastructure"
)

var (
	ErrInvalidPlayerFirstName   = errors.New("player first name is required")
	ErrInvalidPlayerLastName    = errors.New("player last name is required")
	ErrInvalidPlayerCountry     = errors.New("player country is required")
	ErrPlayerNotFoundOrNotOwned = errors.New("player not found or not owned by user")
)

type UpdatePlayerService struct {
	playerRepo OwnedPlayerUpdaterRepository
}

type OwnedPlayerUpdaterRepository interface {
	UpdateOwnedPlayer(
		ctx context.Context,
		userID int64,
		playerID int64,
		firstName string,
		lastName string,
		country string,
	) error
}

func NewUpdatePlayerService(db *sql.DB) *UpdatePlayerService {
	return NewUpdatePlayerServiceWithRepository(infrastructure.NewPlayerRepository(db))
}

func NewUpdatePlayerServiceWithRepository(playerRepo OwnedPlayerUpdaterRepository) *UpdatePlayerService {
	return &UpdatePlayerService{
		playerRepo: playerRepo,
	}
}

func (s *UpdatePlayerService) Execute(
	ctx context.Context,
	userID int64,
	playerID int64,
	firstName string,
	lastName string,
	country string,
) error {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	country = strings.TrimSpace(country)

	if firstName == "" {
		return ErrInvalidPlayerFirstName
	}

	if lastName == "" {
		return ErrInvalidPlayerLastName
	}

	if country == "" {
		return ErrInvalidPlayerCountry
	}

	err := s.playerRepo.UpdateOwnedPlayer(ctx, userID, playerID, firstName, lastName, country)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPlayerNotFoundOrNotOwned
		}

		return err
	}

	return nil
}
