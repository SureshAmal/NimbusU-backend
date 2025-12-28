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

type facultyCourseRepository struct {
	db *pgxpool.Pool
}

func NewFacultyCourseRepository(db *pgxpool.Pool) domain.FacultyCourseRepository {
	return &facultyCourseRepository{db: db}
}

func (r *facultyCourseRepository) Create(ctx context.Context, fc *domain.FacultyCourse) error {
	query := `
		INSERT INTO faculty_courses (faculty_course_id, faculty_id, course_id, role, is_primary, assigned_by, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING assigned_at
	`
	fc.FacultyCourseID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		fc.FacultyCourseID,
		fc.FacultyID,
		fc.CourseID,
		fc.Role,
		fc.IsPrimary,
		fc.AssignedBy,
		fc.IsActive,
	).Scan(&fc.AssignedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrFacultyAlreadyAssigned
		}
		return fmt.Errorf("failed to create faculty course assignment: %w", err)
	}
	return nil
}

func (r *facultyCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.FacultyCourse, error) {
	query := `
		SELECT faculty_course_id, faculty_id, course_id, role, is_primary, assigned_by, assigned_at, removed_at, is_active
		FROM faculty_courses
		WHERE faculty_course_id = $1
	`
	var fc domain.FacultyCourse
	err := r.db.QueryRow(ctx, query, id).Scan(
		&fc.FacultyCourseID, &fc.FacultyID, &fc.CourseID, &fc.Role, &fc.IsPrimary,
		&fc.AssignedBy, &fc.AssignedAt, &fc.RemovedAt, &fc.IsActive,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrAssignmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty course assignment: %w", err)
	}
	return &fc, nil
}

func (r *facultyCourseRepository) Update(ctx context.Context, fc *domain.FacultyCourse) error {
	query := `
		UPDATE faculty_courses
		SET role = $3, is_primary = $4, is_active = $5
		WHERE faculty_id = $1 AND course_id = $2 AND is_active = true
	`
	result, err := r.db.Exec(ctx, query,
		fc.FacultyID,
		fc.CourseID,
		fc.Role,
		fc.IsPrimary,
		fc.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update faculty course assignment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrAssignmentNotFound
	}
	return nil
}

func (r *facultyCourseRepository) Delete(ctx context.Context, facultyID, courseID uuid.UUID) error {
	query := `UPDATE faculty_courses SET is_active = false, removed_at = now() WHERE faculty_id = $1 AND course_id = $2 AND is_active = true`
	result, err := r.db.Exec(ctx, query, facultyID, courseID)
	if err != nil {
		return fmt.Errorf("failed to delete faculty course assignment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrAssignmentNotFound
	}
	return nil
}

func (r *facultyCourseRepository) ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*domain.FacultyCourseWithDetails, error) {
	query := `
		SELECT fc.faculty_course_id, fc.faculty_id, fc.course_id, fc.role, fc.is_primary,
			   fc.assigned_by, fc.assigned_at, fc.removed_at, fc.is_active,
			   f.faculty_id, f.user_id, f.employee_id, f.designation
		FROM faculty_courses fc
		JOIN faculties f ON fc.faculty_id = f.faculty_id
		WHERE fc.course_id = $1 AND fc.is_active = true
		ORDER BY fc.is_primary DESC, fc.assigned_at
	`

	rows, err := r.db.Query(ctx, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list faculty by course: %w", err)
	}
	defer rows.Close()

	var assignments []*domain.FacultyCourseWithDetails
	for rows.Next() {
		var a domain.FacultyCourseWithDetails
		if err := rows.Scan(
			&a.FacultyCourseID, &a.FacultyID, &a.CourseID, &a.Role, &a.IsPrimary,
			&a.AssignedBy, &a.AssignedAt, &a.RemovedAt, &a.IsActive,
			&a.Faculty.FacultyID, &a.Faculty.UserID, &a.Faculty.EmployeeID, &a.Faculty.Designation,
		); err != nil {
			return nil, fmt.Errorf("failed to scan faculty assignment: %w", err)
		}
		assignments = append(assignments, &a)
	}

	return assignments, nil
}

func (r *facultyCourseRepository) ListByFaculty(ctx context.Context, facultyID uuid.UUID, semesterID *uuid.UUID) ([]*domain.FacultyCourse, error) {
	var args []interface{}
	args = append(args, facultyID)

	query := `
		SELECT fc.faculty_course_id, fc.faculty_id, fc.course_id, fc.role, fc.is_primary,
			   fc.assigned_by, fc.assigned_at, fc.removed_at, fc.is_active
		FROM faculty_courses fc
		JOIN courses c ON fc.course_id = c.course_id
		WHERE fc.faculty_id = $1 AND fc.is_active = true
	`

	if semesterID != nil {
		query += " AND c.semester_id = $2"
		args = append(args, *semesterID)
	}

	query += " ORDER BY fc.assigned_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list courses by faculty: %w", err)
	}
	defer rows.Close()

	var assignments []*domain.FacultyCourse
	for rows.Next() {
		var fc domain.FacultyCourse
		if err := rows.Scan(
			&fc.FacultyCourseID, &fc.FacultyID, &fc.CourseID, &fc.Role, &fc.IsPrimary,
			&fc.AssignedBy, &fc.AssignedAt, &fc.RemovedAt, &fc.IsActive,
		); err != nil {
			return nil, fmt.Errorf("failed to scan faculty course: %w", err)
		}
		assignments = append(assignments, &fc)
	}

	return assignments, nil
}

func (r *facultyCourseRepository) GetAssignment(ctx context.Context, facultyID, courseID uuid.UUID) (*domain.FacultyCourse, error) {
	query := `
		SELECT faculty_course_id, faculty_id, course_id, role, is_primary, assigned_by, assigned_at, removed_at, is_active
		FROM faculty_courses
		WHERE faculty_id = $1 AND course_id = $2 AND is_active = true
	`
	var fc domain.FacultyCourse
	err := r.db.QueryRow(ctx, query, facultyID, courseID).Scan(
		&fc.FacultyCourseID, &fc.FacultyID, &fc.CourseID, &fc.Role, &fc.IsPrimary,
		&fc.AssignedBy, &fc.AssignedAt, &fc.RemovedAt, &fc.IsActive,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrAssignmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get faculty assignment: %w", err)
	}
	return &fc, nil
}

func (r *facultyCourseRepository) IsPrimaryFaculty(ctx context.Context, facultyID, courseID uuid.UUID) (bool, error) {
	query := `
		SELECT is_primary FROM faculty_courses
		WHERE faculty_id = $1 AND course_id = $2 AND is_active = true
	`
	var isPrimary bool
	err := r.db.QueryRow(ctx, query, facultyID, courseID).Scan(&isPrimary)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check primary faculty: %w", err)
	}
	return isPrimary, nil
}
