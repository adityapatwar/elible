package internal

import (
	"elible/internal/app/repository"
	"elible/internal/app/services"
	"elible/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
)



// universityHandler := NewUniversityHandler(deps.UniversityService)
// studyProgramHandler := NewStudyProgramHandler(deps.StudyProgramService)
// knowledgeBaseHandler := NewKnowledgeBaseHandler(deps.KnowledgeBaseService)
type Dependencies struct {
	AdminService   *services.AdminService
	StudentService *services.StudentService
	UniversityService *services.UniversityService
	StudyProgramService *services.StudyProgramService
	KnowledgeBaseService *services.KnowledgeBaseService
	// Add your other services here
}

func InitializeDependencies(cfg *config.Config, mongoClient *mongo.Client) (*Dependencies, error) {
	adminRepo := repository.NewAdminRepository(cfg, mongoClient)
	studentRepo := repository.NewStudentRepository(cfg, mongoClient)
	univRepo := repository.NewUniversityRepository(cfg, mongoClient)
	programtRepo := repository.NewStudyProgramRepository(cfg, mongoClient)
	knowRepo := repository.NewKnowledgeBaseRepository(cfg, mongoClient)

	adminService := services.NewAdminService(adminRepo)
	studentService := services.NewStudentService(studentRepo)
	univService := services.NewUniversityService(univRepo)
	programService := services.NewStudyProgramService(programtRepo)
	knowService := services.NewKnowledgeBaseService(knowRepo)

	return &Dependencies{
		AdminService:   adminService,
		StudentService: studentService,
		UniversityService: univService,
		StudyProgramService: programService,
		KnowledgeBaseService: knowService,
	}, nil
}
