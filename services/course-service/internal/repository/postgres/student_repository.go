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

type studentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) domain.StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, student *domain.Student) error {
	query := `
		INSERT INTO students (student_id, user_id, registration_number, roll_number, department_id, program_id,
			current_semester, batch_year, admission_date, current_cgpa, total_credits_earned, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at
	`
	student.StudentID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		student.StudentID,
		student.UserID,
		student.RegistrationNumber,
		student.RollNumber,
		student.DepartmentID,
		student.ProgramID,
		student.CurrentSemester,
		student.BatchYear,
		student.AdmissionDate,
		student.CurrentCGPA,
		student.TotalCreditsEarned,
		student.IsActive,
	).Scan(&student.CreatedAt, &student.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "registration_number") {
				return domain.ErrRegistrationNumberExists
			}
		}
		return fmt.Errorf("failed to create student: %w", err)
	}
	return nil
}

func (r *studentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Student, error) {
	query := `
		SELECT student_id, user_id, registration_number, roll_number, department_id, program_id,
			   current_semester, batch_year, admission_date, current_cgpa, total_credits_earned, is_active, created_at, updated_at
		FROM students
		WHERE student_id = $1
	`
	var s domain.Student
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.StudentID, &s.UserID, &s.RegistrationNumber, &s.RollNumber, &s.DepartmentID, &s.ProgramID,
		&s.CurrentSemester, &s.BatchYear, &s.AdmissionDate, &s.CurrentCGPA, &s.TotalCreditsEarned, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrStudentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}
	return &s, nil
}

func (r *studentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Student, error) {
	query := `
		SELECT student_id, user_id, registration_number, roll_number, department_id, program_id,
			   current_semester, batch_year, admission_date, current_cgpa, total_credits_earned, is_active, created_at, updated_at
		FROM students
		WHERE user_id = $1
	`
	var s domain.Student
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&s.StudentID, &s.UserID, &s.RegistrationNumber, &s.RollNumber, &s.DepartmentID, &s.ProgramID,
		&s.CurrentSemester, &s.BatchYear, &s.AdmissionDate, &s.CurrentCGPA, &s.TotalCreditsEarned, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrStudentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student by user ID: %w", err)
	}
	return &s, nil
}

func (r *studentRepository) GetByRegistrationNumber(ctx context.Context, regNo string) (*domain.Student, error) {
	query := `
		SELECT student_id, user_id, registration_number, roll_number, department_id, program_id,
			   current_semester, batch_year, admission_date, current_cgpa, total_credits_earned, is_active, created_at, updated_at
		FROM students
		WHERE registration_number = $1
	`
	var s domain.Student
	err := r.db.QueryRow(ctx, query, regNo).Scan(
		&s.StudentID, &s.UserID, &s.RegistrationNumber, &s.RollNumber, &s.DepartmentID, &s.ProgramID,
		&s.CurrentSemester, &s.BatchYear, &s.AdmissionDate, &s.CurrentCGPA, &s.TotalCreditsEarned, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrStudentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student by registration number: %w", err)
	}
	return &s, nil
}

