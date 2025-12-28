package service

import (
	"context"
	"testing"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDepartmentService_CreateDepartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDepartmentRepository(ctrl)
	// EventProducer is in service.go, so it's in service_mock.go (same package mocks)
	// But generated mocks for service.go are also in package mocks.
	// We need to check if EventProducer mock is generated. It should be.
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewDepartmentService(mockRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		deptID := uuid.New()
		dept := &domain.Department{
			DepartmentID:   deptID,
			DepartmentName: "Computer Science",
			DepartmentCode: "CS",
		}

		mockRepo.EXPECT().Create(gomock.Any(), dept).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.department.created", deptID.String(), gomock.Any()).Return(nil)

		err := service.CreateDepartment(context.Background(), dept)
		assert.NoError(t, err)
		assert.True(t, dept.IsActive) // Service sets IsActive = true
	})

	t.Run("Repo Error", func(t *testing.T) {
		dept := &domain.Department{
			DepartmentName: "Error Dept",
		}

		mockRepo.EXPECT().Create(gomock.Any(), dept).Return(domain.ErrDepartmentCodeExists)

		err := service.CreateDepartment(context.Background(), dept)
		assert.ErrorIs(t, err, domain.ErrDepartmentCodeExists)
	})
}

func TestDepartmentService_GetDepartment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDepartmentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewDepartmentService(mockRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		deptID := uuid.New()
		expectedDept := &domain.DepartmentWithDetails{
			Department: domain.Department{
				DepartmentID:   deptID,
				DepartmentName: "Computer Science",
			},
		}

		mockRepo.EXPECT().GetWithDetails(gomock.Any(), deptID).Return(expectedDept, nil)

		result, err := service.GetDepartment(context.Background(), deptID)
		assert.NoError(t, err)
		assert.Equal(t, expectedDept, result)
	})

	t.Run("Not Found", func(t *testing.T) {
		deptID := uuid.New()
		mockRepo.EXPECT().GetWithDetails(gomock.Any(), deptID).Return(nil, domain.ErrDepartmentNotFound)

		result, err := service.GetDepartment(context.Background(), deptID)
		assert.ErrorIs(t, err, domain.ErrDepartmentNotFound)
		assert.Nil(t, result)
	})
}
