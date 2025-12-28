package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSemesterService_CreateSemester(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewSemesterService(mockRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		semesterID := uuid.New()

		semester := &domain.Semester{
			SemesterID:   semesterID,
			SemesterName: "Fall 2024",
			SemesterCode: "FALL24",
			AcademicYear: 2024,
			StartDate:    time.Now(),
			EndDate:      time.Now().AddDate(0, 4, 0),
		}

		mockRepo.EXPECT().Create(gomock.Any(), semester).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.semester.created", semesterID.String(), gomock.Any()).Return(nil)

		err := service.CreateSemester(context.Background(), semester)
		assert.NoError(t, err)
	})

	t.Run("Repository Error", func(t *testing.T) {
		semester := &domain.Semester{
			SemesterID:   uuid.New(),
			SemesterName: "Fall 2024",
		}
		repoErr := errors.New("database error")

		mockRepo.EXPECT().Create(gomock.Any(), semester).Return(repoErr)

		err := service.CreateSemester(context.Background(), semester)
		assert.ErrorIs(t, err, repoErr)
	})
}

func TestSemesterService_GetSemester(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)

	service := NewSemesterService(mockRepo, nil)

	t.Run("Success", func(t *testing.T) {
		semesterID := uuid.New()

		semester := &domain.Semester{
			SemesterID:   semesterID,
			SemesterName: "Fall 2024",
			SemesterCode: "FALL24",
			AcademicYear: 2024,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(semester, nil)

		result, err := service.GetSemester(context.Background(), semesterID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Fall 2024", result.SemesterName)
	})

	t.Run("Semester Not Found", func(t *testing.T) {
		semesterID := uuid.New()

		mockRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(nil, domain.ErrSemesterNotFound)

		result, err := service.GetSemester(context.Background(), semesterID)
		assert.ErrorIs(t, err, domain.ErrSemesterNotFound)
		assert.Nil(t, result)
	})
}

func TestSemesterService_GetCurrentSemester(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)

	service := NewSemesterService(mockRepo, nil)

	t.Run("Success", func(t *testing.T) {
		semester := &domain.Semester{
			SemesterID:   uuid.New(),
			SemesterName: "Fall 2024",
			IsCurrent:    true,
		}

		mockRepo.EXPECT().GetCurrent(gomock.Any()).Return(semester, nil)

		result, err := service.GetCurrentSemester(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsCurrent)
	})

	t.Run("No Current Semester", func(t *testing.T) {
		mockRepo.EXPECT().GetCurrent(gomock.Any()).Return(nil, domain.ErrNoCurrentSemester)

		result, err := service.GetCurrentSemester(context.Background())
		assert.ErrorIs(t, err, domain.ErrNoCurrentSemester)
		assert.Nil(t, result)
	})
}

func TestSemesterService_UpdateSemester(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewSemesterService(mockRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		semesterID := uuid.New()

		semester := &domain.Semester{
			SemesterID:   semesterID,
			SemesterName: "Fall 2024",
			AcademicYear: 2024,
		}

		updates := map[string]interface{}{
			"semester_name": "Fall 2024 - Updated",
			"academic_year": 2025,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(semester, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Semester) error {
			assert.Equal(t, "Fall 2024 - Updated", s.SemesterName)
			assert.Equal(t, 2025, s.AcademicYear)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.semester.updated", semesterID.String(), gomock.Any()).Return(nil)

		err := service.UpdateSemester(context.Background(), semesterID, updates)
		assert.NoError(t, err)
	})

	t.Run("Semester Not Found", func(t *testing.T) {
		semesterID := uuid.New()
		updates := map[string]interface{}{"semester_name": "New Name"}

		mockRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(nil, domain.ErrSemesterNotFound)

		err := service.UpdateSemester(context.Background(), semesterID, updates)
		assert.ErrorIs(t, err, domain.ErrSemesterNotFound)
	})

	t.Run("Update With Date Fields", func(t *testing.T) {
		semesterID := uuid.New()
		newStartDate := time.Now()
		newEndDate := time.Now().AddDate(0, 4, 0)
		regStart := time.Now().AddDate(0, -1, 0)

		semester := &domain.Semester{
			SemesterID:   semesterID,
			SemesterName: "Fall 2024",
		}

		updates := map[string]interface{}{
			"start_date":         newStartDate,
			"end_date":           newEndDate,
			"registration_start": regStart,
			"registration_end":   nil,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(semester, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Semester) error {
			assert.Equal(t, newStartDate, s.StartDate)
			assert.Equal(t, newEndDate, s.EndDate)
			assert.NotNil(t, s.RegistrationStart)
			assert.Nil(t, s.RegistrationEnd)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.semester.updated", semesterID.String(), gomock.Any()).Return(nil)

		err := service.UpdateSemester(context.Background(), semesterID, updates)
		assert.NoError(t, err)
	})
}

func TestSemesterService_DeleteSemester(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewSemesterService(mockRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		semesterID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), semesterID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.semester.deleted", semesterID.String(), gomock.Any()).Return(nil)

		err := service.DeleteSemester(context.Background(), semesterID)
		assert.NoError(t, err)
	})

	t.Run("Semester Not Found", func(t *testing.T) {
		semesterID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), semesterID).Return(domain.ErrSemesterNotFound)

		err := service.DeleteSemester(context.Background(), semesterID)
		assert.ErrorIs(t, err, domain.ErrSemesterNotFound)
	})
}

func TestSemesterService_ListSemesters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)

	service := NewSemesterService(mockRepo, nil)

	t.Run("Success", func(t *testing.T) {
		year := 2024
		filter := domain.SemesterFilter{
			AcademicYear: &year,
		}

		semesters := []*domain.Semester{
			{SemesterID: uuid.New(), SemesterName: "Fall 2024"},
			{SemesterID: uuid.New(), SemesterName: "Spring 2024"},
		}

		mockRepo.EXPECT().List(gomock.Any(), filter, 10, 0).Return(semesters, int64(2), nil)

		result, total, err := service.ListSemesters(context.Background(), filter, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("Empty List", func(t *testing.T) {
		filter := domain.SemesterFilter{}

		mockRepo.EXPECT().List(gomock.Any(), filter, 20, 20).Return([]*domain.Semester{}, int64(0), nil)

		result, total, err := service.ListSemesters(context.Background(), filter, 2, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})
}

func TestSemesterService_SetCurrentSemester(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewSemesterService(mockRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		semesterID := uuid.New()

		mockRepo.EXPECT().SetCurrent(gomock.Any(), semesterID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.semester.current_changed", semesterID.String(), gomock.Any()).Return(nil)

		err := service.SetCurrentSemester(context.Background(), semesterID)
		assert.NoError(t, err)
	})

	t.Run("Semester Not Found", func(t *testing.T) {
		semesterID := uuid.New()

		mockRepo.EXPECT().SetCurrent(gomock.Any(), semesterID).Return(domain.ErrSemesterNotFound)

		err := service.SetCurrentSemester(context.Background(), semesterID)
		assert.ErrorIs(t, err, domain.ErrSemesterNotFound)
	})
}
