package service

import (
	"context"
	"errors"
	"testing"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSubjectService_CreateSubject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)
	mockDeptRepo := mocks.NewMockDepartmentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewSubjectService(mockRepo, mockDeptRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()
		deptID := uuid.New()

		subject := &domain.Subject{
			SubjectID:    subjectID,
			SubjectName:  "Data Structures",
			SubjectCode:  "CS201",
			DepartmentID: deptID,
			Credits:      3,
		}

		dept := &domain.Department{DepartmentID: deptID}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Subject) error {
			assert.True(t, s.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.subject.created", subjectID.String(), gomock.Any()).Return(nil)

		err := service.CreateSubject(context.Background(), subject, nil, nil)
		assert.NoError(t, err)
		assert.True(t, subject.IsActive)
	})

	t.Run("With Prerequisites", func(t *testing.T) {
		subjectID := uuid.New()
		deptID := uuid.New()
		prereqID := uuid.New()

		subject := &domain.Subject{
			SubjectID:    subjectID,
			SubjectName:  "Advanced Algorithms",
			SubjectCode:  "CS301",
			DepartmentID: deptID,
			Credits:      3,
		}

		prerequisites := []domain.SubjectPrerequisite{
			{PrerequisiteSubjectID: prereqID, IsMandatory: true},
		}

		dept := &domain.Department{DepartmentID: deptID}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().AddPrerequisite(gomock.Any(), subjectID, prereqID, true).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.subject.created", subjectID.String(), gomock.Any()).Return(nil)

		err := service.CreateSubject(context.Background(), subject, prerequisites, nil)
		assert.NoError(t, err)
	})

	t.Run("With Corequisites", func(t *testing.T) {
		subjectID := uuid.New()
		deptID := uuid.New()
		coreqID := uuid.New()

		subject := &domain.Subject{
			SubjectID:    subjectID,
			SubjectName:  "Physics Lab",
			SubjectCode:  "PHY101L",
			DepartmentID: deptID,
			Credits:      1,
		}

		corequisites := []uuid.UUID{coreqID}

		dept := &domain.Department{DepartmentID: deptID}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().AddCorequisite(gomock.Any(), subjectID, coreqID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.subject.created", subjectID.String(), gomock.Any()).Return(nil)

		err := service.CreateSubject(context.Background(), subject, nil, corequisites)
		assert.NoError(t, err)
	})

	t.Run("Department Not Found", func(t *testing.T) {
		subject := &domain.Subject{
			SubjectID:    uuid.New(),
			DepartmentID: uuid.New(),
		}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), subject.DepartmentID).Return(nil, domain.ErrDepartmentNotFound)

		err := service.CreateSubject(context.Background(), subject, nil, nil)
		assert.ErrorIs(t, err, domain.ErrDepartmentNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		subjectID := uuid.New()
		deptID := uuid.New()

		subject := &domain.Subject{
			SubjectID:    subjectID,
			DepartmentID: deptID,
		}

		dept := &domain.Department{DepartmentID: deptID}
		repoErr := errors.New("database error")

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

		err := service.CreateSubject(context.Background(), subject, nil, nil)
		assert.ErrorIs(t, err, repoErr)
	})
}

func TestSubjectService_GetSubject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)

	service := NewSubjectService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()
		deptID := uuid.New()

		subjectDetails := &domain.SubjectWithDetails{
			Subject: domain.Subject{
				SubjectID:   subjectID,
				SubjectName: "Data Structures",
				SubjectCode: "CS201",
				Credits:     3,
			},
			Department: domain.DepartmentBasic{
				DepartmentID:   deptID,
				DepartmentName: "Computer Science",
			},
		}

		mockRepo.EXPECT().GetWithDetails(gomock.Any(), subjectID).Return(subjectDetails, nil)

		result, err := service.GetSubject(context.Background(), subjectID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Data Structures", result.SubjectName)
	})

	t.Run("Subject Not Found", func(t *testing.T) {
		subjectID := uuid.New()

		mockRepo.EXPECT().GetWithDetails(gomock.Any(), subjectID).Return(nil, domain.ErrSubjectNotFound)

		result, err := service.GetSubject(context.Background(), subjectID)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
		assert.Nil(t, result)
	})
}

