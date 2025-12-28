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

type DepartmentHandler struct {
	service   domain.DepartmentService
	validator *validator.Validate
}

func NewDepartmentHandler(service domain.DepartmentService) *DepartmentHandler {
	v := validator.New()
	v.SetTagName("binding")
	return &DepartmentHandler{
		service:   service,
		validator: v,
	}
}

func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "validation failed", err)
		return
	}

	dept := req.ToDomain()
	if err := h.service.CreateDepartment(r.Context(), dept); err != nil {
		if err == domain.ErrDepartmentCodeExists {
			ErrorResponse(w, http.StatusConflict, "department code already exists", err)
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "failed to create department", err)
		return
	}

	SuccessResponse(w, http.StatusCreated, "department created", dto.DepartmentToResponse(dept))
}

func (h *DepartmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid department ID", err)
		return
	}

	dept, err := h.service.GetDepartment(r.Context(), id)
	if err != nil {
		if err == domain.ErrDepartmentNotFound {
			ErrorResponse(w, http.StatusNotFound, "department not found", err)
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "failed to get department", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "department retrieved", dto.DepartmentWithDetailsToResponse(dept))
}

func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid department ID", err)
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	updates := req.ToUpdates()
	if err := h.service.UpdateDepartment(r.Context(), id, updates); err != nil {
		if err == domain.ErrDepartmentNotFound {
			ErrorResponse(w, http.StatusNotFound, "department not found", err)
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "failed to update department", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "department updated", nil)
}

func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid department ID", err)
		return
	}

	if err := h.service.DeleteDepartment(r.Context(), id); err != nil {
		if err == domain.ErrDepartmentNotFound {
			ErrorResponse(w, http.StatusNotFound, "department not found", err)
			return
		}
		ErrorResponse(w, http.StatusInternalServerError, "failed to delete department", err)
		return
	}

	SuccessResponse(w, http.StatusOK, "department deleted", nil)
}

func (h *DepartmentHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var filter domain.DepartmentFilter
	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		isActive := isActiveStr == "true"
		filter.IsActive = &isActive
	}

	depts, total, err := h.service.ListDepartments(r.Context(), filter, page, limit)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "failed to list departments", err)
		return
	}

	response := make([]*dto.DepartmentResponse, 0, len(depts))
	for _, dept := range depts {
		response = append(response, dto.DepartmentToResponse(dept))
	}

	PaginatedResponse(w, http.StatusOK, "departments retrieved", response, page, limit, total)
}
