package repository

import (
	"database/sql"
	"Gym-StrongCode/internal/models"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(userID, classID int) (int64, error) {
	res, err := r.db.Exec(`
		INSERT INTO bookings (user_id, class_id) VALUES (?, ?)`, userID, classID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *BookingRepository) Exists(userID, classID int) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM bookings WHERE user_id = ? AND class_id = ?`, userID, classID).Scan(&count)
	return count > 0, err
}

func (r *BookingRepository) GetByUser(userID int) ([]models.Booking, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, class_id, status, created_at
		FROM bookings WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.ClassID, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingRepository) ListAll() ([]models.Booking, error) {
	rows, err := r.db.Query(`SELECT id, user_id, class_id, status, created_at FROM bookings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.ClassID, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}

func (r *BookingRepository) Cancel(bookingID, userID int) error {
	_, err := r.db.Exec(`DELETE FROM bookings WHERE id = ? AND user_id = ?`, bookingID, userID)
	return err
}