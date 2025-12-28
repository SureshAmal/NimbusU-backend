package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/dto"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/mocks"
)

func TestAuthHandler_Login(t *testing.T) {
	// Setup Gin to test mode
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)

	t.Run("Success", func(t *testing.T) {
		// Mock request data
		req := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(req)

		// Mock response data from service
		mockUser := &domain.UserWithProfile{
			User: domain.User{
				Email: "test@example.com",
			},
		}
		accessToken := "access_token"
		refreshToken := "refresh_token"

		// Expect Login call
		mockAuthService.EXPECT().
			Login(gomock.Any(), req.Email, req.Password, gomock.Any(), gomock.Any()).
			Return(accessToken, refreshToken, mockUser, nil)

		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonValue))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call handler
		handler.Login(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, accessToken, data["access_token"])
		assert.Equal(t, refreshToken, data["refresh_token"])
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		req := dto.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		jsonValue, _ := json.Marshal(req)

		mockAuthService.EXPECT().
			Login(gomock.Any(), req.Email, req.Password, gomock.Any(), gomock.Any()).
			Return("", "", nil, domain.ErrInvalidCredentials)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonValue))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)

	t.Run("Success", func(t *testing.T) {
		req := dto.RefreshTokenRequest{
			RefreshToken: "valid_refresh_token",
		}
		jsonValue, _ := json.Marshal(req)

		newAccessToken := "new_access_token"
		newRefreshToken := "new_refresh_token"

		mockAuthService.EXPECT().
			RefreshToken(gomock.Any(), req.RefreshToken).
			Return(newAccessToken, newRefreshToken, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(jsonValue))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshToken(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, newAccessToken, data["access_token"])
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := dto.RefreshTokenRequest{
			RefreshToken: "invalid_token",
		}
		jsonValue, _ := json.Marshal(req)

		mockAuthService.EXPECT().
			RefreshToken(gomock.Any(), req.RefreshToken).
			Return("", "", domain.ErrInvalidToken)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(jsonValue))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.RefreshToken(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
