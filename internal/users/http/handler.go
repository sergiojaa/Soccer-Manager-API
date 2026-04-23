package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sergiojaa/soccer-manager-api/internal/users/application"
)

type Handler struct {
	signupService *application.SignupService
}

func NewHandler(signupService *application.SignupService) *Handler {
	return &Handler{signupService: signupService}
}

type signupRequest struct {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		case errors.Is(err, application.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
			return
		case errors.Is(err, application.ErrEmailAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
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
