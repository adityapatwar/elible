package internal

import (
	"elible/internal/app/repository"
	"elible/internal/app/services"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type Dependencies struct {
	AdminService *services.AdminService
	// Add your other services here
}

func InitializeDependencies(cfg *config.Config, mongoClient *mongo.Client) (*Dependencies, error) {
	adminRepo := repository.NewAdminRepository(cfg, mongoClient)
	adminService := services.NewAdminService(adminRepo)

	return &Dependencies{
		AdminService: adminService,
		// Add your other services here
	}, nil
}
