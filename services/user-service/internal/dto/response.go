package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
)

// UserResponse represents a user in API responses
type UserResponse struct {
	UserID     uuid.UUID  `json:"user_id"`
	RegisterNo int64      `json:"register_no"`
	Email      string     `json:"email"`
	RoleID     uuid.UUID  `json:"role_id"`
	RoleName   string     `json:"role_name,omitempty"`
	Status     string     `json:"status"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ProfileResponse represents a user profile in API responses
type ProfileResponse struct {
	ProfileID         uuid.UUID `json:"profile_id"`
	UserID            uuid.UUID `json:"user_id"`
	RegisterNo        int64     `json:"register_no"`
	FirstName         string    `json:"first_name"`
	MiddleName        *string   `json:"middle_name"`
	LastName          string    `json:"last_name"`
	Phone             *string   `json:"phone"`
	Gender            *string   `json:"gender"`
	ProfilePictureURL *string   `json:"profile_picture_url"`
	Bio               *string   `json:"bio"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserWithProfileResponse combines user and profile data
type UserWithProfileResponse struct {
	UserResponse
	Profile ProfileResponse `json:"profile"`
}

// LoginResponse represents login response data
type LoginResponse struct {
	AccessToken  string                  `json:"access_token"`
	RefreshToken string                  `json:"refresh_token"`
	TokenType    string                  `json:"token_type"`
	ExpiresIn    int                     `json:"expires_in"`
	User         UserWithProfileResponse `json:"user"`
}

// RefreshTokenResponse represents token refresh response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// RoleResponse represents a role in API responses
type RoleResponse struct {
	RoleID      uuid.UUID `json:"role_id"`
	RoleName    string    `json:"role_name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// PermissionResponse represents a permission in API responses
type PermissionResponse struct {
	PermissionID   uuid.UUID `json:"permission_id"`
	PermissionName string    `json:"permission_name"`
	Resource       string    `json:"resource"`
	Action         string    `json:"action"`
	Description    *string   `json:"description"`
}

// SessionResponse represents an active session
type SessionResponse struct {
	SessionID  uuid.UUID `json:"session_id"`
	DeviceInfo *string   `json:"device_info"`
	IPAddress  *string   `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// Helper functions to convert domain entities to DTOs

func ToUserResponse(user *domain.User, role *domain.Role) UserResponse {
	response := UserResponse{
		UserID:     user.UserID,
		RegisterNo: user.RegisterNo,
		Email:      user.Email,
		RoleID:     user.RoleID,
		Status:     user.Status,
		LastLogin:  user.LastLogin,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	if role != nil {
		response.RoleName = role.RoleName
	}

	return response
}

func ToProfileResponse(profile *domain.UserProfile) ProfileResponse {
	return ProfileResponse{
		ProfileID:         profile.ProfileID,
		UserID:            profile.UserID,
		RegisterNo:        profile.RegisterNo,
		FirstName:         profile.FirstName,
		MiddleName:        profile.MiddleName,
		LastName:          profile.LastName,
		Phone:             profile.Phone,
		Gender:            profile.Gender,
		ProfilePictureURL: profile.ProfilePictureURL,
		Bio:               profile.Bio,
		CreatedAt:         profile.CreatedAt,
		UpdatedAt:         profile.UpdatedAt,
	}
}

func ToUserWithProfileResponse(userWithProfile *domain.UserWithProfile) UserWithProfileResponse {
	return UserWithProfileResponse{
		UserResponse: ToUserResponse(&userWithProfile.User, userWithProfile.Role),
		Profile:      ToProfileResponse(&userWithProfile.UserProfile),
	}
}

func ToRoleResponse(role *domain.Role) RoleResponse {
	return RoleResponse{
		RoleID:      role.RoleID,
		RoleName:    role.RoleName,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
	}
}

func ToPermissionResponse(permission *domain.Permission) PermissionResponse {
	return PermissionResponse{
		PermissionID:   permission.PermissionID,
		PermissionName: permission.PermissionName,
		Resource:       permission.Resource,
		Action:         permission.Action,
		Description:    permission.Description,
	}
}

func ToSessionResponse(session *domain.ActiveSession) SessionResponse {
	return SessionResponse{
		SessionID:  session.SessionID,
		DeviceInfo: session.DeviceInfo,
		IPAddress:  session.IPAddress,
		CreatedAt:  session.CreatedAt,
		ExpiresAt:  session.ExpiresAt,
	}
}
