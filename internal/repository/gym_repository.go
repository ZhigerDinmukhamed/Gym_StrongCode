package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
)

type GymRepository struct {
	db *sql.DB
}

func NewGymRepository(db *sql.DB) *GymRepository {
	return &GymRepository{db: db}
}

func (r *GymRepository) Create(name, address string) (*models.Gym, error) {
	res, err := r.db.Exec(`INSERT INTO gyms (name, address) VALUES (?, ?)`, name, address)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return r.GetByID(int(id))
}

func (r *GymRepository) GetByID(id int) (*models.Gym, error) {
	g := &models.Gym{}
	err := r.db.QueryRow(`SELECT id, name, address, created_at FROM gyms WHERE id = ?`, id).
		Scan(&g.ID, &g.Name, &g.Address, &g.CreatedAt)
	return g, err
}

func (r *GymRepository) List() ([]models.Gym, error) {
	rows, err := r.db.Query(`SELECT id, name, address, created_at FROM gyms`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gyms []models.Gym
	for rows.Next() {
		var g models.Gym
		if err := rows.Scan(&g.ID, &g.Name, &g.Address, &g.CreatedAt); err != nil {
			return nil, err
		}
		gyms = append(gyms, g)
	}
	return gyms, nil
}

func (r *GymRepository) Update(id int, name, address string) error {
	_, err := r.db.Exec(`UPDATE gyms SET name = ?, address = ? WHERE id = ?`, name, address, id)
	return err
}

func (r *GymRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM gyms WHERE id = ?`, id)
	return err
}
