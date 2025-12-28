package http

import (
	"bytes"
	"context"
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

func TestCourseHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCourseService(ctrl)
	mockAssignService := mocks.NewMockFacultyAssignmentService(ctrl)
	mockEnrollService := mocks.NewMockEnrollmentService(ctrl)
	handler := NewCourseHandler(mockService, mockAssignService, mockEnrollService)

	r := chi.NewRouter()
	r.Post("/courses", handler.Create)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		req := dto.CreateCourseRequest{
			CourseCode:     "CS101",
			CourseName:     "Intro to CS",
			SubjectID:      uuid.New(),
			DepartmentID:   uuid.New(),
			SemesterID:     uuid.New(),
			SemesterNumber: 1,
			AcademicYear:   2025,
		}
		body, _ := json.Marshal(req)

		mockService.EXPECT().CreateCourse(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx interface{}, course *domain.Course) error {
			assert.Equal(t, req.CourseCode, course.CourseCode)
			assert.Equal(t, userID, course.CreatedBy)
			return nil
		})

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))
		// Inject user_id into context
		ctx := context.WithValue(reqHttp.Context(), "user_id", userID)
		reqHttp = reqHttp.WithContext(ctx)

		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := dto.CreateCourseRequest{
			CourseCode:     "CS101",
			CourseName:     "Intro to CS",
			SubjectID:      uuid.New(),
			DepartmentID:   uuid.New(),
			SemesterID:     uuid.New(),
			SemesterNumber: 1,
			AcademicYear:   2025,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewBuffer(body))
		// No user_id in context

		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCourseHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockCourseService(ctrl)
	handler := NewCourseHandler(mockService, nil, nil)

	r := chi.NewRouter()
	r.Get("/courses/{id}", handler.GetByID)

	t.Run("Success", func(t *testing.T) {
		courseID := uuid.New()
		course := &domain.CourseWithDetails{
			Course: domain.Course{
				CourseID:   courseID,
				CourseCode: "CS101",
				CourseName: "Intro to CS",
			},
		}

		mockService.EXPECT().GetCourse(gomock.Any(), courseID).Return(course, nil)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodGet, "/courses/"+courseID.String(), nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		courseID := uuid.New()
		mockService.EXPECT().GetCourse(gomock.Any(), courseID).Return(nil, domain.ErrCourseNotFound)

		w := httptest.NewRecorder()
		reqHttp := httptest.NewRequest(http.MethodGet, "/courses/"+courseID.String(), nil)
		r.ServeHTTP(w, reqHttp)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
