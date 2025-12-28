package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type enrollmentRepository struct {
	db *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) domain.EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Create(ctx context.Context, enrollment *domain.CourseEnrollment) error {
	query := `
		INSERT INTO course_enrollments (enrollment_id, student_id, course_id, enrollment_status, enrolled_by, waitlist_position)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING enrollment_date, created_at, updated_at
	`
	enrollment.EnrollmentID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		enrollment.EnrollmentID,
		enrollment.StudentID,
		enrollment.CourseID,
		enrollment.EnrollmentStatus,
		enrollment.EnrolledBy,
		enrollment.WaitlistPosition,
	).Scan(&enrollment.EnrollmentDate, &enrollment.CreatedAt, &enrollment.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrAlreadyEnrolled
		}
		return fmt.Errorf("failed to create enrollment: %w", err)
	}
	return nil
}

func (r *enrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CourseEnrollment, error) {
	query := `
		SELECT enrollment_id, student_id, course_id, enrollment_status, enrolled_by, enrollment_date,
			   dropped_date, drop_reason, completion_date, grade, grade_points, waitlist_position, created_at, updated_at
		FROM course_enrollments
		WHERE enrollment_id = $1
	`
	var e domain.CourseEnrollment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.EnrollmentID, &e.StudentID, &e.CourseID, &e.EnrollmentStatus, &e.EnrolledBy, &e.EnrollmentDate,
		&e.DroppedDate, &e.DropReason, &e.CompletionDate, &e.Grade, &e.GradePoints, &e.WaitlistPosition, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrEnrollmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	return &e, nil
}

func (r *enrollmentRepository) GetByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*domain.CourseEnrollment, error) {
	query := `
		SELECT enrollment_id, student_id, course_id, enrollment_status, enrolled_by, enrollment_date,
			   dropped_date, drop_reason, completion_date, grade, grade_points, waitlist_position, created_at, updated_at
		FROM course_enrollments
		WHERE student_id = $1 AND course_id = $2
	`
	var e domain.CourseEnrollment
	err := r.db.QueryRow(ctx, query, studentID, courseID).Scan(
		&e.EnrollmentID, &e.StudentID, &e.CourseID, &e.EnrollmentStatus, &e.EnrolledBy, &e.EnrollmentDate,
		&e.DroppedDate, &e.DropReason, &e.CompletionDate, &e.Grade, &e.GradePoints, &e.WaitlistPosition, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrEnrollmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment by student and course: %w", err)
	}
	return &e, nil
}

func (r *enrollmentRepository) Update(ctx context.Context, enrollment *domain.CourseEnrollment) error {
	query := `
		UPDATE course_enrollments
		SET enrollment_status = $2, dropped_date = $3, drop_reason = $4, completion_date = $5,
			grade = $6, grade_points = $7, waitlist_position = $8, updated_at = now()
		WHERE enrollment_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		enrollment.EnrollmentID,
		enrollment.EnrollmentStatus,
		enrollment.DroppedDate,
		enrollment.DropReason,
		enrollment.CompletionDate,
		enrollment.Grade,
		enrollment.GradePoints,
		enrollment.WaitlistPosition,
	).Scan(&enrollment.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrEnrollmentNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update enrollment: %w", err)
	}
	return nil
}

func (r *enrollmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM course_enrollments WHERE enrollment_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrEnrollmentNotFound
	}
	return nil
}

func (r *enrollmentRepository) ListByStudent(ctx context.Context, filter domain.EnrollmentFilter, limit, offset int) ([]*domain.EnrollmentWithDetails, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.StudentID != nil {
		conditions = append(conditions, fmt.Sprintf("e.student_id = $%d", argNum))
		args = append(args, *filter.StudentID)
		argNum++
	}
	if filter.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("c.semester_id = $%d", argNum))
		args = append(args, *filter.SemesterID)
		argNum++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("e.enrollment_status = $%d", argNum))
		args = append(args, *filter.Status)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM course_enrollments e
		JOIN courses c ON e.course_id = c.course_id
		%s
	`, whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count enrollments: %w", err)
	}

	// List query
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT e.enrollment_id, e.student_id, e.course_id, e.enrollment_status, e.enrolled_by, e.enrollment_date,
			   e.dropped_date, e.drop_reason, e.completion_date, e.grade, e.grade_points, e.waitlist_position, e.created_at, e.updated_at,
			   s.student_id, s.registration_number,
			   c.course_id, c.course_code, c.course_name,
			   sem.semester_id, sem.semester_name, sem.semester_code
		FROM course_enrollments e
		JOIN students s ON e.student_id = s.student_id
		JOIN courses c ON e.course_id = c.course_id
		JOIN semesters sem ON c.semester_id = sem.semester_id
		%s
		ORDER BY e.enrollment_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list enrollments: %w", err)
	}
	defer rows.Close()

	var enrollments []*domain.EnrollmentWithDetails
	for rows.Next() {
		var e domain.EnrollmentWithDetails
		var sem domain.SemesterBasic
		if err := rows.Scan(
			&e.EnrollmentID, &e.StudentID, &e.CourseID, &e.EnrollmentStatus, &e.EnrolledBy, &e.EnrollmentDate,
			&e.DroppedDate, &e.DropReason, &e.CompletionDate, &e.Grade, &e.GradePoints, &e.WaitlistPosition, &e.CreatedAt, &e.UpdatedAt,
			&e.Student.StudentID, &e.Student.RegistrationNumber,
			&e.Course.CourseID, &e.Course.CourseCode, &e.Course.CourseName,
			&sem.SemesterID, &sem.SemesterName, &sem.SemesterCode,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		e.Semester = &sem
		enrollments = append(enrollments, &e)
	}

	return enrollments, total, nil
}

