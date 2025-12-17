package unit

import (
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/tests/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentRepository_CreateStandalone(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)
	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)

	payment, err := repo.CreateStandalone(userID, 10000, "USD", "card", "completed")
	require.NoError(t, err)
	assert.NotZero(t, payment.ID)
	assert.Equal(t, userID, payment.UserID)
	assert.Equal(t, 10000, payment.AmountCents)
	assert.Equal(t, "completed", payment.Status)
}

func TestPaymentRepository_GetByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)
	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	created, _ := repo.CreateStandalone(userID, 5000, "USD", "card", "pending")

	payment, err := repo.GetByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, payment.ID)
	assert.Equal(t, "pending", payment.Status)
}

func TestPaymentRepository_GetByUser(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)
	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	// Создаем несколько платежей с разными статусами
	repo.CreateStandalone(userID, 1000, "USD", "card", "completed")
	repo.CreateStandalone(userID, 2000, "USD", "card", "pending")

	// Получаем все платежи пользователя
	payments, err := repo.GetByUser(userID, "")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(payments), 2)

	// Фильтруем только завершенные платежи
	completed, err := repo.GetByUser(userID, "completed")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(completed), 1)
	for _, p := range completed {
		assert.Equal(t, "completed", p.Status)
	}
}

func TestPaymentRepository_GetAll(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	repo := repository.NewPaymentRepository(db)
	userID := testutils.CreateTestUser(t, db, "user@test.com", "password", false)
	repo.CreateStandalone(userID, 1000, "USD", "card", "completed")

	payments, err := repo.GetAll()
	require.NoError(t, err)
	assert.NotEmpty(t, payments)
}
