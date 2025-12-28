-- 004_create_semesters.up.sql
-- Create semesters table

CREATE TABLE IF NOT EXISTS semesters (
    semester_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    semester_name VARCHAR(50) NOT NULL,
    semester_code VARCHAR(20) UNIQUE NOT NULL,
    academic_year INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    registration_start DATE,
    registration_end DATE,
    is_current BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_semesters_current ON semesters(is_current);
CREATE INDEX IF NOT EXISTS idx_semesters_academic_year ON semesters(academic_year);
CREATE INDEX IF NOT EXISTS idx_semesters_dates ON semesters(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_semesters_code ON semesters(semester_code);

-- Create trigger for updated_at
CREATE TRIGGER update_semesters_updated_at
    BEFORE UPDATE ON semesters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
