package models

type Gym struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Address   string `json:"address" db:"address"`
	CreatedAt string `json:"created_at" db:"created_at"`
}