-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Insert default roles
INSERT INTO roles (role_name, description) VALUES
    ('admin', 'System administrators with full access'),
    ('faculty', 'Teaching staff members'),
    ('student', 'Enrolled students'),
    ('staff', 'Non-teaching staff members')
ON CONFLICT (role_name) DO NOTHING;

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(role_name);
