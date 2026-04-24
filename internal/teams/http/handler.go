package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/teams/application"
)

type Handler struct {
	getTeamService    *application.GetTeamService
	updateTeamService *application.UpdateTeamService
}

type updateTeamRequest struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func NewHandler(
	getTeamService *application.GetTeamService,
	updateTeamService *application.UpdateTeamService,
) *Handler {
	return &Handler{
		getTeamService:    getTeamService,
		updateTeamService: updateTeamService,
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

func (h *Handler) UpdateMyTeam(c *gin.Context) {
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

	var req updateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err := h.updateTeamService.Execute(c.Request.Context(), userID, req.Name, req.Country)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidTeamName):
			c.JSON(http.StatusBadRequest, gin.H{"error": "team name is required"})
			return
		case errors.Is(err, application.ErrInvalidTeamCountry):
			c.JSON(http.StatusBadRequest, gin.H{"error": "team country is required"})
			return
		case errors.Is(err, application.ErrTeamNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "team updated successfully",
	})
}
