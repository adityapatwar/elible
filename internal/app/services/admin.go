// internal/app/services/admin_service.go
package services

import (
	"errors"

	"elible/internal/app/models"
	"elible/internal/app/repository"
	"elible/internal/app/utils"

	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

func (s *AdminService) Create(admin *models.Admin) error {
	existingAdmin, err := s.repo.FindByUsername(admin.Username)
	if err != nil {
		return err
	}

	if existingAdmin != nil {
		return errors.New("admin already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.Password = string(hashedPassword)

	if err := s.repo.Create(admin); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) Login(username, password string) (*models.Admin, string, error) {
	admin, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid password")
	}

	tokenDetails, err := utils.CreateToken(username)
	if err != nil {
		return nil, "", err
	}

	token := &models.Token{
		AccessToken: tokenDetails.AccessToken,
		AccessUUID:  tokenDetails.AccessUUID,
		AtExpires:   tokenDetails.AtExpires,
	}

	// Store the token into the database
	if err = s.repo.SaveToken(token); err != nil {
		return nil, "", err
	}

	return admin, token.AccessToken, nil
}
