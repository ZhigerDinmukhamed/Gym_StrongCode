package repository

import (
	"database/sql"
	"log"

	"Gym-StrongCode/internal/utils"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func NewDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Включаем FK
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	// Автоматические миграции
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite",
		driver,
	)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		utils.GetLogger().Fatal("Migration failed", zap.Error(err))
		return nil, err
	}

	log.Println("Database migrated successfully")
	return db, nil
}
