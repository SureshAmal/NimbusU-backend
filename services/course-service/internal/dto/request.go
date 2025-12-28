package dto

import (
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

// ==================== Department Requests ====================

type CreateDepartmentRequest struct {
	DepartmentName   string     `json:"department_name" binding:"required,max=100"`
	DepartmentCode   string     `json:"department_code" binding:"required,max=20"`
	HeadOfDepartment *uuid.UUID `json:"head_of_department"`
	Description      *string    `json:"description"`
}

type UpdateDepartmentRequest struct {
	DepartmentName   *string    `json:"department_name" binding:"omitempty,max=100"`
	HeadOfDepartment *uuid.UUID `json:"head_of_department"`
	Description      *string    `json:"description"`
	IsActive         *bool      `json:"is_active"`
}

// ==================== Program Requests ====================

type CreateProgramRequest struct {
	ProgramName   string    `json:"program_name" binding:"required,max=100"`
	ProgramCode   string    `json:"program_code" binding:"required,max=20"`
	DepartmentID  uuid.UUID `json:"department_id" binding:"required"`
	DegreeType    *string   `json:"degree_type" binding:"omitempty,oneof=Bachelor Master PhD"`
	DurationYears int       `json:"duration_years" binding:"required,min=1,max=10"`
	TotalCredits  *int      `json:"total_credits" binding:"omitempty,min=1"`
	Description   *string   `json:"description"`
}

type UpdateProgramRequest struct {
	ProgramName   *string `json:"program_name" binding:"omitempty,max=100"`
	DegreeType    *string `json:"degree_type" binding:"omitempty,oneof=Bachelor Master PhD"`
	DurationYears *int    `json:"duration_years" binding:"omitempty,min=1,max=10"`
	TotalCredits  *int    `json:"total_credits" binding:"omitempty,min=1"`
	Description   *string `json:"description"`
	IsActive      *bool   `json:"is_active"`
}

// ==================== Subject Requests ====================

type CreateSubjectRequest struct {
	SubjectName   string                `json:"subject_name" binding:"required,max=255"`
	SubjectCode   string                `json:"subject_code" binding:"required,max=20"`
	DepartmentID  uuid.UUID             `json:"department_id" binding:"required"`
	Credits       int                   `json:"credits" binding:"required,min=1,max=10"`
	SubjectType   *string               `json:"subject_type" binding:"omitempty,oneof=theory practical project"`
	Description   *string               `json:"description"`
	Syllabus      *string               `json:"syllabus"`
	Prerequisites []PrerequisiteRequest `json:"prerequisites"`
	Corequisites  []uuid.UUID           `json:"corequisites"`
}

type PrerequisiteRequest struct {
	SubjectID   uuid.UUID `json:"subject_id" binding:"required"`
	IsMandatory bool      `json:"is_mandatory"`
}

type UpdateSubjectRequest struct {
	SubjectName *string `json:"subject_name" binding:"omitempty,max=255"`
	Credits     *int    `json:"credits" binding:"omitempty,min=1,max=10"`
	SubjectType *string `json:"subject_type" binding:"omitempty,oneof=theory practical project"`
	Description *string `json:"description"`
	Syllabus    *string `json:"syllabus"`
	IsActive    *bool   `json:"is_active"`
}

// ==================== Semester Requests ====================

type CreateSemesterRequest struct {
	SemesterName      string     `json:"semester_name" binding:"required,max=50"`
	SemesterCode      string     `json:"semester_code" binding:"required,max=20"`
	AcademicYear      int        `json:"academic_year" binding:"required,min=2000,max=2100"`
	StartDate         time.Time  `json:"start_date" binding:"required"`
	EndDate           time.Time  `json:"end_date" binding:"required,gtfield=StartDate"`
	RegistrationStart *time.Time `json:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end"`
}

type UpdateSemesterRequest struct {
	SemesterName      *string    `json:"semester_name" binding:"omitempty,max=50"`
	StartDate         *time.Time `json:"start_date"`
	EndDate           *time.Time `json:"end_date"`
	RegistrationStart *time.Time `json:"registration_start"`
	RegistrationEnd   *time.Time `json:"registration_end"`
}

// ==================== Course Requests ====================

type CreateCourseRequest struct {
	CourseCode     string     `json:"course_code" binding:"required,max=20"`
	CourseName     string     `json:"course_name" binding:"required,max=255"`
	SubjectID      uuid.UUID  `json:"subject_id" binding:"required"`
	DepartmentID   uuid.UUID  `json:"department_id" binding:"required"`
	ProgramID      *uuid.UUID `json:"program_id"`
	SemesterID     uuid.UUID  `json:"semester_id" binding:"required"`
	SemesterNumber int        `json:"semester_number" binding:"required,min=1,max=8"`
	AcademicYear   int        `json:"academic_year" binding:"required,min=2000,max=2100"`
	MaxStudents    *int       `json:"max_students" binding:"omitempty,min=1"`
	Description    *string    `json:"description"`
}

type UpdateCourseRequest struct {
	CourseName  *string `json:"course_name" binding:"omitempty,max=255"`
	MaxStudents *int    `json:"max_students" binding:"omitempty,min=1"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=draft active completed cancelled"`
}

// ==================== Faculty Assignment Requests ====================

type AssignFacultyRequest struct {
	FacultyID uuid.UUID `json:"faculty_id" binding:"required"`
	Role      string    `json:"role" binding:"required,oneof=instructor co-instructor teaching_assistant"`
	IsPrimary bool      `json:"is_primary"`
}

type UpdateFacultyAssignmentRequest struct {
	Role      *string `json:"role" binding:"omitempty,oneof=instructor co-instructor teaching_assistant"`
	IsPrimary *bool   `json:"is_primary"`
}

// ==================== Enrollment Requests ====================

type EnrollRequest struct {
	StudentID *uuid.UUID `json:"student_id"`
}

type DropCourseRequest struct {
	StudentID *uuid.UUID `json:"student_id"`
	Reason    string     `json:"reason"`
}

type UpdateEnrollmentRequest struct {
	EnrollmentStatus string   `json:"enrollment_status" binding:"required,oneof=enrolled waitlisted dropped completed failed"`
	Grade            *string  `json:"grade" binding:"omitempty,max=5"`
	GradePoints      *float64 `json:"grade_points" binding:"omitempty,min=0,max=10"`
}

type BulkEnrollRequest struct {
	StudentIDs        []uuid.UUID `json:"student_ids" binding:"required,min=1"`
	SkipPrerequisites bool        `json:"skip_prerequisites"`
}

// ==================== Faculty Requests ====================

type CreateFacultyRequest struct {
	UserID         uuid.UUID  `json:"user_id" binding:"required"`
	EmployeeID     string     `json:"employee_id" binding:"required,max=50"`
	DepartmentID   uuid.UUID  `json:"department_id" binding:"required"`
	Designation    *string    `json:"designation" binding:"omitempty,max=100"`
	Qualification  *string    `json:"qualification" binding:"omitempty,max=255"`
	Specialization *string    `json:"specialization"`
	JoiningDate    *time.Time `json:"joining_date"`
	OfficeRoom     *string    `json:"office_room" binding:"omitempty,max=50"`
	OfficeHours    *string    `json:"office_hours"`
}

type UpdateFacultyRequest struct {
	Designation    *string `json:"designation" binding:"omitempty,max=100"`
	Qualification  *string `json:"qualification" binding:"omitempty,max=255"`
	Specialization *string `json:"specialization"`
	OfficeRoom     *string `json:"office_room" binding:"omitempty,max=50"`
	OfficeHours    *string `json:"office_hours"`
	IsActive       *bool   `json:"is_active"`
}

// ==================== Student Requests ====================

type CreateStudentRequest struct {
	UserID             uuid.UUID  `json:"user_id" binding:"required"`
	RegistrationNumber string     `json:"registration_number" binding:"required,max=50"`
	RollNumber         *string    `json:"roll_number" binding:"omitempty,max=50"`
	DepartmentID       uuid.UUID  `json:"department_id" binding:"required"`
	ProgramID          uuid.UUID  `json:"program_id" binding:"required"`
	CurrentSemester    int        `json:"current_semester" binding:"required,min=1,max=8"`
	BatchYear          int        `json:"batch_year" binding:"required,min=2000,max=2100"`
	AdmissionDate      *time.Time `json:"admission_date"`
}

type UpdateStudentRequest struct {
	RollNumber *string `json:"roll_number" binding:"omitempty,max=50"`
	IsActive   *bool   `json:"is_active"`
}

type PromoteStudentRequest struct {
	NewSemester   int      `json:"new_semester" binding:"required,min=1,max=8"`
	UpdateCGPA    *float64 `json:"update_cgpa" binding:"omitempty,min=0,max=10"`
	CreditsEarned int      `json:"credits_earned" binding:"min=0"`
}

// ==================== Calendar Requests ====================

type CreateCalendarEventRequest struct {
	SemesterID  uuid.UUID  `json:"semester_id" binding:"required"`
	EventName   string     `json:"event_name" binding:"required,max=255"`
	EventType   string     `json:"event_type" binding:"required,oneof=holiday exam registration deadline event other"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	Description *string    `json:"description"`
	IsHoliday   bool       `json:"is_holiday"`
}

type UpdateCalendarEventRequest struct {
	EventName   *string    `json:"event_name" binding:"omitempty,max=255"`
	EventType   *string    `json:"event_type" binding:"omitempty,oneof=holiday exam registration deadline event other"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description *string    `json:"description"`
	IsHoliday   *bool      `json:"is_holiday"`
}

// ==================== Query Parameters ====================

type PaginationParams struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
}

