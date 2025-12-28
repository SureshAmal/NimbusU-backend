package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	mockActivityLog := mocks.NewMockActivityLogRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewUserService(mockUserRepo, mockProfileRepo, mockRoleRepo, mockActivityLog, mockProducer)

	t.Run("Success", func(t *testing.T) {
		roleID := uuid.New()
		user := &domain.User{
			Email:        "new@example.com",
			RegisterNo:   12345,
			PasswordHash: "password123",
			RoleID:       roleID,
		}
		profile := &domain.UserProfile{
			FirstName: "New",
			LastName:  "User",
		}

		// Expectations
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), user.Email).Return(nil, domain.ErrUserNotFound)
		mockUserRepo.EXPECT().GetByRegisterNo(gomock.Any(), user.RegisterNo).Return(nil, domain.ErrUserNotFound)
		mockUserRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, u *domain.User) error {
			assert.NotEmpty(t, u.UserID)
			assert.NotEqual(t, "password123", u.PasswordHash) // Should be hashed
			return nil
		})
		mockProfileRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockProducer.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := service.CreateUser(context.Background(), user, profile)
		assert.NoError(t, err)
	})

	t.Run("Email Exists", func(t *testing.T) {
		user := &domain.User{
			Email: "exists@example.com",
		}
		profile := &domain.UserProfile{}

		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), user.Email).Return(&domain.User{}, nil)

		err := service.CreateUser(context.Background(), user, profile)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)

	// Unused for GetUser but required for NewUserService
	mockActivityLog := mocks.NewMockActivityLogRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewUserService(mockUserRepo, mockProfileRepo, mockRoleRepo, mockActivityLog, mockProducer)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		roleID := uuid.New()

		user := &domain.User{
			UserID: userID,
			Email:  "test@example.com",
			RoleID: roleID,
		}
		profile := &domain.UserProfile{
			UserID:    userID,
			FirstName: "Test",
		}
		role := &domain.Role{
			RoleID:   roleID,
			RoleName: "student",
		}

		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
		mockProfileRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(profile, nil)
		mockRoleRepo.EXPECT().GetByID(gomock.Any(), roleID).Return(role, nil)

		result, err := service.GetUser(context.Background(), userID)

		assert.NoError(t, err)
		assert.Equal(t, user.Email, result.User.Email)
		assert.Equal(t, profile.FirstName, result.UserProfile.FirstName)
		assert.Equal(t, role.RoleName, result.Role.RoleName)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := uuid.New()
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(nil, domain.ErrUserNotFound)

		result, err := service.GetUser(context.Background(), userID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, result)
	})
}

// Password hashing helper for tests
func init() {
	// Reduce bcrypt cost for faster tests if configurable, but utility uses DefaultCost.
	// We just rely on it being reasonably fast.
}
