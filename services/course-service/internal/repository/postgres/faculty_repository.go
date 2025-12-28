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

type facultyRepository struct {
	db *pgxpool.Pool
}

func NewFacultyRepository(db *pgxpool.Pool) domain.FacultyRepository {
	return &facultyRepository{db: db}
}

func (r *facultyRepository) Create(ctx context.Context, faculty *domain.Faculty) error {
	query := `
		INSERT INTO faculties (faculty_id, user_id, employee_id, department_id, designation, qualification,
			specialization, joining_date, office_room, office_hours, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at, updated_at
	`
	faculty.FacultyID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		faculty.FacultyID,
		faculty.UserID,
		faculty.EmployeeID,
		faculty.DepartmentID,
		faculty.Designation,
		faculty.Qualification,
		faculty.Specialization,
		faculty.JoiningDate,
		faculty.OfficeRoom,
		faculty.OfficeHours,
		faculty.IsActive,
	).Scan(&faculty.CreatedAt, &faculty.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "employee_id") {
				return domain.ErrEmployeeIDExists
			}
		}
		return fmt.Errorf("failed to create faculty: %w", err)
	}
	return nil
}

func (r *facultyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Faculty, error) {
	query := `
		SELECT faculty_id, user_id, employee_id, department_id, designation, qualification,
			   specialization, joining_date, office_room, office_hours, is_active, created_at, updated_at
		FROM faculties
		WHERE faculty_id = $1
	`
	var f domain.Faculty
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.FacultyID, &f.UserID, &f.EmployeeID, &f.DepartmentID, &f.Designation, &f.Qualification,
		&f.Specialization, &f.JoiningDate, &f.OfficeRoom, &f.OfficeHours, &f.IsActive, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrFacultyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty: %w", err)
	}
	return &f, nil
}

func (r *facultyRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Faculty, error) {
	query := `
		SELECT faculty_id, user_id, employee_id, department_id, designation, qualification,
			   specialization, joining_date, office_room, office_hours, is_active, created_at, updated_at
		FROM faculties
		WHERE user_id = $1
	`
	var f domain.Faculty
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&f.FacultyID, &f.UserID, &f.EmployeeID, &f.DepartmentID, &f.Designation, &f.Qualification,
		&f.Specialization, &f.JoiningDate, &f.OfficeRoom, &f.OfficeHours, &f.IsActive, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrFacultyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty by user ID: %w", err)
	}
	return &f, nil
}

func (r *facultyRepository) GetByEmployeeID(ctx context.Context, employeeID string) (*domain.Faculty, error) {
	query := `
		SELECT faculty_id, user_id, employee_id, department_id, designation, qualification,
			   specialization, joining_date, office_room, office_hours, is_active, created_at, updated_at
		FROM faculties
		WHERE employee_id = $1
	`
	var f domain.Faculty
	err := r.db.QueryRow(ctx, query, employeeID).Scan(
		&f.FacultyID, &f.UserID, &f.EmployeeID, &f.DepartmentID, &f.Designation, &f.Qualification,
		&f.Specialization, &f.JoiningDate, &f.OfficeRoom, &f.OfficeHours, &f.IsActive, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrFacultyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty by employee ID: %w", err)
	}
	return &f, nil
}

func (r *facultyRepository) Update(ctx context.Context, faculty *domain.Faculty) error {
	query := `
		UPDATE faculties
		SET designation = $2, qualification = $3, specialization = $4, joining_date = $5,
			office_room = $6, office_hours = $7, is_active = $8, updated_at = now()
		WHERE faculty_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		faculty.FacultyID,
		faculty.Designation,
		faculty.Qualification,
		faculty.Specialization,
		faculty.JoiningDate,
		faculty.OfficeRoom,
		faculty.OfficeHours,
		faculty.IsActive,
	).Scan(&faculty.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrFacultyNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update faculty: %w", err)
	}
	return nil
}

func (r *facultyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE faculties SET is_active = false, updated_at = now() WHERE faculty_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete faculty: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrFacultyNotFound
	}
	return nil
}

func (r *facultyRepository) List(ctx context.Context, filter domain.FacultyFilter, limit, offset int) ([]*domain.FacultyWithDetails, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.DepartmentID != nil {
		conditions = append(conditions, fmt.Sprintf("f.department_id = $%d", argNum))
		args = append(args, *filter.DepartmentID)
		argNum++
	}
	if filter.Designation != nil {
		conditions = append(conditions, fmt.Sprintf("f.designation = $%d", argNum))
		args = append(args, *filter.Designation)
		argNum++
	}
	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf("(f.employee_id ILIKE $%d)", argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}
	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("f.is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM faculties f %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count faculty: %w", err)
	}

	// List query with department join
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT f.faculty_id, f.user_id, f.employee_id, f.department_id, f.designation, f.qualification,
			   f.specialization, f.joining_date, f.office_room, f.office_hours, f.is_active, f.created_at, f.updated_at,
			   d.department_id, d.department_name, d.department_code
		FROM faculties f
		JOIN departments d ON f.department_id = d.department_id
		%s
		ORDER BY f.employee_id
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list faculty: %w", err)
	}
	defer rows.Close()

	var faculties []*domain.FacultyWithDetails
	for rows.Next() {
		var f domain.FacultyWithDetails
		if err := rows.Scan(
			&f.FacultyID, &f.UserID, &f.EmployeeID, &f.DepartmentID, &f.Designation, &f.Qualification,
			&f.Specialization, &f.JoiningDate, &f.OfficeRoom, &f.OfficeHours, &f.IsActive, &f.CreatedAt, &f.UpdatedAt,
			&f.Department.DepartmentID, &f.Department.DepartmentName, &f.Department.DepartmentCode,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan faculty: %w", err)
		}
		// Note: Name and Email would come from user service
		faculties = append(faculties, &f)
	}

	return faculties, total, nil
}

func (r *facultyRepository) GetWithDetails(ctx context.Context, id uuid.UUID) (*domain.FacultyWithDetails, error) {
	query := `
		SELECT f.faculty_id, f.user_id, f.employee_id, f.department_id, f.designation, f.qualification,
			   f.specialization, f.joining_date, f.office_room, f.office_hours, f.is_active, f.created_at, f.updated_at,
			   d.department_id, d.department_name, d.department_code
		FROM faculties f
		JOIN departments d ON f.department_id = d.department_id
		WHERE f.faculty_id = $1
	`

	var f domain.FacultyWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&f.FacultyID, &f.UserID, &f.EmployeeID, &f.DepartmentID, &f.Designation, &f.Qualification,
		&f.Specialization, &f.JoiningDate, &f.OfficeRoom, &f.OfficeHours, &f.IsActive, &f.CreatedAt, &f.UpdatedAt,
		&f.Department.DepartmentID, &f.Department.DepartmentName, &f.Department.DepartmentCode,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrFacultyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty with details: %w", err)
	}

	// Get current courses
	courseQuery := `
		SELECT c.course_id, c.course_code, c.course_name, s.credits, fc.role
		FROM faculty_courses fc
		JOIN courses c ON fc.course_id = c.course_id
		JOIN subjects s ON c.subject_id = s.subject_id
		JOIN semesters sem ON c.semester_id = sem.semester_id
		WHERE fc.faculty_id = $1 AND fc.is_active = true AND sem.is_current = true
	`
	rows, err := r.db.Query(ctx, courseQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty courses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.CourseBasic
		if err := rows.Scan(&c.CourseID, &c.CourseCode, &c.CourseName, &c.Credits, &c.Role); err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		f.CurrentCourses = append(f.CurrentCourses, c)
	}

	return &f, nil
}
