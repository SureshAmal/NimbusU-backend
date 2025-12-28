-- ============================================
-- NimbusU Database Seed Data
-- ============================================
-- This script creates default data for development and testing
-- Run with: psql "your_database_url" -f seed.sql

\echo '======================================'
\echo 'Seeding NimbusU Database'
\echo '======================================'

-- ----------------
-- Default Admin User
-- ----------------
\echo 'Creating default admin user...'

-- Insert admin user
-- Password: Admin@123
INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    1000,
    'admin@nimbusu.edu',
    '$2a$10$6LzyBpz8Rc0/R3X5HcavSu5WhGO0OiTxkacg6LNvv4h994fAiYCcS',
    (SELECT role_id FROM roles WHERE role_name = 'admin' LIMIT 1),
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Insert admin profile
INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, last_name, phone, gender, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    (SELECT user_id FROM users WHERE email = 'admin@nimbusu.edu'),
    1000,
    'System',
    'Administrator',
    '+1-555-0100',
    'other',
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;

\echo 'Admin user created: admin@nimbusu.edu'

-- ----------------
-- Sample Faculty Users
-- ----------------
\echo 'Creating sample faculty users...'

-- Faculty 1
INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    2001,
    'john.doe@nimbusu.edu',
    '$2a$10$rKq/QWw5xEgD.S6h5EQh2e9CblAtZpy7oxz5hxKCHcy5lEpnTnoZC', -- Faculty@123
    (SELECT role_id FROM roles WHERE role_name = 'faculty' LIMIT 1),
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, last_name, phone, gender, bio, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    (SELECT user_id FROM users WHERE email = 'john.doe@nimbusu.edu'),
    2001,
    'John',
    'Doe',
    '+1-555-0101',
    'male',
    'Professor of Computer Science with 15 years of experience',
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;

-- Faculty 2
INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    2002,
    'jane.smith@nimbusu.edu',
    '$2a$10$rKq/QWw5xEgD.S6h5EQh2e9CblAtZpy7oxz5hxKCHcy5lEpnTnoZC', -- Faculty@123
    (SELECT role_id FROM roles WHERE role_name = 'faculty' LIMIT 1),
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, last_name, phone, gender, bio, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    (SELECT user_id FROM users WHERE email = 'jane.smith@nimbusu.edu'),
    2002,
    'Jane',
    'Smith',
    '+1-555-0102',
    'female',
    'Associate Professor of Mathematics',
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;

\echo 'Sample faculty users created'

-- ----------------
-- Sample Student Users
-- ----------------
\echo 'Creating sample student users...'

-- Student 1
INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    3001,
    'alice.johnson@student.nimbusu.edu',
    '$2a$10$8BRrmD2ix8vWopt1VBWvCejNz0lPROYH5xxS4fWyXZCCpnk11eIWG', -- Student@123
    (SELECT role_id FROM roles WHERE role_name = 'student' LIMIT 1),
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, last_name, phone, gender, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    (SELECT user_id FROM users WHERE email = 'alice.johnson@student.nimbusu.edu'),
    3001,
    'Alice',
    'Johnson',
    '+1-555-0201',
    'female',
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;

-- Student 2
INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    3002,
    'bob.williams@student.nimbusu.edu',
    '$2a$10$8BRrmD2ix8vWopt1VBWvCejNz0lPROYH5xxS4fWyXZCCpnk11eIWG', -- Student@123
    (SELECT role_id FROM roles WHERE role_name = 'student' LIMIT 1),
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, last_name, phone, gender, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    (SELECT user_id FROM users WHERE email = 'bob.williams@student.nimbusu.edu'),
    3002,
    'Bob',
    'Williams',
    '+1-555-0202',
    'male',
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;

-- Student 3
INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    3003,
    'carol.davis@student.nimbusu.edu',
    '$2a$10$8BRrmD2ix8vWopt1VBWvCejNz0lPROYH5xxS4fWyXZCCpnk11eIWG', -- Student@123
    (SELECT role_id FROM roles WHERE role_name = 'student' LIMIT 1),
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, last_name, phone, gender, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    (SELECT user_id FROM users WHERE email = 'carol.davis@student.nimbusu.edu'),
    3003,
    'Carol',
    'Davis',
    '+1-555-0203',
    'female',
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;

\echo 'Sample student users created'

-- ----------------
-- Sample Staff Users
-- ----------------
\echo 'Creating sample staff users...'

-- Staff 1
INSERT INTO users (user_id, register_no, email, password_hash, role_id, status, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    4001,
    'david.brown@nimbusu.edu',
    '$2a$10$3eZsyUlhaz9t3ihmfjzRV.DnnPHx.f.oH1kULBHPa5.YuxTdz5/VW', -- Staff@123
    (SELECT role_id FROM roles WHERE role_name = 'staff' LIMIT 1),
    'active',
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO user_profiles (profile_id, user_id, register_no, first_name, last_name, phone, gender, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    (SELECT user_id FROM users WHERE email = 'david.brown@nimbusu.edu'),
    4001,
    'David',
    'Brown',
    '+1-555-0301',
    'male',
    NOW(),
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;

\echo 'Sample staff users created'

-- ----------------
-- Summary
-- ----------------
\echo ''
\echo '======================================'
\echo 'Seed Data Summary'
\echo '======================================'
\echo 'Default Credentials (for all users):'
\echo ''
\echo 'Admin:'
\echo '  Email: admin@nimbusu.edu'
\echo '  Password: Admin@123'
\echo ''
\echo 'Faculty:'
\echo '  Email: john.doe@nimbusu.edu'
\echo '  Email: jane.smith@nimbusu.edu'
\echo '  Password: Faculty@123'
\echo ''
\echo 'Students:'
\echo '  Email: alice.johnson@student.nimbusu.edu'
\echo '  Email: bob.williams@student.nimbusu.edu'
\echo '  Email: carol.davis@student.nimbusu.edu'
\echo '  Password: Student@123'
\echo ''
\echo 'Staff:'
\echo '  Email: david.brown@nimbusu.edu'
\echo '  Password: Staff@123'
\echo '======================================'
