package models

type UserMembership struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"user_id" db:"user_id"`
	MembershipID int  `json:"membership_id" db:"membership_id"`
	StartDate  string `json:"start_date" db:"start_date"`
	EndDate    string `json:"end_date" db:"end_date"`
	Active     bool   `json:"active" db:"active"`
	CreatedAt  string `json:"created_at" db:"created_at"`
}