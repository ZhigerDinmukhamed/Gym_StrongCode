package main

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
}

type Membership struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	DurationDays int       `json:"duration_days"`
	PriceCents   int       `json:"price_cents"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserMembership struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	MembershipID int       `json:"membership_id"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

type Trainer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type Class struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TrainerID   int       `json:"trainer_id"`
	StartTime   time.Time `json:"start_time"`
	DurationMin int       `json:"duration_min"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
}

type Booking struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ClassID   int       `json:"class_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Payment struct {
	ID         int       `json:"id"`
	UserID     *int      `json:"user_id"`
	AmountCents int      `json:"amount_cents"`
	Currency   string    `json:"currency"`
	Method     string    `json:"method"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
