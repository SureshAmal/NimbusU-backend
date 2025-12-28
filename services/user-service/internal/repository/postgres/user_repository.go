package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
)

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.UserID,
		user.RegisterNo,
		user.Email,
		user.PasswordHash,
		user.RoleID,
		user.Status,
		user.CreatedBy,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	query := `
		SELECT user_id, register_no, email, password_hash, role_id, status, 
		       last_login, created_at, updated_at, created_by
		FROM users
		WHERE user_id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.RegisterNo,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Status,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT user_id, register_no, email, password_hash, role_id, status, 
		       last_login, created_at, updated_at, created_by
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.RegisterNo,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Status,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByRegisterNo(ctx context.Context, registerNo int64) (*domain.User, error) {
	query := `
		SELECT user_id, register_no, email, password_hash, role_id, status, 
		       last_login, created_at, updated_at, created_by
		FROM users
		WHERE register_no = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, registerNo).Scan(
		&user.UserID,
		&user.RegisterNo,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Status,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by register_no: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, role_id = $3, status = $4, updated_at = now()
		WHERE user_id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.RoleID,
		user.Status,
		user.UserID,
	).Scan(&user.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.User, int64, error) {
	// Build dynamic query based on filters
	whereClause := []string{}
	args := []interface{}{}
	argPos := 1

	if roleID, ok := filters["role_id"].(uuid.UUID); ok {
		whereClause = append(whereClause, fmt.Sprintf("role_id = $%d", argPos))
		args = append(args, roleID)
		argPos++
	}

	if status, ok := filters["status"].(string); ok {
		whereClause = append(whereClause, fmt.Sprintf("status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	if search, ok := filters["search"].(string); ok {
		whereClause = append(whereClause, fmt.Sprintf("(email ILIKE $%d OR CAST(register_no AS TEXT) LIKE $%d)", argPos, argPos))
		args = append(args, "%"+search+"%")
		argPos++
	}

	where := ""
	if len(whereClause) > 0 {
		where = "WHERE " + strings.Join(whereClause, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", where)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT user_id, register_no, email, password_hash, role_id, status, 
		       last_login, created_at, updated_at, created_by
		FROM users
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argPos, argPos+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.UserID,
			&user.RegisterNo,
			&user.Email,
			&user.PasswordHash,
			&user.RoleID,
			&user.Status,
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	return users, total, nil
}

func (r *userRepository) UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error {
	query := `UPDATE users SET status = $1, updated_at = now() WHERE user_id = $2`

	result, err := r.db.Exec(ctx, query, status, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login = $1, updated_at = now() WHERE user_id = $2`

	result, err := r.db.Exec(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
