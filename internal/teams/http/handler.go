package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/teams/application"
)

type Handler struct {
	getTeamService    *application.GetTeamService
	updateTeamService *application.UpdateTeamService
	localizer         *i18n.Localizer
}

type updateTeamRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
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

	team, err := h.getTeamService.Execute(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrTeamNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": h.localizer.Msg(locale, "error.team_not_found"),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": h.localizer.Msg(locale, "error.internal_server"),
			})
			return
		}
	}

	c.JSON(http.StatusOK, team)
}

func (h *Handler) UpdateMyTeam(c *gin.Context) {
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

	var req updateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": h.localizer.Msg(locale, "error.invalid_request_body"),
		})
		return
	}

	err := h.updateTeamService.Execute(c.Request.Context(), userID, req.Name, req.Country)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidTeamName):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.team_name_required")})
			return
		case errors.Is(err, application.ErrInvalidTeamCountry):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.team_country_required")})
			return
		case errors.Is(err, application.ErrTeamNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": h.localizer.Msg(locale, "error.team_not_found")})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": h.localizer.Msg(locale, "error.internal_server")})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": h.localizer.Msg(locale, "success.team_updated"),
	})
}
