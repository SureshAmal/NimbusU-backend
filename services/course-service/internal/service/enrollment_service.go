package service

import (
	"context"
	"time"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/google/uuid"
)

type enrollmentService struct {
	repo         domain.EnrollmentRepository
	courseRepo   domain.CourseRepository
	studentRepo  domain.StudentRepository
	subjectRepo  domain.SubjectRepository
	semesterRepo domain.SemesterRepository
	producer     domain.EventProducer
}

func NewEnrollmentService(
	repo domain.EnrollmentRepository,
	courseRepo domain.CourseRepository,
	studentRepo domain.StudentRepository,
	subjectRepo domain.SubjectRepository,
	semesterRepo domain.SemesterRepository,
	producer domain.EventProducer,
) domain.EnrollmentService {
	return &enrollmentService{
		repo:         repo,
		courseRepo:   courseRepo,
		studentRepo:  studentRepo,
		subjectRepo:  subjectRepo,
		semesterRepo: semesterRepo,
		producer:     producer,
	}
}

func (s *enrollmentService) EnrollStudent(ctx context.Context, courseID, studentID uuid.UUID, enrolledBy string) (*domain.CourseEnrollment, error) {
	// Check if student exists
	_, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Check if course exists and is active
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if course.Status != "active" {
		return nil, domain.ErrRegistrationClosed
	}

	// Check if already enrolled
	_, err = s.repo.GetByStudentAndCourse(ctx, studentID, courseID)
	if err == nil {
		return nil, domain.ErrAlreadyEnrolled
	}
	if err != domain.ErrEnrollmentNotFound {
		return nil, err
	}

	// Determine enrollment status (enrolled or waitlisted)
	status := "enrolled"
	var waitlistPosition *int
	if course.MaxStudents != nil && course.CurrentEnrollment >= *course.MaxStudents {
		status = "waitlisted"
		pos, err := s.repo.GetNextWaitlistPosition(ctx, courseID)
		if err != nil {
			return nil, err
		}
		waitlistPosition = &pos
	}

	enrollment := &domain.CourseEnrollment{
		StudentID:        studentID,
		CourseID:         courseID,
		EnrollmentStatus: status,
		EnrolledBy:       enrolledBy,
		WaitlistPosition: waitlistPosition,
	}

	if err := s.repo.Create(ctx, enrollment); err != nil {
		return nil, err
	}

	// Increment course enrollment count if enrolled (not waitlisted)
	if status == "enrolled" {
		if err := s.courseRepo.IncrementEnrollment(ctx, courseID); err != nil {
			// Log error but don't fail
		}
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.enrollment.created", enrollment.EnrollmentID.String(), map[string]interface{}{
			"enrollment_id": enrollment.EnrollmentID,
			"student_id":    studentID,
			"course_id":     courseID,
			"status":        status,
		})
	}

	return enrollment, nil
}

func (s *enrollmentService) DropCourse(ctx context.Context, courseID, studentID uuid.UUID, reason string) error {
	enrollment, err := s.repo.GetByStudentAndCourse(ctx, studentID, courseID)
	if err != nil {
		return err
	}

	if enrollment.EnrollmentStatus == "completed" {
		return domain.ErrCannotDropCompletedCourse
	}

	wasEnrolled := enrollment.EnrollmentStatus == "enrolled"

	now := time.Now()
	enrollment.EnrollmentStatus = "dropped"
	enrollment.DroppedDate = &now
	enrollment.DropReason = &reason

	if err := s.repo.Update(ctx, enrollment); err != nil {
		return err
	}

	// Decrement course enrollment count if was enrolled
	if wasEnrolled {
		if err := s.courseRepo.DecrementEnrollment(ctx, courseID); err != nil {
			// Log error but don't fail
		}

		// Promote next student from waitlist
		promoted, err := s.repo.PromoteFromWaitlist(ctx, courseID)
		if err == nil && promoted != nil {
			// Increment enrollment count for promoted student
			s.courseRepo.IncrementEnrollment(ctx, courseID)

			// Publish promotion event
			if s.producer != nil {
				s.producer.PublishEvent("course.enrollment.promoted", promoted.EnrollmentID.String(), map[string]interface{}{
					"enrollment_id": promoted.EnrollmentID,
					"student_id":    promoted.StudentID,
					"course_id":     courseID,
				})
			}
		}
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.enrollment.dropped", enrollment.EnrollmentID.String(), map[string]interface{}{
			"enrollment_id": enrollment.EnrollmentID,
			"student_id":    studentID,
			"course_id":     courseID,
			"reason":        reason,
		})
	}

	return nil
}

