package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	authApp "github.com/sergiojaa/soccer-manager-api/internal/auth/application"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/config"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/database"
	usersApp "github.com/sergiojaa/soccer-manager-api/internal/users/application"
	usersHttp "github.com/sergiojaa/soccer-manager-api/internal/users/http"
	usersInfra "github.com/sergiojaa/soccer-manager-api/internal/users/infrastructure"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	r := gin.Default()

	r.GET("/health", healthHandler(db))

	signupService := usersApp.NewSignupService(db)

	expiresIn, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		log.Fatalf("failed to parse JWT_EXPIRES_IN: %v", err)
	}

	userRepo := usersInfra.NewUserRepository(db)
	tokenService := authApp.NewTokenService(cfg.JWTSecret, expiresIn)
	loginService := usersApp.NewLoginService(userRepo, tokenService)

	userHandler := usersHttp.NewHandler(signupService, loginService)

	r.POST("/auth/signup", userHandler.Signup)
	r.POST("/auth/login", userHandler.Login)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func healthHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "degraded",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	}
}
