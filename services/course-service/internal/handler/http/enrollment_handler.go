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

type EnrollmentHandler struct {
	service   domain.EnrollmentService
	validator *validator.Validate
}

func NewEnrollmentHandler(service domain.EnrollmentService) *EnrollmentHandler {
	v := validator.New()
	v.SetTagName("binding")
	return &EnrollmentHandler{
		service:   service,
		validator: v,
	}
}

func (h *EnrollmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid enrollment ID", err)
		return
	}

	// For now, just return basic enrollment info
	filter := domain.EnrollmentFilter{}
	enrollments, _, err := h.service.GetStudentEnrollments(r.Context(), id, filter, 1, 1)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to get enrollment", err)
		return
	}

	if len(enrollments) == 0 {
		ErrorResponse(w, http.StatusNotFound, "enrollment not found", nil)
		return
	}

	SuccessResponse(w, http.StatusOK, "enrollment retrieved", dto.EnrollmentWithDetailsToResponse(enrollments[0]))
}

func (h *EnrollmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid enrollment ID", err)
		return
	}

	var req dto.UpdateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	if err := h.service.UpdateEnrollment(r.Context(), id, req.EnrollmentStatus, req.Grade, req.GradePoints); err != nil {
		switch err {
		case domain.ErrEnrollmentNotFound:
			ErrorResponse(w, http.StatusNotFound, "enrollment not found", err)
		case domain.ErrInvalidEnrollmentStatus:
			ErrorResponse(w, http.StatusBadRequest, "invalid status transition", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to update enrollment", err)
		}
		return
	}

	SuccessResponse(w, http.StatusOK, "enrollment updated", nil)
}

func (h *EnrollmentHandler) Drop(w http.ResponseWriter, r *http.Request) {
	courseIDStr := chi.URLParam(r, "courseId")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	studentIDStr := chi.URLParam(r, "studentId")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid student ID", err)
		return
	}

	var req dto.DropCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.service.DropCourse(r.Context(), courseID, studentID, req.Reason); err != nil {
		switch err {
		case domain.ErrEnrollmentNotFound:
			ErrorResponse(w, http.StatusNotFound, "enrollment not found", err)
		case domain.ErrCannotDropCompletedCourse:
			ErrorResponse(w, http.StatusBadRequest, "cannot drop completed course", err)
		default:
			ErrorResponse(w, http.StatusInternalServerError, "failed to drop course", err)
		}
		return
	}

	SuccessResponse(w, http.StatusOK, "course dropped", nil)
}

func (h *EnrollmentHandler) GetStudentEnrollments(w http.ResponseWriter, r *http.Request) {
	studentIDStr := chi.URLParam(r, "studentId")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid student ID", err)
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

	var filter domain.EnrollmentFilter
	if semIDStr := r.URL.Query().Get("semester_id"); semIDStr != "" {
		if semID, err := uuid.Parse(semIDStr); err == nil {
			filter.SemesterID = &semID
		}
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = &status
	}

	enrollments, total, err := h.service.GetStudentEnrollments(r.Context(), studentID, filter, page, limit)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to get enrollments", err)
		return
	}

	response := make([]*dto.EnrollmentWithDetailsResponse, 0, len(enrollments))
	for _, e := range enrollments {
		response = append(response, dto.EnrollmentWithDetailsToResponse(e))
	}

	PaginatedResponse(w, http.StatusOK, "enrollments retrieved", response, page, limit, total)
}

func (h *EnrollmentHandler) BulkEnroll(w http.ResponseWriter, r *http.Request) {
	courseIDStr := chi.URLParam(r, "courseId")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	var req dto.BulkEnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	results, err := h.service.BulkEnroll(r.Context(), courseID, req.StudentIDs, req.SkipPrerequisites)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to bulk enroll", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "bulk enrollment completed", results)
}

func (h *EnrollmentHandler) CheckPrerequisites(w http.ResponseWriter, r *http.Request) {
	courseIDStr := chi.URLParam(r, "courseId")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid course ID", err)
		return
	}

	studentIDStr := chi.URLParam(r, "studentId")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid student ID", err)
		return
	}

	met, missing, err := h.service.CheckPrerequisites(r.Context(), studentID, courseID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to check prerequisites", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "prerequisites checked", map[string]interface{}{
		"met":                   met,
		"missing_prerequisites": missing,
	})
}
