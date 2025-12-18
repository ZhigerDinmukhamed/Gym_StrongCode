package service

import (
	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) List() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) Delete(id int) error {
	return s.userRepo.Delete(id)
}
