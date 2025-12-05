package repository

import (
	"os"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Error("expected non-nil database")
	}
}

func TestNewDatabaseInvalidPath(t *testing.T) {
	_, err := NewDatabase("/invalid/path/database.db")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestInitSchema(t *testing.T) {
	dbPath := "test_schema.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	err = InitSchema(db)
	if err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM memberships").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query memberships: %v", err)
	}
	if count != 4 {
		t.Errorf("expected 4 memberships, got %d", count)
	}
}

func TestInitSchemaCreatesAdminUser(t *testing.T) {
	dbPath := "test_admin.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	err = InitSchema(db)
	if err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	var adminEmail string
	err = db.QueryRow("SELECT email FROM users WHERE is_admin = 1").Scan(&adminEmail)
	if err != nil {
		t.Fatalf("failed to query admin user: %v", err)
	}
	if adminEmail != "admin@strongcode.kz" {
		t.Errorf("expected admin@strongcode.kz, got %s", adminEmail)
	}
}

func TestDatabasePingConnection(t *testing.T) {
	dbPath := "test_ping.db"
	defer os.Remove(dbPath)

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("NewDatabase failed: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}
}
