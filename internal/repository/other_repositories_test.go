// internal/repository/other_repositories_test.go
package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
	"testing"
)

func TestTrainerRepository_Create(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:") // Use an in-memory database for testing
	repo := NewTrainerRepository(db)

	// Test creating a trainer
	id, err := repo.Create("John Doe", "Fitness Trainer")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected valid ID, got %d", id)
	}
}

func TestClassRepository_Create(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	repo := NewClassRepository(db)

	class := &models.Class{Title: "Yoga"} // Remove Details if not needed
	err := repo.Create(class)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestBookingRepository_Create(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	repo := NewBookingRepository(db)

	// Assuming userID and classID are valid
	userID := 1
	classID := 1
	id, err := repo.Create(userID, classID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected valid ID, got %d", id)
	}
}

func TestPaymentRepository_Create(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	repo := NewPaymentRepository(db)

	// Test creating a payment
	// Assuming necessary parameters are defined
	_, err := repo.CreateStandalone(1, 1, "Payment Method", "Transaction ID", "Description")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
