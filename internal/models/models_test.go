package models

import (
	"testing"
	"time"
)

func TestUserStruct(t *testing.T) {
	now := time.Now()
	user := User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashedpassword",
		IsAdmin:      true,
		CreatedAt:    now,
	}

	if user.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("expected Name 'John Doe', got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("expected Email 'john@example.com', got %s", user.Email)
	}
	if !user.IsAdmin {
		t.Errorf("expected IsAdmin true, got false")
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, user.CreatedAt)
	}
}

func TestMembershipStruct(t *testing.T) {
	now := time.Now()
	m := Membership{
		ID:           2,
		Name:         "Premium",
		DurationDays: 30,
		PriceCents:   4999,
		CreatedAt:    now,
	}

	if m.Name != "Premium" {
		t.Errorf("expected Name 'Premium', got %s", m.Name)
	}
	if m.DurationDays != 30 {
		t.Errorf("expected DurationDays 30, got %d", m.DurationDays)
	}
	if m.PriceCents != 4999 {
		t.Errorf("expected PriceCents 4999, got %d", m.PriceCents)
	}
}

func TestUserMembershipStruct(t *testing.T) {
	now := time.Now()
	um := UserMembership{
		ID:           3,
		UserID:       1,
		MembershipID: 2,
		StartDate:    "2024-06-01",
		EndDate:      "2024-07-01",
		Active:       true,
		CreatedAt:    now,
	}

	if !um.Active {
		t.Errorf("expected Active true, got false")
	}
	if um.StartDate != "2024-06-01" {
		t.Errorf("expected StartDate '2024-06-01', got %s", um.StartDate)
	}
	if um.EndDate != "2024-07-01" {
		t.Errorf("expected EndDate '2024-07-01', got %s", um.EndDate)
	}
}

func TestTrainerStruct(t *testing.T) {
	now := time.Now()
	trainer := Trainer{
		ID:        4,
		Name:      "Alice Smith",
		Bio:       "Certified fitness trainer",
		CreatedAt: now,
	}

	if trainer.Name != "Alice Smith" {
		t.Errorf("expected Name 'Alice Smith', got %s", trainer.Name)
	}
	if trainer.Bio != "Certified fitness trainer" {
		t.Errorf("expected Bio 'Certified fitness trainer', got %s", trainer.Bio)
	}
}

func TestClassStruct(t *testing.T) {
	now := time.Now()
	class := Class{
		ID:          5,
		Title:       "Yoga",
		Description: "Morning yoga session",
		TrainerID:   4,
		StartTime:   now,
		DurationMin: 60,
		Capacity:    20,
		CreatedAt:   now,
	}

	if class.Title != "Yoga" {
		t.Errorf("expected Title 'Yoga', got %s", class.Title)
	}
	if class.DurationMin != 60 {
		t.Errorf("expected DurationMin 60, got %d", class.DurationMin)
	}
	if class.Capacity != 20 {
		t.Errorf("expected Capacity 20, got %d", class.Capacity)
	}
}

func TestBookingStruct(t *testing.T) {
	now := time.Now()
	booking := Booking{
		ID:        6,
		UserID:    1,
		ClassID:   5,
		Status:    "confirmed",
		CreatedAt: now,
	}

	if booking.Status != "confirmed" {
		t.Errorf("expected Status 'confirmed', got %s", booking.Status)
	}
}

func TestPaymentStruct(t *testing.T) {
	now := time.Now()
	uid := 1
	payment := Payment{
		ID:          7,
		UserID:      &uid,
		AmountCents: 4999,
		Currency:    "USD",
		Method:      "card",
		Status:      "completed",
		CreatedAt:   now,
	}

	if payment.UserID == nil || *payment.UserID != 1 {
		t.Errorf("expected UserID 1, got %v", payment.UserID)
	}
	if payment.AmountCents != 4999 {
		t.Errorf("expected AmountCents 4999, got %d", payment.AmountCents)
	}
	if payment.Currency != "USD" {
		t.Errorf("expected Currency 'USD', got %s", payment.Currency)
	}
	if payment.Status != "completed" {
		t.Errorf("expected Status 'completed', got %s", payment.Status)
	}
}
