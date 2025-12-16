package main

import (
	"Gym_StrongCode/config"
	"Gym_StrongCode/internal/handler"
	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/internal/service"
	"log"

	_ "Gym_StrongCode/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mattn/go-sqlite3" // SQLite драйвер
)

// Swagger docs...
// @title           Gym StrongCode API
// @version         2.0
// @description     API для управления фитнес-клубом
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey  Bearer
// @in header
// @name Authorization

func main() {
	// Загружаем конфиг
	cfg := config.Load()

	// Подключаемся к БД + автоматически применяются миграции (всё внутри NewDatabase!)
	db, err := repository.NewDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// ← УДАЛЕНА СТРОКА InitSchema! Больше не нужна!

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	classRepo := repository.NewClassRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Сервисы
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db)
	trainerService := service.NewTrainerService(trainerRepo)
	classService := service.NewClassService(classRepo, trainerRepo)
	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo)
	paymentService := service.NewPaymentService(paymentRepo)

	// Хендлеры
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)
	membershipHandler := handler.NewMembershipHandler(membershipService)
	trainerHandler := handler.NewTrainerHandler(trainerService)
	classHandler := handler.NewClassHandler(classService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Gin режим
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)

		// Публичные
		api.POST("/users/register", authHandler.Register)
		api.POST("/users/login", authHandler.Login)
		api.GET("/classes", classHandler.GetClasses)
		api.GET("/memberships", membershipHandler.GetMemberships)

		// Авторизованные
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			authorized.GET("/me", userHandler.GetCurrentUser)
			authorized.POST("/bookings", bookingHandler.CreateBooking)
			authorized.GET("/bookings", bookingHandler.ListBookings)
			authorized.POST("/memberships/buy", membershipHandler.BuyMembership)
			authorized.POST("/payments", paymentHandler.CreatePayment)
		}

		// Админка
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/trainers", trainerHandler.CreateTrainer)
			admin.POST("/classes", classHandler.CreateClass)
		}
	}

	log.Printf("Gym StrongCode Server starting on %s", cfg.ServerAddress)
	log.Printf("Swagger UI: http://localhost%s/swagger/index.html", cfg.ServerAddress)

	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}