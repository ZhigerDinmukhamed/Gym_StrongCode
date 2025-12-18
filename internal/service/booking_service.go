package service

import (
	"fmt"

	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
)

type BookingService struct {
	bookingRepo     *repository.BookingRepository
	classRepo       *repository.ClassRepository
	membershipRepo  *repository.MembershipRepository
	notificationSvc *NotificationService
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	classRepo *repository.ClassRepository,
	membershipRepo *repository.MembershipRepository,
	notificationSvc *NotificationService,
) *BookingService {
	return &BookingService{
		bookingRepo:     bookingRepo,
		classRepo:       classRepo,
		membershipRepo:  membershipRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *BookingService) Create(userID, classID int, userEmail string) error {
	// ... проверки (класс существует, есть места, активная подписка и т.д.)

	// Создаём бронирование
	_, err := s.bookingRepo.Create(userID, classID)
	if err != nil {
		return err
	}

	// Уведомление по email
	class, err := s.classRepo.GetByID(classID)
	if err != nil {
		return err
	}

	body := fmt.Sprintf(`
		<h2>Бронирование подтверждено!</h2>
		<p>Вы успешно забронировали занятие: <strong>%s</strong></p>
		<p>Дата и время: %s</p>
		<p>Спасибо за выбор StrongCode!</p>
	`, class.Title, class.StartTime)

	s.notificationSvc.SendNotification(userEmail, "Бронирование занятия", body)

	return nil
}

func (s *BookingService) ListUser(userID int) ([]models.Booking, error) {
	return s.bookingRepo.GetByUser(userID)
}

func (s *BookingService) ListAll() ([]models.Booking, error) {
	return s.bookingRepo.ListAll()
}

func (s *BookingService) Cancel(bookingID, userID int) error {
	return s.bookingRepo.Cancel(bookingID, userID)
}
