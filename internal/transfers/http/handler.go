package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/transfers/application"
)

type Handler struct {
	listPlayerService *application.ListPlayerService
	listMarketService *application.ListMarketService
	buyPlayerService  *application.BuyPlayerService
}

func NewHandler(
	listPlayerService *application.ListPlayerService,
	listMarketService *application.ListMarketService,
	buyPlayerService *application.BuyPlayerService,
) *Handler {
	return &Handler{
		listPlayerService: listPlayerService,
		listMarketService: listMarketService,
		buyPlayerService:  buyPlayerService,
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

func (h *Handler) ListMarket(c *gin.Context) {
	listings, err := h.listMarketService.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": listings,
	})
}

func (h *Handler) BuyPlayer(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context is missing"})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return
	}

	listingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || listingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transfer listing id"})
		return
	}

	err = h.buyPlayerService.Execute(c.Request.Context(), userID, listingID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrListingNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "transfer listing not found"})
			return
		case errors.Is(err, application.ErrCannotBuyOwnPlayer):
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot buy your own player"})
			return
		case errors.Is(err, application.ErrInsufficientBudget):
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient budget"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "player purchased successfully",
	})
}
