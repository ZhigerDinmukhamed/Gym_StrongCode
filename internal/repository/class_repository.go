package repository

import (
	"database/sql"
	"Gym-StrongCode/internal/models"
)

type ClassRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

func (r *ClassRepository) Create(c *models.Class) (*models.Class, error) {
	res, err := r.db.Exec(`
		INSERT INTO classes (title, description, trainer_id, gym_id, start_time, duration_min, capacity)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		c.Title, c.Description, c.TrainerID, c.GymID, c.StartTime, c.DurationMin, c.Capacity)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return r.GetByID(int(id))
}

func (r *ClassRepository) GetByID(id int) (*models.Class, error) {
	c := &models.Class{}
	err := r.db.QueryRow(`
		SELECT id, title, description, trainer_id, gym_id, start_time, duration_min, capacity, created_at
		FROM classes WHERE id = ?`, id).
		Scan(&c.ID, &c.Title, &c.Description, &c.TrainerID, &c.GymID, &c.StartTime, &c.DurationMin, &c.Capacity, &c.CreatedAt)
	return c, err
}

func (r *ClassRepository) List() ([]models.Class, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, trainer_id, gym_id, start_time, duration_min, capacity, created_at
		FROM classes ORDER BY start_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var c models.Class
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.TrainerID, &c.GymID, &c.StartTime, &c.DurationMin, &c.Capacity, &c.CreatedAt); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, nil
}

func (r *ClassRepository) Update(id int, c *models.Class) error {
	_, err := r.db.Exec(`
		UPDATE classes SET title = ?, description = ?, trainer_id = ?, gym_id = ?, start_time = ?, duration_min = ?, capacity = ?
		WHERE id = ?`,
		c.Title, c.Description, c.TrainerID, c.GymID, c.StartTime, c.DurationMin, c.Capacity, id)
	return err
}

func (r *ClassRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM classes WHERE id = ?`, id)
	return err
}

func (r *ClassRepository) GetBookingCount(classID int) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM bookings WHERE class_id = ?`, classID).Scan(&count)
	return count, err
}