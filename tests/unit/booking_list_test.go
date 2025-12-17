package unit

import (
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingRepository_ListAll(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewBookingRepository(db)
	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)

	_, err := repo.Create(userID, classID)
	require.NoError(t, err)

	bookings, err := repo.ListAll()
	require.NoError(t, err)
	assert.NotEmpty(t, bookings)
}
