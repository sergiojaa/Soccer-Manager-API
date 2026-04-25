package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/players/application"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
)

type Handler struct {
	updatePlayerService *application.UpdatePlayerService
	localizer           *i18n.Localizer
}

type updatePlayerRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
}

func NewHandler(updatePlayerService *application.UpdatePlayerService, localizer *i18n.Localizer) *Handler {
	return &Handler{
		updatePlayerService: updatePlayerService,
		localizer:           localizer,
	}
}

func (h *Handler) UpdatePlayer(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": h.localizer.Msg(locale, "error.user_context_missing"),
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": h.localizer.Msg(locale, "error.user_context_invalid"),
		})
		return
	}

	playerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": h.localizer.Msg(locale, "error.invalid_player_id"),
		})
		return
	}

	var req updatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": h.localizer.Msg(locale, "error.invalid_request_body"),
		})
		return
	}

	err = h.updatePlayerService.Execute(
		c.Request.Context(),
		userID,
		playerID,
		req.FirstName,
		req.LastName,
		req.Country,
	)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidPlayerFirstName):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.player_first_name_required")})
			return
		case errors.Is(err, application.ErrInvalidPlayerLastName):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.player_last_name_required")})
			return
		case errors.Is(err, application.ErrInvalidPlayerCountry):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.player_country_required")})
			return
		case errors.Is(err, application.ErrPlayerNotFoundOrNotOwned):
			c.JSON(http.StatusNotFound, gin.H{"error": h.localizer.Msg(locale, "error.player_not_found")})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": h.localizer.Msg(locale, "error.internal_server")})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": h.localizer.Msg(locale, "success.player_updated"),
	})
}
