package unit

import (
	"Gym-StrongCode/internal/models"
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/internal/service"
	"Gym-StrongCode/tests/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	classRepo := repository.NewClassRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	gymRepo := repository.NewGymRepository(db)
	classService := service.NewClassService(classRepo, trainerRepo, gymRepo)

	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")

	class := &models.Class{
		Title:       "Yoga",
		Description: "Beginner",
		TrainerID:   trainerID,
		GymID:       gymID,
		StartTime:   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		DurationMin: 60,
		Capacity:    20,
	}

	created, err := classService.Create(class)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "Yoga", created.Title)
}

func TestClassService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	classRepo := repository.NewClassRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	gymRepo := repository.NewGymRepository(db)
	classService := service.NewClassService(classRepo, trainerRepo, gymRepo)

	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	testutils.CreateTestClass(t, db, "Class 1", trainerID, gymID, 20)

	classes, err := classService.List()
	require.NoError(t, err)
	assert.NotEmpty(t, classes)
}

func TestClassService_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	classRepo := repository.NewClassRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	gymRepo := repository.NewGymRepository(db)
	classService := service.NewClassService(classRepo, trainerRepo, gymRepo)

	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Old Title", trainerID, gymID, 20)

	updated := &models.Class{
		Title:       "New Title",
		Description: "New Description",
		TrainerID:   trainerID,
		GymID:       gymID,
		StartTime:   time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		DurationMin: 90,
		Capacity:    25,
	}

	err := classService.Update(classID, updated)
	require.NoError(t, err)
}

func TestClassService_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	classRepo := repository.NewClassRepository(db)
	trainerRepo := repository.NewTrainerRepository(db)
	gymRepo := repository.NewGymRepository(db)
	classService := service.NewClassService(classRepo, trainerRepo, gymRepo)

	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)

	err := classService.Delete(classID)
	require.NoError(t, err)
}
