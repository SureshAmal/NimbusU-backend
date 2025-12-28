-- 007_create_courses.down.sql
DROP TRIGGER IF EXISTS update_courses_updated_at ON courses;
DROP INDEX IF EXISTS idx_courses_semester_number;
DROP INDEX IF EXISTS idx_courses_active;
DROP INDEX IF EXISTS idx_courses_academic_year;
DROP INDEX IF EXISTS idx_courses_status;
DROP INDEX IF EXISTS idx_courses_semester;
DROP INDEX IF EXISTS idx_courses_program;
DROP INDEX IF EXISTS idx_courses_department;
DROP INDEX IF EXISTS idx_courses_subject;
DROP INDEX IF EXISTS idx_courses_code;
DROP TABLE IF EXISTS courses CASCADE;