func (s *enrollmentService) UpdateEnrollment(ctx context.Context, enrollmentID uuid.UUID, status string, grade *string, gradePoints *float64) error {
	enrollment, err := s.repo.GetByID(ctx, enrollmentID)
	if err != nil {
		return err
	}

	// Validate status transition
	validTransitions := map[string][]string{
		"enrolled":   {"completed", "dropped"},
		"waitlisted": {"enrolled", "dropped"},
		"dropped":    {},
		"completed":  {},
	}

	valid := false
	for _, allowed := range validTransitions[enrollment.EnrollmentStatus] {
		if status == allowed {
			valid = true
			break
		}
	}
	if !valid && status != enrollment.EnrollmentStatus {
		return domain.ErrInvalidEnrollmentStatus
	}

	enrollment.EnrollmentStatus = status
	if grade != nil {
		enrollment.Grade = grade
	}
	if gradePoints != nil {
		enrollment.GradePoints = gradePoints
	}
	if status == "completed" {
		now := time.Now()
		enrollment.CompletionDate = &now
	}

	if err := s.repo.Update(ctx, enrollment); err != nil {
		return err
	}

	// Publish event
	if s.producer != nil {
		s.producer.PublishEvent("course.enrollment.updated", enrollment.EnrollmentID.String(), map[string]interface{}{
			"enrollment_id": enrollmentID,
			"status":        status,
			"grade":         grade,
		})
	}

	return nil
}

func (s *enrollmentService) GetStudentEnrollments(ctx context.Context, studentID uuid.UUID, filter domain.EnrollmentFilter, page, limit int) ([]*domain.EnrollmentWithDetails, int64, error) {
	filter.StudentID = &studentID
	offset := (page - 1) * limit
	return s.repo.ListByStudent(ctx, filter, limit, offset)
}

func (s *enrollmentService) BulkEnroll(ctx context.Context, courseID uuid.UUID, studentIDs []uuid.UUID, skipPrerequisites bool) ([]domain.BulkEnrollResult, error) {
	results := make([]domain.BulkEnrollResult, 0, len(studentIDs))

	for _, studentID := range studentIDs {
		result := domain.BulkEnrollResult{
			StudentID: studentID,
			Status:    "success",
		}

		// Check prerequisites if not skipping
		if !skipPrerequisites {
			met, _, err := s.CheckPrerequisites(ctx, studentID, courseID)
			if err != nil {
				result.Status = "failed"
				result.Error = err.Error()
				results = append(results, result)
				continue
			}
			if !met {
				result.Status = "failed"
				result.Error = domain.ErrPrerequisitesNotMet.Error()
				results = append(results, result)
				continue
			}
		}

		enrollment, err := s.EnrollStudent(ctx, courseID, studentID, "bulk_enrollment")
		if err != nil {
			result.Status = "failed"
			result.Error = err.Error()
		} else {
			result.EnrollmentID = &enrollment.EnrollmentID
			if enrollment.EnrollmentStatus == "waitlisted" {
				result.Status = "waitlisted"
			}
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *enrollmentService) CheckPrerequisites(ctx context.Context, studentID, courseID uuid.UUID) (bool, []domain.SubjectBasic, error) {
	// Get the course to find its subject
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return false, nil, err
	}

	// Get subject prerequisites
	prereqs, err := s.subjectRepo.GetPrerequisites(ctx, course.SubjectID)
	if err != nil {
		return false, nil, err
	}

	if len(prereqs) == 0 {
		return true, nil, nil
	}

	// Get student's completed courses
	filter := domain.EnrollmentFilter{
		StudentID: &studentID,
		Status:    strPtr("completed"),
	}
	completedEnrollments, _, err := s.repo.ListByStudent(ctx, filter, 1000, 0)
	if err != nil {
		return false, nil, err
	}

	// Build set of completed subject IDs
	completedSubjects := make(map[uuid.UUID]bool)
	for _, e := range completedEnrollments {
		// Need to get course to find subject ID
		c, err := s.courseRepo.GetByID(ctx, e.CourseID)
		if err == nil {
			completedSubjects[c.SubjectID] = true
		}
	}

	// Check which prerequisites are not met
	var missingPrereqs []domain.SubjectBasic
	for _, prereq := range prereqs {
		if prereq.IsMandatory && !completedSubjects[prereq.PrerequisiteSubjectID] {
			missingPrereqs = append(missingPrereqs, domain.SubjectBasic{
				SubjectID:   prereq.PrerequisiteSubjectID,
				SubjectCode: prereq.SubjectCode,
				SubjectName: prereq.SubjectName,
			})
		}
	}

	return len(missingPrereqs) == 0, missingPrereqs, nil
}

func strPtr(s string) *string {
	return &s
}
