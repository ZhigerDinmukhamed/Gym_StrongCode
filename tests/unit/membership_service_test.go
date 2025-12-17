package unit

import (
	"testing"

	"Gym-StrongCode/config"
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/internal/service"
	"Gym-StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMembershipService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	membershipRepo := repository.NewMembershipRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db, notifService)

	// Создаем тестовую подписку
	testutils.CreateTestMembership(t, db, "Basic", 30, 5000)

	memberships, err := membershipService.List()
	require.NoError(t, err)
	assert.NotEmpty(t, memberships)
}

func TestMembershipService_Buy(t *testing.T) {
	t.Skip("MembershipService.Buy требует корректной работы с транзакциями - оставляем для интеграционных тестов")
}

func TestMembershipService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	membershipRepo := repository.NewMembershipRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db, notifService)

	membership, err := membershipService.Create("Gold", 60, 25000)
	require.NoError(t, err)
	assert.NotZero(t, membership.ID)
	assert.Equal(t, "Gold", membership.Name)
}

func TestMembershipService_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	membershipRepo := repository.NewMembershipRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db, notifService)

	membershipID := testutils.CreateTestMembership(t, db, "Old Name", 30, 10000)

	err := membershipService.Update(membershipID, "New Name", 45, 12000)
	require.NoError(t, err)
}

func TestMembershipService_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer db.Close()

	membershipRepo := repository.NewMembershipRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	notifService := service.NewNotificationService(&config.Config{})
	membershipService := service.NewMembershipService(membershipRepo, paymentRepo, db, notifService)

	membershipID := testutils.CreateTestMembership(t, db, "To Delete", 30, 10000)

	err := membershipService.Delete(membershipID)
	require.NoError(t, err)
}
