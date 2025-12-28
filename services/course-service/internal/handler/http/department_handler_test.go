package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/dto"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDepartmentHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDepartmentService(ctrl)
	handler := NewDepartmentHandler(mockService)

	// Setup router for handling
	r := chi.NewRouter()
	r.Post("/departments", handler.Create)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateDepartmentRequest{
			DepartmentName: "Computer Science",
			DepartmentCode: "CS",
		}
		body, _ := json.Marshal(req)

		mockService.EXPECT().CreateDepartment(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx interface{}, dept *domain.Department) error {
			assert.Equal(t, req.DepartmentName, dept.DepartmentName)
			assert.Equal(t, req.DepartmentCode, dept.DepartmentCode)
			return nil
		})

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(body))
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Validation Error", func(t *testing.T) {
		req := dto.CreateDepartmentRequest{
			DepartmentName: "", // Required
			DepartmentCode: "", // Required
		}
		body, _ := json.Marshal(req)

		// Service should NOT be called
		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewBuffer(body))
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDepartmentHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockDepartmentService(ctrl)
	handler := NewDepartmentHandler(mockService)

	r := chi.NewRouter()
	r.Get("/departments/{id}", handler.GetByID)

	t.Run("Success", func(t *testing.T) {
		deptID := uuid.New()
		dept := &domain.DepartmentWithDetails{
			Department: domain.Department{
				DepartmentID:   deptID,
				DepartmentName: "Math",
				DepartmentCode: "MATH",
			},
		}

		mockService.EXPECT().GetDepartment(gomock.Any(), deptID).Return(dept, nil)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodGet, "/departments/"+deptID.String(), nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]interface{})
		assert.Equal(t, "Math", data["department_name"])
	})

	t.Run("Not Found", func(t *testing.T) {
		deptID := uuid.New()
		mockService.EXPECT().GetDepartment(gomock.Any(), deptID).Return(nil, domain.ErrDepartmentNotFound)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodGet, "/departments/"+deptID.String(), nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
