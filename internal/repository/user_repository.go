package repository

import (
	"Gym_StrongCode/internal/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(name, email, passwordHash string, isAdmin bool) (*models.User, error) {
	res, err := r.db.Exec(`
		INSERT INTO users (name, email, password_hash, is_admin) 
		VALUES (?, ?, ?, ?)`, name, email, passwordHash, isAdmin)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	return r.GetByID(int(id))
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`SELECT id, name, email, password_hash, is_admin, created_at FROM users WHERE email = ?`, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, name, email, is_admin, created_at 
		FROM users WHERE id = ?`, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT id, name, email, is_admin, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsAdmin, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Update(id int, name, email string) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET name = ?, email = ? 
		WHERE id = ?`, name, email, id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}
