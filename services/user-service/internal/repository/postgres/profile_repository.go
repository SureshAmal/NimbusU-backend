package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
)

type userProfileRepository struct {
	db *pgxpool.Pool
}

// NewUserProfileRepository creates a new user profile repository
func NewUserProfileRepository(db *pgxpool.Pool) domain.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(ctx context.Context, profile *domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, middle_name, 
		                           last_name, phone, gender, profile_picture_url, bio)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		profile.ProfileID,
		profile.UserID,
		profile.RegisterNo,
		profile.FirstName,
		profile.MiddleName,
		profile.LastName,
		profile.Phone,
		profile.Gender,
		profile.ProfilePictureURL,
		profile.Bio,
	).Scan(&profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	return nil
}

func (r *userProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	query := `
		SELECT profile_id, user_id, register_no, first_name, middle_name, last_name, 
		       phone, gender, profile_picture_url, bio, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	var profile domain.UserProfile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ProfileID,
		&profile.UserID,
		&profile.RegisterNo,
		&profile.FirstName,
		&profile.MiddleName,
		&profile.LastName,
		&profile.Phone,
		&profile.Gender,
		&profile.ProfilePictureURL,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return &profile, nil
}

func (r *userProfileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	query := `
		UPDATE user_profiles
		SET first_name = $1, middle_name = $2, last_name = $3, phone = $4, 
		    gender = $5, profile_picture_url = $6, bio = $7, updated_at = now()
		WHERE user_id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		profile.FirstName,
		profile.MiddleName,
		profile.LastName,
		profile.Phone,
		profile.Gender,
		profile.ProfilePictureURL,
		profile.Bio,
		profile.UserID,
	).Scan(&profile.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrProfileNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

func (r *userProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_profiles WHERE user_id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user profile: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProfileNotFound
	}

	return nil
}
