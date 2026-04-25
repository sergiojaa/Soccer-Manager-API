package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	authApp "github.com/sergiojaa/soccer-manager-api/internal/auth/application"
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": localizer.Msg(locale, "error.authorization_required"),
			})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": localizer.Msg(locale, "error.authorization_bearer_required"),
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": localizer.Msg(locale, "error.token_required"),
			})
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": localizer.Msg(locale, "error.token_invalid_or_expired"),
			})
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextEmailKey, claims.Email)

		c.Next()
	}
}
