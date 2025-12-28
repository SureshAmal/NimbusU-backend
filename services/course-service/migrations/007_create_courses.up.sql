-- 007_create_courses.up.sql
-- Create courses table

CREATE TABLE IF NOT EXISTS courses (
    course_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_code VARCHAR(20) UNIQUE NOT NULL,
    course_name VARCHAR(255) NOT NULL,
    subject_id UUID NOT NULL REFERENCES subjects(subject_id) ON DELETE RESTRICT,
    department_id UUID NOT NULL REFERENCES departments(department_id) ON DELETE RESTRICT,
    program_id UUID REFERENCES programs(program_id) ON DELETE SET NULL,
    semester_id UUID NOT NULL REFERENCES semesters(semester_id) ON DELETE RESTRICT,
    semester_number INTEGER NOT NULL CHECK(semester_number BETWEEN 1 AND 8),
    academic_year INTEGER NOT NULL,
    max_students INTEGER,
    current_enrollment INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft' CHECK(status IN ('draft', 'active', 'completed', 'cancelled')),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_courses_code ON courses(course_code);
CREATE INDEX IF NOT EXISTS idx_courses_subject ON courses(subject_id);
CREATE INDEX IF NOT EXISTS idx_courses_department ON courses(department_id);
CREATE INDEX IF NOT EXISTS idx_courses_program ON courses(program_id);
CREATE INDEX IF NOT EXISTS idx_courses_semester ON courses(semester_id);
CREATE INDEX IF NOT EXISTS idx_courses_status ON courses(status);
CREATE INDEX IF NOT EXISTS idx_courses_academic_year ON courses(academic_year);
CREATE INDEX IF NOT EXISTS idx_courses_active ON courses(is_active);
CREATE INDEX IF NOT EXISTS idx_courses_semester_number ON courses(semester_number);

-- Create trigger for updated_at
CREATE TRIGGER update_courses_updated_at
    BEFORE UPDATE ON courses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
