package unit

import (
	"testing"

	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGymRepository_CreateAndList(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewGymRepository(db)

	// Создаем залы
	gym1, err := repo.Create("Gym Alpha", "Address 1")
	require.NoError(t, err)
	assert.NotZero(t, gym1.ID)

	gym2, err := repo.Create("Gym Beta", "Address 2")
	require.NoError(t, err)
	assert.NotZero(t, gym2.ID)

	// Получаем список
	gyms, err := repo.List()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(gyms), 2)
}

func TestMembershipRepository_CRUD(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewMembershipRepository(db)

	// Создание подписки
	membership, err := repo.Create("Monthly", 30, 1000000)
	require.NoError(t, err)
	assert.NotZero(t, membership.ID)
	assert.Equal(t, "Monthly", membership.Name)

	// Получение по ID
	retrieved, err := repo.GetByID(membership.ID)
	require.NoError(t, err)
	assert.Equal(t, membership.ID, retrieved.ID)

	// Получение всех подписок
	memberships, err := repo.GetAll()
	require.NoError(t, err)
	assert.NotEmpty(t, memberships)

	// Обновление подписки
	err = repo.Update(membership.ID, "Monthly Pro", 30, 1200000)
	require.NoError(t, err)

	updated, _ := repo.GetByID(membership.ID)
	assert.Equal(t, "Monthly Pro", updated.Name)
	assert.Equal(t, 1200000, updated.PriceCents)

	// Удаление подписки
	err = repo.Delete(membership.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(membership.ID)
	assert.Error(t, err) // Должна быть ошибка sql.ErrNoRows
}
