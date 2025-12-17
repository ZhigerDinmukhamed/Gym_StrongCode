package service

import (
	"Gym-StrongCode/internal/models"
	"Gym-StrongCode/internal/repository"
	"fmt"
)

type PaymentService struct {
	paymentRepo *repository.PaymentRepository
}

func NewPaymentService(paymentRepo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo}
}

func (s *PaymentService) Create(userID, amountCents int, method, description, referenceID string) (*models.Payment, error) {
	if amountCents <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	// Валидация метода оплаты
	validMethods := []string{"card", "cash", "bank_transfer", "qr_code"}
	valid := false
	for _, m := range validMethods {
		if method == m {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("invalid payment method: %s", method)
	}

	return s.paymentRepo.CreateStandalone(userID, amountCents, "KZT", method, "completed")
}

func (s *PaymentService) GetByUser(userID int, status string) ([]models.Payment, error) {
	return s.paymentRepo.GetByUser(userID, status)
}

func (s *PaymentService) GetAll() ([]models.Payment, error) {
	return s.paymentRepo.GetAll()
}

func (s *PaymentService) ListAll() ([]models.Payment, error) {
	return s.GetAll()
}

func (s *PaymentService) CreateStandalone(userID, amountCents int, method string) (*models.Payment, error) {
	return s.paymentRepo.CreateStandalone(userID, amountCents, "KZT", method, "completed")
}
