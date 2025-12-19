package unit

import (
	"testing"

	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	user, err := repo.Create("John Doe", "john@example.com", "hashedpassword", false)

	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.False(t, user.IsAdmin)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	// Создаем пользователя
	created, _ := repo.Create("Jane Doe", "jane@example.com", "hashedpass", false)

	// Получаем по email
	user, err := repo.GetByEmail("jane@example.com")

	require.NoError(t, err)
	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, "Jane Doe", user.Name)
	assert.Equal(t, "hashedpass", user.PasswordHash)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	user, err := repo.GetByEmail("nonexistent@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	created, _ := repo.Create("Bob Smith", "bob@example.com", "hashedpass", true)

	user, err := repo.GetByID(created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, "Bob Smith", user.Name)
	assert.True(t, user.IsAdmin)
}

func TestUserRepository_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	// Создаем несколько пользователей
	repo.Create("User 1", "user1@example.com", "pass1", false)
	repo.Create("User 2", "user2@example.com", "pass2", false)
	repo.Create("User 3", "user3@example.com", "pass3", true)

	users, err := repo.List()

	require.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestUserRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	created, _ := repo.Create("Delete Me", "delete@example.com", "pass", false)

	// Удаляем
	err := repo.Delete(created.ID)
	require.NoError(t, err)

	// Проверяем, что пользователь удален
	user, err := repo.GetByID(created.ID)
	assert.Error(t, err)
	assert.Nil(t, user)
}
