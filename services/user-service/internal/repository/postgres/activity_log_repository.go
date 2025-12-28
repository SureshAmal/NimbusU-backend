package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
)

type activityLogRepository struct {
	db *pgxpool.Pool
}

// NewActivityLogRepository creates a new activity log repository
func NewActivityLogRepository(db *pgxpool.Pool) domain.ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, log *domain.UserActivityLog) error {
	query := `
		INSERT INTO user_activity_logs (log_id, user_id, action, resource_type, resource_id, 
		                                 ip_address, user_agent, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query,
		log.LogID,
		log.UserID,
		log.Action,
		log.ResourceType,
		log.ResourceID,
		log.IPAddress,
		log.UserAgent,
		log.Details,
	).Scan(&log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create activity log: %w", err)
	}

	return nil
}

func (r *activityLogRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.UserActivityLog, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM user_activity_logs WHERE user_id = $1`
	var total int64
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count activity logs: %w", err)
	}

	// Get paginated results
	query := `
		SELECT log_id, user_id, action, resource_type, resource_id, 
		       ip_address, user_agent, details, created_at
		FROM user_activity_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get activity logs: %w", err)
	}
	defer rows.Close()

	logs := []*domain.UserActivityLog{}
	for rows.Next() {
		var log domain.UserActivityLog
		err := rows.Scan(
			&log.LogID,
			&log.UserID,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.IPAddress,
			&log.UserAgent,
			&log.Details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan activity log: %w", err)
		}
		logs = append(logs, &log)
	}

	return logs, total, nil
}

func (r *activityLogRepository) GetByAction(ctx context.Context, action string, limit, offset int) ([]*domain.UserActivityLog, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM user_activity_logs WHERE action = $1`
	var total int64
	err := r.db.QueryRow(ctx, countQuery, action).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count activity logs: %w", err)
	}

	// Get paginated results
	query := `
		SELECT log_id, user_id, action, resource_type, resource_id, 
		       ip_address, user_agent, details, created_at
		FROM user_activity_logs
		WHERE action = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, action, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get activity logs: %w", err)
	}
	defer rows.Close()

	logs := []*domain.UserActivityLog{}
	for rows.Next() {
		var log domain.UserActivityLog
		err := rows.Scan(
			&log.LogID,
			&log.UserID,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.IPAddress,
			&log.UserAgent,
			&log.Details,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan activity log: %w", err)
		}
		logs = append(logs, &log)
	}

	return logs, total, nil
}