func (r *studentRepository) Update(ctx context.Context, student *domain.Student) error {
	query := `
		UPDATE students
		SET roll_number = $2, current_semester = $3, current_cgpa = $4, total_credits_earned = $5, is_active = $6, updated_at = now()
		WHERE student_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		student.StudentID,
		student.RollNumber,
		student.CurrentSemester,
		student.CurrentCGPA,
		student.TotalCreditsEarned,
		student.IsActive,
	).Scan(&student.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrStudentNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}
	return nil
}

func (r *studentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE students SET is_active = false, updated_at = now() WHERE student_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrStudentNotFound
	}
	return nil
}

func (r *studentRepository) List(ctx context.Context, filter domain.StudentFilter, limit, offset int) ([]*domain.StudentWithDetails, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.DepartmentID != nil {
		conditions = append(conditions, fmt.Sprintf("s.department_id = $%d", argNum))
		args = append(args, *filter.DepartmentID)
		argNum++
	}
	if filter.ProgramID != nil {
		conditions = append(conditions, fmt.Sprintf("s.program_id = $%d", argNum))
		args = append(args, *filter.ProgramID)
		argNum++
	}
	if filter.CurrentSemester != nil {
		conditions = append(conditions, fmt.Sprintf("s.current_semester = $%d", argNum))
		args = append(args, *filter.CurrentSemester)
		argNum++
	}
	if filter.BatchYear != nil {
		conditions = append(conditions, fmt.Sprintf("s.batch_year = $%d", argNum))
		args = append(args, *filter.BatchYear)
		argNum++
	}
	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf("(s.registration_number ILIKE $%d OR s.roll_number ILIKE $%d)", argNum, argNum))
		args = append(args, "%"+*filter.Search+"%")
		argNum++
	}
	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("s.is_active = $%d", argNum))
		args = append(args, *filter.IsActive)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students s %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count students: %w", err)
	}

	// List query with joins
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT s.student_id, s.user_id, s.registration_number, s.roll_number, s.department_id, s.program_id,
			   s.current_semester, s.batch_year, s.admission_date, s.current_cgpa, s.total_credits_earned, s.is_active, s.created_at, s.updated_at,
			   d.department_id, d.department_name, d.department_code,
			   p.program_id, p.program_name, p.program_code, p.duration_years
		FROM students s
		JOIN departments d ON s.department_id = d.department_id
		JOIN programs p ON s.program_id = p.program_id
		%s
		ORDER BY s.registration_number
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list students: %w", err)
	}
	defer rows.Close()

	var students []*domain.StudentWithDetails
	for rows.Next() {
		var s domain.StudentWithDetails
		if err := rows.Scan(
			&s.StudentID, &s.UserID, &s.RegistrationNumber, &s.RollNumber, &s.DepartmentID, &s.ProgramID,
			&s.CurrentSemester, &s.BatchYear, &s.AdmissionDate, &s.CurrentCGPA, &s.TotalCreditsEarned, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
			&s.Department.DepartmentID, &s.Department.DepartmentName, &s.Department.DepartmentCode,
			&s.Program.ProgramID, &s.Program.ProgramName, &s.Program.ProgramCode, &s.Program.DurationYears,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan student: %w", err)
		}
		// Note: Name and Email would come from user service
		students = append(students, &s)
	}

	return students, total, nil
}

func (r *studentRepository) GetWithDetails(ctx context.Context, id uuid.UUID) (*domain.StudentWithDetails, error) {
	query := `
		SELECT s.student_id, s.user_id, s.registration_number, s.roll_number, s.department_id, s.program_id,
			   s.current_semester, s.batch_year, s.admission_date, s.current_cgpa, s.total_credits_earned, s.is_active, s.created_at, s.updated_at,
			   d.department_id, d.department_name, d.department_code,
			   p.program_id, p.program_name, p.program_code, p.duration_years
		FROM students s
		JOIN departments d ON s.department_id = d.department_id
		JOIN programs p ON s.program_id = p.program_id
		WHERE s.student_id = $1
	`

	var st domain.StudentWithDetails
	err := r.db.QueryRow(ctx, query, id).Scan(
		&st.StudentID, &st.UserID, &st.RegistrationNumber, &st.RollNumber, &st.DepartmentID, &st.ProgramID,
		&st.CurrentSemester, &st.BatchYear, &st.AdmissionDate, &st.CurrentCGPA, &st.TotalCreditsEarned, &st.IsActive, &st.CreatedAt, &st.UpdatedAt,
		&st.Department.DepartmentID, &st.Department.DepartmentName, &st.Department.DepartmentCode,
		&st.Program.ProgramID, &st.Program.ProgramName, &st.Program.ProgramCode, &st.Program.DurationYears,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrStudentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student with details: %w", err)
	}

	// Get current enrollments
	enrollmentQuery := `
		SELECT c.course_id, c.course_code, c.course_name, sub.credits
		FROM course_enrollments e
		JOIN courses c ON e.course_id = c.course_id
		JOIN subjects sub ON c.subject_id = sub.subject_id
		JOIN semesters sem ON c.semester_id = sem.semester_id
		WHERE e.student_id = $1 AND e.enrollment_status = 'enrolled' AND sem.is_current = true
	`
	rows, err := r.db.Query(ctx, enrollmentQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get student enrollments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e domain.EnrollmentBasic
		if err := rows.Scan(&e.CourseID, &e.CourseCode, &e.CourseName, &e.Credits); err != nil {
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		st.CurrentEnrollments = append(st.CurrentEnrollments, e)
	}

	return &st, nil
}

func (r *studentRepository) UpdateSemester(ctx context.Context, id uuid.UUID, semester int, cgpa *float64, credits int) error {
	query := `
		UPDATE students 
		SET current_semester = $2, current_cgpa = $3, total_credits_earned = total_credits_earned + $4, updated_at = now()
		WHERE student_id = $1
	`
	result, err := r.db.Exec(ctx, query, id, semester, cgpa, credits)
	if err != nil {
		return fmt.Errorf("failed to update student semester: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrStudentNotFound
	}
	return nil
}
