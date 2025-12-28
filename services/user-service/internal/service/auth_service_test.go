package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/mocks"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockTokenRepo := mocks.NewMockPasswordResetTokenRepository(ctrl)
	mockActivityRepo := mocks.NewMockActivityLogRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	// We need real JWT manager
	jwtManager := utils.NewJWTManager("secret", 3600, 3600)

	service := NewAuthService(
		mockUserRepo,
		mockProfileRepo,
		mockRoleRepo,
		mockSessionRepo,
		mockTokenRepo,
		mockActivityRepo,
		jwtManager,
		mockProducer,
		3600,
	)

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := utils.HashPassword(password)
		userID := uuid.New()
		roleID := uuid.New()

		user := &domain.User{
			UserID:       userID,
			Email:        email,
			PasswordHash: hashedPassword,
			RoleID:       roleID,
			Status:       "active",
		}

		profile := &domain.UserProfile{
			UserID:    userID,
			FirstName: "Test",
			LastName:  "User",
		}

		role := &domain.Role{
			RoleID:   roleID,
			RoleName: "student",
		}

		// Expectations
		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(user, nil)
		mockProfileRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(profile, nil)
		mockRoleRepo.EXPECT().GetByID(gomock.Any(), roleID).Return(role, nil)
		mockSessionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockUserRepo.EXPECT().UpdateLastLogin(gomock.Any(), userID).Return(nil)
		mockActivityRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)                 // Log success
		mockProducer.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil) // Login event

		// Execute
		accessToken, refreshToken, userWithProfile, err := service.Login(context.Background(), email, password, "127.0.0.1", "Go-Test")

		// Assertions
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.Equal(t, email, userWithProfile.Email)
		assert.Equal(t, "Test", userWithProfile.UserProfile.FirstName)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		email := "test@example.com"
		password := "wrongpassword"
		hashedPassword, _ := utils.HashPassword("correctpassword")
		userID := uuid.New()

		user := &domain.User{
			UserID:       userID,
			Email:        email,
			PasswordHash: hashedPassword,
			Status:       "active",
		}

		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(user, nil)
		mockActivityRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)                 // Log failure
		mockProducer.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil) // Login failed event

		accessToken, refreshToken, _, err := service.Login(context.Background(), email, password, "127.0.0.1", "Go-Test")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})

	t.Run("User Not Found", func(t *testing.T) {
		email := "unknown@example.com"

		mockUserRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, domain.ErrUserNotFound)
		mockActivityRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil) // Log failure
		mockProducer.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		accessToken, refreshToken, _, err := service.Login(context.Background(), email, "any", "127.0.0.1", "Go-Test")

		assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
	})
}

func TestAuthService_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockSessionRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	// Other mocks needed for constructor but unused in Logout
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockRoleRepo := mocks.NewMockRoleRepository(ctrl)
	mockTokenRepo := mocks.NewMockPasswordResetTokenRepository(ctrl)
	mockActivityRepo := mocks.NewMockActivityLogRepository(ctrl)
	jwtManager := utils.NewJWTManager("secret", 3600, 3600)

	service := NewAuthService(
		mockUserRepo,
		mockProfileRepo,
		mockRoleRepo,
		mockSessionRepo,
		mockTokenRepo,
		mockActivityRepo,
		jwtManager,
		mockProducer,
		3600,
	)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		sessionID := uuid.New()
		refreshToken := "valid_refresh"

		session := &domain.ActiveSession{
			SessionID: sessionID,
			UserID:    userID,
		}

		user := &domain.User{
			UserID: userID,
			Email:  "test@example.com",
		}

		mockSessionRepo.EXPECT().GetByRefreshToken(gomock.Any(), refreshToken).Return(session, nil)
		mockSessionRepo.EXPECT().Delete(gomock.Any(), sessionID).Return(nil)
		mockUserRepo.EXPECT().GetByID(gomock.Any(), userID).Return(user, nil)
		mockProducer.EXPECT().PublishEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		err := service.Logout(context.Background(), userID, refreshToken)

		assert.NoError(t, err)
	})
}
