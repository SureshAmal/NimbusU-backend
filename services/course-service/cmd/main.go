package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httphandler "github.com/SureshAmal/NimbusU-backend/services/course-service/internal/handler/http"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/repository/postgres"
	"github.com/SureshAmal/NimbusU-backend/services/course-service/internal/service"
	"github.com/SureshAmal/NimbusU-backend/shared/config"
	"github.com/SureshAmal/NimbusU-backend/shared/database"
	"github.com/SureshAmal/NimbusU-backend/shared/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found, using environment variables\n")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.InitLogger(cfg.Server.Env); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Course Service",
		zap.String("env", cfg.Server.Env),
		zap.String("port", cfg.Server.Port),
	)

	// Connect to PostgreSQL
	logger.Info("Connecting to PostgreSQL", zap.String("url", cfg.Database.URL))
	db, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.ClosePostgresPool(db)
	logger.Info("Connected to PostgreSQL")

	// Initialize repositories
	logger.Info("Initializing repositories")
	deptRepo := postgres.NewDepartmentRepository(db)
	progRepo := postgres.NewProgramRepository(db)
	subjRepo := postgres.NewSubjectRepository(db)
	semRepo := postgres.NewSemesterRepository(db)
	courseRepo := postgres.NewCourseRepository(db)
	facultyRepo := postgres.NewFacultyRepository(db)
	studentRepo := postgres.NewStudentRepository(db)
	fcRepo := postgres.NewFacultyCourseRepository(db)
	enrollRepo := postgres.NewEnrollmentRepository(db)
	calendarRepo := postgres.NewCalendarRepository(db)

	// Initialize services (nil producer for now - can add Kafka later)
	logger.Info("Initializing services")
	deptService := service.NewDepartmentService(deptRepo, nil)
	progService := service.NewProgramService(progRepo, deptRepo, nil)
	subjService := service.NewSubjectService(subjRepo, deptRepo, nil)
	semService := service.NewSemesterService(semRepo, nil)
	courseService := service.NewCourseService(courseRepo, subjRepo, semRepo, enrollRepo, nil)
	facultyService := service.NewFacultyService(facultyRepo, deptRepo, fcRepo, nil)
	studentService := service.NewStudentService(studentRepo, deptRepo, progRepo, nil)
	facultyAssignService := service.NewFacultyAssignmentService(fcRepo, facultyRepo, courseRepo, nil)
	enrollService := service.NewEnrollmentService(enrollRepo, courseRepo, studentRepo, subjRepo, semRepo, nil)
	calendarService := service.NewCalendarService(calendarRepo, semRepo, nil)

	// Unused services - log for documentation
	_ = facultyService
	_ = studentService

	// Setup routes
	logger.Info("Setting up routes")
	router := httphandler.SetupRoutes(
		deptService,
		progService,
		subjService,
		semService,
		courseService,
		facultyAssignService,
		enrollService,
		calendarService,
	)

	// Create HTTP server
	port := cfg.Server.Port
	if port == "" {
		port = "8082"
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Course Service listening", zap.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Course Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Course Service stopped")
}
