package handlers

import (
	"elible/internal/app/models"
	"elible/internal/app/services"
	errors "elible/internal/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UniversityHandler struct {
	service *services.UniversityService
}

func NewUniversityHandler(service *services.UniversityService) *UniversityHandler {
	return &UniversityHandler{service: service}
}

func (h *UniversityHandler) CreateUniversity(c *gin.Context) {
	var request models.University
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, "Invalid request data"))
		return
	}

	id, err := h.service.CreateUniversity(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "Failed to create university"))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "University created successfully", gin.H{"id": id})
	c.JSON(http.StatusOK, response)
}

func (h *UniversityHandler) UpdateUniversity(c *gin.Context) {
	var request UpdateUniversityRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	if err := h.service.UpdateUniversity(objectId, request.University); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "Failed to update university"))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "University updated successfully", request.University)
	c.JSON(http.StatusOK, response)
}

func (h *UniversityHandler) DeleteUniversity(c *gin.Context) {
	var request RequestWithID

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	if err := h.service.DeleteUniversity(objectId); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "Failed to delete university"))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "University deleted successfully", nil)
	c.JSON(http.StatusOK, response)
}

func (h *UniversityHandler) GetUniversity(c *gin.Context) {
	var request RequestWithID

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	u, err := h.service.GetUniversity(objectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "Failed to retrieve university"))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "University fetched successfully", u)
	c.JSON(http.StatusOK, response)
}

func (h *UniversityHandler) GetUniversities(c *gin.Context) {
	us, err := h.service.GetUniversities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "Failed to retrieve universities"))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Universities fetched successfully", us)
	c.JSON(http.StatusOK, response)
}

func (h *UniversityHandler) GetUniversityByName(c *gin.Context) {
	name := c.Param("name")

	u, err := h.service.GetUniversityByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "Failed to retrieve university by name"))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "University fetched successfully", u)
	c.JSON(http.StatusOK, response)
}
