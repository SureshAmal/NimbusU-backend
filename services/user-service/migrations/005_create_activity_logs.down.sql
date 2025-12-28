DROP INDEX IF EXISTS idx_user_activity_logs_resource_type;
DROP INDEX IF EXISTS idx_user_activity_logs_created_at;
DROP INDEX IF EXISTS idx_user_activity_logs_action;
DROP INDEX IF EXISTS idx_user_activity_logs_user_id;
DROP TABLE IF EXISTS user_activity_logs CASCADE;
