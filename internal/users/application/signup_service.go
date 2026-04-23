package application

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/sergiojaa/soccer-manager-api/internal/users/infrastructure"
)

type SignupService struct {
	db       *sql.DB
	userRepo *infrastructure.UserRepository
}

func NewSignupService(db *sql.DB) *SignupService {
	return &SignupService{
		db:       db,
		userRepo: infrastructure.NewUserRepository(db),
	}
}

func (s *SignupService) Execute(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	userID, err := s.userRepo.Create(ctx, tx, email, string(hash))
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}
