package service

import (
	"Gym-StrongCode/internal/models"
	"Gym-StrongCode/internal/repository"
)

type GymService struct {
	gymRepo *repository.GymRepository
}

func NewGymService(gymRepo *repository.GymRepository) *GymService {
	return &GymService{gymRepo: gymRepo}
}

func (s *GymService) Create(name, address string) (*models.Gym, error) {
	return s.gymRepo.Create(name, address)
}

func (s *GymService) List() ([]models.Gym, error) {
	return s.gymRepo.List()
}

func (s *GymService) Update(id int, name, address string) error {
	return s.gymRepo.Update(id, name, address)
}

func (s *GymService) Delete(id int) error {
	return s.gymRepo.Delete(id)
}