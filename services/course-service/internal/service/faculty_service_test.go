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

func TestFacultyService_CreateFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyRepository(ctrl)
	mockDeptRepo := mocks.NewMockDepartmentRepository(ctrl)
	mockFCRepo := mocks.NewMockFacultyCourseRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewFacultyService(mockRepo, mockDeptRepo, mockFCRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()
		userID := uuid.New()
		deptID := uuid.New()

		faculty := &domain.Faculty{
			FacultyID:    facultyID,
			UserID:       userID,
			EmployeeID:   "EMP001",
			DepartmentID: deptID,
		}

		dept := &domain.Department{
			DepartmentID:   deptID,
			DepartmentName: "Computer Science",
		}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f *domain.Faculty) error {
			assert.True(t, f.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.faculty.created", facultyID.String(), gomock.Any()).Return(nil)

		err := service.CreateFaculty(context.Background(), faculty)
		assert.NoError(t, err)
		assert.True(t, faculty.IsActive)
	})

	t.Run("Department Not Found", func(t *testing.T) {
		faculty := &domain.Faculty{
			FacultyID:    uuid.New(),
			DepartmentID: uuid.New(),
		}

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), faculty.DepartmentID).Return(nil, domain.ErrDepartmentNotFound)

		err := service.CreateFaculty(context.Background(), faculty)
		assert.ErrorIs(t, err, domain.ErrDepartmentNotFound)
	})

	t.Run("Repository Error", func(t *testing.T) {
		facultyID := uuid.New()
		deptID := uuid.New()

		faculty := &domain.Faculty{
			FacultyID:    facultyID,
			DepartmentID: deptID,
		}

		dept := &domain.Department{DepartmentID: deptID}
		repoErr := errors.New("database error")

		mockDeptRepo.EXPECT().GetByID(gomock.Any(), deptID).Return(dept, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repoErr)

		err := service.CreateFaculty(context.Background(), faculty)
		assert.ErrorIs(t, err, repoErr)
	})
}

func TestFacultyService_GetFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyRepository(ctrl)

	service := NewFacultyService(mockRepo, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()
		deptID := uuid.New()

		facultyDetails := &domain.FacultyWithDetails{
			Faculty: domain.Faculty{
				FacultyID:    facultyID,
				EmployeeID:   "EMP001",
				DepartmentID: deptID,
				IsActive:     true,
			},
			Department: domain.DepartmentBasic{
				DepartmentID:   deptID,
				DepartmentName: "Computer Science",
			},
			Name:  "John Doe",
			Email: "john@example.com",
		}

		mockRepo.EXPECT().GetWithDetails(gomock.Any(), facultyID).Return(facultyDetails, nil)

		result, err := service.GetFaculty(context.Background(), facultyID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "John Doe", result.Name)
	})

	t.Run("Faculty Not Found", func(t *testing.T) {
		facultyID := uuid.New()

		mockRepo.EXPECT().GetWithDetails(gomock.Any(), facultyID).Return(nil, domain.ErrFacultyNotFound)

		result, err := service.GetFaculty(context.Background(), facultyID)
		assert.ErrorIs(t, err, domain.ErrFacultyNotFound)
		assert.Nil(t, result)
	})
}

func TestFacultyService_GetFacultyByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyRepository(ctrl)

	service := NewFacultyService(mockRepo, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		facultyID := uuid.New()

		faculty := &domain.Faculty{
			FacultyID:  facultyID,
			UserID:     userID,
			EmployeeID: "EMP001",
		}

		facultyDetails := &domain.FacultyWithDetails{
			Faculty: *faculty,
			Name:    "John Doe",
		}

		mockRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(faculty, nil)
		mockRepo.EXPECT().GetWithDetails(gomock.Any(), facultyID).Return(facultyDetails, nil)

		result, err := service.GetFacultyByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, facultyID, result.FacultyID)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := uuid.New()

		mockRepo.EXPECT().GetByUserID(gomock.Any(), userID).Return(nil, domain.ErrFacultyNotFound)

		result, err := service.GetFacultyByUserID(context.Background(), userID)
		assert.ErrorIs(t, err, domain.ErrFacultyNotFound)
		assert.Nil(t, result)
	})
}

