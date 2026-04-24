package application

import (
	"context"
	"database/sql"
	"strings"

	playersApplication "github.com/sergiojaa/soccer-manager-api/internal/players/application"
	playersInfra "github.com/sergiojaa/soccer-manager-api/internal/players/infrastructure"
	teamInfrastructure "github.com/sergiojaa/soccer-manager-api/internal/teams/infrastructure"
	"github.com/sergiojaa/soccer-manager-api/internal/users/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

type SignupService struct {
	db         *sql.DB
	userRepo   *infrastructure.UserRepository
	teamRepo   *teamInfrastructure.TeamRepository
	playerRepo *playersInfra.PlayerRepository
}

func NewSignupService(db *sql.DB) *SignupService {
	return &SignupService{
		db:         db,
		userRepo:   infrastructure.NewUserRepository(db),
		teamRepo:   teamInfrastructure.NewTeamRepository(db),
		playerRepo: playersInfra.NewPlayerRepository(db),
	}
}

func (s *SignupService) Execute(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if !isValidEmail(email) {
		return 0, ErrInvalidEmail
	}

	if len(password) < 6 {
		return 0, ErrInvalidPassword
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	userID, err := s.userRepo.Create(ctx, tx, email, string(hash))
	if err != nil {
		if err == infrastructure.ErrDuplicateEmail {
			return 0, ErrEmailAlreadyUsed
		}
		return 0, err
	}

	teamID, err := s.teamRepo.Create(ctx, tx, userID, "My Team", "Unknown")
	if err != nil {
		return 0, err
	}

	players := playersApplication.GeneratePlayers(teamID)

	if err := s.playerRepo.CreateBatch(ctx, tx, players); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && len(email) >= 3
}
