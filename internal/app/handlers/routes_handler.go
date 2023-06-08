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
	adminHandler   *AdminHandler
	studentHandler *StudentHandler
}

func NewRoutesHandler(adminService *services.AdminService, studentService *services.StudentService) *RoutesHandler {
	return &RoutesHandler{
		adminHandler:   NewAdminHandler(adminService),
		studentHandler: NewStudentHandler(studentService),
	}
}

func Routes(router *gin.Engine, cfg *config.Config, deps *internal.Dependencies) {
	adminHandler := NewAdminHandler(deps.AdminService)
	studentHandler := NewStudentHandler(deps.StudentService)

	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/create", middleware.AdminMiddleware(cfg, deps.AdminService, true, false, adminHandler.RegisterAdmin))
		adminGroup.POST("/login", adminHandler.LoginAdmin)
		adminGroup.POST("/profil", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, adminHandler.GetProfileByToken))
		// adminGroup.POST("/logout", adminHandler.)
	}

	studentGroup := router.Group("/student")
	{
		studentGroup.POST("/create", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.RegisterStudent))
		studentGroup.POST("/all", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.GetAllStudents))
		studentGroup.POST("/id", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.GetIdStudents))
		studentGroup.POST("/delete", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.DeleteStudent))
		studentGroup.POST("/update", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.UpdateStudent))
		studentGroup.POST("/add-service", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.AddServiceToStudent))
		studentGroup.POST("/add-lobby", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.AddLobbyProgressToStudent))
		studentGroup.POST("/upload", middleware.AdminMiddleware(cfg, deps.AdminService, false, true, studentHandler.uploadImage))
	}
}
