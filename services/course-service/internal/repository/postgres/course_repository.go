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

type courseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) domain.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(ctx context.Context, course *domain.Course) error {
	query := `
		INSERT INTO courses (course_id, course_code, course_name, subject_id, department_id, program_id,
			semester_id, semester_number, academic_year, max_students, current_enrollment, status,
			description, is_active, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING created_at, updated_at
	`
	course.CourseID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		course.CourseID,
		course.CourseCode,
		course.CourseName,
		course.SubjectID,
		course.DepartmentID,
		course.ProgramID,
		course.SemesterID,
		course.SemesterNumber,
		course.AcademicYear,
		course.MaxStudents,
		course.CurrentEnrollment,
		course.Status,
		course.Description,
		course.IsActive,
		course.CreatedBy,
	).Scan(&course.CreatedAt, &course.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrCourseCodeExists
		}
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
}

func (r *courseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Course, error) {
	query := `
		SELECT course_id, course_code, course_name, subject_id, department_id, program_id,
			   semester_id, semester_number, academic_year, max_students, current_enrollment, status,
			   description, is_active, created_by, created_at, updated_at
		FROM courses
		WHERE course_id = $1
	`
	var c domain.Course
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.CourseID, &c.CourseCode, &c.CourseName, &c.SubjectID, &c.DepartmentID, &c.ProgramID,
		&c.SemesterID, &c.SemesterNumber, &c.AcademicYear, &c.MaxStudents, &c.CurrentEnrollment, &c.Status,
		&c.Description, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrCourseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}
	return &c, nil
}

func (r *courseRepository) GetByCode(ctx context.Context, code string) (*domain.Course, error) {
	query := `
		SELECT course_id, course_code, course_name, subject_id, department_id, program_id,
			   semester_id, semester_number, academic_year, max_students, current_enrollment, status,
			   description, is_active, created_by, created_at, updated_at
		FROM courses
		WHERE course_code = $1
	`
	var c domain.Course
	err := r.db.QueryRow(ctx, query, code).Scan(
		&c.CourseID, &c.CourseCode, &c.CourseName, &c.SubjectID, &c.DepartmentID, &c.ProgramID,
		&c.SemesterID, &c.SemesterNumber, &c.AcademicYear, &c.MaxStudents, &c.CurrentEnrollment, &c.Status,
		&c.Description, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrCourseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course by code: %w", err)
	}
	return &c, nil
}

