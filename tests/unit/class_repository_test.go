package unit

import (
	"testing"
	"time"

	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewClassRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")

	class := &models.Class{
		Title:       "Yoga Class",
		Description: "Beginner Yoga",
		TrainerID:   trainerID,
		GymID:       gymID,
		StartTime:   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		DurationMin: 60,
		Capacity:    20,
	}

	created, err := repo.Create(class)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "Yoga Class", created.Title)
}

func TestClassRepository_GetByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewClassRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)

	class, err := repo.GetByID(classID)
	require.NoError(t, err)
	assert.Equal(t, classID, class.ID)
	assert.Equal(t, "Test Class", class.Title)
}

func TestClassRepository_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewClassRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	testutils.CreateTestClass(t, db, "Class 1", trainerID, gymID, 20)
	testutils.CreateTestClass(t, db, "Class 2", trainerID, gymID, 20)

	classes, err := repo.List()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(classes), 2)
}

func TestClassRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewClassRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Old Title", trainerID, gymID, 20)

	// Обновляем данные занятия
	updated := &models.Class{
		Title:       "New Title",
		Description: "New Description",
		TrainerID:   trainerID,
		GymID:       gymID,
		StartTime:   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		DurationMin: 90,
		Capacity:    25,
	}

	err := repo.Update(classID, updated)
	require.NoError(t, err)

	class, _ := repo.GetByID(classID)
	assert.Equal(t, "New Title", class.Title)
	assert.Equal(t, 90, class.DurationMin)
}

func TestClassRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewClassRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)

	err := repo.Delete(classID)
	require.NoError(t, err)

	_, err = repo.GetByID(classID)
	assert.Error(t, err)
}

func TestClassRepository_GetBookingCount(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewClassRepository(db)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)

	count, err := repo.GetBookingCount(classID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