func TestSubjectService_UpdateSubject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewSubjectService(mockRepo, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()

		subject := &domain.Subject{
			SubjectID:   subjectID,
			SubjectName: "Old Name",
			Credits:     3,
		}

		updates := map[string]interface{}{
			"subject_name": "New Name",
			"credits":      4,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(subject, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Subject) error {
			assert.Equal(t, "New Name", s.SubjectName)
			assert.Equal(t, 4, s.Credits)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.subject.updated", subjectID.String(), gomock.Any()).Return(nil)

		err := service.UpdateSubject(context.Background(), subjectID, updates)
		assert.NoError(t, err)
	})

	t.Run("Subject Not Found", func(t *testing.T) {
		subjectID := uuid.New()
		updates := map[string]interface{}{"subject_name": "New Name"}

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(nil, domain.ErrSubjectNotFound)

		err := service.UpdateSubject(context.Background(), subjectID, updates)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
	})

	t.Run("Update With Nullable Fields", func(t *testing.T) {
		subjectID := uuid.New()
		desc := "Old description"
		subjectType := "core"

		subject := &domain.Subject{
			SubjectID:   subjectID,
			SubjectName: "Test Subject",
			Description: &desc,
			SubjectType: &subjectType,
		}

		updates := map[string]interface{}{
			"description":  nil,
			"subject_type": "elective",
			"is_active":    false,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(subject, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Subject) error {
			assert.Nil(t, s.Description)
			assert.Equal(t, "elective", *s.SubjectType)
			assert.False(t, s.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.subject.updated", subjectID.String(), gomock.Any()).Return(nil)

		err := service.UpdateSubject(context.Background(), subjectID, updates)
		assert.NoError(t, err)
	})
}

func TestSubjectService_DeleteSubject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewSubjectService(mockRepo, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), subjectID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.subject.deleted", subjectID.String(), gomock.Any()).Return(nil)

		err := service.DeleteSubject(context.Background(), subjectID)
		assert.NoError(t, err)
	})

	t.Run("Subject Not Found", func(t *testing.T) {
		subjectID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), subjectID).Return(domain.ErrSubjectNotFound)

		err := service.DeleteSubject(context.Background(), subjectID)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
	})
}

func TestSubjectService_ListSubjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)

	service := NewSubjectService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		deptID := uuid.New()
		filter := domain.SubjectFilter{
			DepartmentID: &deptID,
		}

		subjects := []*domain.Subject{
			{SubjectID: uuid.New(), SubjectName: "Data Structures"},
			{SubjectID: uuid.New(), SubjectName: "Algorithms"},
		}

		mockRepo.EXPECT().List(gomock.Any(), filter, 10, 0).Return(subjects, int64(2), nil)

		result, total, err := service.ListSubjects(context.Background(), filter, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("Empty List", func(t *testing.T) {
		filter := domain.SubjectFilter{}

		mockRepo.EXPECT().List(gomock.Any(), filter, 20, 20).Return([]*domain.Subject{}, int64(0), nil)

		result, total, err := service.ListSubjects(context.Background(), filter, 2, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})
}

func TestSubjectService_AddPrerequisite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)

	service := NewSubjectService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()
		prereqID := uuid.New()

		subject := &domain.Subject{SubjectID: subjectID}
		prereq := &domain.Subject{SubjectID: prereqID}

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(subject, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), prereqID).Return(prereq, nil)
		mockRepo.EXPECT().AddPrerequisite(gomock.Any(), subjectID, prereqID, true).Return(nil)

		err := service.AddPrerequisite(context.Background(), subjectID, prereqID, true)
		assert.NoError(t, err)
	})

	t.Run("Subject Not Found", func(t *testing.T) {
		subjectID := uuid.New()
		prereqID := uuid.New()

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(nil, domain.ErrSubjectNotFound)

		err := service.AddPrerequisite(context.Background(), subjectID, prereqID, true)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
	})

	t.Run("Prerequisite Not Found", func(t *testing.T) {
		subjectID := uuid.New()
		prereqID := uuid.New()

		subject := &domain.Subject{SubjectID: subjectID}

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(subject, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), prereqID).Return(nil, domain.ErrSubjectNotFound)

		err := service.AddPrerequisite(context.Background(), subjectID, prereqID, true)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
	})
}

func TestSubjectService_RemovePrerequisite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)

	service := NewSubjectService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()
		prereqID := uuid.New()

		mockRepo.EXPECT().RemovePrerequisite(gomock.Any(), subjectID, prereqID).Return(nil)

		err := service.RemovePrerequisite(context.Background(), subjectID, prereqID)
		assert.NoError(t, err)
	})
}

func TestSubjectService_AddCorequisite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)

	service := NewSubjectService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()
		coreqID := uuid.New()

		subject := &domain.Subject{SubjectID: subjectID}
		coreq := &domain.Subject{SubjectID: coreqID}

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(subject, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), coreqID).Return(coreq, nil)
		mockRepo.EXPECT().AddCorequisite(gomock.Any(), subjectID, coreqID).Return(nil)

		err := service.AddCorequisite(context.Background(), subjectID, coreqID)
		assert.NoError(t, err)
	})

	t.Run("Subject Not Found", func(t *testing.T) {
		subjectID := uuid.New()
		coreqID := uuid.New()

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(nil, domain.ErrSubjectNotFound)

		err := service.AddCorequisite(context.Background(), subjectID, coreqID)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
	})

	t.Run("Corequisite Not Found", func(t *testing.T) {
		subjectID := uuid.New()
		coreqID := uuid.New()

		subject := &domain.Subject{SubjectID: subjectID}

		mockRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(subject, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), coreqID).Return(nil, domain.ErrSubjectNotFound)

		err := service.AddCorequisite(context.Background(), subjectID, coreqID)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
	})
}

func TestSubjectService_RemoveCorequisite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)

	service := NewSubjectService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		subjectID := uuid.New()
		coreqID := uuid.New()

		mockRepo.EXPECT().RemoveCorequisite(gomock.Any(), subjectID, coreqID).Return(nil)

		err := service.RemoveCorequisite(context.Background(), subjectID, coreqID)
		assert.NoError(t, err)
	})
}
