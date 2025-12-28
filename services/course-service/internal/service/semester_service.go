package service

import (
	"context"
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type semesterService struct {
	repo     domain.SemesterRepository
	producer domain.EventProducer
}

func NewSemesterService(repo domain.SemesterRepository, producer domain.EventProducer) domain.SemesterService {
	return &semesterService{repo: repo, producer: producer}
}

func (s *semesterService) CreateSemester(ctx context.Context, semester *domain.Semester) error {
	if err := s.repo.Create(ctx, semester); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.semester.created", semester.SemesterID.String(), map[string]interface{}{
			"semester_id":   semester.SemesterID,
			"semester_name": semester.SemesterName,
			"academic_year": semester.AcademicYear,
		})
	}

	return nil
}

func (s *semesterService) GetSemester(ctx context.Context, id uuid.UUID) (*domain.Semester, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *semesterService) GetCurrentSemester(ctx context.Context) (*domain.Semester, error) {
	return s.repo.GetCurrent(ctx)
}

func (s *semesterService) UpdateSemester(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	sem, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if name, ok := updates["semester_name"].(string); ok {
		sem.SemesterName = name
	}
	if year, ok := updates["academic_year"].(int); ok {
		sem.AcademicYear = year
	}
	if startDate, ok := updates["start_date"].(time.Time); ok {
		sem.StartDate = startDate
	}
	if endDate, ok := updates["end_date"].(time.Time); ok {
		sem.EndDate = endDate
	}
	if regStart, ok := updates["registration_start"]; ok {
		if regStart == nil {
			sem.RegistrationStart = nil
		} else if t, ok := regStart.(time.Time); ok {
			sem.RegistrationStart = &t
		}
	}
	if regEnd, ok := updates["registration_end"]; ok {
		if regEnd == nil {
			sem.RegistrationEnd = nil
		} else if t, ok := regEnd.(time.Time); ok {
			sem.RegistrationEnd = &t
		}
	}

	if err := s.repo.Update(ctx, sem); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.semester.updated", sem.SemesterID.String(), map[string]interface{}{
			"semester_id":   sem.SemesterID,
			"semester_name": sem.SemesterName,
		})
	}

	return nil
}

func (s *semesterService) DeleteSemester(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.semester.deleted", id.String(), map[string]interface{}{
			"semester_id": id,
		})
	}

	return nil
}

func (s *semesterService) ListSemesters(ctx context.Context, filter domain.SemesterFilter, page, limit int) ([]*domain.Semester, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *semesterService) SetCurrentSemester(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.SetCurrent(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.semester.current_changed", id.String(), map[string]interface{}{
			"semester_id": id,
		})
	}

	return nil
}
