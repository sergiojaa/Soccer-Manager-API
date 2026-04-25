package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sergiojaa/soccer-manager-api/internal/shared/httpx"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/middleware"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/validation"
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
	PlayerID    int64 `json:"playerId" binding:"required,gt=0"`
	AskingPrice int64 `json:"askingPrice" binding:"required,gt=0"`
}

func (h *Handler) ListPlayer(c *gin.Context) {
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

	var req listPlayerRequest

	if err := validation.BindJSON(c, &req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", h.localizer.Msg(locale, "error.invalid_request_body"))
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
			httpx.Error(c, http.StatusBadRequest, "INVALID_ASKING_PRICE", h.localizer.Msg(locale, "error.invalid_asking_price"))
			return
		case errors.Is(err, application.ErrTransferPlayerNotFoundOrNotOwned):
			httpx.Error(c, http.StatusNotFound, "PLAYER_NOT_FOUND", h.localizer.Msg(locale, "error.player_not_found"))
			return
		case errors.Is(err, application.ErrPlayerAlreadyListed):
			httpx.Error(c, http.StatusConflict, "PLAYER_ALREADY_LISTED", h.localizer.Msg(locale, "error.player_already_listed"))
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
			return
		}
	}

	httpx.Success(c, http.StatusCreated, gin.H{}, h.localizer.Msg(locale, "success.player_listed_for_transfer"))
}

func (h *Handler) ListMarket(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	listings, err := h.listMarketService.Execute(c.Request.Context())
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
		return
	}

	httpx.Success(c, http.StatusOK, gin.H{"items": listings}, "")
}

func (h *Handler) BuyPlayer(c *gin.Context) {
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

	listingID, err := validation.ParsePositiveInt64Param(c, "id")
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_TRANSFER_LISTING_ID", h.localizer.Msg(locale, "error.invalid_transfer_listing_id"))
		return
	}

	err = h.buyPlayerService.Execute(c.Request.Context(), userID, listingID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrListingNotFound):
			httpx.Error(c, http.StatusNotFound, "TRANSFER_LISTING_NOT_FOUND", h.localizer.Msg(locale, "error.transfer_listing_not_found"))
			return
		case errors.Is(err, application.ErrCannotBuyOwnPlayer):
			httpx.Error(c, http.StatusBadRequest, "CANNOT_BUY_OWN_PLAYER", h.localizer.Msg(locale, "error.cannot_buy_own_player"))
			return
		case errors.Is(err, application.ErrInsufficientBudget):
			httpx.Error(c, http.StatusBadRequest, "INSUFFICIENT_BUDGET", h.localizer.Msg(locale, "error.insufficient_budget"))
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
			return
		}
	}

	httpx.Success(c, http.StatusOK, gin.H{}, h.localizer.Msg(locale, "success.player_purchased"))
}
