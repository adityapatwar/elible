package services

import (
	"elible/internal/app/models"
	"elible/internal/app/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KnowledgeBaseService struct {
	repo *repository.KnowledgeBaseRepository
}

func NewKnowledgeBaseService(repo *repository.KnowledgeBaseRepository) *KnowledgeBaseService {
	return &KnowledgeBaseService{
		repo: repo,
	}
}

func (s *KnowledgeBaseService) CreateKnowledgeBase(kb *models.KnowledgeBase) error {
	return s.repo.CreateKnowledgeBase(kb)
}

func (s *KnowledgeBaseService) DeleteKnowledgeBase(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteKnowledgeBase(oid)
}

func (s *KnowledgeBaseService) UpdateKnowledgeBase(id string, kb *models.KnowledgeBase) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	kb.ID = oid
	return s.repo.UpdateKnowledgeBase(kb)
}

func (s *KnowledgeBaseService) ListKnowledgeBases() ([]models.KnowledgeBase, error) {
	return s.repo.ListKnowledgeBase()
}

func (s *KnowledgeBaseService) AddKnowledgeProgram(id string, program models.KnowledgeProgram) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.AddKnowledgeProgram(oid, program)
}

func (s *KnowledgeBaseService) RemoveKnowledgeProgram(id string, programName string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.RemoveKnowledgeProgram(oid, programName)
}

func (s *KnowledgeBaseService) UpdateKnowledgeProgram(id string,oldName string, program models.KnowledgeProgram) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.UpdateKnowledgeProgram(oid, oldName,program)
}

func (s *KnowledgeBaseService) ListKnowledgePrograms(id string) ([]models.KnowledgeProgram, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repo.ListKnowledgePrograms(oid)
}
