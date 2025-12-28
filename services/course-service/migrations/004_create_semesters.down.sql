-- 004_create_semesters.down.sql
DROP TRIGGER IF EXISTS update_semesters_updated_at ON semesters;
DROP INDEX IF EXISTS idx_semesters_code;
DROP INDEX IF EXISTS idx_semesters_dates;
DROP INDEX IF EXISTS idx_semesters_academic_year;
DROP INDEX IF EXISTS idx_semesters_current;
DROP TABLE IF EXISTS semesters CASCADE;
