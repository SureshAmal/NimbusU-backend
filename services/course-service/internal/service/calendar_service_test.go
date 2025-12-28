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

func TestCalendarService_CreateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCalendarRepository(ctrl)
	mockSemesterRepo := mocks.NewMockSemesterRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewCalendarService(mockRepo, mockSemesterRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		eventID := uuid.New()
		semesterID := uuid.New()

		event := &domain.AcademicCalendarEvent{
			EventID:    eventID,
			SemesterID: semesterID,
			EventName:  "Final Exams",
			EventType:  "exam",
			StartDate:  time.Now(),
		}

		semester := &domain.Semester{
			SemesterID:   semesterID,
			SemesterName: "Fall 2024",
		}

		mockSemesterRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(semester, nil)
		mockRepo.EXPECT().Create(gomock.Any(), event).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.calendar.event_created", eventID.String(), gomock.Any()).Return(nil)

		err := service.CreateEvent(context.Background(), event)
		assert.NoError(t, err)
	})

	t.Run("Semester Not Found", func(t *testing.T) {
		event := &domain.AcademicCalendarEvent{
			SemesterID: uuid.New(),
			EventName:  "Test Event",
		}

		mockSemesterRepo.EXPECT().GetByID(gomock.Any(), event.SemesterID).Return(nil, domain.ErrSemesterNotFound)

		err := service.CreateEvent(context.Background(), event)
		assert.ErrorIs(t, err, domain.ErrSemesterNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		eventID := uuid.New()
		semesterID := uuid.New()

		event := &domain.AcademicCalendarEvent{
			EventID:    eventID,
			SemesterID: semesterID,
			EventName:  "Test Event",
		}

		semester := &domain.Semester{SemesterID: semesterID}
		repoErr := errors.New("database error")

		mockSemesterRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(semester, nil)
		mockRepo.EXPECT().Create(gomock.Any(), event).Return(repoErr)

		err := service.CreateEvent(context.Background(), event)
		assert.ErrorIs(t, err, repoErr)
	})
}

func TestCalendarService_GetEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCalendarRepository(ctrl)
	mockSemesterRepo := mocks.NewMockSemesterRepository(ctrl)

	service := NewCalendarService(mockRepo, mockSemesterRepo, nil)

	t.Run("Success", func(t *testing.T) {
		eventID := uuid.New()
		semesterID := uuid.New()

		event := &domain.AcademicCalendarEvent{
			EventID:    eventID,
			SemesterID: semesterID,
			EventName:  "Final Exams",
			EventType:  "exam",
		}

		semester := &domain.Semester{
			SemesterID:   semesterID,
			SemesterName: "Fall 2024",
			SemesterCode: "FALL24",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), eventID).Return(event, nil)
		mockSemesterRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(semester, nil)

		result, err := service.GetEvent(context.Background(), eventID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, eventID, result.EventID)
		assert.Equal(t, "Fall 2024", result.Semester.SemesterName)
	})

	t.Run("Event Not Found", func(t *testing.T) {
		eventID := uuid.New()

		mockRepo.EXPECT().GetByID(gomock.Any(), eventID).Return(nil, domain.ErrCalendarEventNotFound)

		result, err := service.GetEvent(context.Background(), eventID)
		assert.ErrorIs(t, err, domain.ErrCalendarEventNotFound)
		assert.Nil(t, result)
	})
}

func TestCalendarService_UpdateEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCalendarRepository(ctrl)
	mockSemesterRepo := mocks.NewMockSemesterRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewCalendarService(mockRepo, mockSemesterRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		eventID := uuid.New()
		event := &domain.AcademicCalendarEvent{
			EventID:   eventID,
			EventName: "Old Name",
			EventType: "holiday",
		}

		updates := map[string]interface{}{
			"event_name": "New Name",
			"is_holiday": true,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), eventID).Return(event, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, e *domain.AcademicCalendarEvent) error {
			assert.Equal(t, "New Name", e.EventName)
			assert.True(t, e.IsHoliday)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.calendar.event_updated", eventID.String(), gomock.Any()).Return(nil)

		err := service.UpdateEvent(context.Background(), eventID, updates)
		assert.NoError(t, err)
	})

	t.Run("Event Not Found", func(t *testing.T) {
		eventID := uuid.New()
		updates := map[string]interface{}{"event_name": "New Name"}

		mockRepo.EXPECT().GetByID(gomock.Any(), eventID).Return(nil, domain.ErrCalendarEventNotFound)

		err := service.UpdateEvent(context.Background(), eventID, updates)
		assert.ErrorIs(t, err, domain.ErrCalendarEventNotFound)
	})
}

func TestCalendarService_DeleteEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCalendarRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewCalendarService(mockRepo, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		eventID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), eventID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.calendar.event_deleted", eventID.String(), gomock.Any()).Return(nil)

		err := service.DeleteEvent(context.Background(), eventID)
		assert.NoError(t, err)
	})

	t.Run("Event Not Found", func(t *testing.T) {
		eventID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), eventID).Return(domain.ErrCalendarEventNotFound)

		err := service.DeleteEvent(context.Background(), eventID)
		assert.ErrorIs(t, err, domain.ErrCalendarEventNotFound)
	})
}

func TestCalendarService_ListEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCalendarRepository(ctrl)

	service := NewCalendarService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		semesterID := uuid.New()
		filter := domain.CalendarFilter{
			SemesterID: &semesterID,
		}

		events := []*domain.AcademicCalendarEventWithDetails{
			{
				AcademicCalendarEvent: domain.AcademicCalendarEvent{
					EventID:   uuid.New(),
					EventName: "Event 1",
				},
			},
			{
				AcademicCalendarEvent: domain.AcademicCalendarEvent{
					EventID:   uuid.New(),
					EventName: "Event 2",
				},
			},
		}

		mockRepo.EXPECT().List(gomock.Any(), filter, 10, 0).Return(events, int64(2), nil)

		result, total, err := service.ListEvents(context.Background(), filter, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("Empty List", func(t *testing.T) {
		filter := domain.CalendarFilter{}

		mockRepo.EXPECT().List(gomock.Any(), filter, 20, 20).Return([]*domain.AcademicCalendarEventWithDetails{}, int64(0), nil)

		result, total, err := service.ListEvents(context.Background(), filter, 2, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})
}
