package service

import (
	"context"
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type calendarService struct {
	repo         domain.CalendarRepository
	semesterRepo domain.SemesterRepository
	producer     domain.EventProducer
}

func NewCalendarService(repo domain.CalendarRepository, semesterRepo domain.SemesterRepository, producer domain.EventProducer) domain.CalendarService {
	return &calendarService{repo: repo, semesterRepo: semesterRepo, producer: producer}
}

func (s *calendarService) CreateEvent(ctx context.Context, event *domain.AcademicCalendarEvent) error {
	// Validate semester exists
	_, err := s.semesterRepo.GetByID(ctx, event.SemesterID)
	if err != nil {
		return err
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.calendar.event_created", event.EventID.String(), map[string]interface{}{
			"event_id":    event.EventID,
			"event_name":  event.EventName,
			"event_type":  event.EventType,
			"semester_id": event.SemesterID,
			"start_date":  event.StartDate,
		})
	}

	return nil
}

func (s *calendarService) GetEvent(ctx context.Context, id uuid.UUID) (*domain.AcademicCalendarEventWithDetails, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	semester, err := s.semesterRepo.GetByID(ctx, event.SemesterID)
	if err != nil {
		return nil, err
	}

	return &domain.AcademicCalendarEventWithDetails{
		AcademicCalendarEvent: *event,
		Semester: domain.SemesterBasic{
			SemesterID:   semester.SemesterID,
			SemesterName: semester.SemesterName,
			SemesterCode: semester.SemesterCode,
		},
	}, nil
}

func (s *calendarService) UpdateEvent(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if name, ok := updates["event_name"].(string); ok {
		event.EventName = name
	}
	if eventType, ok := updates["event_type"].(string); ok {
		event.EventType = eventType
	}
	if startDate, ok := updates["start_date"].(time.Time); ok {
		event.StartDate = startDate
	}
	if endDate, ok := updates["end_date"]; ok {
		if endDate == nil {
			event.EndDate = nil
		} else if t, ok := endDate.(time.Time); ok {
			event.EndDate = &t
		}
	}
	if desc, ok := updates["description"]; ok {
		if desc == nil {
			event.Description = nil
		} else if d, ok := desc.(string); ok {
			event.Description = &d
		}
	}
	if isHoliday, ok := updates["is_holiday"].(bool); ok {
		event.IsHoliday = isHoliday
	}

	if err := s.repo.Update(ctx, event); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.calendar.event_updated", event.EventID.String(), map[string]interface{}{
			"event_id":   event.EventID,
			"event_name": event.EventName,
		})
	}

	return nil
}

func (s *calendarService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.calendar.event_deleted", id.String(), map[string]interface{}{
			"event_id": id,
		})
	}

	return nil
}

func (s *calendarService) ListEvents(ctx context.Context, filter domain.CalendarFilter, page, limit int) ([]*domain.AcademicCalendarEventWithDetails, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}
