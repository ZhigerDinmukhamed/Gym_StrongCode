package unit

import (
	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/internal/service"
	"Gym_StrongCode/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userID := testutils.CreateTestUser(t, db, "test@example.com", "password", false)

	user, err := userService.GetByID(userID)
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	_, err := userService.GetByID(99999)
	require.Error(t, err)
}

func TestUserService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	// Создаем несколько тестовых пользователей
	testutils.CreateTestUser(t, db, "user1@test.com", "password", false)
	testutils.CreateTestUser(t, db, "user2@test.com", "password", false)

	users, err := userService.List()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 2)
}

func TestUserService_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userID := testutils.CreateTestUser(t, db, "delete@test.com", "password", false)

	err := userService.Delete(userID)
	require.NoError(t, err)

	_, err = userService.GetByID(userID)
	assert.Error(t, err)
}
