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

type departmentRepository struct {
	db *pgxpool.Pool
}

func NewDepartmentRepository(db *pgxpool.Pool) domain.DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, department *domain.Department) error {
	query := `
		INSERT INTO departments (department_id, department_name, department_code, head_of_department, description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`
	department.DepartmentID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		department.DepartmentID,
		department.DepartmentName,
		department.DepartmentCode,
		department.HeadOfDepartment,
		department.Description,
		department.IsActive,
	).Scan(&department.CreatedAt, &department.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrDepartmentCodeExists
		}
		return fmt.Errorf("failed to create department: %w", err)
	}
	return nil
}

func (r *departmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Department, error) {
	query := `
		SELECT department_id, department_name, department_code, head_of_department,
			   description, is_active, created_at, updated_at
		FROM departments
		WHERE department_id = $1
	`
	var d domain.Department
	err := r.db.QueryRow(ctx, query, id).Scan(
		&d.DepartmentID, &d.DepartmentName, &d.DepartmentCode, &d.HeadOfDepartment,
		&d.Description, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrDepartmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	return &d, nil
}

func (r *departmentRepository) GetByCode(ctx context.Context, code string) (*domain.Department, error) {
	query := `
		SELECT department_id, department_name, department_code, head_of_department,
			   description, is_active, created_at, updated_at
		FROM departments
		WHERE department_code = $1
	`
	var d domain.Department
	err := r.db.QueryRow(ctx, query, code).Scan(
		&d.DepartmentID, &d.DepartmentName, &d.DepartmentCode, &d.HeadOfDepartment,
		&d.Description, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrDepartmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get department by code: %w", err)
	}
	return &d, nil
}

func (r *departmentRepository) Update(ctx context.Context, department *domain.Department) error {
	query := `
		UPDATE departments
		SET department_name = $2, head_of_department = $3, description = $4, is_active = $5, updated_at = now()
		WHERE department_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		department.DepartmentID,
		department.DepartmentName,
		department.HeadOfDepartment,
		department.Description,
		department.IsActive,
	).Scan(&department.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrDepartmentNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}
	return nil
}

func (r *departmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE departments SET is_active = false, updated_at = now() WHERE department_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrDepartmentNotFound
	}
	return nil
}

func (r *departmentRepository) List(ctx context.Context, filter domain.DepartmentFilter, limit, offset int) ([]*domain.Department, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM departments %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count departments: %w", err)
	}

	// List query
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT department_id, department_name, department_code, head_of_department,
			   description, is_active, created_at, updated_at
		FROM departments
		%s
		ORDER BY department_name
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list departments: %w", err)
	}
	defer rows.Close()

	var departments []*domain.Department
	for rows.Next() {
		var d domain.Department
		if err := rows.Scan(
			&d.DepartmentID, &d.DepartmentName, &d.DepartmentCode, &d.HeadOfDepartment,
			&d.Description, &d.IsActive, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan department: %w", err)
		}
		departments = append(departments, &d)
	}

	return departments, total, nil
}

func (r *departmentRepository) GetWithDetails(ctx context.Context, id uuid.UUID) (*domain.DepartmentWithDetails, error) {
	query := `
		SELECT d.department_id, d.department_name, d.department_code, d.head_of_department,
			   d.description, d.is_active, d.created_at, d.updated_at,
			   f.faculty_id, f.employee_id, f.designation,
			   (SELECT COUNT(*) FROM programs p WHERE p.department_id = d.department_id AND p.is_active = true) as programs_count,
			   (SELECT COUNT(*) FROM faculties fa WHERE fa.department_id = d.department_id AND fa.is_active = true) as faculty_count,
			   (SELECT COUNT(*) FROM students s WHERE s.department_id = d.department_id AND s.is_active = true) as students_count
		FROM departments d
		LEFT JOIN faculties f ON d.head_of_department = f.faculty_id
		WHERE d.department_id = $1
	`

	var dd domain.DepartmentWithDetails
	var facultyID, employeeID, designation *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&dd.DepartmentID, &dd.DepartmentName, &dd.DepartmentCode, &dd.HeadOfDepartment,
		&dd.Description, &dd.IsActive, &dd.CreatedAt, &dd.UpdatedAt,
		&facultyID, &employeeID, &designation,
		&dd.ProgramsCount, &dd.FacultyCount, &dd.StudentsCount,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrDepartmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get department with details: %w", err)
	}

	// Set head faculty if exists
	if facultyID != nil {
		fid, _ := uuid.Parse(*facultyID)
		dd.HeadFaculty = &domain.FacultyBasic{
			FacultyID:   fid,
			EmployeeID:  *employeeID,
			Designation: designation,
			Name:        "", // Would need to join with user service for name
		}
	}

	return &dd, nil
}
