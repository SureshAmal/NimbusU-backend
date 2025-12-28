package http

import (
	"net/http"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/dto"
	"github.com/SureshAmal/NimbusU-backend/shared/middleware"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// @Summary      User Login
// @Description  Authenticate user and return access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login Credentials"
// @Success      200  {object}  utils.APIResponse{data=dto.LoginResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	accessToken, refreshToken, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password, ipAddress, userAgent)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Login failed", err)
		return
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		User:         dto.ToUserWithProfileResponse(user),
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// RefreshToken handles token refresh
// @Summary      Refresh Access Token
// @Description  Get a new access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh Token"
// @Success      200  {object}  utils.APIResponse{data=dto.RefreshTokenResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed", response)
}

// Logout handles user logout
// @Summary      User Logout
// @Description  Invalidate current session and refresh token
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh Token"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID, req.RefreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Logout failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logout successful", nil)
}

// ChangePassword handles password change
// @Summary      Change Password
// @Description  Change the current user's password
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.ChangePasswordRequest true "Password Change Data"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/password/change [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err == domain.ErrInvalidCredentials {
			utils.ErrorResponse(c, http.StatusBadRequest, "Old password is incorrect", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Password change failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// RequestPasswordReset handles password reset request
// @Summary      Request Password Reset
// @Description  Request a password reset email (generates token)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.PasswordResetRequestRequest true "Email Address"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/password/reset-request [post]
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.PasswordResetRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	token, err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Password reset request failed", err)
		return
	}

	// In production, this token should be sent via email, not in response
	// For now, we return it for testing purposes
	utils.SuccessResponse(c, http.StatusOK, "Password reset email sent", gin.H{"token": token})
}

// ResetPassword handles password reset with token
// @Summary      Reset Password
// @Description  Reset password using a valid token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "Reset Data"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/password/reset [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		if err == domain.ErrInvalidToken {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Password reset failed", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password reset successful", nil)
}

// GetActiveSessions returns all active sessions for the current user
// @Summary      Get Active Sessions
// @Description  Retrieve a list of all active sessions for the current user
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]dto.SessionResponse}
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/sessions [get]
func (h *AuthHandler) GetActiveSessions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	sessions, err := h.authService.GetActiveSessions(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sessions", err)
		return
	}

	sessionResponses := make([]dto.SessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = dto.ToSessionResponse(session)
	}

	utils.SuccessResponse(c, http.StatusOK, "Sessions retrieved", sessionResponses)
}

// RevokeSession revokes a specific session
// @Summary      Revoke Session
// @Description  Revoke a specific session by ID
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        sessionId path string true "Session ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/sessions/{sessionId} [delete]
func (h *AuthHandler) RevokeSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("sessionId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID", err)
		return
	}

	if err := h.authService.RevokeSession(c.Request.Context(), sessionID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to revoke session", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Session revoked", nil)
}

// RevokeAllSessions revokes all sessions for the current user
// @Summary      Revoke All Sessions
// @Description  Revoke all active sessions for the current user (logout everywhere)
// @Tags         auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /auth/sessions [delete]
func (h *AuthHandler) RevokeAllSessions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	if err := h.authService.RevokeAllSessions(c.Request.Context(), userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to revoke sessions", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "All sessions revoked", nil)
}
