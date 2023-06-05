// internal/app/handlers/routes_handler.go
package handlers

import (
	internal "elible/internal/app"
	"elible/internal/app/services"
	"elible/internal/config"

	"github.com/gin-gonic/gin"
)

type RoutesHandler struct {
	adminHandler *AdminHandler
}

func NewRoutesHandler(adminService *services.AdminService) *RoutesHandler {
	return &RoutesHandler{
		adminHandler: NewAdminHandler(adminService),
	}
}

func Routes(router *gin.Engine, cfg *config.Config, deps *internal.Dependencies) {
	adminHandler := NewAdminHandler(deps.AdminService)

	adminGroup := router.Group("/admin")
	{
		adminGroup.POST("/create", adminHandler.RegisterAdmin)
		adminGroup.POST("/login", adminHandler.LoginAdmin)
		// adminGroup.POST("/logout", adminHandler.)
	}

}
