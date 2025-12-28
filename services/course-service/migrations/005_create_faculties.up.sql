-- 005_create_faculties.up.sql
-- Create faculties table (extends User Service users)

CREATE TABLE IF NOT EXISTS faculties (
    faculty_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    department_id UUID NOT NULL REFERENCES departments(department_id) ON DELETE RESTRICT,
    designation VARCHAR(100),
    qualification VARCHAR(255),
    specialization TEXT,
    joining_date DATE,
    office_room VARCHAR(50),
    office_hours TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_faculties_user ON faculties(user_id);
CREATE INDEX IF NOT EXISTS idx_faculties_employee ON faculties(employee_id);
CREATE INDEX IF NOT EXISTS idx_faculties_department ON faculties(department_id);
CREATE INDEX IF NOT EXISTS idx_faculties_active ON faculties(is_active);

-- Add foreign key from departments to faculties for head_of_department
ALTER TABLE departments
    ADD CONSTRAINT fk_departments_head
    FOREIGN KEY (head_of_department)
    REFERENCES faculties(faculty_id)
    ON DELETE SET NULL;

-- Create trigger for updated_at
CREATE TRIGGER update_faculties_updated_at
    BEFORE UPDATE ON faculties
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
