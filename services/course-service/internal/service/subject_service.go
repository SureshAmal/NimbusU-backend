package service

import (
	"context"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type subjectService struct {
	repo     domain.SubjectRepository
	deptRepo domain.DepartmentRepository
	producer domain.EventProducer
}

func NewSubjectService(repo domain.SubjectRepository, deptRepo domain.DepartmentRepository, producer domain.EventProducer) domain.SubjectService {
	return &subjectService{repo: repo, deptRepo: deptRepo, producer: producer}
}

func (s *subjectService) CreateSubject(ctx context.Context, subject *domain.Subject, prerequisites []domain.SubjectPrerequisite, corequisites []uuid.UUID) error {
	// Validate department exists
	_, err := s.deptRepo.GetByID(ctx, subject.DepartmentID)
	if err != nil {
		return err
	}

	subject.IsActive = true
	if err := s.repo.Create(ctx, subject); err != nil {
		return err
	}

	// Add prerequisites
	for _, prereq := range prerequisites {
		if err := s.repo.AddPrerequisite(ctx, subject.SubjectID, prereq.PrerequisiteSubjectID, prereq.IsMandatory); err != nil {
			// Log error but continue - prerequisite addition is non-critical
			continue
		}
	}

	// Add corequisites
	for _, coreqID := range corequisites {
		if err := s.repo.AddCorequisite(ctx, subject.SubjectID, coreqID); err != nil {
			continue
		}
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.subject.created", subject.SubjectID.String(), map[string]interface{}{
			"subject_id":    subject.SubjectID,
			"subject_name":  subject.SubjectName,
			"subject_code":  subject.SubjectCode,
			"department_id": subject.DepartmentID,
			"credits":       subject.Credits,
		})
	}

	return nil
}

func (s *subjectService) GetSubject(ctx context.Context, id uuid.UUID) (*domain.SubjectWithDetails, error) {
	return s.repo.GetWithDetails(ctx, id)
}

func (s *subjectService) UpdateSubject(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	subj, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates
	if name, ok := updates["subject_name"].(string); ok {
		subj.SubjectName = name
	}
	if credits, ok := updates["credits"].(int); ok {
		subj.Credits = credits
	}
	if subjectType, ok := updates["subject_type"]; ok {
		if subjectType == nil {
			subj.SubjectType = nil
		} else if st, ok := subjectType.(string); ok {
			subj.SubjectType = &st
		}
	}
	if desc, ok := updates["description"]; ok {
		if desc == nil {
			subj.Description = nil
		} else if d, ok := desc.(string); ok {
			subj.Description = &d
		}
	}
	if syllabus, ok := updates["syllabus"]; ok {
		if syllabus == nil {
			subj.Syllabus = nil
		} else if syl, ok := syllabus.(string); ok {
			subj.Syllabus = &syl
		}
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		subj.IsActive = isActive
	}

	if err := s.repo.Update(ctx, subj); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.subject.updated", subj.SubjectID.String(), map[string]interface{}{
			"subject_id":   subj.SubjectID,
			"subject_name": subj.SubjectName,
		})
	}

	return nil
}

func (s *subjectService) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.subject.deleted", id.String(), map[string]interface{}{
			"subject_id": id,
		})
	}

	return nil
}

func (s *subjectService) ListSubjects(ctx context.Context, filter domain.SubjectFilter, page, limit int) ([]*domain.Subject, int64, error) {
	offset := (page - 1) * limit
	return s.repo.List(ctx, filter, limit, offset)
}

func (s *subjectService) AddPrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID, isMandatory bool) error {
	// Validate both subjects exist
	if _, err := s.repo.GetByID(ctx, subjectID); err != nil {
		return err
	}
	if _, err := s.repo.GetByID(ctx, prerequisiteID); err != nil {
		return err
	}

	return s.repo.AddPrerequisite(ctx, subjectID, prerequisiteID, isMandatory)
}

func (s *subjectService) RemovePrerequisite(ctx context.Context, subjectID, prerequisiteID uuid.UUID) error {
	return s.repo.RemovePrerequisite(ctx, subjectID, prerequisiteID)
}

func (s *subjectService) AddCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error {
	// Validate both subjects exist
	if _, err := s.repo.GetByID(ctx, subjectID); err != nil {
		return err
	}
	if _, err := s.repo.GetByID(ctx, corequisiteID); err != nil {
		return err
	}

	return s.repo.AddCorequisite(ctx, subjectID, corequisiteID)
}

func (s *subjectService) RemoveCorequisite(ctx context.Context, subjectID, corequisiteID uuid.UUID) error {
	return s.repo.RemoveCorequisite(ctx, subjectID, corequisiteID)
}
