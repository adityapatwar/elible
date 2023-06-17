// internal/app/handlers/routes_handler.go
package handlers

import (
	internal "elible/internal/app"
	"elible/internal/app/middleware"
	"elible/internal/app/services"
	"elible/internal/config"

	"github.com/gin-gonic/gin"
)

type RoutesHandler struct {
	adminHandler         *AdminHandler
	studentHandler       *StudentHandler
	universityHandler    *UniversityHandler
	studyProgramHandler  *StudyProgramHandler
	knowledgeBaseHandler *KnowledgeBaseHandler
}

func NewRoutesHandler(adminService *services.AdminService, studentService *services.StudentService, universityService *services.UniversityService, studyProgramService *services.StudyProgramService, knowledgeBaseService *services.KnowledgeBaseService) *RoutesHandler {
	return &RoutesHandler{
		adminHandler:         NewAdminHandler(adminService),
		studentHandler:       NewStudentHandler(studentService),
		universityHandler:    NewUniversityHandler(universityService),
		studyProgramHandler:  NewStudyProgramHandler(studyProgramService),
		knowledgeBaseHandler: NewKnowledgeBaseHandler(knowledgeBaseService),
	}
}

func Routes(router *gin.Engine, cfg *config.Config, deps *internal.Dependencies) {
	adminHandler := NewAdminHandler(deps.AdminService)
	studentHandler := NewStudentHandler(deps.StudentService)
	universityHandler := NewUniversityHandler(deps.UniversityService)
	studyProgramHandler := NewStudyProgramHandler(deps.StudyProgramService)
	knowledgeBaseHandler := NewKnowledgeBaseHandler(deps.KnowledgeBaseService)

	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/create", middleware.AdminMiddleware(cfg, deps.AdminService, true, false, adminHandler.RegisterAdmin))
		adminGroup.POST("/login", adminHandler.LoginAdmin)
		adminGroup.POST("/profil", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, adminHandler.GetProfileByToken))
	}

	studentGroup := router.Group("/student")
	{
		studentGroup.POST("/create", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.RegisterStudent))
		studentGroup.POST("/all", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.GetAllStudents))
		studentGroup.POST("/id", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.GetIdStudents))
		studentGroup.POST("/delete", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.DeleteStudent))
		studentGroup.POST("/deactivate", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.DeactivateStudent))
		studentGroup.POST("/update", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.UpdateStudent))
		studentGroup.POST("/add-service", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.AddServiceToStudent))
		studentGroup.POST("/add-lobby", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.AddLobbyProgressToStudent))
		studentGroup.POST("/upload", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.uploadImage))
		studentGroup.POST("/activated-all", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.ActivateStudnetAll))
		studentGroup.POST("/upload-excel", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.UploadAndImportDataStudent))
	}

	universityGroup := router.Group("/university")
	{
		universityGroup.POST("/create", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, universityHandler.CreateUniversity))
		universityGroup.POST("/update", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, universityHandler.UpdateUniversity))
		universityGroup.POST("/delete", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, universityHandler.DeleteUniversity))
		universityGroup.POST("/id", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, universityHandler.GetUniversity))
		universityGroup.POST("/all", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, universityHandler.GetUniversities))
	}

	studyProgramGroup := router.Group("/study-program")
	{
		studyProgramGroup.POST("/create", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studyProgramHandler.CreateStudyProgram))
		studyProgramGroup.POST("/update", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studyProgramHandler.UpdateStudyProgram))
		studyProgramGroup.POST("/delete", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studyProgramHandler.DeleteStudyProgram))
		studyProgramGroup.POST("/id", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studyProgramHandler.GetStudyProgram))
		studyProgramGroup.POST("/all", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studyProgramHandler.GetStudyPrograms))
		studyProgramGroup.POST("/upload", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studyProgramHandler.UploadAndImportData))
	}

	knowledgeBaseGroup := router.Group("/knowledge-base")
	{
		knowledgeBaseGroup.POST("/create", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.CreateKnowledgeBase))
		knowledgeBaseGroup.POST("/update", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.UpdateKnowledgeBase))
		knowledgeBaseGroup.POST("/delete", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.DeleteKnowledgeBase))
		knowledgeBaseGroup.POST("/all", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.ListKnowledgeBase))
	}

	knowledgeProgramsGroup := router.Group("/knowledge-programs")
	{
		knowledgeProgramsGroup.POST("/add", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.AddKnowledgeProgram))
		knowledgeProgramsGroup.POST("/update", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.UpdateKnowledgeProgram))
		knowledgeProgramsGroup.POST("/delete", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.RemoveKnowledgeProgram))
		knowledgeProgramsGroup.POST("/get", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, knowledgeBaseHandler.ListKnowledgePrograms))
	}

}
