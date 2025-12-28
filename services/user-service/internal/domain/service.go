package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrProfileNotFound    = errors.New("profile not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrUnauthorized       = errors.New("unauthorized")
)

// UserService defines business logic for user management
type UserService interface {
	// User management
	CreateUser(ctx context.Context, user *User, profile *UserProfile) error
	GetUser(ctx context.Context, userID uuid.UUID) (*UserWithProfile, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*UserWithProfile, int64, error)

	// User status
	ActivateUser(ctx context.Context, userID uuid.UUID) error
	SuspendUser(ctx context.Context, userID uuid.UUID) error

	// Profile management
	UpdateProfile(ctx context.Context, userID uuid.UUID, profile *UserProfile) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error)

	// Bulk operations
	BulkCreateUsers(ctx context.Context, users []*User, profiles []*UserProfile) error
}

// AuthService defines business logic for authentication
type AuthService interface {
	// Authentication
	Login(ctx context.Context, email, password string, ipAddress, userAgent string) (accessToken, refreshToken string, user *UserWithProfile, err error)
	Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)

	// Password management
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	RequestPasswordReset(ctx context.Context, email string) (token string, err error)
	ResetPassword(ctx context.Context, token, newPassword string) error

	// Session management
	GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]*ActiveSession, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeAllSessions(ctx context.Context, userID uuid.UUID) error
}

// RoleService defines business logic for role management
type RoleService interface {
	CreateRole(ctx context.Context, role *Role) error
	GetRole(ctx context.Context, roleID uuid.UUID) (*Role, error)
	GetRoleByName(ctx context.Context, roleName string) (*Role, error)
	ListRoles(ctx context.Context) ([]*Role, error)
	UpdateRole(ctx context.Context, roleID uuid.UUID, updates map[string]interface{}) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	// Permission management
	AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	RevokePermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)
}

// EventProducer defines interface for publishing events
type EventProducer interface {
	PublishEvent(topic string, key string, event interface{}) error
	Close() error
}
