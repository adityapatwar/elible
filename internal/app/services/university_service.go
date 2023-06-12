package services

import (
	"elible/internal/app/models"
	"elible/internal/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UniversityService struct {
	repo *repository.UniversityRepository
}

func NewUniversityService(repo *repository.UniversityRepository) *UniversityService {
	return &UniversityService{repo: repo}
}

func (s *UniversityService) CreateUniversity(u models.University) (primitive.ObjectID, error) {
	return s.repo.CreateUniversity(u)
}

func (s *UniversityService) UpdateUniversity(id primitive.ObjectID, u models.University) error {
	return s.repo.UpdateUniversity(id, u)
}

func (s *UniversityService) DeleteUniversity(id primitive.ObjectID) error {
	return s.repo.DeleteUniversity(id)
}

func (s *UniversityService) GetUniversity(id primitive.ObjectID) (models.University, error) {
	return s.repo.GetUniversity(id)
}

func (s *UniversityService) GetUniversityByName(name string) (models.University, error) {
	return s.repo.GetUniversityByName(name)
}

func (s *UniversityService) GetUniversities() ([]models.University, error) {
	return s.repo.GetUniversities()
}
