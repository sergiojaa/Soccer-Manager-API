package http

import (
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, err := h.signupService.Execute(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId": userID,
	})
}
