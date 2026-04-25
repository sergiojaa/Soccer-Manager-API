package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sergiojaa/soccer-manager-api/internal/shared/i18n"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Signup(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	var req signupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": h.localizer.Msg(locale, "error.invalid_request_body"),
		})
		return
	}

	userID, err := h.signupService.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.invalid_email")})
			return
		case errors.Is(err, application.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": h.localizer.Msg(locale, "error.invalid_password")})
			return
		case errors.Is(err, application.ErrEmailAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": h.localizer.Msg(locale, "error.email_already_exists")})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": h.localizer.Msg(locale, "error.internal_server"),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId": userID,
	})
}

func (h *Handler) Login(c *gin.Context) {
	locale := h.localizer.ResolveLocale(c.GetHeader("Accept-Language"))
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": h.localizer.Msg(locale, "error.invalid_request_body"),
		})
		return
	}

	token, err := h.loginService.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": h.localizer.Msg(locale, "error.invalid_credentials")})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": h.localizer.Msg(locale, "error.internal_server")})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
	})
}
