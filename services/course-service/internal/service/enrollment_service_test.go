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

func TestEnrollmentService_EnrollStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepository(ctrl)
	mockCourseRepo := mocks.NewMockCourseRepository(ctrl)
	mockStudentRepo := mocks.NewMockStudentRepository(ctrl)
	// unused mocks for this test
	mockSubjectRepo := mocks.NewMockSubjectRepository(ctrl)
	mockSemesterRepo := mocks.NewMockSemesterRepository(ctrl)
	mockProducer := mocks.NewMockEventProducer(ctrl)

	service := NewEnrollmentService(mockRepo, mockCourseRepo, mockStudentRepo, mockSubjectRepo, mockSemesterRepo, mockProducer)

	t.Run("Success Enrolled", func(t *testing.T) {
		studentID := uuid.New()
		courseID := uuid.New()
		maxStudents := 50
		course := &domain.Course{
			CourseID:          courseID,
			Status:            "active",
			MaxStudents:       &maxStudents,
			CurrentEnrollment: 10,
		}

		student := &domain.Student{StudentID: studentID}

		mockStudentRepo.EXPECT().GetByID(gomock.Any(), studentID).Return(student, nil)
		mockCourseRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)
		mockRepo.EXPECT().GetByStudentAndCourse(gomock.Any(), studentID, courseID).Return(nil, domain.ErrEnrollmentNotFound)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, e *domain.CourseEnrollment) error {
			assert.Equal(t, "enrolled", e.EnrollmentStatus)
			return nil
		})

		mockCourseRepo.EXPECT().IncrementEnrollment(gomock.Any(), courseID).Return(nil)
		mockProducer.EXPECT().PublishEvent("course.enrollment.created", gomock.Any(), gomock.Any()).Return(nil)

		enrollment, err := service.EnrollStudent(context.Background(), courseID, studentID, "admin")
		assert.NoError(t, err)
		assert.Equal(t, "enrolled", enrollment.EnrollmentStatus)
	})

	t.Run("Waitlisted", func(t *testing.T) {
		studentID := uuid.New()
		courseID := uuid.New()
		maxStudents := 10
		course := &domain.Course{
			CourseID:          courseID,
			Status:            "active",
			MaxStudents:       &maxStudents,
			CurrentEnrollment: 10, // Full
		}

		student := &domain.Student{StudentID: studentID}

		mockStudentRepo.EXPECT().GetByID(gomock.Any(), studentID).Return(student, nil)
		mockCourseRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)
		mockRepo.EXPECT().GetByStudentAndCourse(gomock.Any(), studentID, courseID).Return(nil, domain.ErrEnrollmentNotFound)

		// Expect waitlist position check
		mockRepo.EXPECT().GetNextWaitlistPosition(gomock.Any(), courseID).Return(1, nil)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, e *domain.CourseEnrollment) error {
			assert.Equal(t, "waitlisted", e.EnrollmentStatus)
			assert.Equal(t, 1, *e.WaitlistPosition)
			return nil
		})

		// Should NOT increment enrollment
		mockProducer.EXPECT().PublishEvent("course.enrollment.created", gomock.Any(), gomock.Any()).Return(nil)

		enrollment, err := service.EnrollStudent(context.Background(), courseID, studentID, "admin")
		assert.NoError(t, err)
		assert.Equal(t, "waitlisted", enrollment.EnrollmentStatus)
	})

	t.Run("Already Enrolled", func(t *testing.T) {
		studentID := uuid.New()
		courseID := uuid.New()
		course := &domain.Course{CourseID: courseID, Status: "active"}
		student := &domain.Student{StudentID: studentID}

		mockStudentRepo.EXPECT().GetByID(gomock.Any(), studentID).Return(student, nil)
		mockCourseRepo.EXPECT().GetByID(gomock.Any(), courseID).Return(course, nil)
		// Return existing enrollment
		mockRepo.EXPECT().GetByStudentAndCourse(gomock.Any(), studentID, courseID).Return(&domain.CourseEnrollment{}, nil)

		_, err := service.EnrollStudent(context.Background(), courseID, studentID, "admin")
		assert.ErrorIs(t, err, domain.ErrAlreadyEnrolled)
	})
}
