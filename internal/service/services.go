// internal/service/services.go
package service

import (
	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// ============================================
// AuthService
// ============================================

// AuthService управляет аутентификацией и авторизацией
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

// NewAuthService создает новый AuthService
func NewAuthService(userRepo *repository.UserRepository, jwtSecret []byte) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepo.Create(name, email, string(hash))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login авторизует пользователя и возвращает JWT токен
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := s.createToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

// createToken создает JWT токен
func (s *AuthService) createToken(userID int, email string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ============================================
// MembershipService
// ============================================

// MembershipService управляет подписками
type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	paymentRepo    *repository.PaymentRepository
	db             *sql.DB
}

// NewMembershipService создает новый MembershipService
func NewMembershipService(membershipRepo *repository.MembershipRepository, paymentRepo *repository.PaymentRepository, db *sql.DB) *MembershipService {
	return &MembershipService{
		membershipRepo: membershipRepo,
		paymentRepo:    paymentRepo,
		db:             db,
	}
}

// GetAll возвращает все доступные подписки
func (s *MembershipService) GetAll() ([]models.Membership, error) {
	return s.membershipRepo.GetAll()
}

// BuyMembership покупает подписку для пользователя (атомарная операция)
func (s *MembershipService) BuyMembership(userID, membershipID int, method string) (map[string]interface{}, error) {
	membership, err := s.membershipRepo.GetByID(membershipID)
	if err != nil {
		return nil, err
	}

	// Используем транзакцию для атомарности операции
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Создаем платеж
	paymentID, err := s.paymentRepo.Create(tx, userID, membership.PriceCents, "KZT", method, "done")
	if err != nil {
		return nil, err
	}

	// Создаем подписку
	startDate := time.Now().UTC()
	endDate := startDate.AddDate(0, 0, membership.DurationDays)

	if err := s.membershipRepo.CreateUserMembership(
		tx,
		userID,
		membershipID,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return map[string]interface{}{
		"payment_id": paymentID,
		"start":      startDate.Format("2006-01-02"),
		"end":        endDate.Format("2006-01-02"),
	}, nil
}

// ============================================
// TrainerService
// ============================================

// TrainerService управляет тренерами
type TrainerService struct {
	trainerRepo *repository.TrainerRepository
}

// NewTrainerService создает новый TrainerService
func NewTrainerService(trainerRepo *repository.TrainerRepository) *TrainerService {
	return &TrainerService{trainerRepo: trainerRepo}
}

// Create создает нового тренера
func (s *TrainerService) Create(name, bio string) (int64, error) {
	if name == "" {
		return 0, fmt.Errorf("trainer name is required")
	}
	return s.trainerRepo.Create(name, bio)
}

// ============================================
// ClassService
// ============================================

// ClassService управляет занятиями
type ClassService struct {
	classRepo   *repository.ClassRepository
	trainerRepo *repository.TrainerRepository
}

// NewClassService создает новый ClassService
func NewClassService(classRepo *repository.ClassRepository, trainerRepo *repository.TrainerRepository) *ClassService {
	return &ClassService{
		classRepo:   classRepo,
		trainerRepo: trainerRepo,
	}
}

// GetAll возвращает все занятия
func (s *ClassService) GetAll() ([]models.Class, error) {
	return s.classRepo.GetAll()
}

// Create создает новое занятие
func (s *ClassService) Create(class *models.Class) error {
	if class.Title == "" {
		return fmt.Errorf("class title is required")
	}

	// Валидация времени начала
	if class.StartTime.Before(time.Now()) {
		return fmt.Errorf("start time must be in the future")
	}

	// Проверка существования тренера
	if class.TrainerID > 0 {
		exists, err := s.trainerRepo.Exists(class.TrainerID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("trainer not found")
		}
	}

	return s.classRepo.Create(class)
}

// ============================================
// BookingService
// ============================================

// BookingService управляет бронированиями
type BookingService struct {
	bookingRepo    *repository.BookingRepository
	classRepo      *repository.ClassRepository
	membershipRepo *repository.MembershipRepository
}

// NewBookingService создает новый BookingService
func NewBookingService(
	bookingRepo *repository.BookingRepository,
	classRepo *repository.ClassRepository,
	membershipRepo *repository.MembershipRepository,
) *BookingService {
	return &BookingService{
		bookingRepo:    bookingRepo,
		classRepo:      classRepo,
		membershipRepo: membershipRepo,
	}
}

// Create создает новое бронирование
func (s *BookingService) Create(userID, classID int) (int64, string, error) {
	// Проверяем существование класса
	class, err := s.classRepo.GetByID(classID)
	if err != nil {
		return 0, "", err
	}

	// Проверяем наличие активной подписки
	hasActive, err := s.membershipRepo.HasActiveMembership(userID)
	if err != nil {
		return 0, "", err
	}
	if !hasActive {
		return 0, "", fmt.Errorf("no active membership")
	}

	// Проверяем, не забронировано ли уже
	exists, err := s.bookingRepo.Exists(userID, classID)
	if err != nil {
		return 0, "", err
	}
	if exists {
		return 0, "", fmt.Errorf("already booked")
	}

	// Проверяем наличие мест
	count, err := s.classRepo.GetBookingCount(classID)
	if err != nil {
		return 0, "", err
	}
	if count >= class.Capacity {
		return 0, "", fmt.Errorf("class is full")
	}

	// Создаем бронирование
	bookingID, err := s.bookingRepo.Create(userID, classID)
	if err != nil {
		return 0, "", err
	}

	return bookingID, class.StartTime.Format(time.RFC3339), nil
}

// GetByUser возвращает все бронирования пользователя
func (s *BookingService) GetByUser(userID int) ([]models.Booking, error) {
	return s.bookingRepo.GetByUser(userID)
}

// ============================================
// PaymentService
// ============================================

// PaymentService управляет платежами
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
}

// NewPaymentService создает новый PaymentService
func NewPaymentService(paymentRepo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo}
}

// Create создает новый платеж
func (s *PaymentService) Create(userID, amountCents int, method string) (int64, error) {
	if amountCents <= 0 {
		return 0, fmt.Errorf("amount must be positive")
	}
	return s.paymentRepo.CreateStandalone(userID, amountCents, "KZT", method, "done")
}