func TestFacultyService_UpdateFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewFacultyService(mockRepo, nil, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()
		designation := "Professor"

		faculty := &domain.Faculty{
			FacultyID:  facultyID,
			EmployeeID: "EMP001",
			IsActive:   true,
		}

		updates := map[string]interface{}{
			"designation": designation,
			"is_active":   false,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), facultyID).Return(faculty, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f *domain.Faculty) error {
			assert.Equal(t, designation, *f.Designation)
			assert.False(t, f.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.faculty.updated", facultyID.String(), gomock.Any()).Return(nil)

		err := service.UpdateFaculty(context.Background(), facultyID, updates)
		assert.NoError(t, err)
	})

	t.Run("Faculty Not Found", func(t *testing.T) {
		facultyID := uuid.New()
		updates := map[string]interface{}{"designation": "Professor"}

		mockRepo.EXPECT().GetByID(gomock.Any(), facultyID).Return(nil, domain.ErrFacultyNotFound)

		err := service.UpdateFaculty(context.Background(), facultyID, updates)
		assert.ErrorIs(t, err, domain.ErrFacultyNotFound)
	})
}

func TestFacultyService_DeleteFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewFacultyService(mockRepo, nil, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), facultyID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.faculty.deleted", facultyID.String(), gomock.Any()).Return(nil)

		err := service.DeleteFaculty(context.Background(), facultyID)
		assert.NoError(t, err)
	})

	t.Run("Faculty Not Found", func(t *testing.T) {
		facultyID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), facultyID).Return(domain.ErrFacultyNotFound)

		err := service.DeleteFaculty(context.Background(), facultyID)
		assert.ErrorIs(t, err, domain.ErrFacultyNotFound)
	})
}

func TestFacultyService_ListFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyRepository(ctrl)

	service := NewFacultyService(mockRepo, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		deptID := uuid.New()
		filter := domain.FacultyFilter{
			DepartmentID: &deptID,
		}

		facultyList := []*domain.FacultyWithDetails{
			{Faculty: domain.Faculty{FacultyID: uuid.New()}, Name: "John Doe"},
			{Faculty: domain.Faculty{FacultyID: uuid.New()}, Name: "Jane Smith"},
		}

		mockRepo.EXPECT().List(gomock.Any(), filter, 10, 0).Return(facultyList, int64(2), nil)

		result, total, err := service.ListFaculty(context.Background(), filter, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("Empty List", func(t *testing.T) {
		filter := domain.FacultyFilter{}

		mockRepo.EXPECT().List(gomock.Any(), filter, 20, 20).Return([]*domain.FacultyWithDetails{}, int64(0), nil)

		result, total, err := service.ListFaculty(context.Background(), filter, 2, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})
}

func TestFacultyService_GetFacultyCourses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFCRepo := mocks.NewMockFacultyCourseRepository(ctrl)

	service := NewFacultyService(nil, nil, mockFCRepo, nil)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()
		semesterID := uuid.New()

		courses := []*domain.FacultyCourse{
			{FacultyCourseID: uuid.New(), FacultyID: facultyID, Role: "instructor"},
			{FacultyCourseID: uuid.New(), FacultyID: facultyID, Role: "ta"},
			{FacultyCourseID: uuid.New(), FacultyID: facultyID, Role: "instructor"},
		}

		mockFCRepo.EXPECT().ListByFaculty(gomock.Any(), facultyID, &semesterID).Return(courses, nil)

		result, total, err := service.GetFacultyCourses(context.Background(), facultyID, &semesterID, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, result, 2) // First page with limit 2
	})

	t.Run("Empty Courses", func(t *testing.T) {
		facultyID := uuid.New()

		mockFCRepo.EXPECT().ListByFaculty(gomock.Any(), facultyID, (*uuid.UUID)(nil)).Return([]*domain.FacultyCourse{}, nil)

		result, total, err := service.GetFacultyCourses(context.Background(), facultyID, nil, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, result)
	})

	t.Run("Pagination Beyond Results", func(t *testing.T) {
		facultyID := uuid.New()

		courses := []*domain.FacultyCourse{
			{FacultyCourseID: uuid.New(), FacultyID: facultyID},
		}

		mockFCRepo.EXPECT().ListByFaculty(gomock.Any(), facultyID, (*uuid.UUID)(nil)).Return(courses, nil)

		result, total, err := service.GetFacultyCourses(context.Background(), facultyID, nil, 5, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Empty(t, result) // Page 5 is beyond results
	})
}
