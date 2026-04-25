package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/httpx"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/validation"
	"github.com/sergiojaa/soccer-manager-api/internal/users/application"
)

type Handler struct {
	signupService *application.SignupService
	loginService  *application.LoginService
	localizer     *i18n.Localizer
}

func NewHandler(
	signupService *application.SignupService,
	loginService *application.LoginService,
	localizer *i18n.Localizer,
) *Handler {
	return &Handler{
		signupService: signupService,
		loginService:  loginService,
		localizer:     localizer,
	}
}

type signupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Signup(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	var req signupRequest

	if err := validation.BindJSON(c, &req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", h.localizer.Msg(locale, "error.invalid_request_body"))
		return
	}

	userID, err := h.signupService.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidEmail):
			httpx.Error(c, http.StatusBadRequest, "INVALID_EMAIL", h.localizer.Msg(locale, "error.invalid_email"))
			return
		case errors.Is(err, application.ErrInvalidPassword):
			httpx.Error(c, http.StatusBadRequest, "INVALID_PASSWORD", h.localizer.Msg(locale, "error.invalid_password"))
			return
		case errors.Is(err, application.ErrEmailAlreadyUsed):
			httpx.Error(c, http.StatusConflict, "EMAIL_ALREADY_EXISTS", h.localizer.Msg(locale, "error.email_already_exists"))
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
			return
		}
	}

	httpx.Success(c, http.StatusCreated, gin.H{"userId": userID}, "")
}

func (h *Handler) Login(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	var req loginRequest

	if err := validation.BindJSON(c, &req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "INVALID_REQUEST_BODY", h.localizer.Msg(locale, "error.invalid_request_body"))
		return
	}

	token, err := h.loginService.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidCredentials):
			httpx.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", h.localizer.Msg(locale, "error.invalid_credentials"))
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", h.localizer.Msg(locale, "error.internal_server"))
			return
		}
	}

	httpx.Success(c, http.StatusOK, gin.H{"accessToken": token}, "")
}
