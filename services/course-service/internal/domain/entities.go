package domain

import (
	"time"

	"github.com/google/uuid"
)

// Department represents an academic department
type Department struct {
	DepartmentID     uuid.UUID  `json:"department_id" db:"department_id"`
	DepartmentName   string     `json:"department_name" db:"department_name"`
	DepartmentCode   string     `json:"department_code" db:"department_code"`
	HeadOfDepartment *uuid.UUID `json:"head_of_department,omitempty" db:"head_of_department"`
	Description      *string    `json:"description,omitempty" db:"description"`
	IsActive         bool       `json:"is_active" db:"is_active"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// DepartmentWithDetails includes head faculty info and counts
type DepartmentWithDetails struct {
	Department
	HeadFaculty   *FacultyBasic `json:"head_faculty,omitempty"`
	ProgramsCount int           `json:"programs_count"`
	FacultyCount  int           `json:"faculty_count"`
	StudentsCount int           `json:"students_count"`
}

// Program represents a degree program
type Program struct {
	ProgramID     uuid.UUID `json:"program_id" db:"program_id"`
	ProgramName   string    `json:"program_name" db:"program_name"`
	ProgramCode   string    `json:"program_code" db:"program_code"`
	DepartmentID  uuid.UUID `json:"department_id" db:"department_id"`
	DegreeType    *string   `json:"degree_type,omitempty" db:"degree_type"`
	DurationYears int       `json:"duration_years" db:"duration_years"`
	TotalCredits  *int      `json:"total_credits,omitempty" db:"total_credits"`
	Description   *string   `json:"description,omitempty" db:"description"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ProgramWithDepartment includes department info
type ProgramWithDepartment struct {
	Program
	Department DepartmentBasic `json:"department"`
}

// Subject represents an academic subject
type Subject struct {
	SubjectID    uuid.UUID `json:"subject_id" db:"subject_id"`
	SubjectName  string    `json:"subject_name" db:"subject_name"`
	SubjectCode  string    `json:"subject_code" db:"subject_code"`
	DepartmentID uuid.UUID `json:"department_id" db:"department_id"`
	Credits      int       `json:"credits" db:"credits"`
	SubjectType  *string   `json:"subject_type,omitempty" db:"subject_type"`
	Description  *string   `json:"description,omitempty" db:"description"`
	Syllabus     *string   `json:"syllabus,omitempty" db:"syllabus"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SubjectWithDetails includes department and prerequisites
type SubjectWithDetails struct {
	Subject
	Department    DepartmentBasic       `json:"department"`
	Prerequisites []SubjectPrerequisite `json:"prerequisites,omitempty"`
	Corequisites  []SubjectBasic        `json:"corequisites,omitempty"`
}

// SubjectPrerequisite represents a prerequisite relationship
type SubjectPrerequisite struct {
	PrerequisiteID        uuid.UUID `json:"prerequisite_id" db:"prerequisite_id"`
	SubjectID             uuid.UUID `json:"subject_id" db:"subject_id"`
	PrerequisiteSubjectID uuid.UUID `json:"prerequisite_subject_id" db:"prerequisite_subject_id"`
	IsMandatory           bool      `json:"is_mandatory" db:"is_mandatory"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	// Joined fields
	SubjectCode string `json:"subject_code,omitempty" db:"prereq_subject_code"`
	SubjectName string `json:"subject_name,omitempty" db:"prereq_subject_name"`
}

// SubjectCorequisite represents a corequisite relationship
type SubjectCorequisite struct {
	CorequisiteID        uuid.UUID `json:"corequisite_id" db:"corequisite_id"`
	SubjectID            uuid.UUID `json:"subject_id" db:"subject_id"`
	CorequisiteSubjectID uuid.UUID `json:"corequisite_subject_id" db:"corequisite_subject_id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

// Semester represents an academic semester
type Semester struct {
	SemesterID        uuid.UUID  `json:"semester_id" db:"semester_id"`
	SemesterName      string     `json:"semester_name" db:"semester_name"`
	SemesterCode      string     `json:"semester_code" db:"semester_code"`
	AcademicYear      int        `json:"academic_year" db:"academic_year"`
	StartDate         time.Time  `json:"start_date" db:"start_date"`
	EndDate           time.Time  `json:"end_date" db:"end_date"`
	RegistrationStart *time.Time `json:"registration_start,omitempty" db:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end,omitempty" db:"registration_end"`
	IsCurrent         bool       `json:"is_current" db:"is_current"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// Faculty represents a faculty member
type Faculty struct {
	FacultyID      uuid.UUID  `json:"faculty_id" db:"faculty_id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	EmployeeID     string     `json:"employee_id" db:"employee_id"`
	DepartmentID   uuid.UUID  `json:"department_id" db:"department_id"`
	Designation    *string    `json:"designation,omitempty" db:"designation"`
	Qualification  *string    `json:"qualification,omitempty" db:"qualification"`
	Specialization *string    `json:"specialization,omitempty" db:"specialization"`
	JoiningDate    *time.Time `json:"joining_date,omitempty" db:"joining_date"`
	OfficeRoom     *string    `json:"office_room,omitempty" db:"office_room"`
	OfficeHours    *string    `json:"office_hours,omitempty" db:"office_hours"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// FacultyWithDetails includes department and user info
type FacultyWithDetails struct {
	Faculty
	Department     DepartmentBasic `json:"department"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Phone          *string         `json:"phone,omitempty"`
	CurrentCourses []CourseBasic   `json:"current_courses,omitempty"`
}

// Student represents a student
type Student struct {
	StudentID          uuid.UUID  `json:"student_id" db:"student_id"`
	UserID             uuid.UUID  `json:"user_id" db:"user_id"`
	RegistrationNumber string     `json:"registration_number" db:"registration_number"`
	RollNumber         *string    `json:"roll_number,omitempty" db:"roll_number"`
	DepartmentID       uuid.UUID  `json:"department_id" db:"department_id"`
	ProgramID          uuid.UUID  `json:"program_id" db:"program_id"`
	CurrentSemester    int        `json:"current_semester" db:"current_semester"`
	BatchYear          int        `json:"batch_year" db:"batch_year"`
	AdmissionDate      *time.Time `json:"admission_date,omitempty" db:"admission_date"`
	CurrentCGPA        *float64   `json:"current_cgpa,omitempty" db:"current_cgpa"`
	TotalCreditsEarned int        `json:"total_credits_earned" db:"total_credits_earned"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// StudentWithDetails includes department, program, and user info
type StudentWithDetails struct {
	Student
	Department         DepartmentBasic   `json:"department"`
	Program            ProgramBasic      `json:"program"`
	Name               string            `json:"name"`
	Email              string            `json:"email"`
	Phone              *string           `json:"phone,omitempty"`
	CurrentEnrollments []EnrollmentBasic `json:"current_enrollments,omitempty"`
}

// Course represents a course instance
type Course struct {
	CourseID          uuid.UUID  `json:"course_id" db:"course_id"`
	CourseCode        string     `json:"course_code" db:"course_code"`
	CourseName        string     `json:"course_name" db:"course_name"`
	SubjectID         uuid.UUID  `json:"subject_id" db:"subject_id"`
	DepartmentID      uuid.UUID  `json:"department_id" db:"department_id"`
	ProgramID         *uuid.UUID `json:"program_id,omitempty" db:"program_id"`
	SemesterID        uuid.UUID  `json:"semester_id" db:"semester_id"`
	SemesterNumber    int        `json:"semester_number" db:"semester_number"`
	AcademicYear      int        `json:"academic_year" db:"academic_year"`
	MaxStudents       *int       `json:"max_students,omitempty" db:"max_students"`
	CurrentEnrollment int        `json:"current_enrollment" db:"current_enrollment"`
	Status            string     `json:"status" db:"status"`
	Description       *string    `json:"description,omitempty" db:"description"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	CreatedBy         uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// CourseWithDetails includes all related info
type CourseWithDetails struct {
	Course
	Subject       SubjectBasic         `json:"subject"`
	Department    DepartmentBasic      `json:"department"`
	Program       *ProgramBasic        `json:"program,omitempty"`
	Semester      SemesterBasic        `json:"semester"`
	Faculty       []FacultyCourseBasic `json:"faculty,omitempty"`
	Prerequisites []SubjectBasic       `json:"prerequisites,omitempty"`
}

// FacultyCourse represents a faculty-course assignment
type FacultyCourse struct {
	FacultyCourseID uuid.UUID  `json:"faculty_course_id" db:"faculty_course_id"`
	FacultyID       uuid.UUID  `json:"faculty_id" db:"faculty_id"`
	CourseID        uuid.UUID  `json:"course_id" db:"course_id"`
	Role            string     `json:"role" db:"role"`
	IsPrimary       bool       `json:"is_primary" db:"is_primary"`
	AssignedBy      uuid.UUID  `json:"assigned_by" db:"assigned_by"`
	AssignedAt      time.Time  `json:"assigned_at" db:"assigned_at"`
	RemovedAt       *time.Time `json:"removed_at,omitempty" db:"removed_at"`
	IsActive        bool       `json:"is_active" db:"is_active"`
}

// FacultyCourseWithDetails includes faculty info
type FacultyCourseWithDetails struct {
	FacultyCourse
	Faculty FacultyBasic `json:"faculty"`
}

// CourseEnrollment represents a student enrollment
type CourseEnrollment struct {
	EnrollmentID     uuid.UUID  `json:"enrollment_id" db:"enrollment_id"`
	StudentID        uuid.UUID  `json:"student_id" db:"student_id"`
	CourseID         uuid.UUID  `json:"course_id" db:"course_id"`
	EnrollmentStatus string     `json:"enrollment_status" db:"enrollment_status"`
	EnrolledBy       string     `json:"enrolled_by" db:"enrolled_by"`
	EnrollmentDate   time.Time  `json:"enrollment_date" db:"enrollment_date"`
	DroppedDate      *time.Time `json:"dropped_date,omitempty" db:"dropped_date"`
	DropReason       *string    `json:"drop_reason,omitempty" db:"drop_reason"`
	CompletionDate   *time.Time `json:"completion_date,omitempty" db:"completion_date"`
	Grade            *string    `json:"grade,omitempty" db:"grade"`
	GradePoints      *float64   `json:"grade_points,omitempty" db:"grade_points"`
	WaitlistPosition *int       `json:"waitlist_position,omitempty" db:"waitlist_position"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// EnrollmentWithDetails includes student and course info
type EnrollmentWithDetails struct {
	CourseEnrollment
	Student  StudentBasic   `json:"student"`
	Course   CourseBasic    `json:"course"`
	Semester *SemesterBasic `json:"semester,omitempty"`
}

// AcademicCalendarEvent represents a calendar event
type AcademicCalendarEvent struct {
	EventID     uuid.UUID  `json:"event_id" db:"event_id"`
	SemesterID  uuid.UUID  `json:"semester_id" db:"semester_id"`
	EventName   string     `json:"event_name" db:"event_name"`
	EventType   string     `json:"event_type" db:"event_type"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty" db:"end_date"`
	Description *string    `json:"description,omitempty" db:"description"`
	IsHoliday   bool       `json:"is_holiday" db:"is_holiday"`
	CreatedBy   uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// AcademicCalendarEventWithDetails includes semester info
type AcademicCalendarEventWithDetails struct {
	AcademicCalendarEvent
	Semester SemesterBasic `json:"semester"`
}

// ========== Basic/Summary Types for Embedding ==========

// DepartmentBasic is a minimal department representation
type DepartmentBasic struct {
	DepartmentID   uuid.UUID `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	DepartmentCode string    `json:"department_code"`
}

// ProgramBasic is a minimal program representation
type ProgramBasic struct {
	ProgramID     uuid.UUID `json:"program_id"`
	ProgramName   string    `json:"program_name"`
	ProgramCode   string    `json:"program_code"`
	DurationYears int       `json:"duration_years,omitempty"`
}

// SubjectBasic is a minimal subject representation
type SubjectBasic struct {
	SubjectID   uuid.UUID `json:"subject_id"`
	SubjectCode string    `json:"subject_code"`
	SubjectName string    `json:"subject_name"`
	Credits     int       `json:"credits,omitempty"`
	SubjectType *string   `json:"subject_type,omitempty"`
}

// SemesterBasic is a minimal semester representation
type SemesterBasic struct {
	SemesterID   uuid.UUID  `json:"semester_id"`
	SemesterName string     `json:"semester_name"`
	SemesterCode string     `json:"semester_code,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
}

// FacultyBasic is a minimal faculty representation
type FacultyBasic struct {
	FacultyID   uuid.UUID `json:"faculty_id"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	EmployeeID  string    `json:"employee_id,omitempty"`
	Name        string    `json:"name"`
	Email       string    `json:"email,omitempty"`
	Designation *string   `json:"designation,omitempty"`
}

// FacultyCourseBasic is for embedding in course details
type FacultyCourseBasic struct {
	FacultyID   uuid.UUID `json:"faculty_id"`
	Name        string    `json:"name"`
	Designation *string   `json:"designation,omitempty"`
	Role        string    `json:"role"`
	IsPrimary   bool      `json:"is_primary"`
}

// StudentBasic is a minimal student representation
type StudentBasic struct {
	StudentID          uuid.UUID `json:"student_id"`
	UserID             uuid.UUID `json:"user_id,omitempty"`
	RegistrationNumber string    `json:"registration_number"`
	Name               string    `json:"name"`
	Email              string    `json:"email,omitempty"`
}

// CourseBasic is a minimal course representation
type CourseBasic struct {
	CourseID   uuid.UUID `json:"course_id"`
	CourseCode string    `json:"course_code"`
	CourseName string    `json:"course_name"`
	Credits    int       `json:"credits,omitempty"`
	Role       string    `json:"role,omitempty"`
}

// EnrollmentBasic is a minimal enrollment representation
type EnrollmentBasic struct {
	CourseID   uuid.UUID `json:"course_id"`
	CourseCode string    `json:"course_code"`
	CourseName string    `json:"course_name"`
	Credits    int       `json:"credits"`
}

// ========== Filter Types ==========

// DepartmentFilter for filtering departments
type DepartmentFilter struct {
	IsActive *bool
}

// ProgramFilter for filtering programs
type ProgramFilter struct {
	DepartmentID *uuid.UUID
	DegreeType   *string
	IsActive     *bool
}

// SubjectFilter for filtering subjects
type SubjectFilter struct {
	DepartmentID *uuid.UUID
	SubjectType  *string
	Credits      *int
	Search       *string
	IsActive     *bool
}

// SemesterFilter for filtering semesters
type SemesterFilter struct {
	AcademicYear *int
	IsCurrent    *bool
}

// CourseFilter for filtering courses
type CourseFilter struct {
	DepartmentID   *uuid.UUID
	ProgramID      *uuid.UUID
	SemesterID     *uuid.UUID
	SubjectID      *uuid.UUID
	FacultyID      *uuid.UUID
	SemesterNumber *int
	AcademicYear   *int
	Status         *string
	Search         *string
	IsActive       *bool
}

// FacultyFilter for filtering faculty
type FacultyFilter struct {
	DepartmentID *uuid.UUID
	Designation  *string
	Search       *string
	IsActive     *bool
}

// StudentFilter for filtering students
type StudentFilter struct {
	DepartmentID    *uuid.UUID
	ProgramID       *uuid.UUID
	CurrentSemester *int
	BatchYear       *int
	Search          *string
	IsActive        *bool
}

// EnrollmentFilter for filtering enrollments
type EnrollmentFilter struct {
	StudentID  *uuid.UUID
	CourseID   *uuid.UUID
	SemesterID *uuid.UUID
	Status     *string
}

// CalendarFilter for filtering calendar events
type CalendarFilter struct {
	SemesterID *uuid.UUID
	EventType  *string
	StartDate  *time.Time
	EndDate    *time.Time
	IsHoliday  *bool
}
