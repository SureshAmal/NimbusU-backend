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

type semesterRepository struct {
	db *pgxpool.Pool
}

func NewSemesterRepository(db *pgxpool.Pool) domain.SemesterRepository {
	return &semesterRepository{db: db}
}

func (r *semesterRepository) Create(ctx context.Context, semester *domain.Semester) error {
	query := `
		INSERT INTO semesters (semester_id, semester_name, semester_code, academic_year, start_date, end_date, registration_start, registration_end, is_current)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`
	semester.SemesterID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		semester.SemesterID,
		semester.SemesterName,
		semester.SemesterCode,
		semester.AcademicYear,
		semester.StartDate,
		semester.EndDate,
		semester.RegistrationStart,
		semester.RegistrationEnd,
		semester.IsCurrent,
	).Scan(&semester.CreatedAt, &semester.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrSemesterCodeExists
		}
		return fmt.Errorf("failed to create semester: %w", err)
	}
	return nil
}

func (r *semesterRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Semester, error) {
	query := `
		SELECT semester_id, semester_name, semester_code, academic_year, start_date, end_date,
			   registration_start, registration_end, is_current, created_at, updated_at
		FROM semesters
		WHERE semester_id = $1
	`
	var s domain.Semester
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.SemesterID, &s.SemesterName, &s.SemesterCode, &s.AcademicYear, &s.StartDate, &s.EndDate,
		&s.RegistrationStart, &s.RegistrationEnd, &s.IsCurrent, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrSemesterNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get semester: %w", err)
	}
	return &s, nil
}

func (r *semesterRepository) GetByCode(ctx context.Context, code string) (*domain.Semester, error) {
	query := `
		SELECT semester_id, semester_name, semester_code, academic_year, start_date, end_date,
			   registration_start, registration_end, is_current, created_at, updated_at
		FROM semesters
		WHERE semester_code = $1
	`
	var s domain.Semester
	err := r.db.QueryRow(ctx, query, code).Scan(
		&s.SemesterID, &s.SemesterName, &s.SemesterCode, &s.AcademicYear, &s.StartDate, &s.EndDate,
		&s.RegistrationStart, &s.RegistrationEnd, &s.IsCurrent, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrSemesterNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get semester by code: %w", err)
	}
	return &s, nil
}

func (r *semesterRepository) GetCurrent(ctx context.Context) (*domain.Semester, error) {
	query := `
		SELECT semester_id, semester_name, semester_code, academic_year, start_date, end_date,
			   registration_start, registration_end, is_current, created_at, updated_at
		FROM semesters
		WHERE is_current = true
		LIMIT 1
	`
	var s domain.Semester
	err := r.db.QueryRow(ctx, query).Scan(
		&s.SemesterID, &s.SemesterName, &s.SemesterCode, &s.AcademicYear, &s.StartDate, &s.EndDate,
		&s.RegistrationStart, &s.RegistrationEnd, &s.IsCurrent, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNoCurrentSemester
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get current semester: %w", err)
	}
	return &s, nil
}

func (r *semesterRepository) Update(ctx context.Context, semester *domain.Semester) error {
	query := `
		UPDATE semesters
		SET semester_name = $2, academic_year = $3, start_date = $4, end_date = $5,
			registration_start = $6, registration_end = $7, updated_at = now()
		WHERE semester_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		semester.SemesterID,
		semester.SemesterName,
		semester.AcademicYear,
		semester.StartDate,
		semester.EndDate,
		semester.RegistrationStart,
		semester.RegistrationEnd,
	).Scan(&semester.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrSemesterNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update semester: %w", err)
	}
	return nil
}

func (r *semesterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM semesters WHERE semester_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete semester: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrSemesterNotFound
	}
	return nil
}

func (r *semesterRepository) List(ctx context.Context, filter domain.SemesterFilter, limit, offset int) ([]*domain.Semester, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.AcademicYear != nil {
		conditions = append(conditions, fmt.Sprintf("academic_year = $%d", argNum))
		args = append(args, *filter.AcademicYear)
		argNum++
	}
	if filter.IsCurrent != nil {
		conditions = append(conditions, fmt.Sprintf("is_current = $%d", argNum))
		args = append(args, *filter.IsCurrent)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM semesters %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count semesters: %w", err)
	}

	// List query
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT semester_id, semester_name, semester_code, academic_year, start_date, end_date,
			   registration_start, registration_end, is_current, created_at, updated_at
		FROM semesters
		%s
		ORDER BY academic_year DESC, start_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list semesters: %w", err)
	}
	defer rows.Close()

	var semesters []*domain.Semester
	for rows.Next() {
		var s domain.Semester
		if err := rows.Scan(
			&s.SemesterID, &s.SemesterName, &s.SemesterCode, &s.AcademicYear, &s.StartDate, &s.EndDate,
			&s.RegistrationStart, &s.RegistrationEnd, &s.IsCurrent, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan semester: %w", err)
		}
		semesters = append(semesters, &s)
	}

	return semesters, total, nil
}

func (r *semesterRepository) SetCurrent(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Unset all current semesters
	_, err = tx.Exec(ctx, "UPDATE semesters SET is_current = false WHERE is_current = true")
	if err != nil {
		return fmt.Errorf("failed to unset current semesters: %w", err)
	}

	// Set the new current semester
	result, err := tx.Exec(ctx, "UPDATE semesters SET is_current = true, updated_at = now() WHERE semester_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to set current semester: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrSemesterNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
