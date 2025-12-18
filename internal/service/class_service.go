package service

import (
	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
)

type ClassService struct {
	classRepo   *repository.ClassRepository
	trainerRepo *repository.TrainerRepository
	gymRepo     *repository.GymRepository
}

func NewClassService(classRepo *repository.ClassRepository, trainerRepo *repository.TrainerRepository, gymRepo *repository.GymRepository) *ClassService {
	return &ClassService{classRepo: classRepo, trainerRepo: trainerRepo, gymRepo: gymRepo}
}

func (s *ClassService) Create(c *models.Class) (*models.Class, error) {
	// Валидация trainer_id и gym_id
	if c.TrainerID != 0 {
		if _, err := s.trainerRepo.GetByID(c.TrainerID); err != nil {
			return nil, err
		}
	}
	if _, err := s.gymRepo.GetByID(c.GymID); err != nil {
		return nil, err
	}

	return s.classRepo.Create(c)
}

func (s *ClassService) List() ([]models.Class, error) {
	return s.classRepo.List()
}

func (s *ClassService) Update(id int, c *models.Class) error {
	return s.classRepo.Update(id, c)
}

func (s *ClassService) Delete(id int) error {
	return s.classRepo.Delete(id)
}
