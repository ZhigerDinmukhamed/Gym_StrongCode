// internal/repository/other_repositories.go
package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// ============================================
// TrainerRepository
// ============================================

// TrainerRepository управляет операциями с тренерами в БД
type TrainerRepository struct {
	db *sql.DB
}

// NewTrainerRepository создает новый TrainerRepository
func NewTrainerRepository(db *sql.DB) *TrainerRepository {
	return &TrainerRepository{db: db}
}

// Create создает нового тренера
func (r *TrainerRepository) Create(name, bio string) (int64, error) {
	res, err := r.db.Exec("INSERT INTO trainers(name, bio) VALUES(?, ?)", name, bio)
	if err != nil {
		return 0, fmt.Errorf("failed to create trainer: %w", err)
	}
	return res.LastInsertId()
}

// Exists проверяет существование тренера по ID
func (r *TrainerRepository) Exists(id int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(1) FROM trainers WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check trainer existence: %w", err)
	}
	return count > 0, nil
}

// ============================================
// ClassRepository
// ============================================

// ClassRepository управляет операциями с занятиями в БД
type ClassRepository struct {
	db *sql.DB
}

// NewClassRepository создает новый ClassRepository
func NewClassRepository(db *sql.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

// Create создает новое занятие
func (r *ClassRepository) Create(c *models.Class) error {
	_, err := r.db.Exec(
		`INSERT INTO classes(title, description, trainer_id, start_time, duration_min, capacity) 
		 VALUES(?, ?, ?, ?, ?, ?)`,
		c.Title, c.Description, c.TrainerID, c.StartTime.Format(time.RFC3339), c.DurationMin, c.Capacity,
	)
	if err != nil {
		return fmt.Errorf("failed to create class: %w", err)
	}
	return nil
}

// GetAll возвращает все занятия
func (r *ClassRepository) GetAll() ([]models.Class, error) {
	rows, err := r.db.Query(
		`SELECT id, title, description, trainer_id, start_time, duration_min, capacity, created_at 
		 FROM classes ORDER BY start_time ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query classes: %w", err)
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var c models.Class
		var startTime, createdAt string
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.TrainerID, &startTime, &c.DurationMin, &c.Capacity, &createdAt); err != nil {
			continue
		}
		c.StartTime, _ = time.Parse(time.RFC3339, startTime)
		c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		classes = append(classes, c)
	}

	return classes, nil
}

// GetByID возвращает занятие по ID
func (r *ClassRepository) GetByID(id int) (*models.Class, error) {
	var c models.Class
	var startTime, createdAt string

	err := r.db.QueryRow(
		`SELECT id, title, description, trainer_id, start_time, duration_min, capacity, created_at 
		 FROM classes WHERE id = ?`,
		id,
	).Scan(&c.ID, &c.Title, &c.Description, &c.TrainerID, &startTime, &c.DurationMin, &c.Capacity, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("class not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get class: %w", err)
	}

	c.StartTime, _ = time.Parse(time.RFC3339, startTime)
	c.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &c, nil
}

// GetBookingCount возвращает количество бронирований на занятие
func (r *ClassRepository) GetBookingCount(classID int) (int, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(1) FROM bookings WHERE class_id = ? AND status = 'booked'",
		classID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get booking count: %w", err)
	}
	return count, nil
}

// ============================================
// BookingRepository
// ============================================

// BookingRepository управляет операциями с бронированиями в БД
type BookingRepository struct {
	db *sql.DB
}

// NewBookingRepository создает новый BookingRepository
func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create создает новое бронирование
func (r *BookingRepository) Create(userID, classID int) (int64, error) {
	res, err := r.db.Exec(
		"INSERT INTO bookings(user_id, class_id, status) VALUES(?, ?, 'booked')",
		userID, classID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create booking: %w", err)
	}
	return res.LastInsertId()
}

// GetByUser возвращает все бронирования пользователя
func (r *BookingRepository) GetByUser(userID int) ([]models.Booking, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, class_id, status, created_at FROM bookings WHERE user_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		var createdAt string
		if err := rows.Scan(&b.ID, &b.UserID, &b.ClassID, &b.Status, &createdAt); err != nil {
			continue
		}
		b.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		bookings = append(bookings, b)
	}

	return bookings, nil
}

// Exists проверяет существование бронирования
func (r *BookingRepository) Exists(userID, classID int) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(1) FROM bookings WHERE user_id = ? AND class_id = ? AND status = 'booked'",
		userID, classID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check booking existence: %w", err)
	}
	return count > 0, nil
}

// ============================================
// PaymentRepository
// ============================================

// PaymentRepository управляет операциями с платежами в БД
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository создает новый PaymentRepository
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create создает новый платеж (в рамках транзакции)
func (r *PaymentRepository) Create(tx *sql.Tx, userID, amountCents int, currency, method, status string) (int64, error) {
	res, err := tx.Exec(
		"INSERT INTO payments(user_id, amount_cents, currency, method, status) VALUES(?, ?, ?, ?, ?)",
		userID, amountCents, currency, method, status,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create payment: %w", err)
	}
	return res.LastInsertId()
}

// CreateStandalone создает платеж без транзакции
func (r *PaymentRepository) CreateStandalone(userID, amountCents int, currency, method, status string) (int64, error) {
	res, err := r.db.Exec(
		"INSERT INTO payments(user_id, amount_cents, currency, method, status) VALUES(?, ?, ?, ?, ?)",
		userID, amountCents, currency, method, status,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create payment: %w", err)
	}
	return res.LastInsertId()
}
