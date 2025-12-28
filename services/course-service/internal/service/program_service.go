package service

import (
	"context"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type programService struct {
	repo     domain.ProgramRepository
	deptRepo domain.DepartmentRepository
	producer domain.EventProducer
}

func NewProgramService(repo domain.ProgramRepository, deptRepo domain.DepartmentRepository, producer domain.EventProducer) domain.ProgramService {
	return &programService{repo: repo, deptRepo: deptRepo, producer: producer}
}

func (s *programService) CreateProgram(ctx context.Context, program *domain.Program) error {
	// Validate department exists
	_, err := s.deptRepo.GetByID(ctx, program.DepartmentID)
	if err != nil {
		return err
	}

	program.IsActive = true
	if err := s.repo.Create(ctx, program); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.program.created", program.ProgramID.String(), map[string]interface{}{
			"program_id":    program.ProgramID,
			"program_name":  program.ProgramName,
			"program_code":  program.ProgramCode,
			"department_id": program.DepartmentID,
		})
	}

	return nil
}

func (s *programService) GetProgram(ctx context.Context, id uuid.UUID) (*domain.ProgramWithDepartment, error) {
	return s.repo.GetWithDepartment(ctx, id)
}

func (s *programService) UpdateProgram(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	prog, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if name, ok := updates["program_name"].(string); ok {
		prog.ProgramName = name
	}
	if degreeType, ok := updates["degree_type"]; ok {
		if degreeType == nil {
			prog.DegreeType = nil
		} else if dt, ok := degreeType.(string); ok {
			prog.DegreeType = &dt
		}
	}
	if duration, ok := updates["duration_years"].(int); ok {
		prog.DurationYears = duration
	}
	if credits, ok := updates["total_credits"]; ok {
		if credits == nil {
			prog.TotalCredits = nil
		} else if c, ok := credits.(int); ok {
			prog.TotalCredits = &c
		}
	}
	if desc, ok := updates["description"]; ok {
		if desc == nil {
			prog.Description = nil
		} else if d, ok := desc.(string); ok {
			prog.Description = &d
		}
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		prog.IsActive = isActive
	}

	if err := s.repo.Update(ctx, prog); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.program.updated", prog.ProgramID.String(), map[string]interface{}{
			"program_id":   prog.ProgramID,
			"program_name": prog.ProgramName,
		})
	}

	return nil
}

func (s *programService) DeleteProgram(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.program.deleted", id.String(), map[string]interface{}{
			"program_id": id,
		})
	}

	return nil
}

func (s *programService) ListPrograms(ctx context.Context, filter domain.ProgramFilter, page, limit int) ([]*domain.ProgramWithDepartment, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}
