// internal/repository/user_repository_test.go
package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Create users table
	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT UNIQUE,
		password_hash TEXT,
		is_admin INTEGER,
		created_at DATETIME
	)`)
	if err != nil {
		panic(err)
	}

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupDB()
	defer db.Close()

	repo := NewUserRepository(db)

	user, err := repo.Create("John Doe", "john@example.com", "hashedpassword")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != "John Doe" || user.Email != "john@example.com" {
		t.Fatalf("expected user to be created with correct details, got %+v", user)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupDB()
	defer db.Close()

	repo := NewUserRepository(db)
	_, _ = repo.Create("Jane Doe", "jane@example.com", "hashedpassword")

	user, err := repo.GetByEmail("jane@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != "Jane Doe" {
		t.Fatalf("expected user name to be 'Jane Doe', got %s", user.Name)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupDB()
	defer db.Close()

	repo := NewUserRepository(db)
	createdUser, _ := repo.Create("Alice Smith", "alice@example.com", "hashedpassword")

	user, err := repo.GetByID(createdUser.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != createdUser.ID {
		t.Fatalf("expected user ID to be %d, got %d", createdUser.ID, user.ID)
	}
}

func TestUserRepository_CreateDuplicateEmail(t *testing.T) {
	db := setupDB()
	defer db.Close()

	repo := NewUserRepository(db)
	_, _ = repo.Create("Bob Brown", "bob@example.com", "hashedpassword")

	_, err := repo.Create("Bob Brown", "bob@example.com", "hashedpassword")
	if err == nil {
		t.Fatal("expected error for duplicate email, got none")
	}
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db := setupDB()
	defer db.Close()

	repo := NewUserRepository(db)

	user, err := repo.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Fatal("expected error for non-existent user, got none")
	}
	if user != nil {
		t.Fatalf("expected user to be nil, got %+v", user)
	}
}
