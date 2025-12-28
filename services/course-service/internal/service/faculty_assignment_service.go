package service

import (
	"context"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type facultyAssignmentService struct {
	repo     domain.FacultyCourseRepository
	faculty  domain.FacultyRepository
	course   domain.CourseRepository
	producer domain.EventProducer
}

func NewFacultyAssignmentService(
	repo domain.FacultyCourseRepository,
	faculty domain.FacultyRepository,
	course domain.CourseRepository,
	producer domain.EventProducer,
) domain.FacultyAssignmentService {
	return &facultyAssignmentService{
		repo:     repo,
		faculty:  faculty,
		course:   course,
		producer: producer,
	}
}

func (s *facultyAssignmentService) AssignFaculty(ctx context.Context, courseID, facultyID, assignedBy uuid.UUID, role string, isPrimary bool) (*domain.FacultyCourse, error) {
	// Validate faculty exists
	_, err := s.faculty.GetByID(ctx, facultyID)
	if err != nil {
		return nil, err
	}

	// Validate course exists
	_, err = s.course.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// Check if already assigned
	_, err = s.repo.GetAssignment(ctx, facultyID, courseID)
	if err == nil {
		return nil, domain.ErrFacultyAlreadyAssigned
	}
	if err != domain.ErrAssignmentNotFound {
		return nil, err
	}

	fc := &domain.FacultyCourse{
		FacultyID:  facultyID,
		CourseID:   courseID,
		Role:       role,
		IsPrimary:  isPrimary,
		AssignedBy: assignedBy,
		IsActive:   true,
	}

	if err := s.repo.Create(ctx, fc); err != nil {
		return nil, err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.faculty.assigned", fc.FacultyCourseID.String(), map[string]interface{}{
			"faculty_id": facultyID,
			"course_id":  courseID,
			"role":       role,
			"is_primary": isPrimary,
		})
	}

	return fc, nil
}

func (s *facultyAssignmentService) UpdateAssignment(ctx context.Context, courseID, facultyID uuid.UUID, role string, isPrimary bool) error {
	fc, err := s.repo.GetAssignment(ctx, facultyID, courseID)
	if err != nil {
		return err
	}

	fc.Role = role
	fc.IsPrimary = isPrimary

	if err := s.repo.Update(ctx, fc); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.faculty.assignment_updated", fc.FacultyCourseID.String(), map[string]interface{}{
			"faculty_id": facultyID,
			"course_id":  courseID,
			"role":       role,
			"is_primary": isPrimary,
		})
	}

	return nil
}

func (s *facultyAssignmentService) RemoveFaculty(ctx context.Context, courseID, facultyID uuid.UUID) error {
	if err := s.repo.Delete(ctx, facultyID, courseID); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.faculty.removed", facultyID.String(), map[string]interface{}{
			"faculty_id": facultyID,
			"course_id":  courseID,
		})
	}

	return nil
}

func (s *facultyAssignmentService) ListCourseFaculty(ctx context.Context, courseID uuid.UUID) ([]*domain.FacultyCourseWithDetails, error) {
	return s.repo.ListByCourse(ctx, courseID)
}
