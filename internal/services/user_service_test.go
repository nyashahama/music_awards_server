package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/nyashahama/music-awards/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements UserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

// Implement all other UserRepository methods similarly...

func TestUserService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), nil)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

		user, err := service.Register(context.Background(), "testuser", "test@example.com", "ValidPass123!")

		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "user", user.Role)
		assert.Equal(t, 5, user.AvailableVotes)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		existingUser := &models.User{Email: "test@example.com"}
		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

		_, err := service.Register(context.Background(), "testuser", "test@example.com", "password")

		assert.ErrorIs(t, err, ErrEmailExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		_, err := service.Register(context.Background(), "testuser", "test@example.com", "short")

		assert.ErrorIs(t, err, ErrPasswordValidation)
	})
}

func TestUserService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		hashedPassword, _ := hashPassword("correctpassword")
		user := &models.User{
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
		}

		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

		token, err := service.Login(context.Background(), "test@example.com", "correctpassword")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		hashedPassword, _ := hashPassword("correctpassword")
		user := &models.User{
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
		}

		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

		_, err := service.Login(context.Background(), "test@example.com", "wrongpassword")

		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})
}

// Add similar tests for other methods: UpdateUser, DeleteUser, PromoteToAdmin, etc.
