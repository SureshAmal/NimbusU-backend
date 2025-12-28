package dto

import (
	"testing"
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToUserResponse(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()
	now := time.Now()

	user := &domain.User{
		UserID:     userID,
		RegisterNo: 12345,
		Email:      "test@example.com",
		RoleID:     roleID,
		Status:     "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	role := &domain.Role{
		RoleID:   roleID,
		RoleName: "student",
	}

	response := ToUserResponse(user, role)

	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "student", response.RoleName)
	assert.Equal(t, now, response.CreatedAt)
}

func TestToUserResponse_NoRole(t *testing.T) {
	userID := uuid.New()
	user := &domain.User{
		UserID: userID,
		Email:  "test@example.com",
	}

	response := ToUserResponse(user, nil)

	assert.Equal(t, userID, response.UserID)
	assert.Empty(t, response.RoleName)
}

func TestToProfileResponse(t *testing.T) {
	profileID := uuid.New()
	userID := uuid.New()

	profile := &domain.UserProfile{
		ProfileID: profileID,
		UserID:    userID,
		FirstName: "John",
		LastName:  "Doe",
		Phone:     nil, // Test nil pointer
	}

	response := ToProfileResponse(profile)

	assert.Equal(t, profileID, response.ProfileID)
	assert.Equal(t, "John", response.FirstName)
	assert.Nil(t, response.Phone)
}

func TestToUserWithProfileResponse(t *testing.T) {
	userWithProfile := &domain.UserWithProfile{
		User: domain.User{
			Email: "test@example.com",
		},
		UserProfile: domain.UserProfile{
			FirstName: "John",
		},
		Role: &domain.Role{
			RoleName: "admin",
		},
	}

	response := ToUserWithProfileResponse(userWithProfile)

	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "John", response.Profile.FirstName)
	assert.Equal(t, "admin", response.RoleName)
}
