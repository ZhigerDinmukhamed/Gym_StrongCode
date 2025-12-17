package unit

import (
	"testing"

	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/internal/service"
	"Gym-StrongCode/tests/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	db := testutils.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "test-secret")

	user, err := authService.Register("Test User", "test@example.com", "password123")

	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "test-secret")

	// Первая регистрация
	_, err := authService.Register("Test User", "duplicate@example.com", "password123")
	require.NoError(t, err)

	// Повторная регистрация с тем же email
	user, err := authService.Register("Another User", "duplicate@example.com", "password456")

	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestAuthService_Login_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "test-secret")

	// Регистрируем пользователя
	authService.Register("Test User", "test@example.com", "password123")

	// Логинимся
	token, err := authService.Login("test@example.com", "password123")

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	db := testutils.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "test-secret")

	// Создаем пользователя
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	userRepo.Create("Test", "test@example.com", string(hash), false)

	// Пытаемся логиниться с неверным паролем
	token, err := authService.Login("test@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	db := testutils.SetupTestDB(t)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "test-secret")

	token, err := authService.Login("notfound@example.com", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
}
