package unit

import (
	"testing"

	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewBookingRepository(db)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "pass123", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address 1")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Yoga", trainerID, gymID, 10)

	bookingID, err := repo.Create(userID, classID)

	require.NoError(t, err)
	assert.NotZero(t, bookingID)
}

func TestBookingRepository_Exists(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewBookingRepository(db)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "pass123", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address 1")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Yoga", trainerID, gymID, 10)

	// Создаем бронирование
	repo.Create(userID, classID)

	// Проверяем существование
	exists, err := repo.Exists(userID, classID)
	require.NoError(t, err)
	assert.True(t, exists)

	// Проверяем несуществующее бронирование
	exists, err = repo.Exists(userID, 9999)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestBookingRepository_GetByUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewBookingRepository(db)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "pass123", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address 1")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID1 := testutils.CreateTestClass(t, db, "Yoga", trainerID, gymID, 10)
	classID2 := testutils.CreateTestClass(t, db, "Pilates", trainerID, gymID, 15)

	// Создаем несколько бронирований
	repo.Create(userID, classID1)
	repo.Create(userID, classID2)

	bookings, err := repo.GetByUser(userID)

	require.NoError(t, err)
	assert.Len(t, bookings, 2)
}

func TestBookingRepository_Cancel(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := repository.NewBookingRepository(db)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "pass123", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address 1")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Yoga", trainerID, gymID, 10)

	bookingID, _ := repo.Create(userID, classID)

	// Отменяем бронирование
	err := repo.Cancel(int(bookingID), userID)
	require.NoError(t, err)

	// Проверяем, что оно удалено
	exists, _ := repo.Exists(userID, classID)
	assert.False(t, exists)
}
