package models

type Payment struct {
	ID          int    `json:"id" db:"id"`
	UserID      int    `json:"user_id" db:"user_id"`
	AmountCents int    `json:"amount_cents" db:"amount_cents"`
	Currency    string `json:"currency" db:"currency"`
	Method      string `json:"method" db:"method"`
	Status      string `json:"status" db:"status"`
	Description string `json:"description" db:"description"`
	ReferenceID string `json:"reference_id" db:"reference_id"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}
