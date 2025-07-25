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

		// Use valid password to isolate email exists error
		_, err := service.Register(context.Background(), "testuser", "test@example.com", "ValidPass123!")

		assert.ErrorIs(t, err, ErrEmailExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		// Don't return existing user to trigger password validation
		mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return((*models.User)(nil), nil)

		_, err := service.Register(context.Background(), "testuser", "test@example.com", "weak")

		assert.ErrorIs(t, err, ErrPasswordValidation)
		mockRepo.AssertNotCalled(t, "Create")
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

func TestUserService_UpdateUser(t *testing.T) {
	userID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		existingUser := &models.User{
			UserID: userID,
			Email:  "old@example.com",
		}
		updateData := map[string]interface{}{
			"username": "newusername",
		}

		mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

		user, err := service.UpdateUser(context.Background(), userID, updateData)

		assert.NoError(t, err)
		assert.Equal(t, "newusername", user.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email conflict", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		existingUser := &models.User{UserID: userID, Email: "old@example.com"}
		conflictUser := &models.User{UserID: uuid.New(), Email: "new@example.com"}
		updateData := map[string]interface{}{"email": "new@example.com"}

		mockRepo.On("GetByID", mock.Anything, userID).Return(existingUser, nil)
		mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return(conflictUser, nil)

		_, err := service.UpdateUser(context.Background(), userID, updateData)

		assert.ErrorIs(t, err, ErrEmailExists)
		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestUserService_PromoteToAdmin(t *testing.T) {
	t.Run("successful promotion", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)
		userID := uuid.New()

		user := &models.User{UserID: userID, Role: "user"}
		mockRepo.On("GetByID", mock.Anything, userID).Return(user, nil)
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *models.User) bool {
			return u.Role == "admin"
		})).Return(nil)

		err := service.PromoteToAdmin(context.Background(), userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
