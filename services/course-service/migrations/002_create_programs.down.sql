-- 002_create_programs.down.sql
DROP TRIGGER IF EXISTS update_programs_updated_at ON programs;
DROP INDEX IF EXISTS idx_programs_degree_type;
DROP INDEX IF EXISTS idx_programs_active;
DROP INDEX IF EXISTS idx_programs_code;
DROP INDEX IF EXISTS idx_programs_department;
DROP TABLE IF EXISTS programs CASCADE;
