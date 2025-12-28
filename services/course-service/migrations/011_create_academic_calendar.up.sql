-- 011_create_academic_calendar.up.sql
-- Create academic calendar table

CREATE TABLE IF NOT EXISTS academic_calendar (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    semester_id UUID NOT NULL REFERENCES semesters(semester_id) ON DELETE CASCADE,
    event_name VARCHAR(255) NOT NULL,
    event_type VARCHAR(50) NOT NULL CHECK(event_type IN ('holiday', 'exam', 'registration', 'deadline', 'event', 'other')),
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT,
    is_holiday BOOLEAN DEFAULT false,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_calendar_semester ON academic_calendar(semester_id);
CREATE INDEX IF NOT EXISTS idx_calendar_dates ON academic_calendar(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_calendar_type ON academic_calendar(event_type);
CREATE INDEX IF NOT EXISTS idx_calendar_holiday ON academic_calendar(is_holiday) WHERE is_holiday = true;

-- Create trigger for updated_at
CREATE TRIGGER update_calendar_updated_at
    BEFORE UPDATE ON academic_calendar
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
