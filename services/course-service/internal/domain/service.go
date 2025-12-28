package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Domain errors
var (
	// Not found errors
	ErrDepartmentNotFound    = errors.New("department not found")
	ErrProgramNotFound       = errors.New("program not found")
	ErrSubjectNotFound       = errors.New("subject not found")
	ErrSemesterNotFound      = errors.New("semester not found")
	ErrCourseNotFound        = errors.New("course not found")
	ErrFacultyNotFound       = errors.New("faculty not found")
	ErrStudentNotFound       = errors.New("student not found")
	ErrEnrollmentNotFound    = errors.New("enrollment not found")
	ErrCalendarEventNotFound = errors.New("calendar event not found")
	ErrAssignmentNotFound    = errors.New("faculty assignment not found")

	// Duplicate errors
	ErrDepartmentCodeExists     = errors.New("department code already exists")
	ErrProgramCodeExists        = errors.New("program code already exists")
	ErrSubjectCodeExists        = errors.New("subject code already exists")
	ErrSemesterCodeExists       = errors.New("semester code already exists")
	ErrCourseCodeExists         = errors.New("course code already exists")
	ErrEmployeeIDExists         = errors.New("employee ID already exists")
	ErrRegistrationNumberExists = errors.New("registration number already exists")
	ErrAlreadyEnrolled          = errors.New("student already enrolled in this course")
	ErrFacultyAlreadyAssigned   = errors.New("faculty already assigned to this course")

	// Business logic errors
	ErrCourseFull                  = errors.New("course has reached maximum enrollment")
	ErrRegistrationClosed          = errors.New("course registration is closed")
	ErrPrerequisitesNotMet         = errors.New("course prerequisites not met")
	ErrCannotDropCompletedCourse   = errors.New("cannot drop a completed course")
	ErrCannotModifyCompletedCourse = errors.New("cannot modify a completed course")
	ErrInvalidEnrollmentStatus     = errors.New("invalid enrollment status transition")
	ErrInvalidCourseStatus         = errors.New("invalid course status transition")
	ErrNoCurrentSemester           = errors.New("no current semester is set")
	ErrSelfPrerequisite            = errors.New("subject cannot be its own prerequisite")

	// Permission errors
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("forbidden: insufficient permissions")
)

