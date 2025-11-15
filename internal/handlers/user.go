// Package handlers
package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/dtos"
	"github.com/nyashahama/music-awards/internal/middleware"
	"github.com/nyashahama/music-awards/internal/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.POST("/validate-reset-token", h.ValidateResetToken)
	}

	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", h.ListAllUsers)
		users.GET("/:id", h.GetProfile)
		users.PUT("/:id", h.UpdateProfile)
		users.DELETE("/:id", h.DeleteAccount)
		users.POST("/:id/promote", h.PromoteUser)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.FirstName, req.LastName, req.Email, req.Password, req.Location)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dtos.NewUserResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create context with additional info for email notifications
	ctx := context.WithValue(c.Request.Context(), "User-Agent", c.Request.UserAgent())
	ctx = context.WithValue(ctx, "IP-Address", c.ClientIP())

	token, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.LoginResponse{Token: token})
}

func (h *UserHandler) ListAllUsers(c *gin.Context) {
	currentUserRole := c.MustGet("user_role").(string)
	if currentUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		response[i] = dtos.NewUserResponse(&user)
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	currentUserID := c.MustGet("user_id").(uuid.UUID)
	currentUserRole := c.MustGet("user_role").(string)

	// Add authorization check
	if currentUserRole != "admin" && currentUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	user, err := h.userService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewUserResponse(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	currentUserID := c.MustGet("user_id").(uuid.UUID)
	currentUserRole := c.MustGet("user_role").(string)

	// Authorization check before calling service
	if currentUserRole != "admin" && currentUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req dtos.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := make(map[string]any)
	if req.FirstName != nil {
		updateData["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updateData["last_name"] = *req.LastName
	}
	if req.Email != nil {
		updateData["email"] = *req.Email
	}
	if req.Password != nil {
		updateData["password"] = *req.Password
	}
	if req.Location != nil {
		updateData["location"] = *req.Location
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, updateData)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dtos.NewUserResponse(user))
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	currentUserID := c.MustGet("user_id").(uuid.UUID)
	currentUserRole := c.MustGet("user_role").(string)

	if currentUserRole != "admin" && currentUserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) PromoteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	currentUserRole := c.MustGet("user_role").(string)
	if currentUserRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.userService.PromoteToAdmin(c.Request.Context(), userID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin"})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req dtos.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.userService.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
	})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dtos.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *UserHandler) ValidateResetToken(c *gin.Context) {
	var req dtos.ValidateResetTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.ValidateResetToken(c.Request.Context(), req.Token)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"email": user.Email,
	})
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, services.ErrInvalidToken):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired reset token"})
	case errors.Is(err, services.ErrTokenGeneration):
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate reset token"})
	// ... existing cases
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
