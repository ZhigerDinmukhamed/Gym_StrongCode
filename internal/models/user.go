package models

type User struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"` // Not exposed in JSON
	IsAdmin      bool   `json:"is_admin" db:"is_admin"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}
