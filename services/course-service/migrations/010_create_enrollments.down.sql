-- 010_create_enrollments.down.sql
DROP TRIGGER IF EXISTS trigger_update_enrollment_count ON course_enrollments;
DROP FUNCTION IF EXISTS update_course_enrollment_count();
DROP TRIGGER IF EXISTS update_enrollments_updated_at ON course_enrollments;
DROP INDEX IF EXISTS idx_enrollments_waitlist;
DROP INDEX IF EXISTS idx_enrollments_date;
DROP INDEX IF EXISTS idx_enrollments_status;
DROP INDEX IF EXISTS idx_enrollments_course;
DROP INDEX IF EXISTS idx_enrollments_student;
DROP TABLE IF EXISTS course_enrollments CASCADE;
