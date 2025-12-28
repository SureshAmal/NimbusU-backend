-- Seed data for NimbusU Course Service
-- Run this after applying all migrations

-- ==================== Departments ====================
INSERT INTO departments (department_id, department_name, department_code, description, is_active) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Computer Science and Engineering', 'CSE', 'Department of Computer Science and Engineering focusing on software, algorithms, and systems.', true),
    ('22222222-2222-2222-2222-222222222222', 'Electronics and Communication', 'ECE', 'Department of Electronics and Communication Engineering.', true),
    ('33333333-3333-3333-3333-333333333333', 'Mechanical Engineering', 'ME', 'Department of Mechanical Engineering.', true),
    ('44444444-4444-4444-4444-444444444444', 'Civil Engineering', 'CE', 'Department of Civil Engineering.', true),
    ('55555555-5555-5555-5555-555555555555', 'Electrical Engineering', 'EE', 'Department of Electrical Engineering.', true)
ON CONFLICT (department_code) DO NOTHING;

-- ==================== Programs ====================
INSERT INTO programs (program_id, program_name, program_code, department_id, degree_type, duration_years, total_credits, description, is_active) VALUES
    ('aaaa1111-1111-1111-1111-111111111111', 'Bachelor of Technology in Computer Science', 'BTECH-CSE', '11111111-1111-1111-1111-111111111111', 'Bachelor', 4, 160, 'Undergraduate program in Computer Science and Engineering.', true),
    ('aaaa2222-2222-2222-2222-222222222222', 'Master of Technology in Computer Science', 'MTECH-CSE', '11111111-1111-1111-1111-111111111111', 'Master', 2, 60, 'Postgraduate program in Computer Science.', true),
    ('aaaa3333-3333-3333-3333-333333333333', 'Bachelor of Technology in Electronics', 'BTECH-ECE', '22222222-2222-2222-2222-222222222222', 'Bachelor', 4, 160, 'Undergraduate program in Electronics and Communication.', true),
    ('aaaa4444-4444-4444-4444-444444444444', 'Bachelor of Technology in Mechanical', 'BTECH-ME', '33333333-3333-3333-3333-333333333333', 'Bachelor', 4, 160, 'Undergraduate program in Mechanical Engineering.', true),
    ('aaaa5555-5555-5555-5555-555555555555', 'Bachelor of Technology in Civil', 'BTECH-CE', '44444444-4444-4444-4444-444444444444', 'Bachelor', 4, 160, 'Undergraduate program in Civil Engineering.', true)
ON CONFLICT (program_code) DO NOTHING;

-- ==================== Subjects ====================
INSERT INTO subjects (subject_id, subject_name, subject_code, department_id, credits, subject_type, description, is_active) VALUES
    -- CSE Subjects
    ('bbbb1111-1111-1111-1111-111111111111', 'Introduction to Programming', 'CS101', '11111111-1111-1111-1111-111111111111', 4, 'theory', 'Fundamentals of programming using Python and C.', true),
    ('bbbb1112-1111-1111-1111-111111111111', 'Data Structures', 'CS201', '11111111-1111-1111-1111-111111111111', 4, 'theory', 'Arrays, linked lists, trees, graphs, and algorithms.', true),
    ('bbbb1113-1111-1111-1111-111111111111', 'Algorithms', 'CS301', '11111111-1111-1111-1111-111111111111', 4, 'theory', 'Design and analysis of algorithms.', true),
    ('bbbb1114-1111-1111-1111-111111111111', 'Database Systems', 'CS302', '11111111-1111-1111-1111-111111111111', 4, 'theory', 'Relational databases, SQL, and database design.', true),
    ('bbbb1115-1111-1111-1111-111111111111', 'Operating Systems', 'CS303', '11111111-1111-1111-1111-111111111111', 4, 'theory', 'Process management, memory, and file systems.', true),
    ('bbbb1116-1111-1111-1111-111111111111', 'Computer Networks', 'CS401', '11111111-1111-1111-1111-111111111111', 3, 'theory', 'Network protocols and architectures.', true),
    ('bbbb1117-1111-1111-1111-111111111111', 'Machine Learning', 'CS402', '11111111-1111-1111-1111-111111111111', 3, 'theory', 'Introduction to machine learning algorithms.', true),
    ('bbbb1118-1111-1111-1111-111111111111', 'Programming Lab', 'CS101L', '11111111-1111-1111-1111-111111111111', 2, 'practical', 'Lab work for programming fundamentals.', true),
    -- ECE Subjects
    ('bbbb2111-2222-2222-2222-222222222222', 'Basic Electronics', 'EC101', '22222222-2222-2222-2222-222222222222', 4, 'theory', 'Introduction to electronic devices and circuits.', true),
    ('bbbb2112-2222-2222-2222-222222222222', 'Digital Electronics', 'EC201', '22222222-2222-2222-2222-222222222222', 4, 'theory', 'Digital circuits and logic design.', true),
    -- Common subjects
    ('bbbb9111-9999-9999-9999-999999999999', 'Engineering Mathematics I', 'MA101', '11111111-1111-1111-1111-111111111111', 4, 'theory', 'Calculus and linear algebra.', true),
    ('bbbb9112-9999-9999-9999-999999999999', 'Engineering Mathematics II', 'MA102', '11111111-1111-1111-1111-111111111111', 4, 'theory', 'Differential equations and statistics.', true),
    ('bbbb9113-9999-9999-9999-999999999999', 'Engineering Physics', 'PH101', '11111111-1111-1111-1111-111111111111', 3, 'theory', 'Fundamentals of physics for engineers.', true)
