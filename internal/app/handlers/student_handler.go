// internal/app/handlers/student_handler.go
package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"elible/internal/app/models"
	"elible/internal/app/services"
	utils "elible/internal/app/utils"
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

func (h *StudentHandler) RegisterStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBind(&student); err != nil {
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
	var filter models.StudentFilter
	if err := c.ShouldBind(&filter); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}
	students, err := h.service.GetAll(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Students fetched successfully", students)
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) GetIdStudents(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	students, err := h.service.GetByID(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Students fetched successfully", students)
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBind(&request); err != nil {
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
	if err := c.ShouldBind(&request); err != nil {
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

	if err := c.ShouldBind(&request); err != nil {
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
	response := errors.NewResponseData(http.StatusOK, "Student updated successfully", request.Student)
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) AddServiceToStudent(c *gin.Context) {
	var request AddServiceRequest
	if err := c.ShouldBind(&request); err != nil {
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

func (h *StudentHandler) AddLobbyProgressToStudent(c *gin.Context) {
	var request AddLobbyRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	if err := h.service.AddLobby(objectId.Hex(), &request.Lobby); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Lobby Proggress added to student successfully", map[string]interface{}{
		"Progress": request.Lobby.Progress,
	})
	c.JSON(http.StatusOK, response)
}

func (h *StudentHandler) uploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, "File is not present in the form data"))
		return
	}

	if valid := utils.IsImage(file); !valid {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, "File is not an image"))
		return
	}

	dir := os.Getenv("IMAGE_DIR")
	if dir == "" {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "IMAGE_DIR environment variable is not set"))
		return
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
			return
		}
	}

	newFilename := fmt.Sprintf("%d_%s", rand.Int(), time.Now().Format("20060102"))
	ext := filepath.Ext(file.Filename)
	newFilenameWithExt := newFilename + ext
	newFilenameWithExtDomain := path.Join("images", newFilename+ext)
	dst := filepath.Join(dir, newFilenameWithExt)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	if err := os.Chmod(dst, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	domain := os.Getenv("WEB_DOMAIN")
	if domain == "" {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, "WEB_DOMAIN environment variable is not set"))
		return
	}

	fullPath := path.Join(domain, newFilenameWithExtDomain)

	response := errors.NewResponseData(http.StatusCreated, "Upload Image Success ", fullPath)
	c.JSON(http.StatusCreated, response)
}

func (h *StudentHandler) ActivateStudnetAll(c *gin.Context) {

	if err := h.service.ActivateStudnetAll(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Activate All Student Successfully", nil)
	c.JSON(http.StatusOK, response)
}
