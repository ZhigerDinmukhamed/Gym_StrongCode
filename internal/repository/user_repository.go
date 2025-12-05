// internal/repository/user_repository.go
package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// UserRepository управляет операциями с пользователями в БД
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(name, email, passwordHash string) (*models.User, error) {
	// Проверяем уникальность email
	var count int
	err := r.db.QueryRow("SELECT COUNT(1) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	res, err := r.db.Exec(
		"INSERT INTO users(name, email, password_hash) VALUES(?, ?, ?)",
		name, email, passwordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &models.User{
		ID:        int(id),
		Name:      name,
		Email:     email,
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}, nil
}

// GetByEmail возвращает пользователя по email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	var isAdmin int
	var createdAt string

	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, is_admin, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &isAdmin, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.IsAdmin = isAdmin == 1
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	return &user, nil
}

// GetByID возвращает пользователя по ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	var isAdmin int
	var createdAt string

	err := r.db.QueryRow(
		"SELECT id, name, email, is_admin, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &isAdmin, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.IsAdmin = isAdmin == 1
	user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	return &user, nil
}
