package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
)

type roleRepository struct {
	db *pgxpool.Pool
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *pgxpool.Pool) domain.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (role_id, role_name, description)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query,
		role.RoleID,
		role.RoleName,
		role.Description,
	).Scan(&role.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, roleID uuid.UUID) (*domain.Role, error) {
	query := `
		SELECT role_id, role_name, description, created_at
		FROM roles
		WHERE role_id = $1
	`

	var role domain.Role
	err := r.db.QueryRow(ctx, query, roleID).Scan(
		&role.RoleID,
		&role.RoleName,
		&role.Description,
		&role.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, roleName string) (*domain.Role, error) {
	query := `
		SELECT role_id, role_name, description, created_at
		FROM roles
		WHERE role_name = $1
	`

	var role domain.Role
	err := r.db.QueryRow(ctx, query, roleName).Scan(
		&role.RoleID,
		&role.RoleName,
		&role.Description,
		&role.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]*domain.Role, error) {
	query := `
		SELECT role_id, role_name, description, created_at
		FROM roles
		ORDER BY role_name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	roles := []*domain.Role{}
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.RoleID,
			&role.RoleName,
			&role.Description,
			&role.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `
		UPDATE roles
		SET role_name = $1, description = $2
		WHERE role_id = $3
	`

	result, err := r.db.Exec(ctx, query,
		role.RoleName,
		role.Description,
		role.RoleID,
	)

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRoleNotFound
	}

	return nil
}

func (r *roleRepository) Delete(ctx context.Context, roleID uuid.UUID) error {
	query := `DELETE FROM roles WHERE role_id = $1`

	result, err := r.db.Exec(ctx, query, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRoleNotFound
	}

	return nil
}
