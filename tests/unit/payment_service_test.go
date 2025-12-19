package unit

import (
	"Gym_StrongCode/internal/repository"
	"Gym_StrongCode/internal/service"
	"Gym_StrongCode/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)

	payment, err := paymentService.Create(userID, 5000, "card", "Test payment", "REF123")
	require.NoError(t, err)
	assert.NotZero(t, payment.ID)
}

func TestPaymentService_GetByUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)

	// Создаем пользователя и платеж
	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	paymentRepo.CreateStandalone(userID, 1000, "USD", "card", "completed", "", "")

	payments, err := paymentService.GetByUser(userID, "")
	require.NoError(t, err)
	assert.NotEmpty(t, payments)
}

func TestPaymentService_ListAll(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	paymentRepo.CreateStandalone(userID, 1000, "USD", "card", "completed", "", "")

	payments, err := paymentService.ListAll()
	require.NoError(t, err)
	assert.NotEmpty(t, payments)
}

func TestPaymentService_CreateStandalone(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)

	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)

	payment, err := paymentService.CreateStandalone(userID, 10000, "cash")
	require.NoError(t, err)
	assert.NotZero(t, payment.ID)
	assert.Equal(t, "completed", payment.Status)
}
