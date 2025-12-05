// internal/repository/database.go
package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// NewDatabase создает новое подключение к базе данных
func NewDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// InitSchema инициализирует схему базы данных
func InitSchema(db *sql.DB) error {
	schema := `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	is_admin INTEGER DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS memberships (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	duration_days INTEGER NOT NULL,
	price_cents INTEGER NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_memberships (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	membership_id INTEGER NOT NULL,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	active INTEGER DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(membership_id) REFERENCES memberships(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS trainers (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	bio TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS classes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT,
	trainer_id INTEGER,
	start_time DATETIME NOT NULL,
	duration_min INTEGER NOT NULL,
	capacity INTEGER DEFAULT 20,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(trainer_id) REFERENCES trainers(id)
);

CREATE TABLE IF NOT EXISTS bookings (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	class_id INTEGER NOT NULL,
	status TEXT DEFAULT 'booked',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE,
	UNIQUE(user_id, class_id)
);

CREATE TABLE IF NOT EXISTS payments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	amount_cents INTEGER NOT NULL,
	currency TEXT NOT NULL,
	method TEXT,
	status TEXT DEFAULT 'done',
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_user_memberships_user ON user_memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_user_memberships_active ON user_memberships(active, start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_class ON bookings(class_id);
CREATE INDEX IF NOT EXISTS idx_classes_start ON classes(start_time);
`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	// Seed data
	if err := seedData(db); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	return nil
}

// seedData заполняет базу начальными данными
func seedData(db *sql.DB) error {
	// Seed memberships
	_, err := db.Exec(`
		INSERT OR IGNORE INTO memberships(id, name, duration_days, price_cents) VALUES
		(1, 'Месячная', 30, 15000),
		(2, 'Квартальная', 90, 40000),
		(3, 'Годовая', 365, 150000),
		(4, 'VIP Годовая', 365, 300000)
	`)
	if err != nil {
		return fmt.Errorf("failed to seed memberships: %w", err)
	}

	// Create admin user if not exists
	var adminCount int
	if err := db.QueryRow("SELECT COUNT(1) FROM users WHERE is_admin = 1").Scan(&adminCount); err != nil {
		return fmt.Errorf("failed to check admin count: %w", err)
	}

	if adminCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash admin password: %w", err)
		}

		_, err = db.Exec(
			"INSERT INTO users(name, email, password_hash, is_admin) VALUES(?, ?, ?, 1)",
			"Администратор", "admin@strongcode.kz", string(hash),
		)
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	return nil
}
