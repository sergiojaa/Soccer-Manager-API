package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/transfers/application"
)

type Handler struct {
	listPlayerService *application.ListPlayerService
	listMarketService *application.ListMarketService
	buyPlayerService  *application.BuyPlayerService
	localizer         *i18n.Localizer
}

func NewHandler(
	listPlayerService *application.ListPlayerService,
	listMarketService *application.ListMarketService,
	buyPlayerService *application.BuyPlayerService,
	localizer *i18n.Localizer,
) *Handler {
	return &Handler{
		listPlayerService: listPlayerService,
		listMarketService: listMarketService,
		buyPlayerService:  buyPlayerService,
		localizer:         localizer,
	}
}

type listPlayerRequest struct {
	PlayerID    int64 `json:"playerId"`
	AskingPrice int64 `json:"askingPrice"`
}

func (h *Handler) ListPlayer(c *gin.Context) {
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

	var req listPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": h.localizer.Msg(locale, "error.invalid_request_body"),
		})
		return
	}

	if req.PlayerID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": h.localizer.Msg(locale, "error.invalid_player_id"),
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
				"error": h.localizer.Msg(locale, "error.invalid_asking_price"),
			})
			return
		case errors.Is(err, application.ErrTransferPlayerNotFoundOrNotOwned):
			c.JSON(http.StatusNotFound, gin.H{
				"error": h.localizer.Msg(locale, "error.player_not_found"),
			})
			return
		case errors.Is(err, application.ErrPlayerAlreadyListed):
			c.JSON(http.StatusConflict, gin.H{
				"error": h.localizer.Msg(locale, "error.player_already_listed"),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": h.localizer.Msg(locale, "error.internal_server"),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": h.localizer.Msg(locale, "success.player_listed_for_transfer"),
	})
}

func (h *Handler) ListMarket(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	listings, err := h.listMarketService.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": h.localizer.Msg(locale, "error.internal_server"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": listings,
	})
}

func (h *Handler) BuyPlayer(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": h.localizer.Msg(locale, "error.user_context_missing")})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": h.localizer.Msg(locale, "error.user_context_invalid")})
		return
	}

	listingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || listingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.invalid_transfer_listing_id")})
		return
	}

	err = h.buyPlayerService.Execute(c.Request.Context(), userID, listingID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrListingNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": h.localizer.Msg(locale, "error.transfer_listing_not_found")})
			return
		case errors.Is(err, application.ErrCannotBuyOwnPlayer):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.cannot_buy_own_player")})
			return
		case errors.Is(err, application.ErrInsufficientBudget):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.insufficient_budget")})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": h.localizer.Msg(locale, "error.internal_server")})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": h.localizer.Msg(locale, "success.player_purchased"),
	})
}
