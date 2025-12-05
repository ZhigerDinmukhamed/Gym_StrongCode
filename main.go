package main

import (
	"Gym_StrongCode/config"
	"Gym_StrongCode/internal/cache"
	"Gym_StrongCode/internal/handler"
	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/internal/service"
	"log"

	_ "Gym_StrongCode/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gym StrongCode API
// @version 2.0
// @description API для управления фитнес-клубом: бронирование занятий, подписки, тренеры, админка.
// @contact.name API Support
// @contact.email support@strongcode.kz
// @license.name MIT
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer <ваш_токен>

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализируем БД
	db, err := repository.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Инициализируем схему и сидеры
	if err := repository.InitSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// 🔥 Создаём кэш
	appCache := cache.NewCache()

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	classRepo := repository.NewClassRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Создаем сервисы
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db)
	trainerService := service.NewTrainerService(trainerRepo)

	// 🔥 ClassService теперь с кэшем!
	classService := service.NewClassService(classRepo, trainerRepo, appCache)

	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo)
	paymentService := service.NewPaymentService(paymentRepo)

	// Создаем хендлеры
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	membershipHandler := handler.NewMembershipHandler(membershipService)
	trainerHandler := handler.NewTrainerHandler(trainerService)
	classHandler := handler.NewClassHandler(classService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Настраиваем Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API группа
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", handler.HealthCheck)

		// Публичные эндпоинты
		api.POST("/users/register", authHandler.Register)
		api.POST("/users/login", authHandler.Login)

		// 🔥 классы теперь читаются из кэша
		api.GET("/classes", classHandler.GetClasses)

		api.GET("/memberships", membershipHandler.GetMemberships)

		// Защищенные эндпоинты (требуют авторизации)
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authorized.GET("/me", userHandler.GetCurrentUser)
			authorized.POST("/bookings", bookingHandler.CreateBooking)
			authorized.GET("/bookings", bookingHandler.ListBookings)
			authorized.POST("/memberships/buy", membershipHandler.BuyMembership)
			authorized.POST("/payments", paymentHandler.CreatePayment)
		}

		// Админские эндпоинты
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/trainers", trainerHandler.CreateTrainer)
			admin.POST("/classes", classHandler.CreateClass)
		}
	}

	log.Printf("🚀 Gym StrongCode Server starting on %s", cfg.ServerAddress)
	log.Printf("📚 Swagger UI: http://localhost%s/swagger/index.html", cfg.ServerAddress)
	log.Printf("🏥 Health check: http://localhost%s/api/health", cfg.ServerAddress)

	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
