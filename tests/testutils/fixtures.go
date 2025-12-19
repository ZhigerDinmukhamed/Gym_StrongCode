package testutils

import (
	"database/sql"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// CreateTestUser создает тестового пользователя
func CreateTestUser(t *testing.T, db *sql.DB, email, password string, isAdmin bool) int {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	result, err := db.Exec(
		"INSERT INTO users (name, email, password_hash, is_admin) VALUES (?, ?, ?, ?)",
		"Test User", email, string(hash), isAdmin,
	)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	id, _ := result.LastInsertId()
	return int(id)
}

// CreateTestGym создает тестовый зал
func CreateTestGym(t *testing.T, db *sql.DB, name, address string) int {
	result, err := db.Exec(
		"INSERT INTO gyms (name, address) VALUES (?, ?)",
		name, address,
	)
	if err != nil {
		t.Fatalf("Failed to create test gym: %v", err)
	}

	id, _ := result.LastInsertId()
	return int(id)
}

// CreateTestMembership создает тестовую подписку
func CreateTestMembership(t *testing.T, db *sql.DB, name string, durationDays, priceCents int) int {
	result, err := db.Exec(
		"INSERT INTO memberships (name, duration_days, price_cents) VALUES (?, ?, ?)",
		name, durationDays, priceCents,
	)
	if err != nil {
		t.Fatalf("Failed to create test membership: %v", err)
	}

	id, _ := result.LastInsertId()
	return int(id)
}

// CreateTestTrainer создает тестового тренера
func CreateTestTrainer(t *testing.T, db *sql.DB, name, bio string) int {
	result, err := db.Exec(
		"INSERT INTO trainers (name, bio) VALUES (?, ?)",
		name, bio,
	)
	if err != nil {
		t.Fatalf("Failed to create test trainer: %v", err)
	}

	id, _ := result.LastInsertId()
	return int(id)
}

// CreateTestClass создает тестовое занятие
func CreateTestClass(t *testing.T, db *sql.DB, title string, trainerID, gymID, capacity int) int {
	result, err := db.Exec(
		"INSERT INTO classes (title, description, trainer_id, gym_id, start_time, duration_min, capacity) VALUES (?, ?, ?, ?, datetime('now', '+1 day'), 60, ?)",
		title, "Test Description", trainerID, gymID, capacity,
	)
	if err != nil {
		t.Fatalf("Failed to create test class: %v", err)
	}

	id, _ := result.LastInsertId()
	return int(id)
}
