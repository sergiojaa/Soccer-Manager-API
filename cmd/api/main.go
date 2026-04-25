package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	authApp "github.com/sergiojaa/soccer-manager-api/internal/auth/application"
	playersApp "github.com/sergiojaa/soccer-manager-api/internal/players/application"
	playersHttp "github.com/sergiojaa/soccer-manager-api/internal/players/http"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/config"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/database"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/httpx"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	teamsApp "github.com/sergiojaa/soccer-manager-api/internal/teams/application"
	teamsHttp "github.com/sergiojaa/soccer-manager-api/internal/teams/http"
	transfersApp "github.com/sergiojaa/soccer-manager-api/internal/transfers/application"
	transfersHttp "github.com/sergiojaa/soccer-manager-api/internal/transfers/http"
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

	localizer, err := i18n.New(cfg.DefaultLocale)
	if err != nil {
		log.Fatalf("failed to load locales: %v", err)
	}

	r.GET("/health", healthHandler(db, localizer))

	signupService := usersApp.NewSignupService(db)

	expiresIn, err := time.ParseDuration(cfg.JWTExpiresIn)
	if err != nil {
		log.Fatalf("failed to parse JWT_EXPIRES_IN: %v", err)
	}

	userRepo := usersInfra.NewUserRepository(db)
	tokenService := authApp.NewTokenService(cfg.JWTSecret, expiresIn)
	loginService := usersApp.NewLoginService(userRepo, tokenService)

	userHandler := usersHttp.NewHandler(signupService, loginService, localizer)

	r.POST("/auth/signup", userHandler.Signup)
	r.POST("/auth/login", userHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware(cfg.JWTSecret, localizer))

	authorized.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get(middleware.ContextUserIDKey)
		email, _ := c.Get(middleware.ContextEmailKey)

		httpx.Success(c, http.StatusOK, gin.H{
			"userId": userID,
			"email":  email,
		}, "")
	})

	getTeamService := teamsApp.NewGetTeamService(db)
	updateTeamService := teamsApp.NewUpdateTeamService(db)
	updatePlayerService := playersApp.NewUpdatePlayerService(db)

	teamHandler := teamsHttp.NewHandler(getTeamService, updateTeamService, localizer)
	playerHandler := playersHttp.NewHandler(updatePlayerService, localizer)

	listPlayerService := transfersApp.NewListPlayerService(db)
	listMarketService := transfersApp.NewListMarketService(db)
	buyPlayerService := transfersApp.NewBuyPlayerService(db)

	transferHandler := transfersHttp.NewHandler(listPlayerService, listMarketService, buyPlayerService, localizer)

	authorized.GET("/team", teamHandler.GetMyTeam)
	authorized.GET("/transfers", transferHandler.ListMarket)
	authorized.PATCH("/team", teamHandler.UpdateMyTeam)
	authorized.PATCH("/players/:id", playerHandler.UpdatePlayer)
	authorized.POST("/transfers/:id/buy", transferHandler.BuyPlayer)
	authorized.POST("/transfers/list", transferHandler.ListPlayer)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func healthHandler(db *sql.DB, localizer *i18n.Localizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := localizer.ResolveLocale(c.GetHeader("Accept-Language"))
		if err := db.Ping(); err != nil {
			httpx.Error(c, http.StatusServiceUnavailable, "INTERNAL_SERVER_ERROR", localizer.Msg(locale, "error.internal_server"))
			return
		}

		httpx.Success(c, http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
		}, "")
	}
}
