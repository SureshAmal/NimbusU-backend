package service

import (
	"context"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type departmentService struct {
	repo     domain.DepartmentRepository
	producer domain.EventProducer
}

func NewDepartmentService(repo domain.DepartmentRepository, producer domain.EventProducer) domain.DepartmentService {
	return &departmentService{repo: repo, producer: producer}
}

func (s *departmentService) CreateDepartment(ctx context.Context, department *domain.Department) error {
	department.IsActive = true
	if err := s.repo.Create(ctx, department); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.department.created", department.DepartmentID.String(), map[string]interface{}{
			"department_id":   department.DepartmentID,
			"department_name": department.DepartmentName,
			"department_code": department.DepartmentCode,
		})
	}

	return nil
}

func (s *departmentService) GetDepartment(ctx context.Context, id uuid.UUID) (*domain.DepartmentWithDetails, error) {
	return s.repo.GetWithDetails(ctx, id)
}

func (s *departmentService) UpdateDepartment(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	dept, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if name, ok := updates["department_name"].(string); ok {
		dept.DepartmentName = name
	}
	if hod, ok := updates["head_of_department"]; ok {
		if hod == nil {
			dept.HeadOfDepartment = nil
		} else if hodStr, ok := hod.(string); ok {
			hodID, err := uuid.Parse(hodStr)
			if err == nil {
				dept.HeadOfDepartment = &hodID
			}
		}
	}
	if desc, ok := updates["description"]; ok {
		if desc == nil {
			dept.Description = nil
		} else if descStr, ok := desc.(string); ok {
			dept.Description = &descStr
		}
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		dept.IsActive = isActive
	}

	if err := s.repo.Update(ctx, dept); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.department.updated", dept.DepartmentID.String(), map[string]interface{}{
			"department_id":   dept.DepartmentID,
			"department_name": dept.DepartmentName,
		})
	}

	return nil
}

func (s *departmentService) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.department.deleted", id.String(), map[string]interface{}{
			"department_id": id,
		})
	}

	return nil
}

func (s *departmentService) ListDepartments(ctx context.Context, filter domain.DepartmentFilter, page, limit int) ([]*domain.Department, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}
