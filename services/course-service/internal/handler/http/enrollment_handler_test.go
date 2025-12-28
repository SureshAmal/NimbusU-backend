package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/dto"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEnrollmentHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockEnrollmentService(ctrl)
	handler := NewEnrollmentHandler(mockService)

	r := chi.NewRouter()
	r.Put("/enrollments/{id}", handler.Update)

	t.Run("Success", func(t *testing.T) {
		enrollmentID := uuid.New()
		grade := "A"
		gradePoints := 4.0
		req := dto.UpdateEnrollmentRequest{
			EnrollmentStatus: "completed",
			Grade:            &grade,
			GradePoints:      &gradePoints,
		}
		body, _ := json.Marshal(req)

		mockService.EXPECT().UpdateEnrollment(gomock.Any(), enrollmentID, "completed", &grade, &gradePoints).Return(nil)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPut, "/enrollments/"+enrollmentID.String(), bytes.NewBuffer(body))
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Status", func(t *testing.T) {
		enrollmentID := uuid.New()
		req := dto.UpdateEnrollmentRequest{
			EnrollmentStatus: "invalid_status",
		}
		body, _ := json.Marshal(req)

		// Should fail validation before service call
		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPut, "/enrollments/"+enrollmentID.String(), bytes.NewBuffer(body))
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestEnrollmentHandler_Drop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockEnrollmentService(ctrl)
	handler := NewEnrollmentHandler(mockService)

	r := chi.NewRouter()
	r.Delete("/enrollments/courses/{courseId}/students/{studentId}", handler.Drop)

	t.Run("Success", func(t *testing.T) {
		courseID := uuid.New()
		studentID := uuid.New()
		req := dto.DropCourseRequest{
			Reason: "Too hard",
		}
		body, _ := json.Marshal(req)

		mockService.EXPECT().DropCourse(gomock.Any(), courseID, studentID, "Too hard").Return(nil)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodDelete, "/enrollments/courses/"+courseID.String()+"/students/"+studentID.String(), bytes.NewBuffer(body))
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
