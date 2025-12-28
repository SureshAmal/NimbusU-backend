-- 003_create_subjects.down.sql
DROP TRIGGER IF EXISTS update_subjects_updated_at ON subjects;
DROP INDEX IF EXISTS idx_subjects_type;
DROP INDEX IF EXISTS idx_subjects_active;
DROP INDEX IF EXISTS idx_subjects_department;
DROP INDEX IF EXISTS idx_subjects_code;
DROP TABLE IF EXISTS subjects CASCADE;
