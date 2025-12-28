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

func TestProgramService_CreateProgram(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProgramRepository(ctrl)
	mockDeptRepo := mocks.NewMockDepartmentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewProgramService(mockRepo, mockDeptRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		programID := uuid.New()
		deptID := uuid.New()

		program := &domain.Program{
			ProgramID:     programID,
			ProgramName:   "Computer Science",
			ProgramCode:   "CS",
			DepartmentID:  deptID,
			DurationYears: 4,
		}

		dept := &domain.Department{
			DepartmentID:   deptID,
			DepartmentName: "Engineering",
		}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Program) error {
			assert.True(t, p.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.program.created", programID.String(), gomock.Any()).Return(nil)

		err := service.CreateProgram(context.Background(), program)
		assert.NoError(t, err)
		assert.True(t, program.IsActive)
	})

	t.Run("Department Not Found", func(t *testing.T) {
		program := &domain.Program{
			ProgramID:    uuid.New(),
			DepartmentID: uuid.New(),
		}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), program.DepartmentID).Return(nil, domain.ErrDepartmentNotFound)

		err := service.CreateProgram(context.Background(), program)
		assert.ErrorIs(t, err, domain.ErrDepartmentNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		programID := uuid.New()
		deptID := uuid.New()

		program := &domain.Program{
			ProgramID:    programID,
			DepartmentID: deptID,
		}

		dept := &domain.Department{DepartmentID: deptID}
		repoErr := errors.New("database error")

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

		err := service.CreateProgram(context.Background(), program)
		assert.ErrorIs(t, err, repoErr)
	})
}

func TestProgramService_GetProgram(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProgramRepository(ctrl)

	service := NewProgramService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		programID := uuid.New()
		deptID := uuid.New()

		programWithDept := &domain.ProgramWithDepartment{
			Program: domain.Program{
				ProgramID:     programID,
				ProgramName:   "Computer Science",
				ProgramCode:   "CS",
				DurationYears: 4,
			},
			Department: domain.DepartmentBasic{
				DepartmentID:   deptID,
				DepartmentName: "Engineering",
			},
		}

		mockRepo.EXPECT().GetWithDepartment(gomock.Any(), programID).Return(programWithDept, nil)

		result, err := service.GetProgram(context.Background(), programID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Computer Science", result.ProgramName)
		assert.Equal(t, "Engineering", result.Department.DepartmentName)
	})

	t.Run("Program Not Found", func(t *testing.T) {
		programID := uuid.New()

		mockRepo.EXPECT().GetWithDepartment(gomock.Any(), programID).Return(nil, domain.ErrProgramNotFound)

		result, err := service.GetProgram(context.Background(), programID)
		assert.ErrorIs(t, err, domain.ErrProgramNotFound)
		assert.Nil(t, result)
	})
}

func TestProgramService_UpdateProgram(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProgramRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewProgramService(mockRepo, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		programID := uuid.New()

		program := &domain.Program{
			ProgramID:     programID,
			ProgramName:   "Old Name",
			DurationYears: 4,
			IsActive:      true,
		}

		updates := map[string]interface{}{
			"program_name":   "New Name",
			"duration_years": 5,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), programID).Return(program, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Program) error {
			assert.Equal(t, "New Name", p.ProgramName)
			assert.Equal(t, 5, p.DurationYears)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.program.updated", programID.String(), gomock.Any()).Return(nil)

		err := service.UpdateProgram(context.Background(), programID, updates)
		assert.NoError(t, err)
	})

	t.Run("Program Not Found", func(t *testing.T) {
		programID := uuid.New()
		updates := map[string]interface{}{"program_name": "New Name"}

		mockRepo.EXPECT().GetByID(gomock.Any(), programID).Return(nil, domain.ErrProgramNotFound)

		err := service.UpdateProgram(context.Background(), programID, updates)
		assert.ErrorIs(t, err, domain.ErrProgramNotFound)
	})

	t.Run("Update With Nullable Fields", func(t *testing.T) {
		programID := uuid.New()
		desc := "Old description"

		program := &domain.Program{
			ProgramID:   programID,
			ProgramName: "Test Program",
			Description: &desc,
		}

		updates := map[string]interface{}{
			"description": nil,
			"is_active":   false,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), programID).Return(program, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, p *domain.Program) error {
			assert.Nil(t, p.Description)
			assert.False(t, p.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.program.updated", programID.String(), gomock.Any()).Return(nil)

		err := service.UpdateProgram(context.Background(), programID, updates)
		assert.NoError(t, err)
	})
}

func TestProgramService_DeleteProgram(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProgramRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewProgramService(mockRepo, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		programID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), programID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.program.deleted", programID.String(), gomock.Any()).Return(nil)

		err := service.DeleteProgram(context.Background(), programID)
		assert.NoError(t, err)
	})

	t.Run("Program Not Found", func(t *testing.T) {
		programID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), programID).Return(domain.ErrProgramNotFound)

		err := service.DeleteProgram(context.Background(), programID)
		assert.ErrorIs(t, err, domain.ErrProgramNotFound)
	})
}

func TestProgramService_ListPrograms(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProgramRepository(ctrl)

	service := NewProgramService(mockRepo, nil, nil)

	t.Run("Success", func(t *testing.T) {
		deptID := uuid.New()
		filter := domain.ProgramFilter{
			DepartmentID: &deptID,
		}

		programs := []*domain.ProgramWithDepartment{
			{Program: domain.Program{ProgramID: uuid.New(), ProgramName: "CS"}},
			{Program: domain.Program{ProgramID: uuid.New(), ProgramName: "IT"}},
		}

		mockRepo.EXPECT().List(gomock.Any(), filter, 10, 0).Return(programs, int64(2), nil)

		result, total, err := service.ListPrograms(context.Background(), filter, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("Empty List", func(t *testing.T) {
		filter := domain.ProgramFilter{}

		mockRepo.EXPECT().List(gomock.Any(), filter, 20, 20).Return([]*domain.ProgramWithDepartment{}, int64(0), nil)

		result, total, err := service.ListPrograms(context.Background(), filter, 2, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})
}
