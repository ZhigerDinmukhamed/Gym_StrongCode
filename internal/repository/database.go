// internal/repository/database.go
package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

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

	// Запускаем миграции
	if err := runMigrations(db, path); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Создаём админа, если его нет
	if err := ensureAdminUser(db); err != nil {
		log.Printf("Warning: failed to create admin user: %v", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB, dbPath string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	// ИСПРАВЛЕНО: было "everyone" → стало "err"
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Миграции успешно применены")
	return nil
}

func ensureAdminUser(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", "admin@strongcode.kz").Scan(&count)
	if err != nil || count > 0 {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO users(name, email, password_hash, is_admin) VALUES(?, ?, ?, 1)",
		"Администратор", "admin@strongcode.kz", string(hash),
	)
	return err
}