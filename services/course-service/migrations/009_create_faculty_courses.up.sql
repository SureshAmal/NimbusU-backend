-- 009_create_faculty_courses.up.sql
-- Create faculty courses assignment table

CREATE TABLE IF NOT EXISTS faculty_courses (
    faculty_course_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    faculty_id UUID NOT NULL REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'instructor' CHECK(role IN ('instructor', 'co-instructor', 'teaching_assistant')),
    is_primary BOOLEAN DEFAULT false,
    assigned_by UUID NOT NULL,
    assigned_at TIMESTAMPTZ DEFAULT now(),
    removed_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_faculty_courses_faculty ON faculty_courses(faculty_id);
CREATE INDEX IF NOT EXISTS idx_faculty_courses_course ON faculty_courses(course_id);
CREATE INDEX IF NOT EXISTS idx_faculty_courses_active ON faculty_courses(is_active);
CREATE INDEX IF NOT EXISTS idx_faculty_courses_role ON faculty_courses(role);

-- Unique constraint for active assignments
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_active_faculty_course 
    ON faculty_courses(faculty_id, course_id) 
    WHERE is_active = true;
