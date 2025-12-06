package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Gym_StrongCode/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/health", nil)

	HealthCheck(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Gym StrongCode API", response["service"])
}

func TestRegisterInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&service.AuthService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users/register", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&service.AuthService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users/login", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterMissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&service.AuthService{})

	payload := registerRequest{
		Name:     "Test User",
		Email:    "",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users/register", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTrainerInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTrainerHandler(&service.TrainerService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/admin/trainers", bytes.NewBufferString("invalid"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTrainer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateTrainerMissingName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewTrainerHandler(&service.TrainerService{})

	payload := createTrainerRequest{
		Name: "",
		Bio:  "Test bio",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/admin/trainers", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateTrainer(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateClassInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewClassHandler(&service.ClassService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/admin/classes", bytes.NewBufferString("invalid"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateClass(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateClassMissingTitle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewClassHandler(&service.ClassService{})

	payload := createClassRequest{
		Title:       "",
		Description: "Test",
		TrainerID:   1,
		StartTime:   time.Now().Add(time.Hour),
		DurationMin: 60,
		Capacity:    15,
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/admin/classes", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateClass(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBookingInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewBookingHandler(&service.BookingService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/bookings", bytes.NewBufferString("invalid"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateBooking(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePaymentInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewPaymentHandler(&service.PaymentService{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/payments", bytes.NewBufferString("invalid"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreatePayment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePaymentInvalidAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewPaymentHandler(&service.PaymentService{})

	payload := createPaymentRequest{
		AmountCents: 0,
		Method:      "card",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/payments", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreatePayment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
