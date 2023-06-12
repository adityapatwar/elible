package services

import (
	"errors"

	"elible/internal/app/models"
	"elible/internal/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{
		repo: repo,
	}
}

func (s *StudentService) Create(student *models.Student) error {
	existingStudent, err := s.repo.FindByUsername(student.Name)
	if err != nil {
		return err
	}

	if existingStudent != nil {
		return errors.New("student already exists")
	}

	if err := s.repo.Create(student); err != nil {
		return err
	}

	return nil
}

func (s *StudentService) GetAll(filter *models.StudentFilter) ([]*models.Student, error) {
	return s.repo.GetAll(filter)
}

func (s *StudentService) GetByID(studentID string) (*models.Student, error) {
	objectId, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return nil,err
	}
	return s.repo.GetByID(objectId)
}

func (s *StudentService) Delete(studentID string) error {
	objectId, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return err
	}
	return s.repo.Delete(objectId)
}

func (s *StudentService) Deactivate(studentID string) error {
	objectId, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return err
	}
	return s.repo.Deactivate(objectId)
}

func (s *StudentService) Update(studentID string, student *models.Student) error {
	objectId, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return err
	}
	return s.repo.Update(objectId, student)
}

func (s *StudentService) AddService(studentID string, service *models.TrackRecord) error {
	objectId, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return err
	}
	return s.repo.AddService(objectId, service)
}

func (s *StudentService) AddLobby(studentID string, lobby *models.Student) error {
	objectId, err := primitive.ObjectIDFromHex(studentID)
	if err != nil {
		return err
	}
	return s.repo.AddLobby(objectId, lobby)
}
