package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/httpx"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/validation"
	"github.com/sergiojaa/soccer-manager-api/internal/teams/application"
)

type Handler struct {
	getTeamService    *application.GetTeamService
	updateTeamService *application.UpdateTeamService
	localizer         *i18n.Localizer
}

type updateTeamRequest struct {
	Name    string `json:"name" binding:"required"`
	Country string `json:"country" binding:"required"`
}

func NewHandler(
	getTeamService *application.GetTeamService,
	updateTeamService *application.UpdateTeamService,
	localizer *i18n.Localizer,
) *Handler {
	return &Handler{
		getTeamService:    getTeamService,
		updateTeamService: updateTeamService,
		localizer:         localizer,
	}
}

func (h *Handler) GetMyTeam(c *gin.Context) {
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

	team, err := h.getTeamService.Execute(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrTeamNotFound):
			httpx.Error(c, http.StatusNotFound, "TEAM_NOT_FOUND", h.localizer.Msg(locale, "error.team_not_found"))
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
			return
		}
	}

	httpx.Success(c, http.StatusOK, team, "")
}

func (h *Handler) UpdateMyTeam(c *gin.Context) {
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

	var req updateTeamRequest

	if err := validation.BindJSON(c, &req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", h.localizer.Msg(locale, "error.invalid_request_body"))
		return
	}

	err := h.updateTeamService.Execute(c.Request.Context(), userID, req.Name, req.Country)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidTeamName):
			httpx.Error(c, http.StatusBadRequest, "INVALID_TEAM_NAME", h.localizer.Msg(locale, "error.team_name_required"))
			return
		case errors.Is(err, application.ErrInvalidTeamCountry):
			httpx.Error(c, http.StatusBadRequest, "INVALID_TEAM_COUNTRY", h.localizer.Msg(locale, "error.team_country_required"))
			return
		case errors.Is(err, application.ErrTeamNotFound):
			httpx.Error(c, http.StatusNotFound, "TEAM_NOT_FOUND", h.localizer.Msg(locale, "error.team_not_found"))
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
			return
		}
	}

	httpx.Success(c, http.StatusOK, gin.H{}, h.localizer.Msg(locale, "success.team_updated"))
}
