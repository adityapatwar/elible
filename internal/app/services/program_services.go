package services

import (
	"elible/internal/app/models"
	"elible/internal/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyProgramService struct {
	repo *repository.StudyProgramRepository
}

func NewStudyProgramService(repo *repository.StudyProgramRepository) *StudyProgramService {
	return &StudyProgramService{
		repo: repo,
	}
}

func (s *StudyProgramService) CreateStudyProgram(sp models.StudyProgram, kbYear string, kpName string) (primitive.ObjectID, error) {
	return s.repo.CreateStudyProgram(sp, kbYear, kpName)
}

func (s *StudyProgramService) UpdateStudyProgram(id string, sp models.StudyProgram) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.UpdateStudyProgram(oid, sp)
}

func (s *StudyProgramService) DeleteStudyProgram(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteStudyProgram(oid)
}

func (s *StudyProgramService) GetStudyProgram(id primitive.ObjectID) (models.StudyProgram, error) {
	return s.repo.GetStudyProgram(id)
}

func (s *StudyProgramService) GetStudyPrograms(kbYear string, kpName string) ([]models.StudyProgram, error) {
	return s.repo.GetStudyPrograms(kbYear, kpName)
}
