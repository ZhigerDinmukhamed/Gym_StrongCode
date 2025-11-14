package main

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateClassRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TrainerID   int    `json:"trainer_id"`
	StartTime   string `json:"start_time"`
	DurationMin int    `json:"duration_min"`
	Capacity    int    `json:"capacity"`
}

type BookingRequest struct {
	ClassID int `json:"class_id"`
}

type BuyMembershipRequest struct {
	MembershipID int    `json:"membership_id"`
	Method       string `json:"method"`
}
