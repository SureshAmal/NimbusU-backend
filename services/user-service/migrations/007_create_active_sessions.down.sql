DROP INDEX IF EXISTS idx_active_sessions_expires_at;
DROP INDEX IF EXISTS idx_active_sessions_refresh_token;
DROP INDEX IF EXISTS idx_active_sessions_user_id;
DROP TABLE IF EXISTS active_sessions CASCADE;
