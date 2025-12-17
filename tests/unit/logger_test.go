package unit

import (
	"Gym-StrongCode/internal/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	// Удаляем папку logs если существует для чистоты теста
	os.RemoveAll("logs")
	defer os.RemoveAll("logs")

	// Инициализация не должна паниковать
	assert.NotPanics(t, func() {
		utils.InitLogger()
	})

	// Logger должен быть установлен
	assert.NotNil(t, utils.Logger)

	// Папка logs должна существовать
	_, err := os.Stat("logs")
	assert.NoError(t, err)
}

func TestGetLogger(t *testing.T) {
	os.RemoveAll("logs")
	defer os.RemoveAll("logs")

	// GetLogger должен инициализировать логгер если его нет
	utils.Logger = nil
	logger := utils.GetLogger()
	assert.NotNil(t, logger)
	assert.NotNil(t, utils.Logger)
}
