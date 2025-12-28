-- 008_create_prerequisites.up.sql
-- Create course prerequisites and corequisites tables

CREATE TABLE IF NOT EXISTS course_prerequisites (
    prerequisite_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(subject_id) ON DELETE CASCADE,
    prerequisite_subject_id UUID NOT NULL REFERENCES subjects(subject_id) ON DELETE CASCADE,
    is_mandatory BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT check_different_subjects CHECK (subject_id != prerequisite_subject_id),
    CONSTRAINT unique_prerequisite UNIQUE (subject_id, prerequisite_subject_id)
);

CREATE TABLE IF NOT EXISTS course_corequisites (
    corequisite_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id UUID NOT NULL REFERENCES subjects(subject_id) ON DELETE CASCADE,
    corequisite_subject_id UUID NOT NULL REFERENCES subjects(subject_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT check_different_corequisites CHECK (subject_id != corequisite_subject_id),
    CONSTRAINT unique_corequisite UNIQUE (subject_id, corequisite_subject_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_prerequisites_subject ON course_prerequisites(subject_id);
CREATE INDEX IF NOT EXISTS idx_prerequisites_prereq ON course_prerequisites(prerequisite_subject_id);
CREATE INDEX IF NOT EXISTS idx_corequisites_subject ON course_corequisites(subject_id);
CREATE INDEX IF NOT EXISTS idx_corequisites_coreq ON course_corequisites(corequisite_subject_id);
