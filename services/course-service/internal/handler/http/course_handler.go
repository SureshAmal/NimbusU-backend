package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CourseHandler struct {
	service           domain.CourseService
	assignmentService domain.FacultyAssignmentService
	enrollmentService domain.EnrollmentService
	validator         *validator.Validate
}

func NewCourseHandler(
	service domain.CourseService,
	assignmentService domain.FacultyAssignmentService,
	enrollmentService domain.EnrollmentService,
) *CourseHandler {
	v := validator.New()
	v.SetTagName("binding")
	return &CourseHandler{
		service:           service,
		assignmentService: assignmentService,
		enrollmentService: enrollmentService,
		validator:         v,
	}
}

func (h *CourseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id")
	if userID == nil {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	course := req.ToDomain()
	course.CreatedBy = userID.(uuid.UUID)

	if err := h.service.CreateCourse(r.Context(), course); err != nil {
		switch err {
		case domain.ErrSubjectNotFound:
			ErrorResponse(w, http.StatusBadRequest, "subject not found", err)
		case domain.ErrSemesterNotFound:
			ErrorResponse(w, http.StatusBadRequest, "semester not found", err)
		case domain.ErrCourseCodeExists:
			ErrorResponse(w, http.StatusConflict, "course code already exists", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to create course", err)
		}
		return
	}

	SuccessResponse(w, http.StatusCreated, "course created", dto.CourseToResponse(course))
}

func (h *CourseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	course, err := h.service.GetCourse(r.Context(), id)
	if err != nil {
		if err == domain.ErrCourseNotFound {
			ErrorResponse(w, http.StatusNotFound, "course not found", err)
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "failed to get course", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "course retrieved", dto.CourseWithDetailsToResponse(course))
}

func (h *CourseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	var req dto.UpdateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	updates := req.ToUpdates()
	if err := h.service.UpdateCourse(r.Context(), id, updates); err != nil {
		switch err {
		case domain.ErrCourseNotFound:
			ErrorResponse(w, http.StatusNotFound, "course not found", err)
		case domain.ErrCannotModifyCompletedCourse:
			ErrorResponse(w, http.StatusBadRequest, "cannot modify completed course", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to update course", err)
		}
		return
	}

	SuccessResponse(w, http.StatusOK, "course updated", nil)
}

func (h *CourseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	if err := h.service.DeleteCourse(r.Context(), id); err != nil {
		if err == domain.ErrCourseNotFound {
			ErrorResponse(w, http.StatusNotFound, "course not found", err)
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "failed to delete course", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "course deleted", nil)
}

func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var filter domain.CourseFilter
	if deptIDStr := r.URL.Query().Get("department_id"); deptIDStr != "" {
		if deptID, err := uuid.Parse(deptIDStr); err == nil {
			filter.DepartmentID = &deptID
		}
	}
	if semIDStr := r.URL.Query().Get("semester_id"); semIDStr != "" {
		if semID, err := uuid.Parse(semIDStr); err == nil {
			filter.SemesterID = &semID
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = &status
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = &search
	}
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	courses, total, err := h.service.ListCourses(r.Context(), filter, page, limit)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to list courses", err)
		return
	}

	response := make([]*dto.CourseWithDetailsResponse, 0, len(courses))
	for _, course := range courses {
		response = append(response, dto.CourseWithDetailsToResponse(course))
	}

	PaginatedResponse(w, http.StatusOK, "courses retrieved", response, page, limit, total)
}

func (h *CourseHandler) Activate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	if err := h.service.ActivateCourse(r.Context(), id); err != nil {
		switch err {
		case domain.ErrCourseNotFound:
			ErrorResponse(w, http.StatusNotFound, "course not found", err)
		case domain.ErrCannotModifyCompletedCourse:
			ErrorResponse(w, http.StatusBadRequest, "cannot modify completed course", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to activate course", err)
		}
		return
	}

	SuccessResponse(w, http.StatusOK, "course activated", nil)
}

func (h *CourseHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	if err := h.service.DeactivateCourse(r.Context(), id); err != nil {
		switch err {
		case domain.ErrCourseNotFound:
			ErrorResponse(w, http.StatusNotFound, "course not found", err)
		case domain.ErrCannotModifyCompletedCourse:
			ErrorResponse(w, http.StatusBadRequest, "cannot modify completed course", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to deactivate course", err)
		}
		return
	}

	SuccessResponse(w, http.StatusOK, "course deactivated", nil)
}

func (h *CourseHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var status *string
	if s := r.URL.Query().Get("status"); s != "" {
		status = &s
	}

	enrollments, total, err := h.service.GetCourseStudents(r.Context(), courseID, status, page, limit)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to get course students", err)
		return
	}

	response := make([]*dto.EnrollmentWithDetailsResponse, 0, len(enrollments))
	for _, e := range enrollments {
		response = append(response, dto.EnrollmentWithDetailsToResponse(e))
	}

	PaginatedResponse(w, http.StatusOK, "course students retrieved", response, page, limit, total)
}

func (h *CourseHandler) AssignFaculty(w http.ResponseWriter, r *http.Request) {
	courseIDStr := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	var req dto.AssignFacultyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	// Get user ID from context
	userID := r.Context().Value("user_id")
	if userID == nil {
		ErrorResponse(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	fc, err := h.assignmentService.AssignFaculty(r.Context(), courseID, req.FacultyID, userID.(uuid.UUID), req.Role, req.IsPrimary)
	if err != nil {
		switch err {
		case domain.ErrFacultyNotFound:
			ErrorResponse(w, http.StatusBadRequest, "faculty not found", err)
		case domain.ErrCourseNotFound:
			ErrorResponse(w, http.StatusNotFound, "course not found", err)
		case domain.ErrFacultyAlreadyAssigned:
			ErrorResponse(w, http.StatusConflict, "faculty already assigned to this course", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to assign faculty", err)
		}
		return
	}

	SuccessResponse(w, http.StatusCreated, "faculty assigned", dto.FacultyCourseToResponse(fc))
}

func (h *CourseHandler) RemoveFaculty(w http.ResponseWriter, r *http.Request) {
	courseIDStr := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	facultyIDStr := chi.URLParam(r, "facultyId")
	facultyID, err := uuid.Parse(facultyIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid faculty ID", err)
		return
	}

	if err := h.assignmentService.RemoveFaculty(r.Context(), courseID, facultyID); err != nil {
		if err == domain.ErrAssignmentNotFound {
			ErrorResponse(w, http.StatusNotFound, "assignment not found", err)
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "failed to remove faculty", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "faculty removed", nil)
}

func (h *CourseHandler) GetFaculty(w http.ResponseWriter, r *http.Request) {
	courseIDStr := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	assignments, err := h.assignmentService.ListCourseFaculty(r.Context(), courseID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to get course faculty", err)
		return
	}

	response := make([]*dto.FacultyCourseWithDetailsResponse, 0, len(assignments))
	for _, a := range assignments {
		response = append(response, dto.FacultyCourseWithDetailsToResponse(a))
	}

	SuccessResponse(w, http.StatusOK, "course faculty retrieved", response)
}

func (h *CourseHandler) EnrollStudent(w http.ResponseWriter, r *http.Request) {
	courseIDStr := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	var req dto.EnrollStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	enrollment, err := h.enrollmentService.EnrollStudent(r.Context(), courseID, req.StudentID, "admin")
	if err != nil {
		switch err {
		case domain.ErrStudentNotFound:
			ErrorResponse(w, http.StatusBadRequest, "student not found", err)
		case domain.ErrCourseNotFound:
			ErrorResponse(w, http.StatusNotFound, "course not found", err)
		case domain.ErrAlreadyEnrolled:
			ErrorResponse(w, http.StatusConflict, "student already enrolled", err)
		case domain.ErrRegistrationClosed:
			ErrorResponse(w, http.StatusBadRequest, "registration is closed", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to enroll student", err)
		}
		return
	}

	SuccessResponse(w, http.StatusCreated, "student enrolled", dto.EnrollmentToResponse(enrollment))
}
