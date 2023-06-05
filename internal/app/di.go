package internal

import (
	"elible/internal/app/repository"
	"elible/internal/app/services"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
)

type Dependencies struct {
	AdminService   *services.AdminService
	StudentService *services.StudentService
	// Add your other services here
}

func InitializeDependencies(cfg *config.Config, mongoClient *mongo.Client) (*Dependencies, error) {
	adminRepo := repository.NewAdminRepository(cfg, mongoClient)
	studentRepo := repository.NewStudentRepository(cfg, mongoClient)

	adminService := services.NewAdminService(adminRepo)
	studentService := services.NewStudentService(studentRepo)

	return &Dependencies{
		AdminService:   adminService,
		StudentService: studentService,
		// Add your other services here
	}, nil
}
