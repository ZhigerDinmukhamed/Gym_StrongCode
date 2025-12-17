package unit

import (
	"testing"

	"Gym-StrongCode/config"
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/internal/service"
	"Gym-StrongCode/internal/utils"
	"Gym-StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()
	utils.InitLogger()

	bookingRepo := repository.NewBookingRepository(db)
	classRepo := repository.NewClassRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo, notifService)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)

	err := bookingService.Create(userID, classID, "user@test.com")
	require.NoError(t, err)
}

func TestBookingService_ListUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	bookingRepo := repository.NewBookingRepository(db)
	classRepo := repository.NewClassRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo, notifService)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)
	bookingRepo.Create(userID, classID)

	bookings, err := bookingService.ListUser(userID)
	require.NoError(t, err)
	assert.NotNil(t, bookings)
}

func TestBookingService_ListAll(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	bookingRepo := repository.NewBookingRepository(db)
	classRepo := repository.NewClassRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo, notifService)

	// Создаем тестовое бронирование для проверки
	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)
	bookingRepo.Create(userID, classID)

	bookings, err := bookingService.ListAll()
	require.NoError(t, err)
	assert.NotNil(t, bookings)
}

func TestBookingService_Cancel(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	bookingRepo := repository.NewBookingRepository(db)
	classRepo := repository.NewClassRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	bookingService := service.NewBookingService(bookingRepo, classRepo, membershipRepo, notifService)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	gymID := testutils.CreateTestGym(t, db, "Test Gym", "Address")
	trainerID := testutils.CreateTestTrainer(t, db, "Trainer", "Bio")
	classID := testutils.CreateTestClass(t, db, "Test Class", trainerID, gymID, 20)
	bookingID, _ := bookingRepo.Create(userID, classID)

	err := bookingService.Cancel(int(bookingID), userID)
	require.NoError(t, err)
}
