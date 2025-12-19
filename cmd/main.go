package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Gym_StrongCode/config"
	_ "Gym_StrongCode/docs"
	"Gym_StrongCode/internal/handler"
	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/internal/service"
	"Gym_StrongCode/internal/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           Gym StrongCode API
// @version         2.0
// @description     API для управления фитнес-клубом: пользователи, подписки, залы, тренеры, занятия, бронирования, платежи
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Введите JWT в формате: Bearer <token>
func main() {
	cfg := config.Load()
	utils.InitLogger()
	logger := utils.GetLogger()

	db, err := repository.NewDatabase(cfg.DatabasePath)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	gymRepo := repository.NewGymRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	classRepo := repository.NewClassRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Сервисы
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	notificationService := service.NewNotificationService(cfg)
	gymService := service.NewGymService(gymRepo)
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db, notificationService)
	trainerService := service.NewTrainerService(trainerRepo)
	classService := service.NewClassService(classRepo, trainerRepo, gymRepo)
	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo, notificationService)
	paymentService := service.NewPaymentService(paymentRepo)

	// Запуск background worker для email
	notificationService.StartWorker()
	logger.Info("Email notification worker started")

	// Хендлеры
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	gymHandler := handler.NewGymHandler(gymService)
	membershipHandler := handler.NewMembershipHandler(membershipService)
	trainerHandler := handler.NewTrainerHandler(trainerService)
	classHandler := handler.NewClassHandler(classService)
	bookingHandler := handler.NewBookingHandler(bookingService, userRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.RateLimitMiddleware())

	// Статические файлы фронтенда
	r.Static("/static", "./frontend/static")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		// Публичные
		api.POST("/users/register", authHandler.Register)
		api.POST("/users/login", authHandler.Login)
		api.GET("/classes", classHandler.List)
		api.GET("/gyms", gymHandler.List)
		api.GET("/memberships", membershipHandler.List)
		api.GET("/trainers", trainerHandler.List)

		// Авторизованные
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authorized.GET("/me", userHandler.GetCurrent)
			authorized.PUT("/me", userHandler.Update)

			authorized.POST("/bookings", bookingHandler.Create)
			authorized.GET("/bookings", bookingHandler.ListUser)
			authorized.DELETE("/bookings/:id", bookingHandler.Cancel)

			authorized.POST("/memberships/buy", membershipHandler.Buy)
			authorized.POST("/payments", paymentHandler.CreateStandalone)
		}

		// Админ
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			// Users
			admin.GET("/users", userHandler.List)
			admin.DELETE("/users/:id", userHandler.Delete)

			// Gyms
			admin.POST("/gyms", gymHandler.Create)
			admin.PUT("/gyms/:id", gymHandler.Update)
			admin.DELETE("/gyms/:id", gymHandler.Delete)

			// Memberships
			admin.POST("/memberships", membershipHandler.Create)
			admin.PUT("/memberships/:id", membershipHandler.Update)
			admin.DELETE("/memberships/:id", membershipHandler.Delete)

			// Trainers
			admin.POST("/trainers", trainerHandler.Create)
			admin.PUT("/trainers/:id", trainerHandler.Update)
			admin.DELETE("/trainers/:id", trainerHandler.Delete)

			// Classes
			admin.POST("/classes", classHandler.Create)
			admin.PUT("/classes/:id", classHandler.Update)
			admin.DELETE("/classes/:id", classHandler.Delete)

			// Payments & Bookings (read-only)
			admin.GET("/payments", paymentHandler.ListAll)
			admin.GET("/bookings", bookingHandler.ListAll)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}

		c.File("./frontend/index.html")
	})

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}
	go func() {
		logger.Info("Server starting", zap.String("addr", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Stopping email notification worker...")
	notificationService.StopWorker()

	logger.Info("Server stopped gracefully")
}
