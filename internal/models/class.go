package models

type Class struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	TrainerID   int    `json:"trainer_id" db:"trainer_id"`
	GymID       int    `json:"gym_id" db:"gym_id"`        // новый FK
	StartTime   string `json:"start_time" db:"start_time"`
	DurationMin int    `json:"duration_min" db:"duration_min"`
	Capacity    int    `json:"capacity" db:"capacity"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}