package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreateStandalone(userID, amountCents int, currency, method, status, description, referenceID string) (*models.Payment, error) {
	res, err := r.db.Exec(`
        INSERT INTO payments (user_id, amount_cents, currency, method, status, description, reference_id)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, amountCents, currency, method, status, description, referenceID)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return r.GetByID(int(id))
}

func (r *PaymentRepository) GetByID(id int) (*models.Payment, error) {
	p := &models.Payment{}
	err := r.db.QueryRow(`
		SELECT id, user_id, amount_cents, currency, method, status, created_at
		FROM payments WHERE id = ?`, id).
		Scan(&p.ID, &p.UserID, &p.AmountCents, &p.Currency, &p.Method, &p.Status, &p.CreatedAt)
	return p, err
}

func (r *PaymentRepository) ListAll() ([]models.Payment, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, amount_cents, currency, method, status, created_at
		FROM payments ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.AmountCents, &p.Currency, &p.Method, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}
func (r *PaymentRepository) CreateForMembership(userID, amountCents int, currency, method, description, referenceID string) (*models.Payment, error) {
	return r.CreateStandalone(userID, amountCents, currency, method, "completed", description, referenceID)
}

func (r *PaymentRepository) GetByUser(userID int, status string) ([]models.Payment, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, amount_cents, currency, method, status, created_at
		FROM payments WHERE user_id = ? AND status = ? ORDER BY created_at DESC`,
		userID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.UserID, &p.AmountCents, &p.Currency, &p.Method, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}
