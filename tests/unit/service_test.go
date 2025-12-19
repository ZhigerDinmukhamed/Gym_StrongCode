package unit

import (
	"testing"

	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/internal/service"
	"Gym_StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGymService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	gymRepo := repository.NewGymRepository(db)
	gymService := service.NewGymService(gymRepo)

	gym, err := gymService.Create("Test Gym", "Test Address")

	require.NoError(t, err)
	assert.NotZero(t, gym.ID)
	assert.Equal(t, "Test Gym", gym.Name)
	assert.Equal(t, "Test Address", gym.Address)
}

func TestGymService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	gymRepo := repository.NewGymRepository(db)
	gymService := service.NewGymService(gymRepo)

	// Создаем несколько залов
	gymService.Create("Gym 1", "Address 1")
	gymService.Create("Gym 2", "Address 2")

	gyms, err := gymService.List()

	require.NoError(t, err)
	assert.Len(t, gyms, 2)
}

func TestTrainerService_CRUD(t *testing.T) {
	db := testutils.SetupTestDB(t)
	trainerRepo := repository.NewTrainerRepository(db)
	trainerService := service.NewTrainerService(trainerRepo)

	// Создание тренера
	trainer, err := trainerService.Create("John Coach", "Expert in CrossFit")
	require.NoError(t, err)
	assert.NotZero(t, trainer.ID)

	// Получение списка
	trainers, err := trainerService.List()
	require.NoError(t, err)
	assert.NotEmpty(t, trainers)

	// Обновление данных
	err = trainerService.Update(trainer.ID, "John Pro Coach", "Master in CrossFit")
	require.NoError(t, err)

	// Удаление
	err = trainerService.Delete(trainer.ID)
	require.NoError(t, err)
}
