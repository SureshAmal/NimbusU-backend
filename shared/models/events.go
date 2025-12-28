package models

import (
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of event
type EventType string

const (
	// User events
	EventUserCreated   EventType = "USER_CREATED"
	EventUserUpdated   EventType = "USER_UPDATED"
	EventUserDeleted   EventType = "USER_DELETED"
	EventUserActivated EventType = "USER_ACTIVATED"
	EventUserSuspended EventType = "USER_SUSPENDED"

	// Auth events
	EventLoginSuccess    EventType = "LOGIN_SUCCESS"
	EventLoginFailed     EventType = "LOGIN_FAILED"
	EventLogout          EventType = "LOGOUT"
	EventPasswordChanged EventType = "PASSWORD_CHANGED"
)

// BaseEvent represents common fields for all events
type BaseEvent struct {
	EventID     uuid.UUID              `json:"event_id"`
	EventType   EventType              `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	ServiceName string                 `json:"service_name"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// UserEvent represents user-related events
type UserEvent struct {
	BaseEvent
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	RoleID    uuid.UUID `json:"role_id,omitempty"`
	Status    string    `json:"status,omitempty"`
}

// AuthEvent represents authentication events
type AuthEvent struct {
	BaseEvent
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Success     bool      `json:"success"`
	ErrorReason string    `json:"error_reason,omitempty"`
}

// NewUserEvent creates a new user event
func NewUserEvent(eventType EventType, userID uuid.UUID, email string) *UserEvent {
	return &UserEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New(),
			EventType:   eventType,
			Timestamp:   time.Now(),
			ServiceName: "user-service",
		},
		UserID: userID,
		Email:  email,
	}
}

// NewAuthEvent creates a new auth event
func NewAuthEvent(eventType EventType, userID uuid.UUID, email, ipAddress, userAgent string, success bool) *AuthEvent {
	return &AuthEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New(),
			EventType:   eventType,
			Timestamp:   time.Now(),
			ServiceName: "user-service",
		},
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
	}
}
