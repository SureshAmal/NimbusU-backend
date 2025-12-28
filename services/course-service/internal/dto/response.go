package dto

import (
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

// ==================== Department Responses ====================

type DepartmentResponse struct {
	DepartmentID     uuid.UUID             `json:"department_id"`
	DepartmentName   string                `json:"department_name"`
	DepartmentCode   string                `json:"department_code"`
	HeadOfDepartment *FacultyBasicResponse `json:"head_of_department,omitempty"`
	Description      *string               `json:"description,omitempty"`
	IsActive         bool                  `json:"is_active"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

type DepartmentDetailResponse struct {
	DepartmentResponse
	ProgramsCount int `json:"programs_count"`
	FacultyCount  int `json:"faculty_count"`
	StudentsCount int `json:"students_count"`
}

func ToDepartmentResponse(d *domain.Department) DepartmentResponse {
	return DepartmentResponse{
		DepartmentID:   d.DepartmentID,
		DepartmentName: d.DepartmentName,
		DepartmentCode: d.DepartmentCode,
		Description:    d.Description,
		IsActive:       d.IsActive,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

func ToDepartmentDetailResponse(d *domain.DepartmentWithDetails) DepartmentDetailResponse {
	resp := DepartmentDetailResponse{
		DepartmentResponse: ToDepartmentResponse(&d.Department),
		ProgramsCount:      d.ProgramsCount,
		FacultyCount:       d.FacultyCount,
		StudentsCount:      d.StudentsCount,
	}
	if d.HeadFaculty != nil {
		head := ToFacultyBasicResponse(d.HeadFaculty)
		resp.HeadOfDepartment = &head
	}
	return resp
}

// ==================== Program Responses ====================

type ProgramResponse struct {
	ProgramID     uuid.UUID               `json:"program_id"`
	ProgramName   string                  `json:"program_name"`
	ProgramCode   string                  `json:"program_code"`
	Department    DepartmentBasicResponse `json:"department"`
	DegreeType    *string                 `json:"degree_type,omitempty"`
	DurationYears int                     `json:"duration_years"`
	TotalCredits  *int                    `json:"total_credits,omitempty"`
	Description   *string                 `json:"description,omitempty"`
	IsActive      bool                    `json:"is_active"`
	CreatedAt     time.Time               `json:"created_at"`
}

func ToProgramResponse(p *domain.ProgramWithDepartment) ProgramResponse {
	return ProgramResponse{
		ProgramID:     p.ProgramID,
		ProgramName:   p.ProgramName,
		ProgramCode:   p.ProgramCode,
		Department:    ToDepartmentBasicResponse(&p.Department),
		DegreeType:    p.DegreeType,
		DurationYears: p.DurationYears,
		TotalCredits:  p.TotalCredits,
		Description:   p.Description,
		IsActive:      p.IsActive,
		CreatedAt:     p.CreatedAt,
	}
}

// ==================== Subject Responses ====================

type SubjectResponse struct {
	SubjectID   uuid.UUID               `json:"subject_id"`
	SubjectName string                  `json:"subject_name"`
	SubjectCode string                  `json:"subject_code"`
	Department  DepartmentBasicResponse `json:"department"`
	Credits     int                     `json:"credits"`
	SubjectType *string                 `json:"subject_type,omitempty"`
	IsActive    bool                    `json:"is_active"`
}

type SubjectDetailResponse struct {
	SubjectResponse
	Description   *string                `json:"description,omitempty"`
	Syllabus      *string                `json:"syllabus,omitempty"`
	Prerequisites []PrerequisiteResponse `json:"prerequisites"`
	Corequisites  []SubjectBasicResponse `json:"corequisites"`
	CreatedAt     time.Time              `json:"created_at"`
}

type PrerequisiteResponse struct {
	SubjectID   uuid.UUID `json:"subject_id"`
	SubjectCode string    `json:"subject_code"`
	SubjectName string    `json:"subject_name"`
	IsMandatory bool      `json:"is_mandatory"`
}

func ToSubjectDetailResponse(s *domain.SubjectWithDetails) SubjectDetailResponse {
	prereqs := make([]PrerequisiteResponse, len(s.Prerequisites))
	for i, p := range s.Prerequisites {
		prereqs[i] = PrerequisiteResponse{
			SubjectID:   p.PrerequisiteSubjectID,
			SubjectCode: p.SubjectCode,
			SubjectName: p.SubjectName,
			IsMandatory: p.IsMandatory,
		}
	}

	coreqs := make([]SubjectBasicResponse, len(s.Corequisites))
	for i, c := range s.Corequisites {
		coreqs[i] = ToSubjectBasicResponse(&c)
	}

	return SubjectDetailResponse{
		SubjectResponse: SubjectResponse{
			SubjectID:   s.SubjectID,
			SubjectName: s.SubjectName,
			SubjectCode: s.SubjectCode,
			Department:  ToDepartmentBasicResponse(&s.Department),
			Credits:     s.Credits,
			SubjectType: s.SubjectType,
			IsActive:    s.IsActive,
		},
		Description:   s.Description,
		Syllabus:      s.Syllabus,
		Prerequisites: prereqs,
		Corequisites:  coreqs,
		CreatedAt:     s.CreatedAt,
	}
}

// ==================== Semester Responses ====================

type SemesterResponse struct {
	SemesterID        uuid.UUID  `json:"semester_id"`
	SemesterName      string     `json:"semester_name"`
	SemesterCode      string     `json:"semester_code"`
	AcademicYear      int        `json:"academic_year"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	RegistrationStart *time.Time `json:"registration_start,omitempty"`
	RegistrationEnd   *time.Time `json:"registration_end,omitempty"`
	IsCurrent         bool       `json:"is_current"`
}

type CurrentSemesterResponse struct {
	SemesterResponse
	DaysRemaining    int  `json:"days_remaining"`
	RegistrationOpen bool `json:"registration_open"`
}

func ToSemesterResponse(s *domain.Semester) SemesterResponse {
	return SemesterResponse{
		SemesterID:        s.SemesterID,
		SemesterName:      s.SemesterName,
		SemesterCode:      s.SemesterCode,
		AcademicYear:      s.AcademicYear,
		StartDate:         s.StartDate,
		EndDate:           s.EndDate,
		RegistrationStart: s.RegistrationStart,
		RegistrationEnd:   s.RegistrationEnd,
		IsCurrent:         s.IsCurrent,
	}
}

// ==================== Course Responses ====================

type CourseResponse struct {
	CourseID          uuid.UUID               `json:"course_id"`
	CourseCode        string                  `json:"course_code"`
	CourseName        string                  `json:"course_name"`
	Subject           SubjectBasicResponse    `json:"subject"`
	Department        DepartmentBasicResponse `json:"department"`
	Semester          SemesterBasicResponse   `json:"semester"`
	SemesterNumber    int                     `json:"semester_number"`
	AcademicYear      int                     `json:"academic_year"`
	Faculty           []FacultyCourseResponse `json:"faculty"`
	MaxStudents       *int                    `json:"max_students,omitempty"`
	CurrentEnrollment int                     `json:"current_enrollment"`
	Status            string                  `json:"status"`
	IsActive          bool                    `json:"is_active"`
}

type CourseDetailResponse struct {
	CourseResponse
	Program        *ProgramBasicResponse  `json:"program,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Prerequisites  []SubjectBasicResponse `json:"prerequisites,omitempty"`
	AvailableSeats int                    `json:"available_seats"`
	CreatedBy      uuid.UUID              `json:"created_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func ToCourseResponse(c *domain.CourseWithDetails) CourseResponse {
	faculty := make([]FacultyCourseResponse, len(c.Faculty))
	for i, f := range c.Faculty {
		faculty[i] = FacultyCourseResponse{
			FacultyID:   f.FacultyID,
			Name:        f.Name,
			Designation: f.Designation,
			Role:        f.Role,
			IsPrimary:   f.IsPrimary,
		}
	}

	return CourseResponse{
		CourseID:          c.CourseID,
		CourseCode:        c.CourseCode,
		CourseName:        c.CourseName,
		Subject:           ToSubjectBasicResponse(&c.Subject),
		Department:        ToDepartmentBasicResponse(&c.Department),
		Semester:          ToSemesterBasicResponse(&c.Semester),
		SemesterNumber:    c.SemesterNumber,
		AcademicYear:      c.AcademicYear,
		Faculty:           faculty,
		MaxStudents:       c.MaxStudents,
		CurrentEnrollment: c.CurrentEnrollment,
		Status:            c.Status,
		IsActive:          c.IsActive,
	}
}

func ToCourseDetailResponse(c *domain.CourseWithDetails) CourseDetailResponse {
	prereqs := make([]SubjectBasicResponse, len(c.Prerequisites))
	for i, p := range c.Prerequisites {
		prereqs[i] = ToSubjectBasicResponse(&p)
	}

	availableSeats := 0
	if c.MaxStudents != nil {
		availableSeats = *c.MaxStudents - c.CurrentEnrollment
		if availableSeats < 0 {
			availableSeats = 0
		}
	}

	resp := CourseDetailResponse{
		CourseResponse: ToCourseResponse(c),
		Description:    c.Description,
		Prerequisites:  prereqs,
		AvailableSeats: availableSeats,
		CreatedBy:      c.CreatedBy,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}

	if c.Program != nil {
		prog := ToProgramBasicResponse(c.Program)
		resp.Program = &prog
	}

	return resp
}

// ==================== Faculty Responses ====================

type FacultyResponse struct {
	FacultyID      uuid.UUID               `json:"faculty_id"`
	UserID         uuid.UUID               `json:"user_id"`
	EmployeeID     string                  `json:"employee_id"`
	Name           string                  `json:"name"`
	Email          string                  `json:"email"`
	Department     DepartmentBasicResponse `json:"department"`
	Designation    *string                 `json:"designation,omitempty"`
	Specialization *string                 `json:"specialization,omitempty"`
	IsActive       bool                    `json:"is_active"`
}

type FacultyDetailResponse struct {
	FacultyResponse
	Phone          *string               `json:"phone,omitempty"`
	Qualification  *string               `json:"qualification,omitempty"`
	JoiningDate    *time.Time            `json:"joining_date,omitempty"`
	OfficeRoom     *string               `json:"office_room,omitempty"`
	OfficeHours    *string               `json:"office_hours,omitempty"`
	CurrentCourses []CourseBasicResponse `json:"current_courses,omitempty"`
}

func ToFacultyResponse(f *domain.FacultyWithDetails) FacultyResponse {
	return FacultyResponse{
		FacultyID:      f.FacultyID,
		UserID:         f.UserID,
		EmployeeID:     f.EmployeeID,
		Name:           f.Name,
		Email:          f.Email,
		Department:     ToDepartmentBasicResponse(&f.Department),
		Designation:    f.Designation,
		Specialization: f.Specialization,
		IsActive:       f.IsActive,
	}
}

func ToFacultyDetailResponse(f *domain.FacultyWithDetails) FacultyDetailResponse {
	courses := make([]CourseBasicResponse, len(f.CurrentCourses))
	for i, c := range f.CurrentCourses {
		courses[i] = ToCourseBasicResponse(&c)
	}

	return FacultyDetailResponse{
		FacultyResponse: ToFacultyResponse(f),
		Phone:           f.Phone,
		Qualification:   f.Qualification,
		JoiningDate:     f.JoiningDate,
		OfficeRoom:      f.OfficeRoom,
		OfficeHours:     f.OfficeHours,
		CurrentCourses:  courses,
	}
}

// ==================== Student Responses ====================

type StudentResponse struct {
	StudentID          uuid.UUID               `json:"student_id"`
	UserID             uuid.UUID               `json:"user_id"`
	RegistrationNumber string                  `json:"registration_number"`
	Name               string                  `json:"name"`
	Email              string                  `json:"email"`
	Department         DepartmentBasicResponse `json:"department"`
	Program            ProgramBasicResponse    `json:"program"`
	CurrentSemester    int                     `json:"current_semester"`
	BatchYear          int                     `json:"batch_year"`
	CurrentCGPA        *float64                `json:"current_cgpa,omitempty"`
	IsActive           bool                    `json:"is_active"`
}

type StudentDetailResponse struct {
	StudentResponse
	Phone              *string                   `json:"phone,omitempty"`
	RollNumber         *string                   `json:"roll_number,omitempty"`
	AdmissionDate      *time.Time                `json:"admission_date,omitempty"`
	TotalCreditsEarned int                       `json:"total_credits_earned"`
	CurrentEnrollments []EnrollmentBasicResponse `json:"current_enrollments,omitempty"`
}

func ToStudentResponse(s *domain.StudentWithDetails) StudentResponse {
	return StudentResponse{
		StudentID:          s.StudentID,
		UserID:             s.UserID,
		RegistrationNumber: s.RegistrationNumber,
		Name:               s.Name,
		Email:              s.Email,
		Department:         ToDepartmentBasicResponse(&s.Department),
		Program:            ToProgramBasicResponse(&s.Program),
		CurrentSemester:    s.CurrentSemester,
		BatchYear:          s.BatchYear,
		CurrentCGPA:        s.CurrentCGPA,
		IsActive:           s.IsActive,
	}
}

func ToStudentDetailResponse(s *domain.StudentWithDetails) StudentDetailResponse {
	enrollments := make([]EnrollmentBasicResponse, len(s.CurrentEnrollments))
	for i, e := range s.CurrentEnrollments {
		enrollments[i] = EnrollmentBasicResponse{
			CourseID:   e.CourseID,
			CourseCode: e.CourseCode,
			CourseName: e.CourseName,
			Credits:    e.Credits,
		}
	}

	return StudentDetailResponse{
		StudentResponse:    ToStudentResponse(s),
		Phone:              s.Phone,
		RollNumber:         s.RollNumber,
		AdmissionDate:      s.AdmissionDate,
		TotalCreditsEarned: s.TotalCreditsEarned,
		CurrentEnrollments: enrollments,
	}
}

// ==================== Enrollment Responses ====================

type EnrollmentResponse struct {
	EnrollmentID     uuid.UUID            `json:"enrollment_id"`
	Student          StudentBasicResponse `json:"student"`
	EnrollmentStatus string               `json:"enrollment_status"`
	EnrollmentDate   time.Time            `json:"enrollment_date"`
	Grade            *string              `json:"grade,omitempty"`
	GradePoints      *float64             `json:"grade_points,omitempty"`
	WaitlistPosition *int                 `json:"waitlist_position,omitempty"`
}

type StudentEnrollmentResponse struct {
	EnrollmentID     uuid.UUID             `json:"enrollment_id"`
	Course           CourseBasicResponse   `json:"course"`
	Semester         SemesterBasicResponse `json:"semester"`
	EnrollmentStatus string                `json:"enrollment_status"`
	EnrollmentDate   time.Time             `json:"enrollment_date"`
	Grade            *string               `json:"grade,omitempty"`
	GradePoints      *float64              `json:"grade_points,omitempty"`
}

type EnrollResultResponse struct {
	EnrollmentID     uuid.UUID `json:"enrollment_id"`
	EnrollmentStatus string    `json:"enrollment_status"`
	WaitlistPosition *int      `json:"waitlist_position,omitempty"`
	Message          string    `json:"message"`
}

type BulkEnrollResultResponse struct {
	Successful []BulkEnrollSuccessItem `json:"successful"`
	Failed     []BulkEnrollFailedItem  `json:"failed"`
	Summary    BulkEnrollSummary       `json:"summary"`
}

type BulkEnrollSuccessItem struct {
	StudentID    uuid.UUID `json:"student_id"`
	EnrollmentID uuid.UUID `json:"enrollment_id"`
	Status       string    `json:"status"`
}

type BulkEnrollFailedItem struct {
	StudentID uuid.UUID `json:"student_id"`
	Error     string    `json:"error"`
}

type BulkEnrollSummary struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

type EnrollmentSummaryResponse struct {
	TotalEnrolled   int `json:"total_enrolled"`
	TotalWaitlisted int `json:"total_waitlisted"`
	TotalDropped    int `json:"total_dropped"`
	TotalCompleted  int `json:"total_completed"`
}

func ToEnrollmentResponse(e *domain.EnrollmentWithDetails) EnrollmentResponse {
	return EnrollmentResponse{
		EnrollmentID:     e.EnrollmentID,
		Student:          ToStudentBasicResponse(&e.Student),
		EnrollmentStatus: e.EnrollmentStatus,
		EnrollmentDate:   e.EnrollmentDate,
		Grade:            e.Grade,
		GradePoints:      e.GradePoints,
		WaitlistPosition: e.WaitlistPosition,
	}
}

func ToStudentEnrollmentResponse(e *domain.EnrollmentWithDetails) StudentEnrollmentResponse {
	resp := StudentEnrollmentResponse{
		EnrollmentID:     e.EnrollmentID,
		Course:           ToCourseBasicFromEnrollment(&e.Course),
		EnrollmentStatus: e.EnrollmentStatus,
		EnrollmentDate:   e.EnrollmentDate,
		Grade:            e.Grade,
		GradePoints:      e.GradePoints,
	}
	if e.Semester != nil {
		resp.Semester = ToSemesterBasicResponse(e.Semester)
	}
	return resp
}

// ==================== Calendar Responses ====================

type CalendarEventResponse struct {
	EventID     uuid.UUID             `json:"event_id"`
	Semester    SemesterBasicResponse `json:"semester"`
	EventName   string                `json:"event_name"`
	EventType   string                `json:"event_type"`
	StartDate   time.Time             `json:"start_date"`
	EndDate     *time.Time            `json:"end_date,omitempty"`
	Description *string               `json:"description,omitempty"`
	IsHoliday   bool                  `json:"is_holiday"`
}

func ToCalendarEventResponse(e *domain.AcademicCalendarEventWithDetails) CalendarEventResponse {
	return CalendarEventResponse{
		EventID:     e.EventID,
		Semester:    ToSemesterBasicResponse(&e.Semester),
		EventName:   e.EventName,
		EventType:   e.EventType,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		Description: e.Description,
		IsHoliday:   e.IsHoliday,
	}
}

// ==================== Basic Response Types ====================

type DepartmentBasicResponse struct {
	DepartmentID   uuid.UUID `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	DepartmentCode string    `json:"department_code,omitempty"`
}

type ProgramBasicResponse struct {
	ProgramID     uuid.UUID `json:"program_id"`
	ProgramName   string    `json:"program_name"`
	ProgramCode   string    `json:"program_code"`
	DurationYears int       `json:"duration_years,omitempty"`
}

type SubjectBasicResponse struct {
	SubjectID   uuid.UUID `json:"subject_id"`
	SubjectCode string    `json:"subject_code"`
	SubjectName string    `json:"subject_name"`
	Credits     int       `json:"credits,omitempty"`
}

type SemesterBasicResponse struct {
	SemesterID   uuid.UUID  `json:"semester_id"`
	SemesterName string     `json:"semester_name"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
}

type FacultyBasicResponse struct {
	FacultyID   uuid.UUID `json:"faculty_id"`
	Name        string    `json:"name"`
	Designation *string   `json:"designation,omitempty"`
}

type FacultyCourseResponse struct {
	FacultyID   uuid.UUID `json:"faculty_id"`
	Name        string    `json:"name"`
	Designation *string   `json:"designation,omitempty"`
	Role        string    `json:"role"`
	IsPrimary   bool      `json:"is_primary"`
}

type StudentBasicResponse struct {
	StudentID          uuid.UUID `json:"student_id"`
	RegistrationNumber string    `json:"registration_number"`
	Name               string    `json:"name"`
	Email              string    `json:"email,omitempty"`
}

type CourseBasicResponse struct {
	CourseID   uuid.UUID `json:"course_id"`
	CourseCode string    `json:"course_code"`
	CourseName string    `json:"course_name"`
	Role       string    `json:"role,omitempty"`
}

type EnrollmentBasicResponse struct {
	CourseID   uuid.UUID `json:"course_id"`
	CourseCode string    `json:"course_code"`
	CourseName string    `json:"course_name"`
	Credits    int       `json:"credits"`
}

// ==================== Converter Helpers ====================

func ToDepartmentBasicResponse(d *domain.DepartmentBasic) DepartmentBasicResponse {
	return DepartmentBasicResponse{
		DepartmentID:   d.DepartmentID,
		DepartmentName: d.DepartmentName,
		DepartmentCode: d.DepartmentCode,
	}
}

func ToProgramBasicResponse(p *domain.ProgramBasic) ProgramBasicResponse {
	return ProgramBasicResponse{
		ProgramID:     p.ProgramID,
		ProgramName:   p.ProgramName,
		ProgramCode:   p.ProgramCode,
		DurationYears: p.DurationYears,
	}
}

func ToSubjectBasicResponse(s *domain.SubjectBasic) SubjectBasicResponse {
	return SubjectBasicResponse{
		SubjectID:   s.SubjectID,
		SubjectCode: s.SubjectCode,
		SubjectName: s.SubjectName,
		Credits:     s.Credits,
	}
}

func ToSemesterBasicResponse(s *domain.SemesterBasic) SemesterBasicResponse {
	return SemesterBasicResponse{
		SemesterID:   s.SemesterID,
		SemesterName: s.SemesterName,
		StartDate:    s.StartDate,
		EndDate:      s.EndDate,
	}
}

func ToFacultyBasicResponse(f *domain.FacultyBasic) FacultyBasicResponse {
	return FacultyBasicResponse{
		FacultyID:   f.FacultyID,
		Name:        f.Name,
		Designation: f.Designation,
	}
}

func ToStudentBasicResponse(s *domain.StudentBasic) StudentBasicResponse {
	return StudentBasicResponse{
		StudentID:          s.StudentID,
		RegistrationNumber: s.RegistrationNumber,
		Name:               s.Name,
		Email:              s.Email,
	}
}

func ToCourseBasicResponse(c *domain.CourseBasic) CourseBasicResponse {
	return CourseBasicResponse{
		CourseID:   c.CourseID,
		CourseCode: c.CourseCode,
		CourseName: c.CourseName,
		Role:       c.Role,
	}
}

func ToCourseBasicFromEnrollment(c *domain.CourseBasic) CourseBasicResponse {
	return CourseBasicResponse{
		CourseID:   c.CourseID,
		CourseCode: c.CourseCode,
		CourseName: c.CourseName,
	}
}

// Pointer-returning converter functions for handlers

func DepartmentToResponse(d *domain.Department) *DepartmentResponse {
	resp := ToDepartmentResponse(d)
	return &resp
}

func DepartmentWithDetailsToResponse(d *domain.DepartmentWithDetails) *DepartmentDetailResponse {
	resp := ToDepartmentDetailResponse(d)
	return &resp
}

func CourseToResponse(c *domain.Course) *CourseResponse {
	// Create a basic response from a Course (not CourseWithDetails)
	return &CourseResponse{
		CourseID:          c.CourseID,
		CourseCode:        c.CourseCode,
		CourseName:        c.CourseName,
		SemesterNumber:    c.SemesterNumber,
		AcademicYear:      c.AcademicYear,
		MaxStudents:       c.MaxStudents,
		CurrentEnrollment: c.CurrentEnrollment,
		Status:            c.Status,
		IsActive:          c.IsActive,
	}
}

func CourseWithDetailsToResponse(c *domain.CourseWithDetails) *CourseWithDetailsResponse {
	faculty := make([]FacultyCourseResponse, len(c.Faculty))
	for i, f := range c.Faculty {
		faculty[i] = FacultyCourseResponse{
			FacultyID:   f.FacultyID,
			Name:        f.Name,
			Designation: f.Designation,
			Role:        f.Role,
			IsPrimary:   f.IsPrimary,
		}
	}

	prereqs := make([]SubjectBasicResponse, len(c.Prerequisites))
	for i, p := range c.Prerequisites {
		prereqs[i] = ToSubjectBasicResponse(&p)
	}

	availableSeats := 0
	if c.MaxStudents != nil {
		availableSeats = *c.MaxStudents - c.CurrentEnrollment
		if availableSeats < 0 {
			availableSeats = 0
		}
	}

	resp := &CourseWithDetailsResponse{
		CourseID:          c.CourseID,
		CourseCode:        c.CourseCode,
		CourseName:        c.CourseName,
		Subject:           ToSubjectBasicResponse(&c.Subject),
		Department:        ToDepartmentBasicResponse(&c.Department),
		Semester:          ToSemesterBasicResponse(&c.Semester),
		SemesterNumber:    c.SemesterNumber,
		AcademicYear:      c.AcademicYear,
		Faculty:           faculty,
		MaxStudents:       c.MaxStudents,
		CurrentEnrollment: c.CurrentEnrollment,
		Status:            c.Status,
		IsActive:          c.IsActive,
		Description:       c.Description,
		Prerequisites:     prereqs,
		AvailableSeats:    availableSeats,
		CreatedBy:         c.CreatedBy,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}

	if c.Program != nil {
		prog := ToProgramBasicResponse(c.Program)
		resp.Program = &prog
	}

	return resp
}

func EnrollmentToResponse(e *domain.CourseEnrollment) *EnrollmentSimpleResponse {
	return &EnrollmentSimpleResponse{
		EnrollmentID:     e.EnrollmentID,
		StudentID:        e.StudentID,
		CourseID:         e.CourseID,
		EnrollmentStatus: e.EnrollmentStatus,
		EnrollmentDate:   e.EnrollmentDate,
		WaitlistPosition: e.WaitlistPosition,
	}
}

func EnrollmentWithDetailsToResponse(e *domain.EnrollmentWithDetails) *EnrollmentWithDetailsResponse {
	resp := &EnrollmentWithDetailsResponse{
		EnrollmentID:     e.EnrollmentID,
		Student:          ToStudentBasicResponse(&e.Student),
		Course:           ToCourseBasicFromEnrollment(&e.Course),
		EnrollmentStatus: e.EnrollmentStatus,
		EnrollmentDate:   e.EnrollmentDate,
		Grade:            e.Grade,
		GradePoints:      e.GradePoints,
		WaitlistPosition: e.WaitlistPosition,
	}
	if e.Semester != nil {
		sem := ToSemesterBasicResponse(e.Semester)
		resp.Semester = &sem
	}
	return resp
}

func FacultyCourseToResponse(fc *domain.FacultyCourse) *FacultyCourseAssignmentResponse {
	return &FacultyCourseAssignmentResponse{
		FacultyCourseID: fc.FacultyCourseID,
		FacultyID:       fc.FacultyID,
		CourseID:        fc.CourseID,
		Role:            fc.Role,
		IsPrimary:       fc.IsPrimary,
		AssignedAt:      fc.AssignedAt,
	}
}

func FacultyCourseWithDetailsToResponse(fc *domain.FacultyCourseWithDetails) *FacultyCourseWithDetailsResponse {
	return &FacultyCourseWithDetailsResponse{
		FacultyCourseID: fc.FacultyCourseID,
		Faculty:         ToFacultyBasicResponse(&fc.Faculty),
		CourseID:        fc.CourseID,
		Role:            fc.Role,
		IsPrimary:       fc.IsPrimary,
		AssignedAt:      fc.AssignedAt,
	}
}

// Additional response types needed by handlers

type CourseWithDetailsResponse struct {
	CourseID          uuid.UUID               `json:"course_id"`
	CourseCode        string                  `json:"course_code"`
	CourseName        string                  `json:"course_name"`
	Subject           SubjectBasicResponse    `json:"subject"`
	Department        DepartmentBasicResponse `json:"department"`
	Semester          SemesterBasicResponse   `json:"semester"`
	Program           *ProgramBasicResponse   `json:"program,omitempty"`
	SemesterNumber    int                     `json:"semester_number"`
	AcademicYear      int                     `json:"academic_year"`
	Faculty           []FacultyCourseResponse `json:"faculty"`
	MaxStudents       *int                    `json:"max_students,omitempty"`
	CurrentEnrollment int                     `json:"current_enrollment"`
	Status            string                  `json:"status"`
	IsActive          bool                    `json:"is_active"`
	Description       *string                 `json:"description,omitempty"`
	Prerequisites     []SubjectBasicResponse  `json:"prerequisites,omitempty"`
	AvailableSeats    int                     `json:"available_seats"`
	CreatedBy         uuid.UUID               `json:"created_by"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}

type EnrollmentSimpleResponse struct {
	EnrollmentID     uuid.UUID `json:"enrollment_id"`
	StudentID        uuid.UUID `json:"student_id"`
	CourseID         uuid.UUID `json:"course_id"`
	EnrollmentStatus string    `json:"enrollment_status"`
	EnrollmentDate   time.Time `json:"enrollment_date"`
	WaitlistPosition *int      `json:"waitlist_position,omitempty"`
}

type EnrollmentWithDetailsResponse struct {
	EnrollmentID     uuid.UUID              `json:"enrollment_id"`
	Student          StudentBasicResponse   `json:"student"`
	Course           CourseBasicResponse    `json:"course"`
	Semester         *SemesterBasicResponse `json:"semester,omitempty"`
	EnrollmentStatus string                 `json:"enrollment_status"`
	EnrollmentDate   time.Time              `json:"enrollment_date"`
	Grade            *string                `json:"grade,omitempty"`
	GradePoints      *float64               `json:"grade_points,omitempty"`
	WaitlistPosition *int                   `json:"waitlist_position,omitempty"`
}

type FacultyCourseAssignmentResponse struct {
	FacultyCourseID uuid.UUID `json:"faculty_course_id"`
	FacultyID       uuid.UUID `json:"faculty_id"`
	CourseID        uuid.UUID `json:"course_id"`
	Role            string    `json:"role"`
	IsPrimary       bool      `json:"is_primary"`
	AssignedAt      time.Time `json:"assigned_at"`
}

type FacultyCourseWithDetailsResponse struct {
	FacultyCourseID uuid.UUID            `json:"faculty_course_id"`
	Faculty         FacultyBasicResponse `json:"faculty"`
	CourseID        uuid.UUID            `json:"course_id"`
	Role            string               `json:"role"`
	IsPrimary       bool                 `json:"is_primary"`
	AssignedAt      time.Time            `json:"assigned_at"`
}

// ==================== Pagination Response ====================

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationResponse(page, limit int, total int64) PaginationResponse {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}
	return PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
