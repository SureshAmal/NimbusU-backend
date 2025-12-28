-- 006_create_students.down.sql
DROP TRIGGER IF EXISTS update_students_updated_at ON students;
DROP INDEX IF EXISTS idx_students_active;
DROP INDEX IF EXISTS idx_students_semester;
DROP INDEX IF EXISTS idx_students_batch;
DROP INDEX IF EXISTS idx_students_program;
DROP INDEX IF EXISTS idx_students_department;
DROP INDEX IF EXISTS idx_students_registration;
DROP INDEX IF EXISTS idx_students_user;
DROP TABLE IF EXISTS students CASCADE;
