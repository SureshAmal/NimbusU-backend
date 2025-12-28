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

func TestStudentService_CreateStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	mockDeptRepo := mocks.NewMockDepartmentRepository(ctrl)
	mockProgRepo := mocks.NewMockProgramRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewStudentService(mockRepo, mockDeptRepo, mockProgRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		studentID := uuid.New()
		userID := uuid.New()
		deptID := uuid.New()
		programID := uuid.New()

		student := &domain.Student{
			StudentID:          studentID,
			UserID:             userID,
			RegistrationNumber: "STU001",
			DepartmentID:       deptID,
			ProgramID:          programID,
			BatchYear:          2024,
		}

		dept := &domain.Department{DepartmentID: deptID}
		program := &domain.Program{ProgramID: programID}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockProgRepo.EXPECT().GetByID(gomock.Any(), programID).Return(program, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Student) error {
			assert.True(t, s.IsActive)
			assert.Equal(t, 1, s.CurrentSemester)
			assert.Equal(t, 0, s.TotalCreditsEarned)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.student.created", studentID.String(), gomock.Any()).Return(nil)

		err := service.CreateStudent(context.Background(), student)
		assert.NoError(t, err)
		assert.True(t, student.IsActive)
		assert.Equal(t, 1, student.CurrentSemester)
	})

	t.Run("Department Not Found", func(t *testing.T) {
		student := &domain.Student{
			StudentID:    uuid.New(),
			DepartmentID: uuid.New(),
			ProgramID:    uuid.New(),
		}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), student.DepartmentID).Return(nil, domain.ErrDepartmentNotFound)

		err := service.CreateStudent(context.Background(), student)
		assert.ErrorIs(t, err, domain.ErrDepartmentNotFound)
	})

	t.Run("Program Not Found", func(t *testing.T) {
		student := &domain.Student{
			StudentID:    uuid.New(),
			DepartmentID: uuid.New(),
			ProgramID:    uuid.New(),
		}

		dept := &domain.Department{DepartmentID: student.DepartmentID}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), student.DepartmentID).Return(dept, nil)
		mockProgRepo.EXPECT().GetByID(gomock.Any(), student.ProgramID).Return(nil, domain.ErrProgramNotFound)

		err := service.CreateStudent(context.Background(), student)
		assert.ErrorIs(t, err, domain.ErrProgramNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		studentID := uuid.New()
		deptID := uuid.New()
		programID := uuid.New()

		student := &domain.Student{
			StudentID:    studentID,
			DepartmentID: deptID,
			ProgramID:    programID,
		}

		dept := &domain.Department{DepartmentID: deptID}
		program := &domain.Program{ProgramID: programID}
		repoErr := errors.New("database error")

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockProgRepo.EXPECT().GetByID(gomock.Any(), programID).Return(program, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

		err := service.CreateStudent(context.Background(), student)
		assert.ErrorIs(t, err, repoErr)
	})

	t.Run("With Existing Semester", func(t *testing.T) {
		studentID := uuid.New()
		deptID := uuid.New()
		programID := uuid.New()

		student := &domain.Student{
			StudentID:       studentID,
			DepartmentID:    deptID,
			ProgramID:       programID,
			CurrentSemester: 3, // Already set
		}

		dept := &domain.Department{DepartmentID: deptID}
		program := &domain.Program{ProgramID: programID}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockProgRepo.EXPECT().GetByID(gomock.Any(), programID).Return(program, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Student) error {
			assert.Equal(t, 3, s.CurrentSemester) // Should keep existing value
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.student.created", studentID.String(), gomock.Any()).Return(nil)

		err := service.CreateStudent(context.Background(), student)
		assert.NoError(t, err)
	})
}

func TestStudentService_GetStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)

	service := NewStudentService(mockRepo, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		studentID := uuid.New()

		studentDetails := &domain.StudentWithDetails{
			Student: domain.Student{
				StudentID:          studentID,
				RegistrationNumber: "STU001",
				IsActive:           true,
			},
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.EXPECT().GetWithDetails(gomock.Any(), studentID).Return(studentDetails, nil)

		result, err := service.GetStudent(context.Background(), studentID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "John Doe", result.Name)
	})

	t.Run("Student Not Found", func(t *testing.T) {
		studentID := uuid.New()

		mockRepo.EXPECT().GetWithDetails(gomock.Any(), studentID).Return(nil, domain.ErrStudentNotFound)

		result, err := service.GetStudent(context.Background(), studentID)
		assert.ErrorIs(t, err, domain.ErrStudentNotFound)
		assert.Nil(t, result)
	})
}

