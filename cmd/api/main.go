package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/config"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/database"

	"github.com/sergiojaa/soccer-manager-api/internal/users/application"
	usersHttp "github.com/sergiojaa/soccer-manager-api/internal/users/http"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	r := gin.Default()

	// Health check
	r.GET("/health", healthHandler(db))

	// Signup wiring
	signupService := application.NewSignupService(db)
	userHandler := usersHttp.NewHandler(signupService)

	r.POST("/auth/signup", userHandler.Signup)

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
