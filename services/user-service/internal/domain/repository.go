package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// User CRUD
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByRegisterNo(ctx context.Context, registerNo int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID uuid.UUID) error

	// List and search
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*User, int64, error)

	// Status management
	UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

// UserProfileRepository defines the interface for user profile operations
type UserProfileRepository interface {
	Create(ctx context.Context, profile *UserProfile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	Update(ctx context.Context, profile *UserProfile) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

// RoleRepository defines the interface for role operations
type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, roleID uuid.UUID) (*Role, error)
	GetByName(ctx context.Context, roleName string) (*Role, error)
	List(ctx context.Context) ([]*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, roleID uuid.UUID) error
}

// PermissionRepository defines the interface for permission operations
type PermissionRepository interface {
	Create(ctx context.Context, permission *Permission) error
	GetByID(ctx context.Context, permissionID uuid.UUID) (*Permission, error)
	List(ctx context.Context) ([]*Permission, error)
	GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)
}

// RolePermissionRepository defines the interface for role-permission mapping
type RolePermissionRepository interface {
	AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	RevokePermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)
}

// ActivityLogRepository defines the interface for activity logging
type ActivityLogRepository interface {
	Create(ctx context.Context, log *UserActivityLog) error
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*UserActivityLog, int64, error)
	GetByAction(ctx context.Context, action string, limit, offset int) ([]*UserActivityLog, int64, error)
}

// PasswordResetTokenRepository defines the interface for password reset tokens
type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *PasswordResetToken) error
	GetByToken(ctx context.Context, token string) (*PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// SessionRepository defines the interface for session management
type SessionRepository interface {
	Create(ctx context.Context, session *ActiveSession) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*ActiveSession, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*ActiveSession, error)
	Delete(ctx context.Context, sessionID uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
