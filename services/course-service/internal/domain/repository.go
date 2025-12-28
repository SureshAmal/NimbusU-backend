package domain

import (
	"context"

	"github.com/google/uuid"
)

// DepartmentRepository defines the interface for department data access
type DepartmentRepository interface {
	Create(ctx context.Context, department *Department) error
	GetByID(ctx context.Context, id uuid.UUID) (*Department, error)
	GetByCode(ctx context.Context, code string) (*Department, error)
	Update(ctx context.Context, department *Department) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter DepartmentFilter, limit, offset int) ([]*Department, int64, error)
	GetWithDetails(ctx context.Context, id uuid.UUID) (*DepartmentWithDetails, error)
}

// ProgramRepository defines the interface for program data access
type ProgramRepository interface {
	Create(ctx context.Context, program *Program) error
	GetByID(ctx context.Context, id uuid.UUID) (*Program, error)
	GetByCode(ctx context.Context, code string) (*Program, error)
	Update(ctx context.Context, program *Program) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ProgramFilter, limit, offset int) ([]*ProgramWithDepartment, int64, error)
	GetWithDepartment(ctx context.Context, id uuid.UUID) (*ProgramWithDepartment, error)
}

// SubjectRepository defines the interface for subject data access
type SubjectRepository interface {
	Create(ctx context.Context, subject *Subject) error
	GetByID(ctx context.Context, id uuid.UUID) (*Subject, error)
	GetByCode(ctx context.Context, code string) (*Subject, error)
	Update(ctx context.Context, subject *Subject) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter SubjectFilter, limit, offset int) ([]*Subject, int64, error)
	GetWithDetails(ctx context.Context, id uuid.UUID) (*SubjectWithDetails, error)
	AddPrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID, isMandatory bool) error
	RemovePrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID) error
	GetPrerequisites(ctx context.Context, subjectID uuid.UUID) ([]SubjectPrerequisite, error)
	AddCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error
	RemoveCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error
	GetCorequisites(ctx context.Context, subjectID uuid.UUID) ([]SubjectBasic, error)
}

// SemesterRepository defines the interface for semester data access
type SemesterRepository interface {
	Create(ctx context.Context, semester *Semester) error
	GetByID(ctx context.Context, id uuid.UUID) (*Semester, error)
	GetByCode(ctx context.Context, code string) (*Semester, error)
	GetCurrent(ctx context.Context) (*Semester, error)
	Update(ctx context.Context, semester *Semester) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter SemesterFilter, limit, offset int) ([]*Semester, int64, error)
	SetCurrent(ctx context.Context, id uuid.UUID) error
}

// CourseRepository defines the interface for course data access
type CourseRepository interface {
	Create(ctx context.Context, course *Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*Course, error)
	GetByCode(ctx context.Context, code string) (*Course, error)
	Update(ctx context.Context, course *Course) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter CourseFilter, limit, offset int) ([]*CourseWithDetails, int64, error)
	GetWithDetails(ctx context.Context, id uuid.UUID) (*CourseWithDetails, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	IncrementEnrollment(ctx context.Context, id uuid.UUID) error
	DecrementEnrollment(ctx context.Context, id uuid.UUID) error
}

// FacultyRepository defines the interface for faculty data access
type FacultyRepository interface {
	Create(ctx context.Context, faculty *Faculty) error
	GetByID(ctx context.Context, id uuid.UUID) (*Faculty, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Faculty, error)
	GetByEmployeeID(ctx context.Context, employeeID string) (*Faculty, error)
	Update(ctx context.Context, faculty *Faculty) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter FacultyFilter, limit, offset int) ([]*FacultyWithDetails, int64, error)
	GetWithDetails(ctx context.Context, id uuid.UUID) (*FacultyWithDetails, error)
}

// StudentRepository defines the interface for student data access
type StudentRepository interface {
	Create(ctx context.Context, student *Student) error
	GetByID(ctx context.Context, id uuid.UUID) (*Student, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Student, error)
	GetByRegistrationNumber(ctx context.Context, regNo string) (*Student, error)
	Update(ctx context.Context, student *Student) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter StudentFilter, limit, offset int) ([]*StudentWithDetails, int64, error)
	GetWithDetails(ctx context.Context, id uuid.UUID) (*StudentWithDetails, error)
	UpdateSemester(ctx context.Context, id uuid.UUID, semester int, cgpa *float64, credits int) error
}

// FacultyCourseRepository defines the interface for faculty-course assignments
type FacultyCourseRepository interface {
	Create(ctx context.Context, fc *FacultyCourse) error
	GetByID(ctx context.Context, id uuid.UUID) (*FacultyCourse, error)
	Update(ctx context.Context, fc *FacultyCourse) error
	Delete(ctx context.Context, facultyID, courseID uuid.UUID) error
	ListByCourse(ctx context.Context, courseID uuid.UUID) ([]*FacultyCourseWithDetails, error)
	ListByFaculty(ctx context.Context, facultyID uuid.UUID, semesterID *uuid.UUID) ([]*FacultyCourse, error)
	GetAssignment(ctx context.Context, facultyID, courseID uuid.UUID) (*FacultyCourse, error)
	IsPrimaryFaculty(ctx context.Context, facultyID, courseID uuid.UUID) (bool, error)
}

// EnrollmentRepository defines the interface for enrollment data access
type EnrollmentRepository interface {
	Create(ctx context.Context, enrollment *CourseEnrollment) error
	GetByID(ctx context.Context, id uuid.UUID) (*CourseEnrollment, error)
	GetByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*CourseEnrollment, error)
	Update(ctx context.Context, enrollment *CourseEnrollment) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByStudent(ctx context.Context, filter EnrollmentFilter, limit, offset int) ([]*EnrollmentWithDetails, int64, error)
	ListByCourse(ctx context.Context, courseID uuid.UUID, status *string, limit, offset int) ([]*EnrollmentWithDetails, int64, error)
	GetEnrollmentSummary(ctx context.Context, courseID uuid.UUID) (enrolled, waitlisted, dropped, completed int, err error)
	GetNextWaitlistPosition(ctx context.Context, courseID uuid.UUID) (int, error)
	PromoteFromWaitlist(ctx context.Context, courseID uuid.UUID) (*CourseEnrollment, error)
}

// CalendarRepository defines the interface for academic calendar data access
type CalendarRepository interface {
	Create(ctx context.Context, event *AcademicCalendarEvent) error
	GetByID(ctx context.Context, id uuid.UUID) (*AcademicCalendarEvent, error)
	Update(ctx context.Context, event *AcademicCalendarEvent) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter CalendarFilter, limit, offset int) ([]*AcademicCalendarEventWithDetails, int64, error)
}
