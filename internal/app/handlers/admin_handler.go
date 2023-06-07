package handlers

import (
	"net/http"
	"strings"

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
	if err := c.ShouldBind(&admin); err != nil {
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
	if err := c.ShouldBind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	admin, token, err := h.service.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewResponseError(http.StatusUnauthorized, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Login successful", gin.H{
		"admin": gin.H{
			"Username": admin.Username,
			"Email":    admin.Email,
			"FullName": admin.FullName,
		},
		"token": token,
	})

	c.JSON(http.StatusOK, response)
}

func (h *AdminHandler) GetProfileByToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, errors.NewResponseError(http.StatusUnauthorized, "No Authorization header provided"))
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		errors.WriteErrorResponse(c.Writer, 401, "Invalid Authorization header format")
		c.Abort()
		return
	}
	token := tokenParts[1]
	admin, err := h.service.GetAdminByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.NewResponseError(http.StatusUnauthorized, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Fetched admin profile successfully", gin.H{
		"admin": gin.H{
			"Username": admin.Username,
			"Email":    admin.Email,
			"FullName": admin.FullName,
		},
		"token": token,
	})

	c.JSON(http.StatusOK, response)
}
