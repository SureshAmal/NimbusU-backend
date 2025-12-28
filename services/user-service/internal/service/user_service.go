package service

import (
	"context"
	"fmt"

	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/shared/models"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
	"github.com/google/uuid"
)

type userService struct {
	userRepo    domain.UserRepository
	profileRepo domain.UserProfileRepository
	roleRepo    domain.RoleRepository
	activityLog domain.ActivityLogRepository
	producer    domain.EventProducer
}

// NewUserService creates a new user service
func NewUserService(
	userRepo domain.UserRepository,
	profileRepo domain.UserProfileRepository,
	roleRepo domain.RoleRepository,
	activityLog domain.ActivityLogRepository,
	producer domain.EventProducer,
) domain.UserService {
	return &userService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		roleRepo:    roleRepo,
		activityLog: activityLog,
		producer:    producer,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User, profile *domain.UserProfile) error {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, user.Email)
	if existingUser != nil {
		return domain.ErrUserAlreadyExists
	}

	existingUser, _ = s.userRepo.GetByRegisterNo(ctx, user.RegisterNo)
	if existingUser != nil {
		return domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashedPassword

	// Generate IDs
	user.UserID = uuid.New()
	profile.ProfileID = uuid.New()
	profile.UserID = user.UserID

	// Set defaults
	if user.Status == "" {
		user.Status = "active"
	}

	// Create user and profile
	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	if err := s.profileRepo.Create(ctx, profile); err != nil {
		// Rollback user creation if profile fails
		s.userRepo.Delete(ctx, user.UserID)
		return err
	}

	// Publish user created event
	event := models.NewUserEvent(models.EventUserCreated, user.UserID, user.Email)
	event.FirstName = profile.FirstName
	event.LastName = profile.LastName
	event.RoleID = user.RoleID
	event.Status = user.Status

	s.producer.PublishEvent("user.events", user.UserID.String(), event)

	return nil
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserWithProfile, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.GetByID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	return &domain.UserWithProfile{
		User:        *user,
		UserProfile: *profile,
		Role:        role,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Apply updates
	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if roleID, ok := updates["role_id"].(uuid.UUID); ok {
		user.RoleID = roleID
	}
	if status, ok := updates["status"].(string); ok {
		user.Status = status
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Publish user updated event
	event := models.NewUserEvent(models.EventUserUpdated, user.UserID, user.Email)
	event.RoleID = user.RoleID
	event.Status = user.Status
	s.producer.PublishEvent("user.events", user.UserID.String(), event)

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Delete profile first (foreign key constraint)
	if err := s.profileRepo.Delete(ctx, userID); err != nil {
		return err
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	// Publish user deleted event
	event := models.NewUserEvent(models.EventUserDeleted, userID, user.Email)
	s.producer.PublishEvent("user.events", userID.String(), event)

	return nil
}

func (s *userService) ListUsers(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*domain.UserWithProfile, int64, error) {
	offset := (page - 1) * limit

	users, total, err := s.userRepo.List(ctx, filters, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	usersWithProfiles := make([]*domain.UserWithProfile, 0, len(users))
	for _, user := range users {
		profile, err := s.profileRepo.GetByUserID(ctx, user.UserID)
		if err != nil {
			continue
		}

		role, err := s.roleRepo.GetByID(ctx, user.RoleID)
		if err != nil {
			continue
		}

		usersWithProfiles = append(usersWithProfiles, &domain.UserWithProfile{
			User:        *user,
			UserProfile: *profile,
			Role:        role,
		})
	}

	return usersWithProfiles, total, nil
}

func (s *userService) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdateStatus(ctx, userID, "active"); err != nil {
		return err
	}

	// Publish user activated event
	event := models.NewUserEvent(models.EventUserActivated, userID, user.Email)
	event.Status = "active"
	s.producer.PublishEvent("user.events", userID.String(), event)

	return nil
}

func (s *userService) SuspendUser(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdateStatus(ctx, userID, "suspended"); err != nil {
		return err
	}

	// Publish user suspended event
	event := models.NewUserEvent(models.EventUserSuspended, userID, user.Email)
	event.Status = "suspended"
	s.producer.PublishEvent("user.events", userID.String(), event)

	return nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, profile *domain.UserProfile) error {
	existingProfile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Update fields
	existingProfile.FirstName = profile.FirstName
	existingProfile.MiddleName = profile.MiddleName
	existingProfile.LastName = profile.LastName
	existingProfile.Phone = profile.Phone
	existingProfile.Gender = profile.Gender
	existingProfile.ProfilePictureURL = profile.ProfilePictureURL
	existingProfile.Bio = profile.Bio

	return s.profileRepo.Update(ctx, existingProfile)
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}

func (s *userService) BulkCreateUsers(ctx context.Context, users []*domain.User, profiles []*domain.UserProfile) error {
	if len(users) != len(profiles) {
		return fmt.Errorf("users and profiles count mismatch")
	}

	for i, user := range users {
		if err := s.CreateUser(ctx, user, profiles[i]); err != nil {
			return fmt.Errorf("failed to create user %s: %w", user.Email, err)
		}
	}

	return nil
}
