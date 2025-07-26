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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userID uuid.UUID, updateData map[string]interface{}) (*models.User, error) {
	args := m.Called(ctx, userID, updateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) PromoteToAdmin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
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

func TestUserHandler_UpdateProfile(t *testing.T) {
	userID := uuid.New()

	t.Run("authorized update", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)

		updateReq := dtos.UpdateProfileRequest{
			Username: ptr("newusername"),
		}
		user := &models.User{UserID: userID, Username: "newusername"}

		mockService.On("UpdateUser", mock.Anything, userID, mock.Anything).Return(user, nil)

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", userID)
			c.Set("user_role", "user")
		})
		router.PUT("/users/:id", handler.UpdateProfile)

		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest(http.MethodPut, "/users/"+userID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_DeleteAccount(t *testing.T) {
	t.Run("admin deletion", func(t *testing.T) {
		mockService := new(MockUserService)
		handler := NewUserHandler(mockService)
		userID := uuid.New()

		mockService.On("DeleteUser", mock.Anything, userID).Return(nil)

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uuid.New()) // Different admin user
			c.Set("user_role", "admin")
		})
		router.DELETE("/users/:id", handler.DeleteAccount)

		req, _ := http.NewRequest(http.MethodDelete, "/users/"+userID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// Helper function for pointers
func ptr(s string) *string {
	return &s
}
