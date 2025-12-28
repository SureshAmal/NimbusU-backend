package http

import (
	"net/http"
	"strconv"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/dto"
	"github.com/SureshAmal/NimbusU-backend/shared/middleware"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetMe returns current user profile
// @Summary      Get My Profile
// @Description  Get the profile of the currently authenticated user
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=dto.UserWithProfileResponse}
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved", dto.ToUserWithProfileResponse(user))
}

// UpdateMe updates current user profile
// @Summary      Update My Profile
// @Description  Update the profile of the currently authenticated user
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.UpdateProfileRequest true "Profile Update Data"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	profile := &domain.UserProfile{
		FirstName:         req.FirstName,
		MiddleName:        req.MiddleName,
		LastName:          req.LastName,
		Phone:             req.Phone,
		Gender:            req.Gender,
		ProfilePictureURL: req.ProfilePictureURL,
		Bio:               req.Bio,
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), userID, profile); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated", nil)
}

// CreateUser creates a new user (admin only)
// @Summary      Create User
// @Description  Create a new user (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserRequest true "User Creation Data"
// @Success      201  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      409  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID", err)
		return
	}

	user := &domain.User{
		RegisterNo:   req.RegisterNo,
		Email:        req.Email,
		PasswordHash: req.Password, // Will be hashed by service
		RoleID:       roleID,
		Status:       "active",
	}

	profile := &domain.UserProfile{
		RegisterNo: req.RegisterNo,
		FirstName:  req.FirstName,
		MiddleName: &req.MiddleName,
		LastName:   req.LastName,
		Phone:      &req.Phone,
		Gender:     &req.Gender,
	}

	if err := h.userService.CreateUser(c.Request.Context(), user, profile); err != nil {
		if err == domain.ErrUserAlreadyExists {
			utils.ErrorResponse(c, http.StatusConflict, "User already exists", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", gin.H{"user_id": user.UserID})
}

// GetUser returns user by ID (admin only)
// @Summary      Get User by ID
// @Description  Get user details by ID (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  utils.APIResponse{data=dto.UserWithProfileResponse}
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved", dto.ToUserWithProfileResponse(user))
}

// ListUsers returns paginated list of users (admin only)
// @Summary      List Users
// @Description  Get a paginated list of users with filtering (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page number" default(1)
// @Param        limit    query     int     false  "Items per page" default(20)
// @Param        role_id  query     string  false  "Filter by Role ID"
// @Param        status   query     string  false  "Filter by Status"
// @Param        search   query     string  false  "Search term"
// @Success      200  {object}  utils.PaginatedResponse{data=[]dto.UserWithProfileResponse}
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filters := make(map[string]interface{})

	if roleID := c.Query("role_id"); roleID != "" {
		if id, err := uuid.Parse(roleID); err == nil {
			filters["role_id"] = id
		}
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), filters, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list users", err)
		return
	}

	userResponses := make([]dto.UserWithProfileResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.ToUserWithProfileResponse(user)
	}

	utils.PaginatedSuccessResponse(c, userResponses, page, limit, total)
}

// UpdateUser updates user by ID (admin only)
// @Summary      Update User
// @Description  Update user details by ID (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "User ID"
// @Param        request  body      dto.UpdateUserRequest  true  "Update Data"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.RoleID != "" {
		roleID, err := uuid.Parse(req.RoleID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID", err)
			return
		}
		updates["role_id"] = roleID
	}

	if err := h.userService.UpdateUser(c.Request.Context(), userID, updates); err != nil {
		if err == domain.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", nil)
}

// DeleteUser deletes user by ID (admin only)
// @Summary      Delete User
// @Description  Delete user by ID (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		if err == domain.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// ActivateUser activates a user (admin only)
// @Summary      Activate User
// @Description  Activate a suspended or inactive user (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.ActivateUser(c.Request.Context(), userID); err != nil {
		if err == domain.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to activate user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User activated successfully", nil)
}

// SuspendUser suspends a user (admin only)
// @Summary      Suspend User
// @Description  Suspend an active user (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users/{id}/suspend [post]
func (h *UserHandler) SuspendUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	if err := h.userService.SuspendUser(c.Request.Context(), userID); err != nil {
		if err == domain.ErrUserNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to suspend user", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User suspended successfully", nil)
}

// BulkImportUsers imports users in bulk (admin only)
// @Summary      Bulk Import Users
// @Description  Import multiple users at once (Admin only)
// @Tags         admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.BulkUserImportRequest true "Bulk Import Data"
// @Success      201  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      500  {object}  utils.APIResponse
// @Router       /admin/users/bulk-import [post]
func (h *UserHandler) BulkImportUsers(c *gin.Context) {
	var req dto.BulkUserImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	users := make([]*domain.User, len(req.Users))
	profiles := make([]*domain.UserProfile, len(req.Users))

	for i, userReq := range req.Users {
		roleID, err := uuid.Parse(userReq.RoleID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role ID", err)
			return
		}

		users[i] = &domain.User{
			RegisterNo:   userReq.RegisterNo,
			Email:        userReq.Email,
			PasswordHash: userReq.Password,
			RoleID:       roleID,
			Status:       "active",
		}

		profiles[i] = &domain.UserProfile{
			RegisterNo: userReq.RegisterNo,
			FirstName:  userReq.FirstName,
			MiddleName: &userReq.MiddleName,
			LastName:   userReq.LastName,
			Phone:      &userReq.Phone,
			Gender:     &userReq.Gender,
		}
	}

	if err := h.userService.BulkCreateUsers(c.Request.Context(), users, profiles); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to import users", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Users imported successfully", gin.H{"count": len(users)})
}
