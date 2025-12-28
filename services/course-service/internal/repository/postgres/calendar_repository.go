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

type calendarRepository struct {
	db *pgxpool.Pool
}

func NewCalendarRepository(db *pgxpool.Pool) domain.CalendarRepository {
	return &calendarRepository{db: db}
}

func (r *calendarRepository) Create(ctx context.Context, event *domain.AcademicCalendarEvent) error {
	query := `
		INSERT INTO academic_calendar (event_id, semester_id, event_name, event_type, start_date, end_date, description, is_holiday, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`
	event.EventID = uuid.New()
	err := r.db.QueryRow(ctx, query,
		event.EventID,
		event.SemesterID,
		event.EventName,
		event.EventType,
		event.StartDate,
		event.EndDate,
		event.Description,
		event.IsHoliday,
		event.CreatedBy,
	).Scan(&event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create calendar event: %w", err)
	}
	return nil
}

func (r *calendarRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AcademicCalendarEvent, error) {
	query := `
		SELECT event_id, semester_id, event_name, event_type, start_date, end_date, description, is_holiday, created_by, created_at, updated_at
		FROM academic_calendar
		WHERE event_id = $1
	`
	var e domain.AcademicCalendarEvent
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.EventID, &e.SemesterID, &e.EventName, &e.EventType, &e.StartDate, &e.EndDate,
		&e.Description, &e.IsHoliday, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrCalendarEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar event: %w", err)
	}
	return &e, nil
}

func (r *calendarRepository) Update(ctx context.Context, event *domain.AcademicCalendarEvent) error {
	query := `
		UPDATE academic_calendar
		SET event_name = $2, event_type = $3, start_date = $4, end_date = $5, description = $6, is_holiday = $7, updated_at = now()
		WHERE event_id = $1
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		event.EventID,
		event.EventName,
		event.EventType,
		event.StartDate,
		event.EndDate,
		event.Description,
		event.IsHoliday,
	).Scan(&event.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrCalendarEventNotFound
	}
	if err != nil {
		return fmt.Errorf("failed to update calendar event: %w", err)
	}
	return nil
}

func (r *calendarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM academic_calendar WHERE event_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete calendar event: %w", err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrCalendarEventNotFound
	}
	return nil
}

func (r *calendarRepository) List(ctx context.Context, filter domain.CalendarFilter, limit, offset int) ([]*domain.AcademicCalendarEventWithDetails, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filter.SemesterID != nil {
		conditions = append(conditions, fmt.Sprintf("e.semester_id = $%d", argNum))
		args = append(args, *filter.SemesterID)
		argNum++
	}
	if filter.EventType != nil {
		conditions = append(conditions, fmt.Sprintf("e.event_type = $%d", argNum))
		args = append(args, *filter.EventType)
		argNum++
	}
	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("e.start_date >= $%d", argNum))
		args = append(args, *filter.StartDate)
		argNum++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("e.start_date <= $%d", argNum))
		args = append(args, *filter.EndDate)
		argNum++
	}
	if filter.IsHoliday != nil {
		conditions = append(conditions, fmt.Sprintf("e.is_holiday = $%d", argNum))
		args = append(args, *filter.IsHoliday)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM academic_calendar e %s", whereClause)
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count calendar events: %w", err)
	}

	// List query with semester join
	args = append(args, limit, offset)
	listQuery := fmt.Sprintf(`
		SELECT e.event_id, e.semester_id, e.event_name, e.event_type, e.start_date, e.end_date,
			   e.description, e.is_holiday, e.created_by, e.created_at, e.updated_at,
			   s.semester_id, s.semester_name, s.semester_code
		FROM academic_calendar e
		JOIN semesters s ON e.semester_id = s.semester_id
		%s
		ORDER BY e.start_date
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list calendar events: %w", err)
	}
	defer rows.Close()

	var events []*domain.AcademicCalendarEventWithDetails
	for rows.Next() {
		var e domain.AcademicCalendarEventWithDetails
		if err := rows.Scan(
			&e.EventID, &e.SemesterID, &e.EventName, &e.EventType, &e.StartDate, &e.EndDate,
			&e.Description, &e.IsHoliday, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
			&e.Semester.SemesterID, &e.Semester.SemesterName, &e.Semester.SemesterCode,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan calendar event: %w", err)
		}
		events = append(events, &e)
	}

	return events, total, nil
}
