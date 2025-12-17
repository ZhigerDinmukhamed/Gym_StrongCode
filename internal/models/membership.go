package models

type Membership struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	DurationDays int    `json:"duration_days" db:"duration_days"`
	PriceCents   int    `json:"price_cents" db:"price_cents"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}