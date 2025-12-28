-- 001_create_departments.down.sql
DROP TRIGGER IF EXISTS update_departments_updated_at ON departments;
DROP INDEX IF EXISTS idx_departments_head;
DROP INDEX IF EXISTS idx_departments_active;
DROP INDEX IF EXISTS idx_departments_code;
DROP TABLE IF EXISTS departments CASCADE;
