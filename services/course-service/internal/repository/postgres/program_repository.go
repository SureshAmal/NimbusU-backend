package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type programRepository struct {
	db *pgxpool.Pool
}

func NewProgramRepository(db *pgxpool.Pool) domain.ProgramRepository {
	return &programRepository{db: db}
}

func (r *programRepository) Create(ctx context.Context, program *domain.Program) error {
	query := `
		INSERT INTO programs (program_id, program_name, program_code, department_id, degree_type, duration_years, total_credits, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`
	program.ProgramID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		program.ProgramID,
		program.ProgramName,
		program.ProgramCode,
		program.DepartmentID,
		program.DegreeType,
		program.DurationYears,
		program.TotalCredits,
		program.Description,
		program.IsActive,
	).Scan(&program.CreatedAt, &program.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrProgramCodeExists
		}
		return fmt.Errorf("failed to create program: %w", err)
	}
	return nil
}

func (r *programRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	query := `
		SELECT program_id, program_name, program_code, department_id, degree_type,
			   duration_years, total_credits, description, is_active, created_at, updated_at
		FROM programs
		WHERE program_id = $1
	`
	var p domain.Program
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ProgramID, &p.ProgramName, &p.ProgramCode, &p.DepartmentID, &p.DegreeType,
		&p.DurationYears, &p.TotalCredits, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrProgramNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}
	return &p, nil
}

func (r *programRepository) GetByCode(ctx context.Context, code string) (*domain.Program, error) {
	query := `
		SELECT program_id, program_name, program_code, department_id, degree_type,
			   duration_years, total_credits, description, is_active, created_at, updated_at
		FROM programs
		WHERE program_code = $1
	`
	var p domain.Program
	err := r.db.QueryRow(ctx, query, code).Scan(
		&p.ProgramID, &p.ProgramName, &p.ProgramCode, &p.DepartmentID, &p.DegreeType,
		&p.DurationYears, &p.TotalCredits, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrProgramNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get program by code: %w", err)
	}
	return &p, nil
}

func (r *programRepository) Update(ctx context.Context, program *domain.Program) error {
	query := `
		UPDATE programs
		SET program_name = $2, degree_type = $3, duration_years = $4, total_credits = $5, description = $6, is_active = $7, updated_at = now()
		WHERE program_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		program.ProgramID,
		program.ProgramName,
		program.DegreeType,
		program.DurationYears,
		program.TotalCredits,
		program.Description,
		program.IsActive,
	).Scan(&program.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrProgramNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update program: %w", err)
	}
	return nil
}

func (r *programRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE programs SET is_active = false, updated_at = now() WHERE program_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete program: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrProgramNotFound
	}
	return nil
}

func (r *programRepository) List(ctx context.Context, filter domain.ProgramFilter, limit, offset int) ([]*domain.ProgramWithDepartment, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.DepartmentID != nil {
		conditions = append(conditions, fmt.Sprintf("p.department_id = $%d", argNum))
		args = append(args, *filter.DepartmentID)
		argNum++
	}
	if filter.DegreeType != nil {
		conditions = append(conditions, fmt.Sprintf("p.degree_type = $%d", argNum))
		args = append(args, *filter.DegreeType)
		argNum++
	}
	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("p.is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM programs p %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count programs: %w", err)
	}

	// List query with department join
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT p.program_id, p.program_name, p.program_code, p.department_id, p.degree_type,
			   p.duration_years, p.total_credits, p.description, p.is_active, p.created_at, p.updated_at,
			   d.department_id, d.department_name, d.department_code
		FROM programs p
		JOIN departments d ON p.department_id = d.department_id
		%s
		ORDER BY p.program_name
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list programs: %w", err)
	}
	defer rows.Close()

	var programs []*domain.ProgramWithDepartment
	for rows.Next() {
		var p domain.ProgramWithDepartment
		if err := rows.Scan(
			&p.ProgramID, &p.ProgramName, &p.ProgramCode, &p.DepartmentID, &p.DegreeType,
			&p.DurationYears, &p.TotalCredits, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
			&p.Department.DepartmentID, &p.Department.DepartmentName, &p.Department.DepartmentCode,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan program: %w", err)
		}
		programs = append(programs, &p)
	}

	return programs, total, nil
}

func (r *programRepository) GetWithDepartment(ctx context.Context, id uuid.UUID) (*domain.ProgramWithDepartment, error) {
	query := `
		SELECT p.program_id, p.program_name, p.program_code, p.department_id, p.degree_type,
			   p.duration_years, p.total_credits, p.description, p.is_active, p.created_at, p.updated_at,
			   d.department_id, d.department_name, d.department_code
		FROM programs p
		JOIN departments d ON p.department_id = d.department_id
		WHERE p.program_id = $1
	`

	var p domain.ProgramWithDepartment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ProgramID, &p.ProgramName, &p.ProgramCode, &p.DepartmentID, &p.DegreeType,
		&p.DurationYears, &p.TotalCredits, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		&p.Department.DepartmentID, &p.Department.DepartmentName, &p.Department.DepartmentCode,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrProgramNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get program with department: %w", err)
	}

	return &p, nil
}
