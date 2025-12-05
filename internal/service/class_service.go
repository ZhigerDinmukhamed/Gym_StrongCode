package service

import (
	"Gym_StrongCode/internal/cache"
	"Gym_StrongCode/internal/models"
	"Gym_StrongCode/internal/repository"
	"time"
)

type ClassService struct {
	classRepo   *repository.ClassRepository
	trainerRepo *repository.TrainerRepository
	cache       *cache.Cache
}

func NewClassService(classRepo *repository.ClassRepository, trainerRepo *repository.TrainerRepository) *ClassService {
	return &ClassService{
		classRepo:   classRepo,
		trainerRepo: trainerRepo,
		cache:       cache.NewCache(10 * time.Second), // TTL = 10 сек
	}
}

func (s *ClassService) GetClasses() ([]models.ClassInfo, error) {
	if cached, ok := s.cache.Get("classes"); ok {
		return cached.([]models.ClassInfo), nil
	}

	classes, err := s.classRepo.GetAll()
	if err != nil {
		return nil, err
	}

	s.cache.Set("classes", classes)
	return classes, nil
}

func (s *ClassService) CreateClass(class *models.Class) error {
	s.cache.Delete("classes")
	return s.classRepo.Create(class)
}
