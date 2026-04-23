package application

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret    []byte
	expiresIn time.Duration
}

func NewTokenService(secret string, expiresIn time.Duration) *TokenService {
	return &TokenService{
		secret:    []byte(secret),
		expiresIn: expiresIn,
	}
}

type AccessTokenClaims struct {
	UserID int64  `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (s *TokenService) Generate(userID int64, email string) (string, error) {
	now := time.Now()

	claims := AccessTokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}