func TestStudentService_GetStudentByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)

	service := NewStudentService(mockRepo, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		studentID := uuid.New()

		student := &domain.Student{
			StudentID: studentID,
			UserID:    userID,
		}

		studentDetails := &domain.StudentWithDetails{
			Student: *student,
			Name:    "John Doe",
		}

		mockRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(student, nil)
		mockRepo.EXPECT().GetWithDetails(gomock.Any(), studentID).Return(studentDetails, nil)

		result, err := service.GetStudentByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, studentID, result.StudentID)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := uuid.New()

		mockRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(nil, domain.ErrStudentNotFound)

		result, err := service.GetStudentByUserID(context.Background(), userID)
		assert.ErrorIs(t, err, domain.ErrStudentNotFound)
		assert.Nil(t, result)
	})
}

func TestStudentService_UpdateStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewStudentService(mockRepo, nil, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		studentID := uuid.New()
		rollNumber := "ROLL001"

		student := &domain.Student{
			StudentID:          studentID,
			RegistrationNumber: "STU001",
			IsActive:           true,
		}

		updates := map[string]interface{}{
			"roll_number": rollNumber,
			"is_active":   false,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), studentID).Return(student, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Student) error {
			assert.Equal(t, rollNumber, *s.RollNumber)
			assert.False(t, s.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.student.updated", studentID.String(), gomock.Any()).Return(nil)

		err := service.UpdateStudent(context.Background(), studentID, updates)
		assert.NoError(t, err)
	})

	t.Run("Student Not Found", func(t *testing.T) {
		studentID := uuid.New()
		updates := map[string]interface{}{"roll_number": "ROLL001"}

		mockRepo.EXPECT().GetByID(gomock.Any(), studentID).Return(nil, domain.ErrStudentNotFound)

		err := service.UpdateStudent(context.Background(), studentID, updates)
		assert.ErrorIs(t, err, domain.ErrStudentNotFound)
	})

	t.Run("Clear Roll Number", func(t *testing.T) {
		studentID := uuid.New()
		rollNumber := "ROLL001"

		student := &domain.Student{
			StudentID:  studentID,
			RollNumber: &rollNumber,
		}

		updates := map[string]interface{}{
			"roll_number": nil,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), studentID).Return(student, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *domain.Student) error {
			assert.Nil(t, s.RollNumber)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.student.updated", studentID.String(), gomock.Any()).Return(nil)

		err := service.UpdateStudent(context.Background(), studentID, updates)
		assert.NoError(t, err)
	})
}

func TestStudentService_DeleteStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewStudentService(mockRepo, nil, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		studentID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), studentID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.student.deleted", studentID.String(), gomock.Any()).Return(nil)

		err := service.DeleteStudent(context.Background(), studentID)
		assert.NoError(t, err)
	})

	t.Run("Student Not Found", func(t *testing.T) {
		studentID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), studentID).Return(domain.ErrStudentNotFound)

		err := service.DeleteStudent(context.Background(), studentID)
		assert.ErrorIs(t, err, domain.ErrStudentNotFound)
	})
}

func TestStudentService_ListStudents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)

	service := NewStudentService(mockRepo, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		deptID := uuid.New()
		filter := domain.StudentFilter{
			DepartmentID: &deptID,
		}

		students := []*domain.StudentWithDetails{
			{Student: domain.Student{StudentID: uuid.New()}, Name: "John Doe"},
			{Student: domain.Student{StudentID: uuid.New()}, Name: "Jane Smith"},
		}

		mockRepo.EXPECT().List(gomock.Any(), filter, 10, 0).Return(students, int64(2), nil)

		result, total, err := service.ListStudents(context.Background(), filter, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("Empty List", func(t *testing.T) {
		filter := domain.StudentFilter{}

		mockRepo.EXPECT().List(gomock.Any(), filter, 20, 20).Return([]*domain.StudentWithDetails{}, int64(0), nil)

		result, total, err := service.ListStudents(context.Background(), filter, 2, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})
}

func TestStudentService_PromoteStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewStudentService(mockRepo, nil, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		studentID := uuid.New()
		cgpa := 3.5

		mockRepo.EXPECT().UpdateSemester(gomock.Any(), studentID, 4, &cgpa, 60).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.student.promoted", studentID.String(), gomock.Any()).Return(nil)

		err := service.PromoteStudent(context.Background(), studentID, 4, &cgpa, 60)
		assert.NoError(t, err)
	})

	t.Run("Student Not Found", func(t *testing.T) {
		studentID := uuid.New()

		mockRepo.EXPECT().UpdateSemester(gomock.Any(), studentID, 4, (*float64)(nil), 60).Return(domain.ErrStudentNotFound)

		err := service.PromoteStudent(context.Background(), studentID, 4, nil, 60)
		assert.ErrorIs(t, err, domain.ErrStudentNotFound)
	})
}
