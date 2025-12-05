// internal/handler/handlers.go
package handler

import (
	"Gym_StrongCode/internal/middleware"
	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================
// Health Check
// ============================================

// HealthCheck godoc
// @Summary Проверка здоровья сервера
// @Description Возвращает статус OK
// @Tags health
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "Gym StrongCode API"})
}

// ============================================
// AuthHandler
// ============================================

// AuthHandler обрабатывает запросы аутентификации
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler создает новый AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required" example:"Асем Нурова"`
	Email    string `json:"email" binding:"required,email" example:"asem@example.kz"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param body body registerRequest true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"asem@example.kz"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// Login godoc
// @Summary Авторизация
// @Description Возвращает JWT токен для авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "Логин и пароль"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ============================================
// UserHandler
// ============================================

// UserHandler обрабатывает запросы пользователей
type UserHandler struct {
	userRepo *repository.UserRepository
}

// NewUserHandler создает новый UserHandler
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetCurrentUser godoc
// @Summary Получить текущего пользователя
// @Description Возвращает данные авторизованного пользователя
// @Tags user
// @Security Bearer
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ============================================
// MembershipHandler
// ============================================

// MembershipHandler обрабатывает запросы подписок
type MembershipHandler struct {
	membershipService *service.MembershipService
}

// NewMembershipHandler создает новый MembershipHandler
func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

// GetMemberships godoc
// @Summary Список подписок
// @Description Возвращает все доступные типы подписок
// @Tags memberships
// @Success 200 {array} models.Membership
// @Failure 500 {object} map[string]string
// @Router /memberships [get]
func (h *MembershipHandler) GetMemberships(c *gin.Context) {
	memberships, err := h.membershipService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, memberships)
}

type buyMembershipRequest struct {
	MembershipID int    `json:"membership_id" binding:"required" example:"1"`
	Method       string `json:"method" binding:"required" example:"card"`
}

// BuyMembership godoc
// @Summary Купить подписку
// @Description Создаёт платеж и активирует подписку для пользователя
// @Tags memberships
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body buyMembershipRequest true "ID подписки и метод оплаты"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /memberships/buy [post]
func (h *MembershipHandler) BuyMembership(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req buyMembershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.membershipService.BuyMembership(userID, req.MembershipID, req.Method)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============================================
// TrainerHandler
// ============================================

// TrainerHandler обрабатывает запросы тренеров
type TrainerHandler struct {
	trainerService *service.TrainerService
}

// NewTrainerHandler создает новый TrainerHandler
func NewTrainerHandler(trainerService *service.TrainerService) *TrainerHandler {
	return &TrainerHandler{trainerService: trainerService}
}

type createTrainerRequest struct {
	Name string `json:"name" binding:"required" example:"Данияр Смагулов"`
	Bio  string `json:"bio" example:"Мастер спорта международного класса по кроссфиту"`
}

// CreateTrainer godoc
// @Summary Создать тренера
// @Description Создаёт нового тренера (только для администратора)
// @Tags admin
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body createTrainerRequest true "Данные тренера"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /admin/trainers [post]
func (h *TrainerHandler) CreateTrainer(c *gin.Context) {
	var req createTrainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trainerID, err := h.trainerService.Create(req.Name, req.Bio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"trainer_id": trainerID})
}

// ============================================
// ClassHandler
// ============================================

// ClassHandler обрабатывает запросы занятий
type ClassHandler struct {
	classService *service.ClassService
}

// NewClassHandler создает новый ClassHandler
func NewClassHandler(classService *service.ClassService) *ClassHandler {
	return &ClassHandler{classService: classService}
}

// GetClasses godoc
// @Summary Список занятий
// @Description Возвращает все доступные групповые занятия
// @Tags classes
// @Success 200 {array} models.Class
// @Failure 500 {object} map[string]string
// @Router /classes [get]
func (h *ClassHandler) GetClasses(c *gin.Context) {
	classes, err := h.classService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, classes)
}

type createClassRequest struct {
	Title       string    `json:"title" binding:"required" example:"Йога для начинающих"`
	Description string    `json:"description" example:"Спокойная практика для новичков"`
	TrainerID   int       `json:"trainer_id" example:"1"`
	StartTime   time.Time `json:"start_time" binding:"required" example:"2025-12-20T10:00:00Z"`
	DurationMin int       `json:"duration_min" binding:"required" example:"60"`
	Capacity    int       `json:"capacity" binding:"required" example:"15"`
}

// CreateClass godoc
// @Summary Создать занятие
// @Description Создаёт новое групповое занятие (только для администратора)
// @Tags admin
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body createClassRequest true "Данные занятия"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /admin/classes [post]
func (h *ClassHandler) CreateClass(c *gin.Context) {
	var req createClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class := &models.Class{
		Title:       req.Title,
		Description: req.Description,
		TrainerID:   req.TrainerID,
		StartTime:   req.StartTime,
		DurationMin: req.DurationMin,
		Capacity:    req.Capacity,
	}

	if err := h.classService.Create(class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "class created successfully"})
}

// ============================================
// BookingHandler
// ============================================

// BookingHandler обрабатывает запросы бронирований
type BookingHandler struct {
	bookingService *service.BookingService
}

// NewBookingHandler создает новый BookingHandler
func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

type createBookingRequest struct {
	ClassID int `json:"class_id" binding:"required" example:"1"`
}

// CreateBooking godoc
// @Summary Забронировать занятие
// @Description Создаёт бронирование на занятие (требует активную подписку)
// @Tags bookings
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body createBookingRequest true "ID занятия"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req createBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingID, startTime, err := h.bookingService.Create(userID, req.ClassID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "no active membership" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"booking_id": bookingID,
		"start_time": startTime,
	})
}

// ListBookings godoc
// @Summary Мои бронирования
// @Description Возвращает список всех бронирований текущего пользователя
// @Tags bookings
// @Security Bearer
// @Success 200 {array} models.Booking
// @Failure 500 {object} map[string]string
// @Router /bookings [get]
func (h *BookingHandler) ListBookings(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	bookings, err := h.bookingService.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// ============================================
// PaymentHandler
// ============================================

// PaymentHandler обрабатывает запросы платежей
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler создает новый PaymentHandler
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

type createPaymentRequest struct {
	AmountCents int    `json:"amount_cents" binding:"required,gt=0" example:"10000"`
	Method      string `json:"method" binding:"required" example:"card"`
}

// CreatePayment godoc
// @Summary Произвольная оплата
// @Description Создаёт запись об оплате произвольной суммы
// @Tags payments
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body createPaymentRequest true "Сумма и метод оплаты"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentID, err := h.paymentService.Create(userID, req.AmountCents, req.Method)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_id": paymentID})
}
