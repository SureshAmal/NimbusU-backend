package service

import (
	"context"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type courseService struct {
	repo           domain.CourseRepository
	subjectRepo    domain.SubjectRepository
	semesterRepo   domain.SemesterRepository
	enrollmentRepo domain.EnrollmentRepository
	producer       domain.EventProducer
}

func NewCourseService(
	repo domain.CourseRepository,
	subjectRepo domain.SubjectRepository,
	semesterRepo domain.SemesterRepository,
	enrollmentRepo domain.EnrollmentRepository,
	producer domain.EventProducer,
) domain.CourseService {
	return &courseService{
		repo:           repo,
		subjectRepo:    subjectRepo,
		semesterRepo:   semesterRepo,
		enrollmentRepo: enrollmentRepo,
		producer:       producer,
	}
}

func (s *courseService) CreateCourse(ctx context.Context, course *domain.Course) error {
	// Validate subject exists
	subject, err := s.subjectRepo.GetByID(ctx, course.SubjectID)
	if err != nil {
		return err
	}

	// Validate semester exists
	_, err = s.semesterRepo.GetByID(ctx, course.SemesterID)
	if err != nil {
		return err
	}

	// Set defaults
	course.IsActive = true
	if course.Status == "" {
		course.Status = "draft"
	}
	course.CurrentEnrollment = 0

	// Use subject name if course name not provided
	if course.CourseName == "" {
		course.CourseName = subject.SubjectName
	}

	if err := s.repo.Create(ctx, course); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.course.created", course.CourseID.String(), map[string]interface{}{
			"course_id":     course.CourseID,
			"course_code":   course.CourseCode,
			"course_name":   course.CourseName,
			"subject_id":    course.SubjectID,
			"semester_id":   course.SemesterID,
			"department_id": course.DepartmentID,
		})
	}

	return nil
}

func (s *courseService) GetCourse(ctx context.Context, id uuid.UUID) (*domain.CourseWithDetails, error) {
	return s.repo.GetWithDetails(ctx, id)
}

func (s *courseService) UpdateCourse(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if course can be modified
	if course.Status == "completed" {
		return domain.ErrCannotModifyCompletedCourse
	}

	// Apply updates
	if name, ok := updates["course_name"].(string); ok {
		course.CourseName = name
	}
	if maxStudents, ok := updates["max_students"]; ok {
		if maxStudents == nil {
			course.MaxStudents = nil
		} else if ms, ok := maxStudents.(int); ok {
			course.MaxStudents = &ms
		}
	}
	if desc, ok := updates["description"]; ok {
		if desc == nil {
			course.Description = nil
		} else if d, ok := desc.(string); ok {
			course.Description = &d
		}
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		course.IsActive = isActive
	}

	if err := s.repo.Update(ctx, course); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.course.updated", course.CourseID.String(), map[string]interface{}{
			"course_id":   course.CourseID,
			"course_name": course.CourseName,
		})
	}

	return nil
}

func (s *courseService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.course.deleted", id.String(), map[string]interface{}{
			"course_id": id,
		})
	}

	return nil
}

func (s *courseService) ListCourses(ctx context.Context, filter domain.CourseFilter, page, limit int) ([]*domain.CourseWithDetails, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *courseService) ActivateCourse(ctx context.Context, id uuid.UUID) error {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if course.Status == "completed" {
		return domain.ErrCannotModifyCompletedCourse
	}

	if err := s.repo.UpdateStatus(ctx, id, "active"); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.course.activated", id.String(), map[string]interface{}{
			"course_id": id,
		})
	}

	return nil
}

func (s *courseService) DeactivateCourse(ctx context.Context, id uuid.UUID) error {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if course.Status == "completed" {
		return domain.ErrCannotModifyCompletedCourse
	}

	if err := s.repo.UpdateStatus(ctx, id, "cancelled"); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.course.deactivated", id.String(), map[string]interface{}{
			"course_id": id,
		})
	}

	return nil
}

func (s *courseService) GetCourseStudents(ctx context.Context, courseID uuid.UUID, status *string, page, limit int) ([]*domain.EnrollmentWithDetails, int64, error) {
	offset := (page - 1) * limit
	return s.enrollmentRepo.ListByCourse(ctx, courseID, status, limit, offset)
}
