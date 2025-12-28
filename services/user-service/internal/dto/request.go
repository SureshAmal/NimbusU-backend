package dto

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	RegisterNo int64  `json:"register_no" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	FirstName  string `json:"first_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" binding:"required"`
	Phone      string `json:"phone"`
	Gender     string `json:"gender"`
	RoleName   string `json:"role_name" binding:"required,oneof=admin faculty student staff"`
}

// CreateUserRequest represents admin user creation
type CreateUserRequest struct {
	RegisterNo int64  `json:"register_no" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	FirstName  string `json:"first_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name" binding:"required"`
	Phone      string `json:"phone"`
	Gender     string `json:"gender"`
	RoleID     string `json:"role_id" binding:"required,uuid"`
}

// UpdateUserRequest represents user update data
type UpdateUserRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
	RoleID string `json:"role_id" binding:"omitempty,uuid"`
}

// UpdateProfileRequest represents profile update data
type UpdateProfileRequest struct {
	FirstName         string  `json:"first_name"`
	MiddleName        *string `json:"middle_name"`
	LastName          string  `json:"last_name"`
	Phone             *string `json:"phone"`
	Gender            *string `json:"gender"`
	ProfilePictureURL *string `json:"profile_picture_url"`
	Bio               *string `json:"bio"`
}

// ChangePasswordRequest represents password change data
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// PasswordResetRequestRequest represents password reset request
type PasswordResetRequestRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents password reset with token
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateRoleRequest represents role creation data
type CreateRoleRequest struct {
	RoleName    string  `json:"role_name" binding:"required"`
	Description *string `json:"description"`
}

// BulkUserImportRequest represents bulk user import
type BulkUserImportRequest struct {
	Users []CreateUserRequest `json:"users" binding:"required,min=1"`
}

// AssignPermissionRequest represents permission assignment
type AssignPermissionRequest struct {
	PermissionID string `json:"permission_id" binding:"required,uuid"`
}
