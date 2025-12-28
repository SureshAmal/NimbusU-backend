-- 006_create_students.up.sql
-- Create students table (extends User Service users)

CREATE TABLE IF NOT EXISTS students (
    student_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL,
    registration_number VARCHAR(50) UNIQUE NOT NULL,
    roll_number VARCHAR(50),
    department_id UUID NOT NULL REFERENCES departments(department_id) ON DELETE RESTRICT,
    program_id UUID NOT NULL REFERENCES programs(program_id) ON DELETE RESTRICT,
    current_semester INTEGER NOT NULL CHECK(current_semester BETWEEN 1 AND 8),
    batch_year INTEGER NOT NULL,
    admission_date DATE,
    current_cgpa NUMERIC(4,2) CHECK(current_cgpa BETWEEN 0 AND 10),
    total_credits_earned INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_students_user ON students(user_id);
CREATE INDEX IF NOT EXISTS idx_students_registration ON students(registration_number);
CREATE INDEX IF NOT EXISTS idx_students_department ON students(department_id);
CREATE INDEX IF NOT EXISTS idx_students_program ON students(program_id);
CREATE INDEX IF NOT EXISTS idx_students_batch ON students(batch_year);
CREATE INDEX IF NOT EXISTS idx_students_semester ON students(current_semester);
CREATE INDEX IF NOT EXISTS idx_students_active ON students(is_active);

-- Create trigger for updated_at
CREATE TRIGGER update_students_updated_at
    BEFORE UPDATE ON students
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
