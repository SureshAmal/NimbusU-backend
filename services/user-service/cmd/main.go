package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "github.com/SureshAmal/NimbusU-backend/services/user-service/internal/handler/http"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/repository/postgres"
	"github.com/SureshAmal/NimbusU-backend/services/user-service/internal/service"
	"github.com/SureshAmal/NimbusU-backend/shared/config"
	"github.com/SureshAmal/NimbusU-backend/shared/database"
	"github.com/SureshAmal/NimbusU-backend/shared/kafka"
	"github.com/SureshAmal/NimbusU-backend/shared/logger"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @title           NimbusU User Service API
// @version         1.0
// @description     User Management Service for NimbusU University Platform
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.nimbusu.edu/support
// @contact.email  support@nimbusu.edu

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found, using environment variables\n")
	}

	// Initialize logger
	cfg := config.LoadConfig()
	if err := logger.InitLogger(cfg.Server.Env); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting User Service",
		zap.String("env", cfg.Server.Env),
		zap.String("port", cfg.Server.Port),
	)

	// Connect to PostgreSQL
	logger.Info("Connecting to PostgreSQL", zap.String("url", cfg.Database.URL))
	pgPool, err := database.NewPostgresPool(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer database.ClosePostgresPool(pgPool)

	// Connect to Redis
	logger.Info("Connecting to Redis", zap.String("url", cfg.Redis.URL))
	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer database.CloseRedisClient(redisClient)

	// Initialize Kafka producer
	logger.Info("Initializing Kafka producer", zap.Strings("brokers", cfg.Kafka.Brokers))
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		logger.Fatal("Failed to initialize Kafka producer", zap.Error(err))
	}
	defer kafkaProducer.Close()

	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)

	// Initialize repositories
	logger.Info("Initializing repositories")
	userRepo := postgres.NewUserRepository(pgPool)
	profileRepo := postgres.NewUserProfileRepository(pgPool)
	roleRepo := postgres.NewRoleRepository(pgPool)
	sessionRepo := postgres.NewSessionRepository(pgPool)
	passwordTokenRepo := postgres.NewPasswordResetTokenRepository(pgPool)
	activityLogRepo := postgres.NewActivityLogRepository(pgPool)

	// Initialize services
	logger.Info("Initializing services")
	userSvc := service.NewUserService(
		userRepo,
		profileRepo,
		roleRepo,
		activityLogRepo,
		kafkaProducer,
	)

	authSvc := service.NewAuthService(
		userRepo,
		profileRepo,
		roleRepo,
		sessionRepo,
		passwordTokenRepo,
		activityLogRepo,
		jwtManager,
		kafkaProducer,
		cfg.JWT.RefreshTokenExpiry,
	)

	// Initialize HTTP handlers
	logger.Info("Initializing HTTP handlers")
	authHandler := httpHandler.NewAuthHandler(authSvc)
	userHandler := httpHandler.NewUserHandler(userSvc)

	// Setup Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Setup routes
	httpHandler.SetupRoutes(router, authHandler, userHandler, jwtManager, redisClient)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited successfully")
}
