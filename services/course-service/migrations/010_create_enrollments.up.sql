-- 010_create_enrollments.up.sql
-- Create course enrollments table

CREATE TABLE IF NOT EXISTS course_enrollments (
    enrollment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    enrollment_status VARCHAR(20) DEFAULT 'enrolled' 
        CHECK(enrollment_status IN ('enrolled', 'waitlisted', 'dropped', 'completed', 'failed')),
    enrolled_by VARCHAR(20) DEFAULT 'self' CHECK(enrolled_by IN ('self', 'admin', 'system')),
    enrollment_date TIMESTAMPTZ DEFAULT now(),
    dropped_date TIMESTAMPTZ,
    drop_reason TEXT,
    completion_date TIMESTAMPTZ,
    grade VARCHAR(5),
    grade_points NUMERIC(3,2) CHECK(grade_points IS NULL OR grade_points BETWEEN 0 AND 10),
    waitlist_position INTEGER,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT unique_student_course UNIQUE (student_id, course_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_enrollments_student ON course_enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_course ON course_enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_status ON course_enrollments(enrollment_status);
CREATE INDEX IF NOT EXISTS idx_enrollments_date ON course_enrollments(enrollment_date);
CREATE INDEX IF NOT EXISTS idx_enrollments_waitlist ON course_enrollments(waitlist_position) WHERE waitlist_position IS NOT NULL;

-- Create trigger for updated_at
CREATE TRIGGER update_enrollments_updated_at
    BEFORE UPDATE ON course_enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create function to update course enrollment count
CREATE OR REPLACE FUNCTION update_course_enrollment_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        IF NEW.enrollment_status = 'enrolled' THEN
            UPDATE courses SET current_enrollment = current_enrollment + 1 WHERE course_id = NEW.course_id;
        END IF;
    ELSIF TG_OP = 'UPDATE' THEN
        IF OLD.enrollment_status != 'enrolled' AND NEW.enrollment_status = 'enrolled' THEN
            UPDATE courses SET current_enrollment = current_enrollment + 1 WHERE course_id = NEW.course_id;
        ELSIF OLD.enrollment_status = 'enrolled' AND NEW.enrollment_status != 'enrolled' THEN
            UPDATE courses SET current_enrollment = current_enrollment - 1 WHERE course_id = NEW.course_id;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        IF OLD.enrollment_status = 'enrolled' THEN
            UPDATE courses SET current_enrollment = current_enrollment - 1 WHERE course_id = OLD.course_id;
        END IF;
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_enrollment_count
    AFTER INSERT OR UPDATE OR DELETE ON course_enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_course_enrollment_count();
