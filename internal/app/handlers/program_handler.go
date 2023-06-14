package handlers

import (
	"net/http"
	"os"
	"path"

	"elible/internal/app/models"
	"elible/internal/app/services"
	errors "elible/internal/pkg"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudyProgramHandler struct {
	service *services.StudyProgramService
}

func NewStudyProgramHandler(service *services.StudyProgramService) *StudyProgramHandler {
	return &StudyProgramHandler{
		service: service,
	}
}

func (h *StudyProgramHandler) CreateStudyProgram(c *gin.Context) {

	var request AddProgramRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	spID, err := h.service.CreateStudyProgram(request.StudyProgram, request.KbYear, request.KpName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Study Program created successfully", gin.H{"id": spID})
	c.JSON(http.StatusOK, response)
}

func (h *StudyProgramHandler) UpdateStudyProgram(c *gin.Context) {
	var request UpdateProgramRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.UpdateStudyProgram(request.ID, request.StudyProgram); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Study Program updated successfully", request.StudyProgram)
	c.JSON(http.StatusOK, response)
}

func (h *StudyProgramHandler) DeleteStudyProgram(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.DeleteStudyProgram(request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Study Program deleted successfully", gin.H{"id": request.ID})
	c.JSON(http.StatusOK, response)
}

func (h *StudyProgramHandler) GetStudyProgram(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	oid, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	sp, err := h.service.GetStudyProgram(oid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Fetched study program successfully", sp)
	c.JSON(http.StatusOK, response)
}

func (h *StudyProgramHandler) GetStudyPrograms(c *gin.Context) {
	var filter models.GetStudyProgramsFilter
	if err := c.ShouldBind(&filter); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	programs, err := h.service.GetStudyPrograms(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Fetched study programs successfully", programs)
	c.JSON(http.StatusOK, response)
}

func (h *StudyProgramHandler) UploadAndImportData(c *gin.Context) {
	// I assume that you use 'file' as the field name in your form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// You can change the destination path as you wish
	dst := path.Join("./tempFile", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Always delete the file after we're done
	defer os.Remove(dst)

	// Get knowledgeBaseYear and knowledgeProgramName from the form data
	knowledgeBaseYear := c.PostForm("knowledgeBaseYear")
	knowledgeProgramName := c.PostForm("knowledgeProgramName")

	// Import data from the uploaded Excel file
	err = h.service.ImportDataFromExcelStudyPrograms(knowledgeBaseYear, knowledgeProgramName, dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	// Respond to the client
	response := errors.NewResponseData(http.StatusOK, "Data imported successfully", nil)
	c.JSON(http.StatusOK, response)
}
