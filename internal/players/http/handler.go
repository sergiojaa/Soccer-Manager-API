package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/players/application"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/httpx"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/validation"
)

type Handler struct {
	updatePlayerService *application.UpdatePlayerService
	localizer           *i18n.Localizer
}

type updatePlayerRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Country   string `json:"country" binding:"required"`
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
		httpx.Error(c, http.StatusUnauthorized, "USER_CONTEXT_MISSING", h.localizer.Msg(locale, "error.user_context_missing"))
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		httpx.Error(c, http.StatusUnauthorized, "USER_CONTEXT_INVALID", h.localizer.Msg(locale, "error.user_context_invalid"))
		return
	}

	playerID, err := validation.ParsePositiveInt64Param(c, "id")
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_PLAYER_ID", h.localizer.Msg(locale, "error.invalid_player_id"))
		return
	}

	var req updatePlayerRequest
	if err := validation.BindJSON(c, &req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", h.localizer.Msg(locale, "error.invalid_request_body"))
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
			httpx.Error(c, http.StatusBadRequest, "INVALID_PLAYER_FIRST_NAME", h.localizer.Msg(locale, "error.player_first_name_required"))
			return
		case errors.Is(err, application.ErrInvalidPlayerLastName):
			httpx.Error(c, http.StatusBadRequest, "INVALID_PLAYER_LAST_NAME", h.localizer.Msg(locale, "error.player_last_name_required"))
			return
		case errors.Is(err, application.ErrInvalidPlayerCountry):
			httpx.Error(c, http.StatusBadRequest, "INVALID_PLAYER_COUNTRY", h.localizer.Msg(locale, "error.player_country_required"))
			return
		case errors.Is(err, application.ErrPlayerNotFoundOrNotOwned):
			httpx.Error(c, http.StatusNotFound, "PLAYER_NOT_FOUND", h.localizer.Msg(locale, "error.player_not_found"))
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
			return
		}
	}

	httpx.Success(c, http.StatusOK, gin.H{}, h.localizer.Msg(locale, "success.player_updated"))
}
