package handlers

import (
	"net/http"

	"elible/internal/app/models"
	"elible/internal/app/services"
	errors "elible/internal/pkg"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *services.AdminService
}

func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{
		service: service,
	}
}

func (h *AdminHandler) RegisterAdmin(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.Create(&admin); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusCreated, "Admin created successfully", admin)
	c.JSON(http.StatusCreated, response)
}

func (h *AdminHandler) LoginAdmin(c *gin.Context) {
	var credentials models.Credentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	admin, token, err := h.service.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewResponseError(http.StatusUnauthorized, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Login successful", gin.H{
		"admin": admin,
		"token": token,
	})

	c.JSON(http.StatusOK, response)
}
