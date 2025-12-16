package repository

import (
	"database/sql"
	"time"

	"Gym-StrongCode/internal/models"
)

type MembershipRepository struct {
	db *sql.DB
}

func NewMembershipRepository(db *sql.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) GetAll() ([]models.Membership, error) {
	rows, err := r.db.Query(`SELECT id, name, duration_days, price_cents, created_at FROM memberships`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Membership
	for rows.Next() {
		var m models.Membership
		if err := rows.Scan(&m.ID, &m.Name, &m.DurationDays, &m.PriceCents, &m.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *MembershipRepository) GetByID(id int) (*models.Membership, error) {
	m := &models.Membership{}
	err := r.db.QueryRow(`SELECT id, name, duration_days, price_cents FROM memberships WHERE id = ?`, id).
		Scan(&m.ID, &m.Name, &m.DurationDays, &m.PriceCents)
	return m, err
}

func (r *MembershipRepository) Create(name string, durationDays, priceCents int) (*models.Membership, error) {
	res, err := r.db.Exec(`INSERT INTO memberships (name, duration_days, price_cents) VALUES (?, ?, ?)`,
		name, durationDays, priceCents)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return r.GetByID(int(id))
}

func (r *MembershipRepository) Update(id int, name string, durationDays, priceCents int) error {
	_, err := r.db.Exec(`UPDATE memberships SET name = ?, duration_days = ?, price_cents = ? WHERE id = ?`,
		name, durationDays, priceCents, id)
	return err
}

func (r *MembershipRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM memberships WHERE id = ?`, id)
	return err
}

func (r *MembershipRepository) HasActiveMembership(userID int) (bool, error) {
	var count int
	current := time.Now().Format("2006-01-02")
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM user_memberships 
		WHERE user_id = ? AND active = 1 AND end_date >= ?`, userID, current).
		Scan(&count)
	return count > 0, err
}

func (r *MembershipRepository) Activate(userID, membershipID int, durationDays int) error {
	start := time.Now()
	end := start.AddDate(0, 0, durationDays)
	_, err := r.db.Exec(`
		INSERT INTO user_memberships (user_id, membership_id, start_date, end_date, active)
		VALUES (?, ?, ?, ?, 1)`, userID, membershipID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	return err
}