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

func TestCourseService_CreateCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCourseRepository(ctrl)
	mockSubjectRepo := mocks.NewMockSubjectRepository(ctrl)
	mockSemesterRepo := mocks.NewMockSemesterRepository(ctrl)
	mockEnrollmentRepo := mocks.NewMockEnrollmentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewCourseService(mockRepo, mockSubjectRepo, mockSemesterRepo, mockEnrollmentRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		courseID := uuid.New()
		subjectID := uuid.New()
		semesterID := uuid.New()

		course := &domain.Course{
			CourseID:   courseID,
			SubjectID:  subjectID,
			SemesterID: semesterID,
			CourseCode: "CS101",
		}

		subject := &domain.Subject{
			SubjectID:   subjectID,
			SubjectName: "Intro to CS",
		}

		semester := &domain.Semester{
			SemesterID: semesterID,
		}

		mockSubjectRepo.EXPECT().GetByID(gomock.Any(), subjectID).Return(subject, nil)
		mockSemesterRepo.EXPECT().GetByID(gomock.Any(), semesterID).Return(semester, nil)
		mockRepo.EXPECT().Create(gomock.Any(), course).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.course.created", courseID.String(), gomock.Any()).Return(nil)

		err := service.CreateCourse(context.Background(), course)
		assert.NoError(t, err)
		assert.Equal(t, "Intro to CS", course.CourseName) // Should use subject name if empty
		assert.Equal(t, "draft", course.Status)           // Default status
		assert.Equal(t, 0, course.CurrentEnrollment)      // Default enrollment
		assert.True(t, course.IsActive)                   // Default active
	})

	t.Run("Subject Not Found", func(t *testing.T) {
		course := &domain.Course{SubjectID: uuid.New()}
		mockSubjectRepo.EXPECT().GetByID(gomock.Any(), course.SubjectID).Return(nil, domain.ErrSubjectNotFound)

		err := service.CreateCourse(context.Background(), course)
		assert.ErrorIs(t, err, domain.ErrSubjectNotFound)
	})
}

func TestCourseService_UpdateCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCourseRepository(ctrl)
	mockSubjectRepo := mocks.NewMockSubjectRepository(ctrl)
	mockSemesterRepo := mocks.NewMockSemesterRepository(ctrl)
	mockEnrollmentRepo := mocks.NewMockEnrollmentRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewCourseService(mockRepo, mockSubjectRepo, mockSemesterRepo, mockEnrollmentRepo, mockProducer)

	t.Run("Success", func(t *testing.T) {
		courseID := uuid.New()
		course := &domain.Course{
			CourseID:   courseID,
			CourseName: "Old Name",
			Status:     "active",
		}

		updates := map[string]interface{}{
			"course_name": "New Name",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)
		mockRepo.EXPECT().Update(gomock.Any(), course).DoAndReturn(func(ctx context.Context, c *domain.Course) error {
			assert.Equal(t, "New Name", c.CourseName)
			return nil
		})
		mockProducer.EXPECT().PublishEvent("course.course.updated", courseID.String(), gomock.Any()).Return(nil)

		err := service.UpdateCourse(context.Background(), courseID, updates)
		assert.NoError(t, err)
	})

	t.Run("Cannot Modify Completed Course", func(t *testing.T) {
		courseID := uuid.New()
		course := &domain.Course{
			CourseID: courseID,
			Status:   "completed",
		}

		updates := map[string]interface{}{"course_name": "New Name"}

		mockRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)

		err := service.UpdateCourse(context.Background(), courseID, updates)
		assert.ErrorIs(t, err, domain.ErrCannotModifyCompletedCourse)
	})
}

func TestCourseService_ActivateCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCourseRepository(ctrl)
	// other mocks unused
	service := NewCourseService(mockRepo, nil, nil, nil, mocks.NewMockEventProducer(ctrl))

	t.Run("Success", func(t *testing.T) {
		courseID := uuid.New()
		course := &domain.Course{
			CourseID: courseID,
			Status:   "draft",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), courseID, "active").Return(nil)
		// We need the producer here because ActivateCourse calls PublishEvent
		service.(*courseService).producer.(*mocks.MockEventProducer).EXPECT().
			PublishEvent("course.course.activated", courseID.String(), gomock.Any()).Return(nil)

		err := service.ActivateCourse(context.Background(), courseID)
		assert.NoError(t, err)
	})
}