// DepartmentService defines the interface for department business logic
type DepartmentService interface {
	CreateDepartment(ctx context.Context, department *Department) error
	GetDepartment(ctx context.Context, id uuid.UUID) (*DepartmentWithDetails, error)
	UpdateDepartment(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	ListDepartments(ctx context.Context, filter DepartmentFilter, page, limit int) ([]*Department, int64, error)
}

// ProgramService defines the interface for program business logic
type ProgramService interface {
	CreateProgram(ctx context.Context, program *Program) error
	GetProgram(ctx context.Context, id uuid.UUID) (*ProgramWithDepartment, error)
	UpdateProgram(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteProgram(ctx context.Context, id uuid.UUID) error
	ListPrograms(ctx context.Context, filter ProgramFilter, page, limit int) ([]*ProgramWithDepartment, int64, error)
}

// SubjectService defines the interface for subject business logic
type SubjectService interface {
	CreateSubject(ctx context.Context, subject *Subject, prerequisites []SubjectPrerequisite, corequisites []uuid.UUID) error
	GetSubject(ctx context.Context, id uuid.UUID) (*SubjectWithDetails, error)
	UpdateSubject(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteSubject(ctx context.Context, id uuid.UUID) error
	ListSubjects(ctx context.Context, filter SubjectFilter, page, limit int) ([]*Subject, int64, error)
	AddPrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID, isMandatory bool) error
	RemovePrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID) error
	AddCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error
	RemoveCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error
}

// SemesterService defines the interface for semester business logic
type SemesterService interface {
	CreateSemester(ctx context.Context, semester *Semester) error
	GetSemester(ctx context.Context, id uuid.UUID) (*Semester, error)
	GetCurrentSemester(ctx context.Context) (*Semester, error)
	UpdateSemester(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteSemester(ctx context.Context, id uuid.UUID) error
	ListSemesters(ctx context.Context, filter SemesterFilter, page, limit int) ([]*Semester, int64, error)
	SetCurrentSemester(ctx context.Context, id uuid.UUID) error
}

// CourseService defines the interface for course business logic
type CourseService interface {
	CreateCourse(ctx context.Context, course *Course) error
	GetCourse(ctx context.Context, id uuid.UUID) (*CourseWithDetails, error)
	UpdateCourse(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	ListCourses(ctx context.Context, filter CourseFilter, page, limit int) ([]*CourseWithDetails, int64, error)
	ActivateCourse(ctx context.Context, id uuid.UUID) error
	DeactivateCourse(ctx context.Context, id uuid.UUID) error
	GetCourseStudents(ctx context.Context, courseID uuid.UUID, status *string, page, limit int) ([]*EnrollmentWithDetails, int64, error)
}

// FacultyService defines the interface for faculty business logic
type FacultyService interface {
	CreateFaculty(ctx context.Context, faculty *Faculty) error
	GetFaculty(ctx context.Context, id uuid.UUID) (*FacultyWithDetails, error)
	GetFacultyByUserID(ctx context.Context, userID uuid.UUID) (*FacultyWithDetails, error)
	UpdateFaculty(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteFaculty(ctx context.Context, id uuid.UUID) error
	ListFaculty(ctx context.Context, filter FacultyFilter, page, limit int) ([]*FacultyWithDetails, int64, error)
	GetFacultyCourses(ctx context.Context, facultyID uuid.UUID, semesterID *uuid.UUID, page, limit int) ([]*FacultyCourse, int64, error)
}

// StudentService defines the interface for student business logic
type StudentService interface {
	CreateStudent(ctx context.Context, student *Student) error
	GetStudent(ctx context.Context, id uuid.UUID) (*StudentWithDetails, error)
	GetStudentByUserID(ctx context.Context, userID uuid.UUID) (*StudentWithDetails, error)
	UpdateStudent(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteStudent(ctx context.Context, id uuid.UUID) error
	ListStudents(ctx context.Context, filter StudentFilter, page, limit int) ([]*StudentWithDetails, int64, error)
	PromoteStudent(ctx context.Context, id uuid.UUID, newSemester int, cgpa *float64, creditsEarned int) error
}

// FacultyAssignmentService defines the interface for faculty-course assignment business logic
type FacultyAssignmentService interface {
	AssignFaculty(ctx context.Context, courseID, facultyID, assignedBy uuid.UUID, role string, isPrimary bool) (*FacultyCourse, error)
	UpdateAssignment(ctx context.Context, courseID, facultyID uuid.UUID, role string, isPrimary bool) error
	RemoveFaculty(ctx context.Context, courseID, facultyID uuid.UUID) error
	ListCourseFaculty(ctx context.Context, courseID uuid.UUID) ([]*FacultyCourseWithDetails, error)
}

// EnrollmentService defines the interface for enrollment business logic
type EnrollmentService interface {
	EnrollStudent(ctx context.Context, courseID, studentID uuid.UUID, enrolledBy string) (*CourseEnrollment, error)
	DropCourse(ctx context.Context, courseID, studentID uuid.UUID, reason string) error
	UpdateEnrollment(ctx context.Context, enrollmentID uuid.UUID, status string, grade *string, gradePoints *float64) error
	GetStudentEnrollments(ctx context.Context, studentID uuid.UUID, filter EnrollmentFilter, page, limit int) ([]*EnrollmentWithDetails, int64, error)
	BulkEnroll(ctx context.Context, courseID uuid.UUID, studentIDs []uuid.UUID, skipPrerequisites bool) ([]BulkEnrollResult, error)
	CheckPrerequisites(ctx context.Context, studentID, courseID uuid.UUID) (bool, []SubjectBasic, error)
}

// BulkEnrollResult represents the result of a bulk enrollment operation
type BulkEnrollResult struct {
	StudentID    uuid.UUID  `json:"student_id"`
	EnrollmentID *uuid.UUID `json:"enrollment_id,omitempty"`
	Status       string     `json:"status"`
	Error        string     `json:"error,omitempty"`
}

// CalendarService defines the interface for academic calendar business logic
type CalendarService interface {
	CreateEvent(ctx context.Context, event *AcademicCalendarEvent) error
	GetEvent(ctx context.Context, id uuid.UUID) (*AcademicCalendarEventWithDetails, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error
	ListEvents(ctx context.Context, filter CalendarFilter, page, limit int) ([]*AcademicCalendarEventWithDetails, int64, error)
}

// EventProducer defines the interface for publishing events to Kafka
type EventProducer interface {
	PublishEvent(topic string, key string, event interface{}) error
	Close() error
}
