// internal/app/handlers/routes_handler.go
package handlers

import (
	internal "elible/internal/app"
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
		adminGroup.POST("/create", adminHandler.RegisterAdmin)
		adminGroup.POST("/login", adminHandler.LoginAdmin)
		// adminGroup.POST("/logout", adminHandler.)
	}

	studentGroup := router.Group("/student")
	{
		studentGroup.POST("/create", studentHandler.RegisterStudent)
		studentGroup.POST("/all", studentHandler.GetAllStudents)
		studentGroup.POST("/id", studentHandler.GetIdStudents)
		studentGroup.POST("/delete", studentHandler.DeleteStudent)
		studentGroup.POST("/update", studentHandler.UpdateStudent)
		studentGroup.POST("/add-service", studentHandler.AddServiceToStudent)
		studentGroup.POST("/add-lobby", studentHandler.AddLobbyProgressToStudent)
	}
}
