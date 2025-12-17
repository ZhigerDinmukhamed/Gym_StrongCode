package unit

import (
	"testing"

	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/internal/service"
	"Gym-StrongCode/tests/testutils"

	"github.com/stretchr/testify/require"
)

func TestGymService_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	gymRepo := repository.NewGymRepository(db)
	gymService := service.NewGymService(gymRepo)

	gymID := testutils.CreateTestGym(t, db, "Old Name", "Old Address")

	err := gymService.Update(gymID, "New Name", "New Address")
	require.NoError(t, err)
}

func TestGymService_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	gymRepo := repository.NewGymRepository(db)
	gymService := service.NewGymService(gymRepo)

	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")

	err := gymService.Delete(gymID)
	require.NoError(t, err)
}
