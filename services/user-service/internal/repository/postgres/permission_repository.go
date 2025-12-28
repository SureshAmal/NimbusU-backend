package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
)

type permissionRepository struct {
	db *pgxpool.Pool
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *pgxpool.Pool) domain.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	query := `
		INSERT INTO permissions (permission_id, permission_name, resource, action, description)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		permission.PermissionID,
		permission.PermissionName,
		permission.Resource,
		permission.Action,
		permission.Description,
	)

	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}

func (r *permissionRepository) GetByID(ctx context.Context, permissionID uuid.UUID) (*domain.Permission, error) {
	query := `
		SELECT permission_id, permission_name, resource, action, description
		FROM permissions
		WHERE permission_id = $1
	`

	var permission domain.Permission
	err := r.db.QueryRow(ctx, query, permissionID).Scan(
		&permission.PermissionID,
		&permission.PermissionName,
		&permission.Resource,
		&permission.Action,
		&permission.Description,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrPermissionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &permission, nil
}

func (r *permissionRepository) List(ctx context.Context) ([]*domain.Permission, error) {
	query := `
		SELECT permission_id, permission_name, resource, action, description
		FROM permissions
		ORDER BY resource, action
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	permissions := []*domain.Permission{}
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(
			&permission.PermissionID,
			&permission.PermissionName,
			&permission.Resource,
			&permission.Action,
			&permission.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

func (r *permissionRepository) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]*domain.Permission, error) {
	query := `
		SELECT p.permission_id, p.permission_name, p.resource, p.action, p.description
		FROM permissions p
		INNER JOIN role_permissions rp ON p.permission_id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action
	`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by role: %w", err)
	}
	defer rows.Close()

	permissions := []*domain.Permission{}
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(
			&permission.PermissionID,
			&permission.PermissionName,
			&permission.Resource,
			&permission.Action,
			&permission.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

type rolePermissionRepository struct {
	db *pgxpool.Pool
}

// NewRolePermissionRepository creates a new role-permission repository
func NewRolePermissionRepository(db *pgxpool.Pool) domain.RolePermissionRepository {
	return &rolePermissionRepository{db: db}
}

func (r *rolePermissionRepository) AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `
		INSERT INTO role_permissions (role_permission_id, role_id, permission_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, uuid.New(), roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}

	return nil
}

func (r *rolePermissionRepository) RevokePermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	result, err := r.db.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPermissionNotFound
	}

	return nil
}

func (r *rolePermissionRepository) GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*domain.Permission, error) {
	query := `
		SELECT p.permission_id, p.permission_name, p.resource, p.action, p.description
		FROM permissions p
		INNER JOIN role_permissions rp ON p.permission_id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action
	`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer rows.Close()

	permissions := []*domain.Permission{}
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(
			&permission.PermissionID,
			&permission.PermissionName,
			&permission.Resource,
			&permission.Action,
			&permission.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}
