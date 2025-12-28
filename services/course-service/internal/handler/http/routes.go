package http

import (
	"net/http"

	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/domain"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRoutes(
	deptService domain.DepartmentService,
	progService domain.ProgramService,
	subjService domain.SubjectService,
	semService domain.SemesterService,
	courseService domain.CourseService,
	facultyAssignService domain.FacultyAssignmentService,
	enrollService domain.EnrollmentService,
	calendarService domain.CalendarService,
) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Unused services - acknowledge for linter
	_ = progService
	_ = subjService
	_ = semService
	_ = calendarService

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Department routes
		deptHandler := NewDepartmentHandler(deptService)
		r.Route("/departments", func(r chi.Router) {
			r.Get("/", deptHandler.List)
			r.Post("/", deptHandler.Create)
			r.Get("/{id}", deptHandler.GetByID)
			r.Put("/{id}", deptHandler.Update)
			r.Delete("/{id}", deptHandler.Delete)
		})

		// Course routes
		courseHandler := NewCourseHandler(courseService, facultyAssignService, enrollService)
		r.Route("/courses", func(r chi.Router) {
			r.Get("/", courseHandler.List)
			r.Post("/", courseHandler.Create)
			r.Get("/{id}", courseHandler.GetByID)
			r.Put("/{id}", courseHandler.Update)
			r.Delete("/{id}", courseHandler.Delete)
			r.Post("/{id}/activate", courseHandler.Activate)
			r.Post("/{id}/deactivate", courseHandler.Deactivate)
			r.Get("/{id}/students", courseHandler.GetStudents)
			r.Get("/{id}/faculty", courseHandler.GetFaculty)
			r.Post("/{id}/faculty", courseHandler.AssignFaculty)
			r.Delete("/{id}/faculty/{facultyId}", courseHandler.RemoveFaculty)
			r.Post("/{id}/enroll", courseHandler.EnrollStudent)
		})

		// Enrollment routes
		enrollHandler := NewEnrollmentHandler(enrollService)
		r.Route("/enrollments", func(r chi.Router) {
			r.Get("/{id}", enrollHandler.GetByID)
			r.Put("/{id}", enrollHandler.Update)
			r.Post("/courses/{courseId}/bulk", enrollHandler.BulkEnroll)
			r.Delete("/courses/{courseId}/students/{studentId}", enrollHandler.Drop)
			r.Get("/students/{studentId}", enrollHandler.GetStudentEnrollments)
			r.Get("/courses/{courseId}/students/{studentId}/prerequisites", enrollHandler.CheckPrerequisites)
		})
	})

	return r
}
