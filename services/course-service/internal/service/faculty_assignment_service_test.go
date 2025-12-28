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

func TestFacultyAssignmentService_AssignFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyCourseRepository(ctrl)
	mockFacultyRepo := mocks.NewMockFacultyRepository(ctrl)
	mockCourseRepo := mocks.NewMockCourseRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewFacultyAssignmentService(mockRepo, mockFacultyRepo, mockCourseRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()
		assignedBy := uuid.New()

		faculty := &domain.Faculty{FacultyID: facultyID}
		course := &domain.Course{CourseID: courseID}

		mockFacultyRepo.EXPECT().GetByID(gomock.Any(), facultyID).Return(faculty, nil)
		mockCourseRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)
		mockRepo.EXPECT().GetAssignment(gomock.Any(), facultyID, courseID).Return(nil, domain.ErrAssignmentNotFound)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fc *domain.FacultyCourse) error {
			assert.Equal(t, facultyID, fc.FacultyID)
			assert.Equal(t, courseID, fc.CourseID)
			assert.Equal(t, "instructor", fc.Role)
			assert.True(t, fc.IsPrimary)
			assert.True(t, fc.IsActive)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.faculty.assigned", gomock.Any(), gomock.Any()).Return(nil)

		result, err := service.AssignFaculty(context.Background(), courseID, facultyID, assignedBy, "instructor", true)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, facultyID, result.FacultyID)
		assert.Equal(t, courseID, result.CourseID)
	})

	t.Run("Faculty Not Found", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()
		assignedBy := uuid.New()

		mockFacultyRepo.EXPECT().GetByID(gomock.Any(), facultyID).Return(nil, domain.ErrFacultyNotFound)

		result, err := service.AssignFaculty(context.Background(), courseID, facultyID, assignedBy, "instructor", true)
		assert.ErrorIs(t, err, domain.ErrFacultyNotFound)
		assert.Nil(t, result)
	})

	t.Run("Course Not Found", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()
		assignedBy := uuid.New()

		faculty := &domain.Faculty{FacultyID: facultyID}

		mockFacultyRepo.EXPECT().GetByID(gomock.Any(), facultyID).Return(faculty, nil)
		mockCourseRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(nil, domain.ErrCourseNotFound)

		result, err := service.AssignFaculty(context.Background(), courseID, facultyID, assignedBy, "instructor", true)
		assert.ErrorIs(t, err, domain.ErrCourseNotFound)
		assert.Nil(t, result)
	})

	t.Run("Faculty Already Assigned", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()
		assignedBy := uuid.New()

		faculty := &domain.Faculty{FacultyID: facultyID}
		course := &domain.Course{CourseID: courseID}
		existingAssignment := &domain.FacultyCourse{FacultyID: facultyID, CourseID: courseID}

		mockFacultyRepo.EXPECT().GetByID(gomock.Any(), facultyID).Return(faculty, nil)
		mockCourseRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)
		mockRepo.EXPECT().GetAssignment(gomock.Any(), facultyID, courseID).Return(existingAssignment, nil)

		result, err := service.AssignFaculty(context.Background(), courseID, facultyID, assignedBy, "instructor", true)
		assert.ErrorIs(t, err, domain.ErrFacultyAlreadyAssigned)
		assert.Nil(t, result)
	})
}

func TestFacultyAssignmentService_UpdateAssignment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyCourseRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewFacultyAssignmentService(mockRepo, nil, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()
		fcID := uuid.New()

		fc := &domain.FacultyCourse{
			FacultyCourseID: fcID,
			FacultyID:       facultyID,
			CourseID:        courseID,
			Role:            "ta",
			IsPrimary:       false,
		}

		mockRepo.EXPECT().GetAssignment(gomock.Any(), facultyID, courseID).Return(fc, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, f *domain.FacultyCourse) error {
			assert.Equal(t, "instructor", f.Role)
			assert.True(t, f.IsPrimary)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.faculty.assignment_updated", fcID.String(), gomock.Any()).Return(nil)

		err := service.UpdateAssignment(context.Background(), courseID, facultyID, "instructor", true)
		assert.NoError(t, err)
	})

	t.Run("Assignment Not Found", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()

		mockRepo.EXPECT().GetAssignment(gomock.Any(), facultyID, courseID).Return(nil, domain.ErrAssignmentNotFound)

		err := service.UpdateAssignment(context.Background(), courseID, facultyID, "instructor", true)
		assert.ErrorIs(t, err, domain.ErrAssignmentNotFound)
	})
}

func TestFacultyAssignmentService_RemoveFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyCourseRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewFacultyAssignmentService(mockRepo, nil, nil, mockProducer)

	t.Run("Success", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), facultyID, courseID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.faculty.removed", facultyID.String(), gomock.Any()).Return(nil)

		err := service.RemoveFaculty(context.Background(), courseID, facultyID)
		assert.NoError(t, err)
	})

	t.Run("Assignment Not Found", func(t *testing.T) {
		facultyID := uuid.New()
		courseID := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), facultyID, courseID).Return(domain.ErrAssignmentNotFound)

		err := service.RemoveFaculty(context.Background(), courseID, facultyID)
		assert.ErrorIs(t, err, domain.ErrAssignmentNotFound)
	})
}

func TestFacultyAssignmentService_ListCourseFaculty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFacultyCourseRepository(ctrl)

	service := NewFacultyAssignmentService(mockRepo, nil, nil, nil)

	t.Run("Success", func(t *testing.T) {
		courseID := uuid.New()

		facultyList := []*domain.FacultyCourseWithDetails{
			{
				FacultyCourse: domain.FacultyCourse{
					FacultyID: uuid.New(),
					CourseID:  courseID,
					Role:      "instructor",
					IsPrimary: true,
				},
				Faculty: domain.FacultyBasic{Name: "John Doe"},
			},
			{
				FacultyCourse: domain.FacultyCourse{
					FacultyID: uuid.New(),
					CourseID:  courseID,
					Role:      "ta",
					IsPrimary: false,
				},
				Faculty: domain.FacultyBasic{Name: "Jane Smith"},
			},
		}

		mockRepo.EXPECT().ListByCourse(gomock.Any(), courseID).Return(facultyList, nil)

		result, err := service.ListCourseFaculty(context.Background(), courseID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "John Doe", result[0].Faculty.Name)
	})

	t.Run("Empty List", func(t *testing.T) {
		courseID := uuid.New()

		mockRepo.EXPECT().ListByCourse(gomock.Any(), courseID).Return([]*domain.FacultyCourseWithDetails{}, nil)

		result, err := service.ListCourseFaculty(context.Background(), courseID)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}
