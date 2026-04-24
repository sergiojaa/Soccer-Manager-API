package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/transfers/application"
)

type Handler struct {
	listPlayerService *application.ListPlayerService
}

func NewHandler(listPlayerService *application.ListPlayerService) *Handler {
	return &Handler{
		listPlayerService: listPlayerService,
	}
}

type listPlayerRequest struct {
	PlayerID    int64 `json:"playerId"`
	AskingPrice int64 `json:"askingPrice"`
}

func (h *Handler) ListPlayer(c *gin.Context) {
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

	var req listPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.PlayerID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid player id",
		})
		return
	}

	err := h.listPlayerService.Execute(
		c.Request.Context(),
		userID,
		req.PlayerID,
		req.AskingPrice,
	)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidAskingPrice):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "asking price must be greater than zero",
			})
			return
		case errors.Is(err, application.ErrTransferPlayerNotFoundOrNotOwned):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "player not found",
			})
			return
		case errors.Is(err, application.ErrPlayerAlreadyListed):
			c.JSON(http.StatusConflict, gin.H{
				"error": "player is already listed for transfer",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "player listed for transfer successfully",
	})
}
