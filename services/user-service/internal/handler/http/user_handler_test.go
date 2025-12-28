package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/dto"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/mocks"
)

func TestUserHandler_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()

		mockUser := &domain.UserWithProfile{
			User: domain.User{
				UserID: userID,
				Email:  "test@example.com",
			},
			UserProfile: domain.UserProfile{
				FirstName: "Test",
				LastName:  "User",
			},
		}

		mockUserService.EXPECT().
			GetUser(gomock.Any(), userID).
			Return(mockUser, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/users/me", nil)
		c.Set("user_id", userID) // Simulate Auth Middleware

		handler.GetMe(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, "test@example.com", data["email"])
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/users/me", nil)
		// Missing userID in context

		handler.GetMe(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUserHandler_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	t.Run("Success", func(t *testing.T) {
		roleID := uuid.New()
		req := dto.CreateUserRequest{
			RegisterNo: 12345,
			Email:      "new@example.com",
			Password:   "password123",
			FirstName:  "New",
			LastName:   "User",
			RoleID:     roleID.String(),
		}
		jsonValue, _ := json.Marshal(req)

		mockUserService.EXPECT().
			CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx interface{}, user *domain.User, profile *domain.UserProfile) error {
				assert.Equal(t, req.Email, user.Email)
				assert.Equal(t, req.RegisterNo, user.RegisterNo)
				assert.Equal(t, req.FirstName, profile.FirstName)
				user.UserID = uuid.New() // Simulate ID generation
				return nil
			})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/admin/users", bytes.NewBuffer(jsonValue))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateUser(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		roleID := uuid.New()
		req := dto.CreateUserRequest{
			RegisterNo: 12345,
			Email:      "existing@example.com",
			Password:   "password123",
			FirstName:  "Existing",
			LastName:   "User",
			RoleID:     roleID.String(),
		}
		jsonValue, _ := json.Marshal(req)

		mockUserService.EXPECT().
			CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(domain.ErrUserAlreadyExists)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/admin/users", bytes.NewBuffer(jsonValue))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateUser(c)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}
