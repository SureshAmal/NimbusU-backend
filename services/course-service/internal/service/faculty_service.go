package service

import (
	"context"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type facultyService struct {
	repo     domain.FacultyRepository
	deptRepo domain.DepartmentRepository
	fcRepo   domain.FacultyCourseRepository
	producer domain.EventProducer
}

func NewFacultyService(
	repo domain.FacultyRepository,
	deptRepo domain.DepartmentRepository,
	fcRepo domain.FacultyCourseRepository,
	producer domain.EventProducer,
) domain.FacultyService {
	return &facultyService{
		repo:     repo,
		deptRepo: deptRepo,
		fcRepo:   fcRepo,
		producer: producer,
	}
}

func (s *facultyService) CreateFaculty(ctx context.Context, faculty *domain.Faculty) error {
	// Validate department exists
	_, err := s.deptRepo.GetByID(ctx, faculty.DepartmentID)
	if err != nil {
		return err
	}

	faculty.IsActive = true
	if err := s.repo.Create(ctx, faculty); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.faculty.created", faculty.FacultyID.String(), map[string]interface{}{
			"faculty_id":    faculty.FacultyID,
			"user_id":       faculty.UserID,
			"employee_id":   faculty.EmployeeID,
			"department_id": faculty.DepartmentID,
		})
	}

	return nil
}

func (s *facultyService) GetFaculty(ctx context.Context, id uuid.UUID) (*domain.FacultyWithDetails, error) {
	return s.repo.GetWithDetails(ctx, id)
}

func (s *facultyService) GetFacultyByUserID(ctx context.Context, userID uuid.UUID) (*domain.FacultyWithDetails, error) {
	faculty, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetWithDetails(ctx, faculty.FacultyID)
}

func (s *facultyService) UpdateFaculty(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	faculty, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if designation, ok := updates["designation"]; ok {
		if designation == nil {
			faculty.Designation = nil
		} else if d, ok := designation.(string); ok {
			faculty.Designation = &d
		}
	}
	if qualification, ok := updates["qualification"]; ok {
		if qualification == nil {
			faculty.Qualification = nil
		} else if q, ok := qualification.(string); ok {
			faculty.Qualification = &q
		}
	}
	if specialization, ok := updates["specialization"]; ok {
		if specialization == nil {
			faculty.Specialization = nil
		} else if sp, ok := specialization.(string); ok {
			faculty.Specialization = &sp
		}
	}
	if officeRoom, ok := updates["office_room"]; ok {
		if officeRoom == nil {
			faculty.OfficeRoom = nil
		} else if or, ok := officeRoom.(string); ok {
			faculty.OfficeRoom = &or
		}
	}
	if officeHours, ok := updates["office_hours"]; ok {
		if officeHours == nil {
			faculty.OfficeHours = nil
		} else if oh, ok := officeHours.(string); ok {
			faculty.OfficeHours = &oh
		}
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		faculty.IsActive = isActive
	}

	if err := s.repo.Update(ctx, faculty); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.faculty.updated", faculty.FacultyID.String(), map[string]interface{}{
			"faculty_id":  faculty.FacultyID,
			"employee_id": faculty.EmployeeID,
		})
	}

	return nil
}

func (s *facultyService) DeleteFaculty(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.faculty.deleted", id.String(), map[string]interface{}{
			"faculty_id": id,
		})
	}

	return nil
}

func (s *facultyService) ListFaculty(ctx context.Context, filter domain.FacultyFilter, page, limit int) ([]*domain.FacultyWithDetails, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *facultyService) GetFacultyCourses(ctx context.Context, facultyID uuid.UUID, semesterID *uuid.UUID, page, limit int) ([]*domain.FacultyCourse, int64, error) {
	courses, err := s.fcRepo.ListByFaculty(ctx, facultyID, semesterID)
	if err != nil {
		return nil, 0, err
	}
	// Simple pagination for now
	total := int64(len(courses))
	offset := (page - 1) * limit
	end := offset + limit
	if end > len(courses) {
		end = len(courses)
	}
	if offset > len(courses) {
		return []*domain.FacultyCourse{}, total, nil
	}
	return courses[offset:end], total, nil
}
