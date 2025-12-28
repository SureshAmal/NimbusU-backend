package http

import (
	"time"

	_ "github.com/SureshAmal/NimbusU-backend/services/user-service/docs"
	"github.com/SureshAmal/NimbusU-backend/shared/middleware"
	"github.com/SureshAmal/NimbusU-backend/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all HTTP routes for the user service
func SetupRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	jwtManager *utils.JWTManager,
	redisClient *redis.Client,
) {
	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Rate limiting (100 requests per minute per IP)
	rateLimiter := middleware.NewRateLimiter(redisClient, 100, time.Minute)
	router.Use(rateLimiter.RateLimitMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, 200, "User service is healthy", gin.H{
			"service": "user-service",
			"status":  "healthy",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes (no authentication required)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
		authRoutes.POST("/password/reset-request", authHandler.RequestPasswordReset)
		authRoutes.POST("/password/reset", authHandler.ResetPassword)
	}

	// Protected routes (authentication required)
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// Authenticated auth routes
	authProtected := router.Group("/auth")
	authProtected.Use(authMiddleware)
	{
		authProtected.POST("/logout", authHandler.Logout)
		authProtected.POST("/password/change", authHandler.ChangePassword)
		authProtected.GET("/sessions", authHandler.GetActiveSessions)
		authProtected.DELETE("/sessions", authHandler.RevokeAllSessions)
		authProtected.DELETE("/sessions/:sessionId", authHandler.RevokeSession)
	}

	// User routes (self-service)
	userRoutes := router.Group("/users")
	userRoutes.Use(authMiddleware)
	{
		userRoutes.GET("/me", userHandler.GetMe)
		userRoutes.PUT("/me", userHandler.UpdateMe)
	}

	// Admin routes (admin and faculty only)
	adminMiddleware := middleware.RoleMiddleware("admin", "faculty")

	adminUserRoutes := router.Group("/admin/users")
	adminUserRoutes.Use(authMiddleware, adminMiddleware)
	{
		adminUserRoutes.POST("", userHandler.CreateUser)
		adminUserRoutes.GET("", userHandler.ListUsers)
		adminUserRoutes.GET("/:id", userHandler.GetUser)
		adminUserRoutes.PUT("/:id", userHandler.UpdateUser)
		adminUserRoutes.DELETE("/:id", userHandler.DeleteUser)
		adminUserRoutes.POST("/:id/activate", userHandler.ActivateUser)
		adminUserRoutes.POST("/:id/suspend", userHandler.SuspendUser)
		adminUserRoutes.POST("/bulk-import", userHandler.BulkImportUsers)
	}
}