type DepartmentQueryParams struct {
	PaginationParams
	IsActive *bool `form:"is_active"`
}

type ProgramQueryParams struct {
	PaginationParams
	DepartmentID *uuid.UUID `form:"department_id"`
	DegreeType   *string    `form:"degree_type"`
	IsActive     *bool      `form:"is_active"`
}

type SubjectQueryParams struct {
	PaginationParams
	DepartmentID *uuid.UUID `form:"department_id"`
	SubjectType  *string    `form:"subject_type"`
	Credits      *int       `form:"credits"`
	Search       *string    `form:"search"`
	IsActive     *bool      `form:"is_active"`
}

type SemesterQueryParams struct {
	PaginationParams
	AcademicYear *int  `form:"academic_year"`
	IsCurrent    *bool `form:"is_current"`
}

type CourseQueryParams struct {
	PaginationParams
	DepartmentID   *uuid.UUID `form:"department_id"`
	ProgramID      *uuid.UUID `form:"program_id"`
	SemesterID     *uuid.UUID `form:"semester_id"`
	SubjectID      *uuid.UUID `form:"subject_id"`
	FacultyID      *uuid.UUID `form:"faculty_id"`
	SemesterNumber *int       `form:"semester_number"`
	AcademicYear   *int       `form:"academic_year"`
	Status         *string    `form:"status"`
	Search         *string    `form:"search"`
	IsActive       *bool      `form:"is_active"`
}

