package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	authApp "github.com/sergiojaa/soccer-manager-api/internal/auth/application"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/httpx"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
)

const (
	ContextUserIDKey = "userId"
	ContextEmailKey  = "email"
)

func AuthMiddleware(secret string, localizer *i18n.Localizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := localizer.ResolveLocale(c.GetHeader("Accept-Language"))
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpx.Error(c, http.StatusUnauthorized, "AUTHORIZATION_REQUIRED", localizer.Msg(locale, "error.authorization_required"))
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			httpx.Error(c, http.StatusUnauthorized, "AUTHORIZATION_BEARER_REQUIRED", localizer.Msg(locale, "error.authorization_bearer_required"))
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if tokenString == "" {
			httpx.Error(c, http.StatusUnauthorized, "TOKEN_REQUIRED", localizer.Msg(locale, "error.token_required"))
			c.Abort()
			return
		}

		claims := &authApp.AccessTokenClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			httpx.Error(c, http.StatusUnauthorized, "TOKEN_INVALID_OR_EXPIRED", localizer.Msg(locale, "error.token_invalid_or_expired"))
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextEmailKey, claims.Email)

		c.Next()
	}
}