ON CONFLICT (subject_code) DO NOTHING;

-- ==================== Subject Prerequisites ====================
INSERT INTO subject_prerequisites (prerequisite_id, subject_id, prerequisite_subject_id, is_mandatory) VALUES
    -- Data Structures requires Intro to Programming
    (gen_random_uuid(), 'bbbb1112-1111-1111-1111-111111111111', 'bbbb1111-1111-1111-1111-111111111111', true),
    -- Algorithms requires Data Structures
    (gen_random_uuid(), 'bbbb1113-1111-1111-1111-111111111111', 'bbbb1112-1111-1111-1111-111111111111', true),
    -- Database Systems requires Intro to Programming
    (gen_random_uuid(), 'bbbb1114-1111-1111-1111-111111111111', 'bbbb1111-1111-1111-1111-111111111111', true),
    -- Machine Learning requires Algorithms
    (gen_random_uuid(), 'bbbb1117-1111-1111-1111-111111111111', 'bbbb1113-1111-1111-1111-111111111111', true)
ON CONFLICT DO NOTHING;

-- ==================== Semesters ====================
INSERT INTO semesters (semester_id, semester_name, semester_code, academic_year, start_date, end_date, registration_start, registration_end, is_current) VALUES
    ('cccc1111-1111-1111-1111-111111111111', 'Fall 2024', 'FALL-2024', 2024, '2024-08-01', '2024-12-15', '2024-07-01', '2024-07-31', false),
    ('cccc2222-2222-2222-2222-222222222222', 'Spring 2025', 'SPRING-2025', 2025, '2025-01-15', '2025-05-30', '2024-12-15', '2025-01-10', true),
    ('cccc3333-3333-3333-3333-333333333333', 'Fall 2025', 'FALL-2025', 2025, '2025-08-01', '2025-12-15', '2025-07-01', '2025-07-31', false)
ON CONFLICT (semester_code) DO NOTHING;

-- ==================== Academic Calendar Events ====================
INSERT INTO academic_calendar (event_id, semester_id, event_name, event_type, start_date, end_date, description, is_holiday, created_by) VALUES
    -- Spring 2025 events
    (gen_random_uuid(), 'cccc2222-2222-2222-2222-222222222222', 'Semester Start', 'event', '2025-01-15', NULL, 'First day of Spring 2025 semester.', false, '00000000-0000-0000-0000-000000000001'),
    (gen_random_uuid(), 'cccc2222-2222-2222-2222-222222222222', 'Republic Day', 'holiday', '2025-01-26', NULL, 'National holiday.', true, '00000000-0000-0000-0000-000000000001'),
    (gen_random_uuid(), 'cccc2222-2222-2222-2222-222222222222', 'Mid-semester Exams', 'exam', '2025-03-10', '2025-03-20', 'Mid-semester examination period.', false, '00000000-0000-0000-0000-000000000001'),
    (gen_random_uuid(), 'cccc2222-2222-2222-2222-222222222222', 'Spring Break', 'holiday', '2025-03-21', '2025-03-30', 'Spring vacation.', true, '00000000-0000-0000-0000-000000000001'),
    (gen_random_uuid(), 'cccc2222-2222-2222-2222-222222222222', 'End-semester Exams', 'exam', '2025-05-15', '2025-05-28', 'Final examination period.', false, '00000000-0000-0000-0000-000000000001'),
    (gen_random_uuid(), 'cccc2222-2222-2222-2222-222222222222', 'Semester End', 'event', '2025-05-30', NULL, 'Last day of Spring 2025 semester.', false, '00000000-0000-0000-0000-000000000001'),
    (gen_random_uuid(), 'cccc2222-2222-2222-2222-222222222222', 'Course Registration Deadline', 'deadline', '2025-01-25', NULL, 'Last date to add/drop courses.', false, '00000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

-- Note: Faculties, Students, Courses, and Enrollments should be created through the API
-- or inserted after the corresponding users are created in the user-service.

-- Sample message
SELECT 'Seed data inserted successfully!' as message;
