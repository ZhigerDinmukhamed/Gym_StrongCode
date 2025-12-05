package repository

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	schema := `
	CREATE TABLE memberships (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		duration_days INTEGER NOT NULL,
		price_cents INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE user_memberships (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		membership_id INTEGER NOT NULL,
		start_date DATE NOT NULL,
		end_date DATE NOT NULL,
		active INTEGER DEFAULT 1
	);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestGetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.Exec("INSERT INTO memberships (name, duration_days, price_cents, created_at) VALUES (?, ?, ?, ?)",
		"Premium", 30, 9999, "2024-01-01 10:00:00")
	db.Exec("INSERT INTO memberships (name, duration_days, price_cents, created_at) VALUES (?, ?, ?, ?)",
		"Standard", 7, 2999, "2024-01-02 10:00:00")

	repo := NewMembershipRepository(db)
	memberships, err := repo.GetAll()

	if err != nil {
		t.Errorf("GetAll failed: %v", err)
	}
	if len(memberships) != 2 {
		t.Errorf("expected 2 memberships, got %d", len(memberships))
	}
	if memberships[0].Name != "Premium" {
		t.Errorf("expected Premium, got %s", memberships[0].Name)
	}
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	db.Exec("INSERT INTO memberships (id, name, duration_days, price_cents, created_at) VALUES (?, ?, ?, ?, ?)",
		1, "Gold", 90, 19999, "2024-01-01 10:00:00")

	repo := NewMembershipRepository(db)
	membership, err := repo.GetByID(1)

	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}
	if membership.Name != "Gold" {
		t.Errorf("expected Gold, got %s", membership.Name)
	}
	if membership.PriceCents != 19999 {
		t.Errorf("expected 19999, got %d", membership.PriceCents)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMembershipRepository(db)
	_, err := repo.GetByID(999)

	if err == nil {
		t.Error("expected error for non-existent membership")
	}
}

func TestCreateUserMembership(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tx, _ := db.Begin()
	defer tx.Rollback()

	repo := NewMembershipRepository(db)
	err := repo.CreateUserMembership(tx, 1, 1, "2024-01-01", "2024-01-31")

	if err != nil {
		t.Errorf("CreateUserMembership failed: %v", err)
	}
}

func TestHasActiveMembership(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	db.Exec("INSERT INTO user_memberships (user_id, membership_id, start_date, end_date, active) VALUES (?, ?, ?, ?, ?)",
		1, 1, today, tomorrow, 1)

	repo := NewMembershipRepository(db)
	has, err := repo.HasActiveMembership(1)

	if err != nil {
		t.Errorf("HasActiveMembership failed: %v", err)
	}
	if !has {
		t.Error("expected user to have active membership")
	}
}

func TestHasActiveMembershipExpired(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	db.Exec("INSERT INTO user_memberships (user_id, membership_id, start_date, end_date, active) VALUES (?, ?, ?, ?, ?)",
		2, 1, "2024-01-01", yesterday, 1)

	repo := NewMembershipRepository(db)
	has, err := repo.HasActiveMembership(2)

	if err != nil {
		t.Errorf("HasActiveMembership failed: %v", err)
	}
	if has {
		t.Error("expected user to not have active membership")
	}
}