func (r *enrollmentRepository) ListByCourse(ctx context.Context, courseID uuid.UUID, status *string, limit, offset int) ([]*domain.EnrollmentWithDetails, int64, error) {
	var conditions []string
	var args []interface{}
	args = append(args, courseID)
	conditions = append(conditions, "e.course_id = $1")
	argNum := 2

	if status != nil {
		conditions = append(conditions, fmt.Sprintf("e.enrollment_status = $%d", argNum))
		args = append(args, *status)
		argNum++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM course_enrollments e %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count enrollments: %w", err)
	}

	// List query
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT e.enrollment_id, e.student_id, e.course_id, e.enrollment_status, e.enrolled_by, e.enrollment_date,
			   e.dropped_date, e.drop_reason, e.completion_date, e.grade, e.grade_points, e.waitlist_position, e.created_at, e.updated_at,
			   s.student_id, s.registration_number,
			   c.course_id, c.course_code, c.course_name
		FROM course_enrollments e
		JOIN students s ON e.student_id = s.student_id
		JOIN courses c ON e.course_id = c.course_id
		%s
		ORDER BY e.enrollment_date
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list enrollments by course: %w", err)
	}
	defer rows.Close()

	var enrollments []*domain.EnrollmentWithDetails
	for rows.Next() {
		var e domain.EnrollmentWithDetails
		if err := rows.Scan(
			&e.EnrollmentID, &e.StudentID, &e.CourseID, &e.EnrollmentStatus, &e.EnrolledBy, &e.EnrollmentDate,
			&e.DroppedDate, &e.DropReason, &e.CompletionDate, &e.Grade, &e.GradePoints, &e.WaitlistPosition, &e.CreatedAt, &e.UpdatedAt,
			&e.Student.StudentID, &e.Student.RegistrationNumber,
			&e.Course.CourseID, &e.Course.CourseCode, &e.Course.CourseName,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		enrollments = append(enrollments, &e)
	}

	return enrollments, total, nil
}

func (r *enrollmentRepository) GetEnrollmentSummary(ctx context.Context, courseID uuid.UUID) (enrolled, waitlisted, dropped, completed int, err error) {
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE enrollment_status = 'enrolled') as enrolled,
			COUNT(*) FILTER (WHERE enrollment_status = 'waitlisted') as waitlisted,
			COUNT(*) FILTER (WHERE enrollment_status = 'dropped') as dropped,
			COUNT(*) FILTER (WHERE enrollment_status = 'completed') as completed
		FROM course_enrollments
		WHERE course_id = $1
	`
	err = r.db.QueryRow(ctx, query, courseID).Scan(&enrolled, &waitlisted, &dropped, &completed)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get enrollment summary: %w", err)
	}
	return enrolled, waitlisted, dropped, completed, nil
}

func (r *enrollmentRepository) GetNextWaitlistPosition(ctx context.Context, courseID uuid.UUID) (int, error) {
	query := `
		SELECT COALESCE(MAX(waitlist_position), 0) + 1
		FROM course_enrollments
		WHERE course_id = $1 AND enrollment_status = 'waitlisted'
	`
	var position int
	err := r.db.QueryRow(ctx, query, courseID).Scan(&position)
	if err != nil {
		return 0, fmt.Errorf("failed to get next waitlist position: %w", err)
	}
	return position, nil
}

func (r *enrollmentRepository) PromoteFromWaitlist(ctx context.Context, courseID uuid.UUID) (*domain.CourseEnrollment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get the first waitlisted student
	query := `
		SELECT enrollment_id, student_id, course_id, enrollment_status, enrolled_by, enrollment_date,
			   dropped_date, drop_reason, completion_date, grade, grade_points, waitlist_position, created_at, updated_at
		FROM course_enrollments
		WHERE course_id = $1 AND enrollment_status = 'waitlisted'
		ORDER BY waitlist_position
		LIMIT 1
		FOR UPDATE
	`
	var e domain.CourseEnrollment
	err = tx.QueryRow(ctx, query, courseID).Scan(
		&e.EnrollmentID, &e.StudentID, &e.CourseID, &e.EnrollmentStatus, &e.EnrolledBy, &e.EnrollmentDate,
		&e.DroppedDate, &e.DropReason, &e.CompletionDate, &e.Grade, &e.GradePoints, &e.WaitlistPosition, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil // No one on waitlist
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get waitlisted enrollment: %w", err)
	}

	// Update to enrolled
	now := time.Now()
	updateQuery := `
		UPDATE course_enrollments
		SET enrollment_status = 'enrolled', waitlist_position = NULL, enrollment_date = $2, updated_at = $2
		WHERE enrollment_id = $1
	`
	_, err = tx.Exec(ctx, updateQuery, e.EnrollmentID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to promote from waitlist: %w", err)
	}

	// Decrement waitlist positions for remaining
	decrementQuery := `
		UPDATE course_enrollments
		SET waitlist_position = waitlist_position - 1
		WHERE course_id = $1 AND enrollment_status = 'waitlisted' AND waitlist_position > $2
	`
	_, err = tx.Exec(ctx, decrementQuery, courseID, *e.WaitlistPosition)
	if err != nil {
		return nil, fmt.Errorf("failed to decrement waitlist positions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	e.EnrollmentStatus = "enrolled"
	e.WaitlistPosition = nil
	e.EnrollmentDate = now
	e.UpdatedAt = now

	return &e, nil
}
