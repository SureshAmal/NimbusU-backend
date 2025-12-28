package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	RegisterNo   int64      `json:"register_no" db:"register_no"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"` // Never expose in JSON
	RoleID       uuid.UUID  `json:"role_id" db:"role_id"`
	Status       string     `json:"status" db:"status"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy    *uuid.UUID `json:"created_by" db:"created_by"`
}

// Role represents a user role
type Role struct {
	RoleID      uuid.UUID `json:"role_id" db:"role_id"`
	RoleName    string    `json:"role_name" db:"role_name"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Permission represents a granular permission
type Permission struct {
	PermissionID   uuid.UUID `json:"permission_id" db:"permission_id"`
	PermissionName string    `json:"permission_name" db:"permission_name"`
	Resource       string    `json:"resource" db:"resource"`
	Action         string    `json:"action" db:"action"`
	Description    *string   `json:"description" db:"description"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RolePermissionID uuid.UUID `json:"role_permission_id" db:"role_permission_id"`
	RoleID           uuid.UUID `json:"role_id" db:"role_id"`
	PermissionID     uuid.UUID `json:"permission_id" db:"permission_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// UserProfile represents extended user information
type UserProfile struct {
	ProfileID         uuid.UUID `json:"profile_id" db:"profile_id"`
	UserID            uuid.UUID `json:"user_id" db:"user_id"`
	RegisterNo        int64     `json:"register_no" db:"register_no"`
	FirstName         string    `json:"first_name" db:"first_name"`
	MiddleName        *string   `json:"middle_name" db:"middle_name"`
	LastName          string    `json:"last_name" db:"last_name"`
	Phone             *string   `json:"phone" db:"phone"`
	Gender            *string   `json:"gender" db:"gender"`
	ProfilePictureURL *string   `json:"profile_picture_url" db:"profile_picture_url"`
	Bio               *string   `json:"bio" db:"bio"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// UserActivityLog represents an audit log entry
type UserActivityLog struct {
	LogID        uuid.UUID  `json:"log_id" db:"log_id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Action       string     `json:"action" db:"action"`
	ResourceType *string    `json:"resource_type" db:"resource_type"`
	ResourceID   *uuid.UUID `json:"resource_id" db:"resource_id"`
	IPAddress    *string    `json:"ip_address" db:"ip_address"`
	UserAgent    *string    `json:"user_agent" db:"user_agent"`
	Details      *string    `json:"details" db:"details"` // JSON string
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	TokenID   uuid.UUID `json:"token_id" db:"token_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ActiveSession represents an active user session
type ActiveSession struct {
	SessionID    uuid.UUID `json:"session_id" db:"session_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	DeviceInfo   *string   `json:"device_info" db:"device_info"`
	IPAddress    *string   `json:"ip_address" db:"ip_address"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// UserWithProfile combines User and UserProfile
type UserWithProfile struct {
	User
	UserProfile
	Role *Role `json:"role,omitempty"`
}
