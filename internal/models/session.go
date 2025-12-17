package models

import "time"

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	UserEmail string    `json:"user_email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
