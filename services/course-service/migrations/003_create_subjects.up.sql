-- 003_create_subjects.up.sql
-- Create subjects table

CREATE TABLE IF NOT EXISTS subjects (
    subject_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_name VARCHAR(255) NOT NULL,
    subject_code VARCHAR(20) UNIQUE NOT NULL,
    department_id UUID NOT NULL REFERENCES departments(department_id) ON DELETE RESTRICT,
    credits INTEGER NOT NULL,
    subject_type VARCHAR(50) CHECK(subject_type IN ('theory', 'practical', 'project')),
    description TEXT,
    syllabus TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_subjects_code ON subjects(subject_code);
CREATE INDEX IF NOT EXISTS idx_subjects_department ON subjects(department_id);
CREATE INDEX IF NOT EXISTS idx_subjects_active ON subjects(is_active);
CREATE INDEX IF NOT EXISTS idx_subjects_type ON subjects(subject_type);

-- Create trigger for updated_at
CREATE TRIGGER update_subjects_updated_at
    BEFORE UPDATE ON subjects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
