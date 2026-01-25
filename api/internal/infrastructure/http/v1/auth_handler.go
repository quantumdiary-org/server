package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"netschool-proxy/api/api/internal/domain/auth"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(authService *auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	SchoolID    int    `json:"school_id" binding:"required"`
	InstanceURL string `json:"instance_url" binding:"required"`
	APIType     string `json:"api_type" binding:"required"` 
}

type LoginResponse struct {
	Token string `json:"token"`
}













func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.LoginWithAPIType(c.Request.Context(), req.Username, req.Password, req.SchoolID, req.InstanceURL, req.APIType)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: token})
}










func (h *AuthHandler) Logout(c *gin.Context) {
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}