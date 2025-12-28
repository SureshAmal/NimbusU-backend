package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/shared/models"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
	"github.com/google/uuid"
)

type authService struct {
	userRepo           domain.UserRepository
	profileRepo        domain.UserProfileRepository
	roleRepo           domain.RoleRepository
	sessionRepo        domain.SessionRepository
	passwordTokenRepo  domain.PasswordResetTokenRepository
	activityLogRepo    domain.ActivityLogRepository
	jwtManager         *utils.JWTManager
	producer           domain.EventProducer
	refreshTokenExpiry time.Duration
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo domain.UserRepository,
	profileRepo domain.UserProfileRepository,
	roleRepo domain.RoleRepository,
	sessionRepo domain.SessionRepository,
	passwordTokenRepo domain.PasswordResetTokenRepository,
	activityLogRepo domain.ActivityLogRepository,
	jwtManager *utils.JWTManager,
	producer domain.EventProducer,
	refreshTokenExpiry int,
) domain.AuthService {
	return &authService{
		userRepo:           userRepo,
		profileRepo:        profileRepo,
		roleRepo:           roleRepo,
		sessionRepo:        sessionRepo,
		passwordTokenRepo:  passwordTokenRepo,
		activityLogRepo:    activityLogRepo,
		jwtManager:         jwtManager,
		producer:           producer,
		refreshTokenExpiry: time.Duration(refreshTokenExpiry) * time.Second,
	}
}

func (s *authService) Login(ctx context.Context, email, password string, ipAddress, userAgent string) (accessToken, refreshToken string, user *domain.UserWithProfile, err error) {
	// Get user by email
	foundUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Log failed login attempt
		s.logAuthEvent(ctx, uuid.Nil, email, ipAddress, userAgent, false, "user not found")
		return "", "", nil, domain.ErrInvalidCredentials
	}

	// Check if user is active
	if foundUser.Status != "active" {
		s.logAuthEvent(ctx, foundUser.UserID, email, ipAddress, userAgent, false, "user not active")
		return "", "", nil, domain.ErrInvalidCredentials
	}

	// Verify password
	if err := utils.VerifyPassword(foundUser.PasswordHash, password); err != nil {
		s.logAuthEvent(ctx, foundUser.UserID, email, ipAddress, userAgent, false, "invalid password")
		return "", "", nil, domain.ErrInvalidCredentials
	}

	// Get user profile and role
	profile, err := s.profileRepo.GetByUserID(ctx, foundUser.UserID)
	if err != nil {
		return "", "", nil, err
	}

	role, err := s.roleRepo.GetByID(ctx, foundUser.RoleID)
	if err != nil {
		return "", "", nil, err
	}

	// Generate tokens
	accessToken, err = s.jwtManager.GenerateAccessToken(foundUser.UserID, foundUser.Email, role.RoleID, role.RoleName)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = s.jwtManager.GenerateRefreshToken(foundUser.UserID)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session
	session := &domain.ActiveSession{
		SessionID:    uuid.New(),
		UserID:       foundUser.UserID,
		RefreshToken: refreshToken,
		DeviceInfo:   &userAgent,
		IPAddress:    &ipAddress,
		ExpiresAt:    time.Now().Add(s.refreshTokenExpiry),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", "", nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	s.userRepo.UpdateLastLogin(ctx, foundUser.UserID)

	// Log successful login
	s.logAuthEvent(ctx, foundUser.UserID, email, ipAddress, userAgent, true, "")

	// Publish login success event
	event := models.NewAuthEvent(models.EventLoginSuccess, foundUser.UserID, email, ipAddress, userAgent, true)
	s.producer.PublishEvent("auth.events", foundUser.UserID.String(), event)

	userWithProfile := &domain.UserWithProfile{
		User:        *foundUser,
		UserProfile: *profile,
		Role:        role,
	}

	return accessToken, refreshToken, userWithProfile, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	// Get session by refresh token
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	// Delete session
	if err := s.sessionRepo.Delete(ctx, session.SessionID); err != nil {
		return err
	}

	// Get user for event
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Publish logout event
	event := models.NewAuthEvent(models.EventLogout, userID, user.Email, "", "", true)
	s.producer.PublishEvent("auth.events", userID.String(), event)

	return nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	// Get session
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", domain.ErrInvalidToken
	}

	// Get user and role
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}

	if user.Status != "active" {
		return "", "", domain.ErrUnauthorized
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	accessToken, err = s.jwtManager.GenerateAccessToken(user.UserID, user.Email, role.RoleID, role.RoleName)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = s.jwtManager.GenerateRefreshToken(user.UserID)
	if err != nil {
		return "", "", err
	}

	// Delete old session and create new one
	s.sessionRepo.Delete(ctx, session.SessionID)

	newSession := &domain.ActiveSession{
		SessionID:    uuid.New(),
		UserID:       user.UserID,
		RefreshToken: newRefreshToken,
		DeviceInfo:   session.DeviceInfo,
		IPAddress:    session.IPAddress,
		ExpiresAt:    time.Now().Add(s.refreshTokenExpiry),
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := utils.VerifyPassword(user.PasswordHash, oldPassword); err != nil {
		return domain.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Revoke all sessions (force re-login)
	s.sessionRepo.DeleteByUserID(ctx, userID)

	// Publish password changed event
	event := models.NewAuthEvent(models.EventPasswordChanged, userID, user.Email, "", "", true)
	s.producer.PublishEvent("auth.events", userID.String(), event)

	return nil
}

func (s *authService) RequestPasswordReset(ctx context.Context, email string) (token string, err error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists
		return "", nil
	}

	// Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token = hex.EncodeToString(tokenBytes)

	// Create reset token
	resetToken := &domain.PasswordResetToken{
		TokenID:   uuid.New(),
		UserID:    user.UserID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		Used:      false,
	}

	if err := s.passwordTokenRepo.Create(ctx, resetToken); err != nil {
		return "", err
	}

	// TODO: Send email with reset link (would integrate with notification service)

	return token, nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get token
	resetToken, err := s.passwordTokenRepo.GetByToken(ctx, token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Mark token as used
	s.passwordTokenRepo.MarkAsUsed(ctx, resetToken.TokenID)

	// Revoke all sessions
	s.sessionRepo.DeleteByUserID(ctx, user.UserID)

	// Publish password changed event
	event := models.NewAuthEvent(models.EventPasswordChanged, user.UserID, user.Email, "", "", true)
	s.producer.PublishEvent("auth.events", user.UserID.String(), event)

	return nil
}

func (s *authService) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]*domain.ActiveSession, error) {
	return s.sessionRepo.GetByUserID(ctx, userID)
}

func (s *authService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *authService) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	return s.sessionRepo.DeleteByUserID(ctx, userID)
}

// Helper to log auth events
func (s *authService) logAuthEvent(ctx context.Context, userID uuid.UUID, email, ipAddress, userAgent string, success bool, errorReason string) {
	log := &domain.UserActivityLog{
		LogID:     uuid.New(),
		UserID:    userID,
		Action:    "login",
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
	}

	if !success {
		details := fmt.Sprintf(`{"success": false, "error": "%s"}`, errorReason)
		log.Details = &details
	}

	s.activityLogRepo.Create(ctx, log)

	// Publish failed login event
	if !success {
		event := models.NewAuthEvent(models.EventLoginFailed, userID, email, ipAddress, userAgent, false)
		event.ErrorReason = errorReason
		s.producer.PublishEvent("auth.events", email, event)
	}
}
