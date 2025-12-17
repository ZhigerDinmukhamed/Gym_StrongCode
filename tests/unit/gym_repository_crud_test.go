package unit

import (
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGymRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Old Name", "Old Address")

	err := repo.Update(gymID, "New Name", "New Address")
	require.NoError(t, err)

	gym, _ := repo.GetByID(gymID)
	assert.Equal(t, "New Name", gym.Name)
	assert.Equal(t, "New Address", gym.Address)
}

func TestGymRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewGymRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")

	err := repo.Delete(gymID)
	require.NoError(t, err)

	_, err = repo.GetByID(gymID)
	assert.Error(t, err)
}
