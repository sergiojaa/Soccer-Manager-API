package application

import (
	"context"
	"strings"

	"github.com/sergiojaa/soccer-manager-api/internal/users/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepo     UserAuthRepository
	tokenService AccessTokenGenerator
}

type UserAuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*infrastructure.UserAuthRecord, error)
}

type AccessTokenGenerator interface {
	Generate(userID int64, email string) (string, error)
}

func NewLoginService(
	userRepo UserAuthRepository,
	tokenService AccessTokenGenerator,
) *LoginService {
	return &LoginService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *LoginService) Execute(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if !isValidEmail(email) {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == infrastructure.ErrUserNotFound {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokenService.Generate(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
