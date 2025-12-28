-- 009_create_faculty_courses.down.sql
DROP INDEX IF EXISTS idx_unique_active_faculty_course;
DROP INDEX IF EXISTS idx_faculty_courses_role;
DROP INDEX IF EXISTS idx_faculty_courses_active;
DROP INDEX IF EXISTS idx_faculty_courses_course;
DROP INDEX IF EXISTS idx_faculty_courses_faculty;
DROP TABLE IF EXISTS faculty_courses CASCADE;
