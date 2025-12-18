package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
)

type TrainerRepository struct {
	db *sql.DB
}

func NewTrainerRepository(db *sql.DB) *TrainerRepository {
	return &TrainerRepository{db: db}
}

func (r *TrainerRepository) Create(name, bio string) (*models.Trainer, error) {
	res, err := r.db.Exec(`INSERT INTO trainers (name, bio) VALUES (?, ?)`, name, bio)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return r.GetByID(int(id))
}

func (r *TrainerRepository) GetByID(id int) (*models.Trainer, error) {
	t := &models.Trainer{}
	err := r.db.QueryRow(`SELECT id, name, bio, created_at FROM trainers WHERE id = ?`, id).
		Scan(&t.ID, &t.Name, &t.Bio, &t.CreatedAt)
	return t, err
}

func (r *TrainerRepository) List() ([]models.Trainer, error) {
	rows, err := r.db.Query(`SELECT id, name, bio, created_at FROM trainers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainers []models.Trainer
	for rows.Next() {
		var t models.Trainer
		if err := rows.Scan(&t.ID, &t.Name, &t.Bio, &t.CreatedAt); err != nil {
			return nil, err
		}
		trainers = append(trainers, t)
	}
	return trainers, nil
}

func (r *TrainerRepository) Update(id int, name, bio string) error {
	_, err := r.db.Exec(`UPDATE trainers SET name = ?, bio = ? WHERE id = ?`, name, bio, id)
	return err
}

func (r *TrainerRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM trainers WHERE id = ?`, id)
	return err
}
