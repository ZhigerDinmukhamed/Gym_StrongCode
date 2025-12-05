package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// --- Mocks ---

type mockAuthService struct {
	registerFunc func(name, email, password string) (*models.User, error)
	loginFunc    func(email, password string) (string, error)
}

func (m *mockAuthService) Register(name, email, password string) (*models.User, error) {
	if m.registerFunc != nil {
		return m.registerFunc(name, email, password)
	}
	return nil, nil
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(email, password)
	}
	return "", nil
}

var _ service.AuthService = (*mockAuthService)(nil)

type mockUserRepo struct {
	getByIDFunc func(id int) (*models.User, error)
}

func (m *mockUserRepo) GetByID(id int) (*models.User, error) {
	return m.getByIDFunc(id)
}

type mockMembershipService struct {
	getAllFunc        func() ([]models.Membership, error)
	buyMembershipFunc func(userID, membershipID int, method string) (map[string]interface{}, error)
}

func (m *mockMembershipService) GetAll() ([]models.Membership, error) {
	return m.getAllFunc()
}
func (m *mockMembershipService) BuyMembership(userID, membershipID int, method string) (map[string]interface{}, error) {
	return m.buyMembershipFunc(userID, membershipID, method)
}

type mockTrainerService struct {
	createFunc func(name, bio string) (int, error)
}

func (m *mockTrainerService) Create(name, bio string) (int, error) {
	return m.createFunc(name, bio)
}

type mockClassService struct {
	getAllFunc func() ([]models.Class, error)
	createFunc func(class *models.Class) error
}

func (m *mockClassService) GetAll() ([]models.Class, error) {
	return m.getAllFunc()
}
func (m *mockClassService) Create(class *models.Class) error {
	return m.createFunc(class)
}

type mockBookingService struct {
	createFunc    func(userID, classID int) (int, time.Time, error)
	getByUserFunc func(userID int) ([]models.Booking, error)
}

func (m *mockBookingService) Create(userID, classID int) (int, time.Time, error) {
	return m.createFunc(userID, classID)
}
func (m *mockBookingService) GetByUser(userID int) ([]models.Booking, error) {
	return m.getByUserFunc(userID)
}

type mockPaymentService struct {
	createFunc func(userID, amountCents int, method string) (int, error)
}

func (m *mockPaymentService) Create(userID, amountCents int, method string) (int, error) {
	return m.createFunc(userID, amountCents, method)
}

var _ service.PaymentService = (*mockPaymentService)(nil)

// --- Tests ---

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{
		registerFunc: func(name, email, password string) (*models.User, error) {
			return &models.User{ID: 1, Name: name, Email: email}, nil
		},
	}
	h := NewAuthHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"name":"Test","email":"test@test.kz","password":"123456"}`
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test")
}

func TestAuthHandler_Register_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{
		registerFunc: func(name, email, password string) (*models.User, error) {
			return nil, errors.New("email exists")
		},
	}
	h := NewAuthHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"name":"Test","email":"test@test.kz","password":"123456"}`
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{
		loginFunc: func(email, password string) (string, error) {
			return "token123", nil
		},
	}
	h := NewAuthHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"email":"test@test.kz","password":"123456"}`
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token123")
}

func TestAuthHandler_Login_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockAuthService{
		loginFunc: func(email, password string) (string, error) {
			return "", errors.New("invalid credentials")
		},
	}
	h := NewAuthHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"email":"test@test.kz","password":"wrong"}`
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_GetCurrentUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockUserRepo{
		getByIDFunc: func(id int) (*models.User, error) {
			return &models.User{ID: id, Name: "Test"}, nil
		},
	}
	h := NewUserHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)

	h.GetCurrentUser(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test")
}

func TestMembershipHandler_GetMemberships_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockMembershipService{
		getAllFunc: func() ([]models.Membership, error) {
			return []models.Membership{{ID: 1, Name: "Gold"}}, nil
		},
	}
	h := NewMembershipHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetMemberships(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Gold")
}

func TestTrainerHandler_CreateTrainer_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockTrainerService{
		createFunc: func(name, bio string) (int, error) {
			return 42, nil
		},
	}
	h := NewTrainerHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"name":"Trainer","bio":"Bio"}`
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateTrainer(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "trainer_id")
}

func TestClassHandler_GetClasses_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockClassService{
		getAllFunc: func() ([]models.Class, error) {
			return []models.Class{{ID: 1, Title: "Yoga"}}, nil
		},
	}
	h := NewClassHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetClasses(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Yoga")
}

func TestBookingHandler_ListBookings_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockBookingService{
		getByUserFunc: func(userID int) ([]models.Booking, error) {
			return []models.Booking{{ID: 1, ClassID: 2}}, nil
		},
	}
	h := NewBookingHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)

	h.ListBookings(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ClassID")
}

func TestPaymentHandler_CreatePayment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock := &mockPaymentService{
		createFunc: func(userID, amountCents int, method string) (int, error) {
			return 99, nil
		},
	}
	h := NewPaymentHandler(mock)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user_id", 1)
	body := `{"amount_cents":1000,"method":"card"}`
	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreatePayment(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "payment_id")
}
