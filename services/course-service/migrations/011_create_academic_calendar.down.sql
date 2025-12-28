-- 011_create_academic_calendar.down.sql
DROP TRIGGER IF EXISTS update_calendar_updated_at ON academic_calendar;
DROP INDEX IF EXISTS idx_calendar_holiday;
DROP INDEX IF EXISTS idx_calendar_type;
DROP INDEX IF EXISTS idx_calendar_dates;
DROP INDEX IF EXISTS idx_calendar_semester;
DROP TABLE IF EXISTS academic_calendar CASCADE;
