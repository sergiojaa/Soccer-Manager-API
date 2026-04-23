package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sergiojaa/soccer-manager-api/internal/users/application"
)

type Handler struct {
	signupService *application.SignupService
	loginService  *application.LoginService
}

func NewHandler(
	signupService *application.SignupService,
	loginService *application.LoginService,
) *Handler {
	return &Handler{
		signupService: signupService,
		loginService:  loginService,
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
	var req signupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userID, err := h.signupService.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "email must be a valid email address"})
			return
		case errors.Is(err, application.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters long"})
			return
		case errors.Is(err, application.ErrEmailAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": "an account with this email already exists"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId": userID,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	token, err := h.loginService.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token,
	})
}
