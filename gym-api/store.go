package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	DB *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1)
	return &Store{DB: db}, nil
}

func (s *Store) InitSchema() error {
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
	FOREIGN KEY(class_id) REFERENCES classes(id) ON DELETE CASCADE
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
`
	_, err := s.DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}

	// seed some sample memberships and a sample admin user if not exists
	_, _ = s.DB.Exec(`INSERT OR IGNORE INTO memberships(id, name, duration_days, price_cents) VALUES
	(1,'Monthly',30,3000),
	(2,'Yearly',365,30000),
	(3,'VIP',365,90000)`)

	// create admin if not exists
	var adminCount int
	_ = s.DB.QueryRow("SELECT COUNT(1) FROM users WHERE is_admin = 1").Scan(&adminCount)
	if adminCount == 0 {
		pw, _ := HashPassword("adminpass")
		_, _ = s.DB.Exec("INSERT INTO users(name,email,password_hash,is_admin) VALUES(?,?,?,1)", "Admin", "admin@example.com", pw)
	}

	return nil
}
