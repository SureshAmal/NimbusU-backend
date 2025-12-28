-- 005_create_faculties.down.sql
ALTER TABLE departments DROP CONSTRAINT IF EXISTS fk_departments_head;
DROP TRIGGER IF EXISTS update_faculties_updated_at ON faculties;
DROP INDEX IF EXISTS idx_faculties_active;
DROP INDEX IF EXISTS idx_faculties_department;
DROP INDEX IF EXISTS idx_faculties_employee;
DROP INDEX IF EXISTS idx_faculties_user;
DROP TABLE IF EXISTS faculties CASCADE;
