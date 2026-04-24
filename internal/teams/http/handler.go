package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/teams/application"
)

type Handler struct {
	getTeamService *application.GetTeamService
}

func NewHandler(getTeamService *application.GetTeamService) *Handler {
	return &Handler{
		getTeamService: getTeamService,
	}
}

func (h *Handler) GetMyTeam(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user context is missing",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user context",
		})
		return
	}

	team, err := h.getTeamService.Execute(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrTeamNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "team not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, team)
}
