package service

import (
	"Gym-StrongCode/internal/models"
	"Gym-StrongCode/internal/repository"
	"Gym-StrongCode/internal/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) List() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) Delete(id int) error {
	return s.userRepo.Delete(id)
}