type FacultyQueryParams struct {
	PaginationParams
	DepartmentID *uuid.UUID `form:"department_id"`
	Designation  *string    `form:"designation"`
	Search       *string    `form:"search"`
	IsActive     *bool      `form:"is_active"`
}

type StudentQueryParams struct {
	PaginationParams
	DepartmentID    *uuid.UUID `form:"department_id"`
	ProgramID       *uuid.UUID `form:"program_id"`
	CurrentSemester *int       `form:"current_semester"`
	BatchYear       *int       `form:"batch_year"`
	Search          *string    `form:"search"`
	IsActive        *bool      `form:"is_active"`
}

type EnrollmentQueryParams struct {
	PaginationParams
	SemesterID *uuid.UUID `form:"semester_id"`
	Status     *string    `form:"status"`
}

type CalendarQueryParams struct {
	PaginationParams
	SemesterID *uuid.UUID `form:"semester_id"`
	EventType  *string    `form:"event_type"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
	IsHoliday  *bool      `form:"is_holiday"`
}

// ==================== Enroll Student Request ====================

type EnrollStudentRequest struct {
	StudentID uuid.UUID `json:"student_id" validate:"required"`
}

// ==================== ToDomain Methods ====================

func (r *CreateDepartmentRequest) ToDomain() *domain.Department {
	return &domain.Department{
		DepartmentName:   r.DepartmentName,
		DepartmentCode:   r.DepartmentCode,
		HeadOfDepartment: r.HeadOfDepartment,
		Description:      r.Description,
		IsActive:         true,
	}
}

func (r *UpdateDepartmentRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.DepartmentName != nil {
		updates["department_name"] = *r.DepartmentName
	}
	if r.HeadOfDepartment != nil {
		updates["head_of_department"] = r.HeadOfDepartment.String()
	}
	if r.Description != nil {
		updates["description"] = *r.Description
	}
	if r.IsActive != nil {
		updates["is_active"] = *r.IsActive
	}
	return updates
}

func (r *CreateProgramRequest) ToDomain() *domain.Program {
	return &domain.Program{
		ProgramName:   r.ProgramName,
		ProgramCode:   r.ProgramCode,
		DepartmentID:  r.DepartmentID,
		DegreeType:    r.DegreeType,
		DurationYears: r.DurationYears,
		TotalCredits:  r.TotalCredits,
		Description:   r.Description,
		IsActive:      true,
	}
}

func (r *UpdateProgramRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.ProgramName != nil {
		updates["program_name"] = *r.ProgramName
	}
	if r.DegreeType != nil {
		updates["degree_type"] = *r.DegreeType
	}
	if r.DurationYears != nil {
		updates["duration_years"] = *r.DurationYears
	}
	if r.TotalCredits != nil {
		updates["total_credits"] = *r.TotalCredits
	}
	if r.Description != nil {
		updates["description"] = *r.Description
	}
	if r.IsActive != nil {
		updates["is_active"] = *r.IsActive
	}
	return updates
}

func (r *CreateSubjectRequest) ToDomain() *domain.Subject {
	return &domain.Subject{
		SubjectName:  r.SubjectName,
		SubjectCode:  r.SubjectCode,
		DepartmentID: r.DepartmentID,
		Credits:      r.Credits,
		SubjectType:  r.SubjectType,
		Description:  r.Description,
		Syllabus:     r.Syllabus,
		IsActive:     true,
	}
}

func (r *UpdateSubjectRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.SubjectName != nil {
		updates["subject_name"] = *r.SubjectName
	}
	if r.Credits != nil {
		updates["credits"] = *r.Credits
	}
	if r.SubjectType != nil {
		updates["subject_type"] = *r.SubjectType
	}
	if r.Description != nil {
		updates["description"] = *r.Description
	}
	if r.Syllabus != nil {
		updates["syllabus"] = *r.Syllabus
	}
	if r.IsActive != nil {
		updates["is_active"] = *r.IsActive
	}
	return updates
}

func (r *CreateSemesterRequest) ToDomain() *domain.Semester {
	return &domain.Semester{
		SemesterName:      r.SemesterName,
		SemesterCode:      r.SemesterCode,
		AcademicYear:      r.AcademicYear,
		StartDate:         r.StartDate,
		EndDate:           r.EndDate,
		RegistrationStart: r.RegistrationStart,
		RegistrationEnd:   r.RegistrationEnd,
		IsCurrent:         false,
	}
}

func (r *UpdateSemesterRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.SemesterName != nil {
		updates["semester_name"] = *r.SemesterName
	}
	if r.StartDate != nil {
		updates["start_date"] = *r.StartDate
	}
	if r.EndDate != nil {
		updates["end_date"] = *r.EndDate
	}
	if r.RegistrationStart != nil {
		updates["registration_start"] = *r.RegistrationStart
	}
	if r.RegistrationEnd != nil {
		updates["registration_end"] = *r.RegistrationEnd
	}
	return updates
}

func (r *CreateCourseRequest) ToDomain() *domain.Course {
	return &domain.Course{
		CourseCode:        r.CourseCode,
		CourseName:        r.CourseName,
		SubjectID:         r.SubjectID,
		DepartmentID:      r.DepartmentID,
		ProgramID:         r.ProgramID,
		SemesterID:        r.SemesterID,
		SemesterNumber:    r.SemesterNumber,
		AcademicYear:      r.AcademicYear,
		MaxStudents:       r.MaxStudents,
		CurrentEnrollment: 0,
		Status:            "draft",
		Description:       r.Description,
		IsActive:          true,
	}
}

func (r *UpdateCourseRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.CourseName != nil {
		updates["course_name"] = *r.CourseName
	}
	if r.MaxStudents != nil {
		updates["max_students"] = *r.MaxStudents
	}
	if r.Description != nil {
		updates["description"] = *r.Description
	}
	if r.Status != nil {
		updates["status"] = *r.Status
	}
	return updates
}

func (r *CreateFacultyRequest) ToDomain() *domain.Faculty {
	return &domain.Faculty{
		UserID:         r.UserID,
		EmployeeID:     r.EmployeeID,
		DepartmentID:   r.DepartmentID,
		Designation:    r.Designation,
		Qualification:  r.Qualification,
		Specialization: r.Specialization,
		JoiningDate:    r.JoiningDate,
		OfficeRoom:     r.OfficeRoom,
		OfficeHours:    r.OfficeHours,
		IsActive:       true,
	}
}

func (r *UpdateFacultyRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.Designation != nil {
		updates["designation"] = *r.Designation
	}
	if r.Qualification != nil {
		updates["qualification"] = *r.Qualification
	}
	if r.Specialization != nil {
		updates["specialization"] = *r.Specialization
	}
	if r.OfficeRoom != nil {
		updates["office_room"] = *r.OfficeRoom
	}
	if r.OfficeHours != nil {
		updates["office_hours"] = *r.OfficeHours
	}
	if r.IsActive != nil {
		updates["is_active"] = *r.IsActive
	}
	return updates
}

func (r *CreateStudentRequest) ToDomain() *domain.Student {
	return &domain.Student{
		UserID:             r.UserID,
		RegistrationNumber: r.RegistrationNumber,
		RollNumber:         r.RollNumber,
		DepartmentID:       r.DepartmentID,
		ProgramID:          r.ProgramID,
		CurrentSemester:    r.CurrentSemester,
		BatchYear:          r.BatchYear,
		AdmissionDate:      r.AdmissionDate,
		TotalCreditsEarned: 0,
		IsActive:           true,
	}
}

func (r *UpdateStudentRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.RollNumber != nil {
		updates["roll_number"] = *r.RollNumber
	}
	if r.IsActive != nil {
		updates["is_active"] = *r.IsActive
	}
	return updates
}

func (r *CreateCalendarEventRequest) ToDomain() *domain.AcademicCalendarEvent {
	return &domain.AcademicCalendarEvent{
		SemesterID:  r.SemesterID,
		EventName:   r.EventName,
		EventType:   r.EventType,
		StartDate:   r.StartDate,
		EndDate:     r.EndDate,
		Description: r.Description,
		IsHoliday:   r.IsHoliday,
	}
}

func (r *UpdateCalendarEventRequest) ToUpdates() map[string]interface{} {
	updates := make(map[string]interface{})
	if r.EventName != nil {
		updates["event_name"] = *r.EventName
	}
	if r.EventType != nil {
		updates["event_type"] = *r.EventType
	}
	if r.StartDate != nil {
		updates["start_date"] = *r.StartDate
	}
	if r.EndDate != nil {
		updates["end_date"] = *r.EndDate
	}
	if r.Description != nil {
		updates["description"] = *r.Description
	}
	if r.IsHoliday != nil {
		updates["is_holiday"] = *r.IsHoliday
	}
	return updates
}
