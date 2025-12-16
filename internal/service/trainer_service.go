package service

import (
	"Gym-StrongCode/internal/models"
	"Gym-StrongCode/internal/repository"
)

type TrainerService struct {
	trainerRepo *repository.TrainerRepository
}

func NewTrainerService(trainerRepo *repository.TrainerRepository) *TrainerService {
	return &TrainerService{trainerRepo: trainerRepo}
}

func (s *TrainerService) Create(name, bio string) (*models.Trainer, error) {
	return s.trainerRepo.Create(name, bio)
}

func (s *TrainerService) List() ([]models.Trainer, error) {
	return s.trainerRepo.List()
}

func (s *TrainerService) Update(id int, name, bio string) error {
	return s.trainerRepo.Update(id, name, bio)
}

func (s *TrainerService) Delete(id int) error {
	return s.trainerRepo.Delete(id)
}