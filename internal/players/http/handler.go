package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/players/application"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
)

type Handler struct {
	updatePlayerService *application.UpdatePlayerService
}

type updatePlayerRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Country   string `json:"country"`
}

func NewHandler(updatePlayerService *application.UpdatePlayerService) *Handler {
	return &Handler{
		updatePlayerService: updatePlayerService,
	}
}

func (h *Handler) UpdatePlayer(c *gin.Context) {
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

	playerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid player id",
		})
		return
	}

	var req updatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "player first name is required"})
			return
		case errors.Is(err, application.ErrInvalidPlayerLastName):
			c.JSON(http.StatusBadRequest, gin.H{"error": "player last name is required"})
			return
		case errors.Is(err, application.ErrInvalidPlayerCountry):
			c.JSON(http.StatusBadRequest, gin.H{"error": "player country is required"})
			return
		case errors.Is(err, application.ErrPlayerNotFoundOrNotOwned):
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "player updated successfully",
	})
}
