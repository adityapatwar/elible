// internal/app/handlers/student_handler.go
package handlers

import (
	"log"
	"net/http"

	"elible/internal/app/models"
	"elible/internal/app/services"
	errors "elible/internal/pkg"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(service *services.StudentService) *StudentHandler {
	return &StudentHandler{
		service: service,
	}
}

type RequestWithID struct {
	ID string `json:"id" binding:"required"`
}

type UpdateStudentRequest struct {
	ID      string         `json:"id"`
	Student models.Student `json:"student"`
}

type AddServiceRequest struct {
	ID      string             `json:"id" `
	Service models.TrackRecord `json:"service"`
}

func (h *StudentHandler) RegisterStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.Create(&student); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusCreated, "Student created successfully", student)
	c.JSON(http.StatusCreated, response)
}

func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	students, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Students fetched successfully", students)
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.Delete(request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Student deleted successfully", nil)
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) DeactivateStudent(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}
	if err := h.service.Deactivate(request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Student deactivated successfully", nil)
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	var request UpdateStudentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	if err := h.service.Update(objectId.Hex(), &request.Student); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	log.Printf(request.Student.Name)

	response := errors.NewResponseData(http.StatusOK, "Student updated successfully", request.Student)
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) AddServiceToStudent(c *gin.Context) {
	var request AddServiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	if err := h.service.AddService(objectId.Hex(), &request.Service); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Service added to student successfully", request.Service)
	c.JSON(http.StatusOK, response)
}
