package service

import (
	"context"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type studentService struct {
	repo     domain.StudentRepository
	deptRepo domain.DepartmentRepository
	progRepo domain.ProgramRepository
	producer domain.EventProducer
}

func NewStudentService(
	repo domain.StudentRepository,
	deptRepo domain.DepartmentRepository,
	progRepo domain.ProgramRepository,
	producer domain.EventProducer,
) domain.StudentService {
	return &studentService{
		repo:     repo,
		deptRepo: deptRepo,
		progRepo: progRepo,
		producer: producer,
	}
}

func (s *studentService) CreateStudent(ctx context.Context, student *domain.Student) error {
	// Validate department exists
	_, err := s.deptRepo.GetByID(ctx, student.DepartmentID)
	if err != nil {
		return err
	}

	// Validate program exists
	_, err = s.progRepo.GetByID(ctx, student.ProgramID)
	if err != nil {
		return err
	}

	student.IsActive = true
	if student.CurrentSemester == 0 {
		student.CurrentSemester = 1
	}
	student.TotalCreditsEarned = 0

	if err := s.repo.Create(ctx, student); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.student.created", student.StudentID.String(), map[string]interface{}{
			"student_id":          student.StudentID,
			"user_id":             student.UserID,
			"registration_number": student.RegistrationNumber,
			"department_id":       student.DepartmentID,
			"program_id":          student.ProgramID,
		})
	}

	return nil
}

func (s *studentService) GetStudent(ctx context.Context, id uuid.UUID) (*domain.StudentWithDetails, error) {
	return s.repo.GetWithDetails(ctx, id)
}

func (s *studentService) GetStudentByUserID(ctx context.Context, userID uuid.UUID) (*domain.StudentWithDetails, error) {
	student, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetWithDetails(ctx, student.StudentID)
}

func (s *studentService) UpdateStudent(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if rollNumber, ok := updates["roll_number"]; ok {
		if rollNumber == nil {
			student.RollNumber = nil
		} else if rn, ok := rollNumber.(string); ok {
			student.RollNumber = &rn
		}
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		student.IsActive = isActive
	}

	if err := s.repo.Update(ctx, student); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.student.updated", student.StudentID.String(), map[string]interface{}{
			"student_id":          student.StudentID,
			"registration_number": student.RegistrationNumber,
		})
	}

	return nil
}

func (s *studentService) DeleteStudent(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.student.deleted", id.String(), map[string]interface{}{
			"student_id": id,
		})
	}

	return nil
}

func (s *studentService) ListStudents(ctx context.Context, filter domain.StudentFilter, page, limit int) ([]*domain.StudentWithDetails, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *studentService) PromoteStudent(ctx context.Context, id uuid.UUID, newSemester int, cgpa *float64, creditsEarned int) error {
	if err := s.repo.UpdateSemester(ctx, id, newSemester, cgpa, creditsEarned); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.student.promoted", id.String(), map[string]interface{}{
			"student_id":   id,
			"new_semester": newSemester,
			"cgpa":         cgpa,
		})
	}

	return nil
}
