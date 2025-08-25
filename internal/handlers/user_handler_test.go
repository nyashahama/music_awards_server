package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/dtos"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/nyashahama/music-awards/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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

// Test helper functions
func setupHandlerTest() (*MockUserService, *UserHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	router := gin.Default()
	return mockService, handler, router
}

func setupAuthContext(router *gin.Engine, userID uuid.UUID, role string) *gin.Engine {
	router.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("user_role", role)
	})
	return router
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name: "successful registration",
			payload: dtos.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "ValidPass123!",
			},
			mockSetup: func(m *MockUserService) {
				user := &models.User{
					UserID:   uuid.New(),
					Username: "testuser",
					Email:    "test@example.com",
					Role:     "user",
				}
				m.On("Register", mock.Anything, "testuser", "test@example.com", "ValidPass123!").Return(user, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid payload",
			payload: map[string]interface{}{
				"username": "testuser",
				// Missing email and password
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "email already exists",
			payload: dtos.RegisterRequest{
				Username: "testuser",
				Email:    "existing@example.com",
				Password: "ValidPass123!",
			},
			mockSetup: func(m *MockUserService) {
				m.On("Register", mock.Anything, "testuser", "existing@example.com", "ValidPass123!").Return(nil, services.ErrEmailExists)
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupHandlerTest()
			router.POST("/auth/register", handler.Register)
			tt.mockSetup(mockService)

			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name: "successful login",
			payload: dtos.LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(m *MockUserService) {
				m.On("Login", mock.Anything, "test@example.com", "password").Return("valid_token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			payload: dtos.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockUserService) {
				m.On("Login", mock.Anything, "test@example.com", "wrongpassword").Return("", services.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid payload",
			payload: map[string]interface{}{
				"email": "test@example.com",
				// Missing password
			},
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupHandlerTest()
			router.POST("/auth/login", handler.Login)
			tt.mockSetup(mockService)

			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_ListAllUsers(t *testing.T) {
	tests := []struct {
		name           string
		authUserRole   string
		mockSetup      func(*MockUserService)
		expectedStatus int
	}{
		{
			name:         "admin can list all users",
			authUserRole: "admin",
			mockSetup: func(m *MockUserService) {
				users := []models.User{*createTestUser(), *createTestUser()}
				m.On("GetAllUsers", mock.Anything).Return(users, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-admin cannot list all users",
			authUserRole:   "user",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "repository error",
			authUserRole: "admin",
			mockSetup: func(m *MockUserService) {
				m.On("GetAllUsers", mock.Anything).Return([]models.User{}, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupHandlerTest()
			setupAuthContext(router, uuid.New(), tt.authUserRole)
			router.GET("/users", handler.ListAllUsers)
			tt.mockSetup(mockService)

			req, _ := http.NewRequest(http.MethodGet, "/users", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name           string
		pathUserID     string
		authUserID     uuid.UUID
		authUserRole   string
		mockSetup      func(*MockUserService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:         "successful get own profile",
			pathUserID:   userID.String(),
			authUserID:   userID,
			authUserRole: "user",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				user := &models.User{UserID: targetID, Username: "testuser"}
				m.On("GetUserProfile", mock.Anything, targetID).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "admin can get any profile",
			pathUserID:   userID.String(),
			authUserID:   otherUserID, // Different user
			authUserRole: "admin",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				user := &models.User{UserID: targetID, Username: "testuser"}
				m.On("GetUserProfile", mock.Anything, targetID).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "unauthorized access to other profile",
			pathUserID:     userID.String(),
			authUserID:     otherUserID, // Different user
			authUserRole:   "user",      // Not admin
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid user ID",
			pathUserID:     "invalid-uuid",
			authUserID:     userID,
			authUserRole:   "user",
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "user not found",
			pathUserID:   userID.String(),
			authUserID:   userID,
			authUserRole: "user",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				m.On("GetUserProfile", mock.Anything, targetID).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupHandlerTest()
			setupAuthContext(router, tt.authUserID, tt.authUserRole)
			router.GET("/users/:id", handler.GetProfile)
			tt.mockSetup(mockService, userID) // userID is the target ID in the path

			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.pathUserID, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name           string
		pathUserID     string
		authUserID     uuid.UUID
		authUserRole   string
		payload        interface{}
		mockSetup      func(*MockUserService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:         "successful update own profile",
			pathUserID:   userID.String(),
			authUserID:   userID,
			authUserRole: "user",
			payload: dtos.UpdateProfileRequest{
				Username: ptr("newusername"),
			},
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				user := &models.User{UserID: targetID, Username: "newusername"}
				m.On("UpdateUser", mock.Anything, targetID, mock.Anything).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "admin can update any profile",
			pathUserID:   userID.String(),
			authUserID:   otherUserID, // Different user
			authUserRole: "admin",
			payload: dtos.UpdateProfileRequest{
				Username: ptr("newusername"),
			},
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				user := &models.User{UserID: targetID, Username: "newusername"}
				m.On("UpdateUser", mock.Anything, targetID, mock.Anything).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "unauthorized update of other profile",
			pathUserID:   userID.String(),
			authUserID:   otherUserID, // Different user
			authUserRole: "user",      // Not admin
			payload: dtos.UpdateProfileRequest{
				Username: ptr("newusername"),
			},
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "invalid user ID",
			pathUserID:   "invalid-uuid",
			authUserID:   userID,
			authUserRole: "user",
			payload: dtos.UpdateProfileRequest{
				Username: ptr("newusername"),
			},
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "email conflict",
			pathUserID:   userID.String(),
			authUserID:   userID,
			authUserRole: "user",
			payload: dtos.UpdateProfileRequest{
				Email: ptr("existing@example.com"),
			},
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				m.On("UpdateUser", mock.Anything, targetID, mock.Anything).Return(nil, services.ErrEmailExists)
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupHandlerTest()
			// Set up auth context first
			router.Use(func(c *gin.Context) {
				c.Set("user_id", tt.authUserID)
				c.Set("user_role", tt.authUserRole)
			})
			router.PUT("/users/:id", handler.UpdateProfile)
			tt.mockSetup(mockService, userID) // userID is the target ID in the path

			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPut, "/users/"+tt.pathUserID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			// Only assert expectations if we expect the service to be called
			if tt.expectedStatus == http.StatusOK || tt.expectedStatus == http.StatusConflict {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUserHandler_DeleteAccount(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name           string
		pathUserID     string
		authUserID     uuid.UUID
		authUserRole   string
		mockSetup      func(*MockUserService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:         "successful delete own account",
			pathUserID:   userID.String(),
			authUserID:   userID,
			authUserRole: "user",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				m.On("DeleteUser", mock.Anything, targetID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:         "admin can delete any account",
			pathUserID:   userID.String(),
			authUserID:   otherUserID, // Different user
			authUserRole: "admin",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				m.On("DeleteUser", mock.Anything, targetID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "unauthorized delete of other account",
			pathUserID:     userID.String(),
			authUserID:     otherUserID, // Different user
			authUserRole:   "user",      // Not admin
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid user ID",
			pathUserID:     "invalid-uuid",
			authUserID:     userID,
			authUserRole:   "user",
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "user not found",
			pathUserID:   userID.String(),
			authUserID:   userID,
			authUserRole: "user",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				m.On("DeleteUser", mock.Anything, targetID).Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupHandlerTest()
			setupAuthContext(router, tt.authUserID, tt.authUserRole)
			router.DELETE("/users/:id", handler.DeleteAccount)
			tt.mockSetup(mockService, userID) // userID is the target ID in the path

			req, _ := http.NewRequest(http.MethodDelete, "/users/"+tt.pathUserID, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_PromoteUser(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name           string
		pathUserID     string
		authUserID     uuid.UUID
		authUserRole   string
		mockSetup      func(*MockUserService, uuid.UUID)
		expectedStatus int
	}{
		{
			name:         "admin can promote user",
			pathUserID:   userID.String(),
			authUserID:   otherUserID, // Different user (admin)
			authUserRole: "admin",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				m.On("PromoteToAdmin", mock.Anything, targetID).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-admin cannot promote",
			pathUserID:     userID.String(),
			authUserID:     otherUserID, // Different user
			authUserRole:   "user",      // Not admin
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "invalid user ID",
			pathUserID:     "invalid-uuid",
			authUserID:     otherUserID,
			authUserRole:   "admin",
			mockSetup:      func(m *MockUserService, targetID uuid.UUID) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "user not found",
			pathUserID:   userID.String(),
			authUserID:   otherUserID,
			authUserRole: "admin",
			mockSetup: func(m *MockUserService, targetID uuid.UUID) {
				m.On("PromoteToAdmin", mock.Anything, targetID).Return(gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler, router := setupHandlerTest()
			setupAuthContext(router, tt.authUserID, tt.authUserRole)
			router.POST("/users/:id/promote", handler.PromoteUser)
			tt.mockSetup(mockService, userID) // userID is the target ID in the path

			req, _ := http.NewRequest(http.MethodPost, "/users/"+tt.pathUserID+"/promote", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockService.AssertExpectations(t)
		})
	}
}

// Helper function for pointers
func ptr(s string) *string {
	return &s
}

// Helper function to create test user
func createTestUser() *models.User {
	return &models.User{
		UserID:         uuid.New(),
		Username:       "testuser",
		Email:          "test@example.com",
		PasswordHash:   "hashedpassword",
		Role:           "user",
		AvailableVotes: 5,
	}
}
