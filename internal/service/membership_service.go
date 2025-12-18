package service

import (
	"database/sql"
	"fmt"

	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
)

type MembershipService struct {
	membershipRepo  *repository.MembershipRepository
	paymentRepo     *repository.PaymentRepository
	db              *sql.DB
	notificationSvc *NotificationService
}

func NewMembershipService(membershipRepo *repository.MembershipRepository, paymentRepo *repository.PaymentRepository, db *sql.DB, notificationSvc *NotificationService) *MembershipService {
	return &MembershipService{
		membershipRepo:  membershipRepo,
		paymentRepo:     paymentRepo,
		db:              db,
		notificationSvc: notificationSvc,
	}
}

func (s *MembershipService) List() ([]models.Membership, error) {
	return s.membershipRepo.GetAll()
}

func (s *MembershipService) Buy(userID, membershipID int, method string) (map[string]interface{}, error) {
	membership, err := s.membershipRepo.GetByID(membershipID)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	payment, err := s.paymentRepo.CreateForMembership(userID, membership.PriceCents, "KZT", method, "membership purchase", fmt.Sprintf("membership_%d", membershipID))
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Активируем подписку
	if err := s.membershipRepo.Activate(userID, membershipID, membership.DurationDays); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Уведомление
	s.notificationSvc.SendNotification("", "Подписка активирована", "Ваша подписка успешно куплена!")

	return map[string]interface{}{
		"payment":    payment,
		"membership": membership,
		"message":    "membership activated",
	}, nil
}

func (s *MembershipService) Create(name string, durationDays, priceCents int) (*models.Membership, error) {
	return s.membershipRepo.Create(name, durationDays, priceCents)
}

func (s *MembershipService) Update(id int, name string, durationDays, priceCents int) error {
	return s.membershipRepo.Update(id, name, durationDays, priceCents)
}

func (s *MembershipService) Delete(id int) error {
	return s.membershipRepo.Delete(id)
}
