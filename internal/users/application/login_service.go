package application

import (
	"context"
	"strings"

	authApp "github.com/sergiojaa/soccer-manager-api/internal/auth/application"
	"github.com/sergiojaa/soccer-manager-api/internal/users/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepo     *infrastructure.UserRepository
	tokenService *authApp.TokenService
}

func NewLoginService(
	userRepo *infrastructure.UserRepository,
	tokenService *authApp.TokenService,
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
