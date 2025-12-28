-- 008_create_prerequisites.down.sql
DROP INDEX IF EXISTS idx_corequisites_coreq;
DROP INDEX IF EXISTS idx_corequisites_subject;
DROP INDEX IF EXISTS idx_prerequisites_prereq;
DROP INDEX IF EXISTS idx_prerequisites_subject;
DROP TABLE IF EXISTS course_corequisites CASCADE;
DROP TABLE IF EXISTS course_prerequisites CASCADE;
