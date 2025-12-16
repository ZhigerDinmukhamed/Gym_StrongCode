package models

type Trainer struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Bio       string `json:"bio" db:"bio"`
	CreatedAt string `json:"created_at" db:"created_at"`
}