func (r *courseRepository) Update(ctx context.Context, course *domain.Course) error {
	query := `
		UPDATE courses
		SET course_name = $2, max_students = $3, status = $4, description = $5, is_active = $6, updated_at = now()
		WHERE course_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		course.CourseID,
		course.CourseName,
		course.MaxStudents,
		course.Status,
		course.Description,
		course.IsActive,
	).Scan(&course.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrCourseNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
	}
	return nil
}

func (r *courseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE courses SET is_active = false, updated_at = now() WHERE course_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrCourseNotFound
	}
	return nil
}

func (r *courseRepository) List(ctx context.Context, filter domain.CourseFilter, limit, offset int) ([]*domain.CourseWithDetails, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.DepartmentID != nil {
		conditions = append(conditions, fmt.Sprintf("c.department_id = $%d", argNum))
		args = append(args, *filter.DepartmentID)
		argNum++
	}
	if filter.ProgramID != nil {
		conditions = append(conditions, fmt.Sprintf("c.program_id = $%d", argNum))
		args = append(args, *filter.ProgramID)
		argNum++
	}
	if filter.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("c.semester_id = $%d", argNum))
		args = append(args, *filter.SemesterID)
		argNum++
	}
	if filter.SubjectID != nil {
		conditions = append(conditions, fmt.Sprintf("c.subject_id = $%d", argNum))
		args = append(args, *filter.SubjectID)
		argNum++
	}
	if filter.SemesterNumber != nil {
		conditions = append(conditions, fmt.Sprintf("c.semester_number = $%d", argNum))
		args = append(args, *filter.SemesterNumber)
		argNum++
	}
	if filter.AcademicYear != nil {
		conditions = append(conditions, fmt.Sprintf("c.academic_year = $%d", argNum))
		args = append(args, *filter.AcademicYear)
		argNum++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("c.status = $%d", argNum))
		args = append(args, *filter.Status)
		argNum++
	}
	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf("(c.course_name ILIKE $%d OR c.course_code ILIKE $%d)", argNum, argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}
	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("c.is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM courses c %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	// List query with joins
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT c.course_id, c.course_code, c.course_name, c.subject_id, c.department_id, c.program_id,
			   c.semester_id, c.semester_number, c.academic_year, c.max_students, c.current_enrollment, c.status,
			   c.description, c.is_active, c.created_by, c.created_at, c.updated_at,
			   s.subject_id, s.subject_code, s.subject_name, s.credits, s.subject_type,
			   d.department_id, d.department_name, d.department_code,
			   sem.semester_id, sem.semester_name, sem.semester_code
		FROM courses c
		JOIN subjects s ON c.subject_id = s.subject_id
		JOIN departments d ON c.department_id = d.department_id
		JOIN semesters sem ON c.semester_id = sem.semester_id
		%s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list courses: %w", err)
	}
	defer rows.Close()

	var courses []*domain.CourseWithDetails
	for rows.Next() {
		var c domain.CourseWithDetails
		if err := rows.Scan(
			&c.CourseID, &c.CourseCode, &c.CourseName, &c.SubjectID, &c.DepartmentID, &c.ProgramID,
			&c.SemesterID, &c.SemesterNumber, &c.AcademicYear, &c.MaxStudents, &c.CurrentEnrollment, &c.Status,
			&c.Description, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
			&c.Subject.SubjectID, &c.Subject.SubjectCode, &c.Subject.SubjectName, &c.Subject.Credits, &c.Subject.SubjectType,
			&c.Department.DepartmentID, &c.Department.DepartmentName, &c.Department.DepartmentCode,
			&c.Semester.SemesterID, &c.Semester.SemesterName, &c.Semester.SemesterCode,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, &c)
	}

	return courses, total, nil
}

func (r *courseRepository) GetWithDetails(ctx context.Context, id uuid.UUID) (*domain.CourseWithDetails, error) {
	query := `
		SELECT c.course_id, c.course_code, c.course_name, c.subject_id, c.department_id, c.program_id,
			   c.semester_id, c.semester_number, c.academic_year, c.max_students, c.current_enrollment, c.status,
			   c.description, c.is_active, c.created_by, c.created_at, c.updated_at,
			   s.subject_id, s.subject_code, s.subject_name, s.credits, s.subject_type,
			   d.department_id, d.department_name, d.department_code,
			   sem.semester_id, sem.semester_name, sem.semester_code,
			   p.program_id, p.program_name, p.program_code, p.duration_years
		FROM courses c
		JOIN subjects s ON c.subject_id = s.subject_id
		JOIN departments d ON c.department_id = d.department_id
		JOIN semesters sem ON c.semester_id = sem.semester_id
		LEFT JOIN programs p ON c.program_id = p.program_id
		WHERE c.course_id = $1
	`

	var c domain.CourseWithDetails
	var programID, programName, programCode *string
	var durationYears *int

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.CourseID, &c.CourseCode, &c.CourseName, &c.SubjectID, &c.DepartmentID, &c.ProgramID,
		&c.SemesterID, &c.SemesterNumber, &c.AcademicYear, &c.MaxStudents, &c.CurrentEnrollment, &c.Status,
		&c.Description, &c.IsActive, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
		&c.Subject.SubjectID, &c.Subject.SubjectCode, &c.Subject.SubjectName, &c.Subject.Credits, &c.Subject.SubjectType,
		&c.Department.DepartmentID, &c.Department.DepartmentName, &c.Department.DepartmentCode,
		&c.Semester.SemesterID, &c.Semester.SemesterName, &c.Semester.SemesterCode,
		&programID, &programName, &programCode, &durationYears,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrCourseNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course with details: %w", err)
	}

	// Set program if exists
	if programID != nil {
		pid, _ := uuid.Parse(*programID)
		c.Program = &domain.ProgramBasic{
			ProgramID:     pid,
			ProgramName:   *programName,
			ProgramCode:   *programCode,
			DurationYears: *durationYears,
		}
	}

	// Get faculty assignments
	facultyQuery := `
		SELECT f.faculty_id, fc.role, fc.is_primary
		FROM faculty_courses fc
		JOIN faculties f ON fc.faculty_id = f.faculty_id
		WHERE fc.course_id = $1 AND fc.is_active = true
	`
	rows, err := r.db.Query(ctx, facultyQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get course faculty: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var fb domain.FacultyCourseBasic
		if err := rows.Scan(&fb.FacultyID, &fb.Role, &fb.IsPrimary); err != nil {
			return nil, fmt.Errorf("failed to scan faculty: %w", err)
		}
		c.Faculty = append(c.Faculty, fb)
	}

	return &c, nil
}

func (r *courseRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE courses SET status = $2, updated_at = now() WHERE course_id = $1`
	result, err := r.db.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update course status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrCourseNotFound
	}
	return nil
}

func (r *courseRepository) IncrementEnrollment(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE courses SET current_enrollment = current_enrollment + 1, updated_at = now() WHERE course_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment enrollment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrCourseNotFound
	}
	return nil
}

func (r *courseRepository) DecrementEnrollment(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE courses SET current_enrollment = GREATEST(current_enrollment - 1, 0), updated_at = now() WHERE course_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to decrement enrollment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrCourseNotFound
	}
	return nil
}
