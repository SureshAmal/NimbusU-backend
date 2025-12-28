package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/domain"
)

type sessionRepository struct {
	db *pgxpool.Pool
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *pgxpool.Pool) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.ActiveSession) error {
	query := `
		INSERT INTO active_sessions (session_id, user_id, refresh_token, device_info, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`

	err := r.db.QueryRow(ctx, query,
		session.SessionID,
		session.UserID,
		session.RefreshToken,
		session.DeviceInfo,
		session.IPAddress,
		session.ExpiresAt,
	).Scan(&session.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.ActiveSession, error) {
	query := `
		SELECT session_id, user_id, refresh_token, device_info, ip_address, expires_at, created_at
		FROM active_sessions
		WHERE refresh_token = $1 AND expires_at > $2
	`

	var session domain.ActiveSession
	err := r.db.QueryRow(ctx, query, refreshToken, time.Now()).Scan(
		&session.SessionID,
		&session.UserID,
		&session.RefreshToken,
		&session.DeviceInfo,
		&session.IPAddress,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.ActiveSession, error) {
	query := `
		SELECT session_id, user_id, refresh_token, device_info, ip_address, expires_at, created_at
		FROM active_sessions
		WHERE user_id = $1 AND expires_at > $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}
	defer rows.Close()

	sessions := []*domain.ActiveSession{}
	for rows.Next() {
		var session domain.ActiveSession
		err := rows.Scan(
			&session.SessionID,
			&session.UserID,
			&session.RefreshToken,
			&session.DeviceInfo,
			&session.IPAddress,
			&session.ExpiresAt,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

func (r *sessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM active_sessions WHERE session_id = $1`

	result, err := r.db.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSessionNotFound
	}

	return nil
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM active_sessions WHERE user_id = $1`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM active_sessions WHERE expires_at <= $1`

	_, err := r.db.Exec(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}
