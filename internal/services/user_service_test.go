package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) DecrementAvailableVotes(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) IncrementAvailableVotes(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Test helper functions to reduce duplication
func setupTest() (*MockUserRepository, UserService) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	return mockRepo, service
}

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

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		mockSetup   func(*MockUserRepository)
		expectedErr error
	}{
		{
			name:     "successful registration",
			email:    "test@example.com",
			password: "ValidPass123!",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), nil)
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:     "email already exists",
			email:    "test@example.com",
			password: "ValidPass123!",
			mockSetup: func(m *MockUserRepository) {
				existingUser := &models.User{Email: "test@example.com"}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)
			},
			expectedErr: ErrEmailExists,
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "weak",
			mockSetup: func(m *MockUserRepository) {
				// No repository calls expected - password validation fails first
			},
			expectedErr: ErrPasswordValidation,
		},
		{
			name:     "invalid email format",
			email:    "invalid-email",
			password: "ValidPass123!",
			mockSetup: func(m *MockUserRepository) {
				// No repository calls expected for email validation
			},
			expectedErr: errors.New("invalid email format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setupTest()
			tt.mockSetup(mockRepo)

			user, err := service.Register(context.Background(), "testuser", tt.email, tt.password)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				if tt.expectedErr == ErrEmailExists {
					mockRepo.AssertNotCalled(t, "Create")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "user", user.Role)
				assert.Equal(t, 5, user.AvailableVotes)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		password    string
		mockSetup   func(*MockUserRepository)
		expectToken bool
		expectedErr error
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "correctpassword",
			mockSetup: func(m *MockUserRepository) {
				hashedPassword, _ := hashPassword("correctpassword")
				user := &models.User{
					Email:        "test@example.com",
					PasswordHash: hashedPassword,
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectToken: true,
			expectedErr: nil,
		},
		{
			name:     "invalid credentials - wrong password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockSetup: func(m *MockUserRepository) {
				hashedPassword, _ := hashPassword("correctpassword")
				user := &models.User{
					Email:        "test@example.com",
					PasswordHash: hashedPassword,
				}
				m.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
			},
			expectToken: false,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "anypassword",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return((*models.User)(nil), nil)
			},
			expectToken: false,
			expectedErr: ErrInvalidCredentials,
		},
		{
			name:     "repository error",
			email:    "test@example.com",
			password: "anypassword",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), errors.New("db error"))
			},
			expectToken: false,
			expectedErr: errors.New("failed to get user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setupTest()
			tt.mockSetup(mockRepo)

			token, err := service.Login(context.Background(), tt.email, tt.password)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.expectToken {
					assert.NotEmpty(t, token)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserProfile(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func(*MockUserRepository, uuid.UUID)
		expectedErr error
	}{
		{
			name:   "successful get",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
			},
			expectedErr: nil,
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				m.On("GetByID", mock.Anything, userID).Return((*models.User)(nil), nil)
			},
			expectedErr: ErrInvalidID,
		},
		{
			name:   "repository error",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				m.On("GetByID", mock.Anything, userID).Return((*models.User)(nil), errors.New("db error"))
			},
			expectedErr: errors.New("failed to get user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setupTest()
			tt.mockSetup(mockRepo, tt.userID)

			user, err := service.GetUserProfile(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.UserID)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name        string
		updateData  map[string]interface{}
		mockSetup   func(*MockUserRepository, uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful update - username",
			updateData: map[string]interface{}{
				"username": "newusername",
			},
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "successful update - email",
			updateData: map[string]interface{}{
				"email": "new@example.com",
			},
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
				m.On("GetByEmail", mock.Anything, "new@example.com").Return((*models.User)(nil), nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "email conflict",
			updateData: map[string]interface{}{
				"email": "existing@example.com",
			},
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				conflictUser := &models.User{UserID: uuid.New(), Email: "existing@example.com"}
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
				m.On("GetByEmail", mock.Anything, "existing@example.com").Return(conflictUser, nil)
			},
			expectedErr: ErrEmailExists,
		},
		{
			name: "invalid email format",
			updateData: map[string]interface{}{
				"email": "invalid-email",
			},
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
			},
			expectedErr: errors.New("invalid email format"),
		},
		{
			name: "invalid password",
			updateData: map[string]interface{}{
				"password": "weak",
			},
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
			},
			expectedErr: ErrPasswordValidation,
		},
		{
			name: "user not found",
			updateData: map[string]interface{}{
				"username": "newusername",
			},
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				m.On("GetByID", mock.Anything, userID).Return((*models.User)(nil), nil)
			},
			expectedErr: ErrInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setupTest()
			tt.mockSetup(mockRepo, userID)

			user, err := service.UpdateUser(context.Background(), userID, tt.updateData)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
				if tt.expectedErr == ErrEmailExists || tt.expectedErr == ErrInvalidID {
					mockRepo.AssertNotCalled(t, "Update")
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func(*MockUserRepository, uuid.UUID)
		expectedErr error
	}{
		{
			name:   "successful delete",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				m.On("Delete", mock.Anything, userID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "repository error",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				m.On("Delete", mock.Anything, userID).Return(errors.New("db error"))
			},
			expectedErr: errors.New("failed to delete user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setupTest()
			tt.mockSetup(mockRepo, tt.userID)

			err := service.DeleteUser(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_PromoteToAdmin(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func(*MockUserRepository, uuid.UUID)
		expectedErr error
	}{
		{
			name:   "successful promotion",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
					return u.Role == "admin"
				})).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				m.On("GetByID", mock.Anything, userID).Return((*models.User)(nil), nil)
			},
			expectedErr: ErrInvalidID,
		},
		{
			name:   "repository error on get",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				m.On("GetByID", mock.Anything, userID).Return((*models.User)(nil), errors.New("db error"))
			},
			expectedErr: errors.New("failed to get user"),
		},
		{
			name:   "repository error on update",
			userID: uuid.New(),
			mockSetup: func(m *MockUserRepository, userID uuid.UUID) {
				user := createTestUser()
				user.UserID = userID
				m.On("GetByID", mock.Anything, userID).Return(user, nil)
				m.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("db error"))
			},
			expectedErr: errors.New("failed to promote user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setupTest()
			tt.mockSetup(mockRepo, tt.userID)

			err := service.PromoteToAdmin(context.Background(), tt.userID)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockUserRepository)
		expectedErr error
	}{
		{
			name: "successful get all",
			mockSetup: func(m *MockUserRepository) {
				users := []models.User{*createTestUser(), *createTestUser()}
				m.On("GetAll", mock.Anything).Return(users, nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetAll", mock.Anything).Return([]models.User{}, errors.New("db error"))
			},
			expectedErr: errors.New(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setupTest()
			tt.mockSetup(mockRepo)

			users, err := service.GetAllUsers(context.Background())

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Empty(t, users)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, users)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
