package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/dtos"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService implements UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	args := m.Called(ctx, username, email, password)
	return args.Get(0).(*models.User), args.Error(1)
}

// Implement all other UserService methods similarly...

func TestUserHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful registration", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		user := &models.User{
			UserID:   uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}

		mockService.On("Register", mock.Anything, "testuser", "test@example.com", "ValidPass123!").
			Return(user, nil)

		router := gin.Default()
		router.POST("/auth/register", handler.Register)

		payload := dtos.RegisterRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "ValidPass123!",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response dtos.UserResponse
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.Equal(t, user.UserID, response.UserID)
		assert.Equal(t, user.Username, response.Username)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid payload", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		router := gin.Default()
		router.POST("/auth/register", handler.Register)

		// Missing email field
		invalidPayload := `{"username": "testuser", "password": "password"}`

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(invalidPayload))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestUserHandler_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("Login", mock.Anything, "test@example.com", "password").
			Return("valid_token", nil)

		router := gin.Default()
		router.POST("/auth/login", handler.Login)

		payload := dtos.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response dtos.LoginResponse
		json.Unmarshal(resp.Body.Bytes(), &response)

		assert.Equal(t, "valid_token", response.Token)
		mockService.AssertExpectations(t)
	})
}

// Add tests for other endpoints: GetProfile, UpdateProfile, DeleteAccount, etc.
