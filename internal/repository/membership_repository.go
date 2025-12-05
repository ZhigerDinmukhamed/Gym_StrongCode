// internal/repository/membership_repository.go
package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// MembershipRepository управляет операциями с подписками в БД
type MembershipRepository struct {
	db *sql.DB
}

// NewMembershipRepository создает новый MembershipRepository
func NewMembershipRepository(db *sql.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

// GetAll возвращает все доступные подписки
func (r *MembershipRepository) GetAll() ([]models.Membership, error) {
	rows, err := r.db.Query(
		"SELECT id, name, duration_days, price_cents, created_at FROM memberships ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query memberships: %w", err)
	}
	defer rows.Close()

	var memberships []models.Membership
	for rows.Next() {
		var m models.Membership
		var createdAt string
		if err := rows.Scan(&m.ID, &m.Name, &m.DurationDays, &m.PriceCents, &createdAt); err != nil {
			continue
		}
		m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		memberships = append(memberships, m)
	}

	return memberships, nil
}

// GetByID возвращает подписку по ID
func (r *MembershipRepository) GetByID(id int) (*models.Membership, error) {
	var m models.Membership
	var createdAt string

	err := r.db.QueryRow(
		"SELECT id, name, duration_days, price_cents, created_at FROM memberships WHERE id = ?",
		id,
	).Scan(&m.ID, &m.Name, &m.DurationDays, &m.PriceCents, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("membership not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get membership: %w", err)
	}

	m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	return &m, nil
}

// CreateUserMembership создает подписку для пользователя (в рамках транзакции)
func (r *MembershipRepository) CreateUserMembership(tx *sql.Tx, userID, membershipID int, startDate, endDate string) error {
	_, err := tx.Exec(
		`INSERT INTO user_memberships(user_id, membership_id, start_date, end_date, active) 
		 VALUES(?, ?, ?, ?, 1)`,
		userID, membershipID, startDate, endDate,
	)
	if err != nil {
		return fmt.Errorf("failed to create user membership: %w", err)
	}
	return nil
}

// HasActiveMembership проверяет наличие активной подписки у пользователя
func (r *MembershipRepository) HasActiveMembership(userID int) (bool, error) {
	var count int
	now := time.Now().Format("2006-01-02")

	err := r.db.QueryRow(
		`SELECT COUNT(1) FROM user_memberships 
		 WHERE user_id = ? AND active = 1 AND start_date <= ? AND end_date >= ?`,
		userID, now, now,
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check active membership: %w", err)
	}

	return count > 0, nil
}
