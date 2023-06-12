package handlers

import (
	"net/http"

	"elible/internal/app/models"
	"elible/internal/app/services"
	errors "elible/internal/pkg"

	"github.com/gin-gonic/gin"
)

type KnowledgeBaseHandler struct {
	service *services.KnowledgeBaseService
}

func NewKnowledgeBaseHandler(service *services.KnowledgeBaseService) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{
		service: service,
	}
}

func (h *KnowledgeBaseHandler) CreateKnowledgeBase(c *gin.Context) {
	var kb models.KnowledgeBase
	if err := c.ShouldBindJSON(&kb); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.CreateKnowledgeBase(&kb); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusCreated, "KnowledgeBase created successfully", kb)
	c.JSON(http.StatusCreated, response)
}

func (h *KnowledgeBaseHandler) DeleteKnowledgeBase(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.DeleteKnowledgeBase(request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KnowledgeBase deleted successfully"})
}

func (h *KnowledgeBaseHandler) UpdateKnowledgeBase(c *gin.Context) {
	var request UpdateKnowledgeBaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.UpdateKnowledgeBase(request.ID, &request.KnowledgeBase); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Knowledge Base updated successfully", request.KnowledgeBase)
	c.JSON(http.StatusOK, response)
}

func (h *KnowledgeBaseHandler) ListKnowledgeBase(c *gin.Context) {
	knowledgeBases, err := h.service.ListKnowledgeBases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Knowledge Base fetched successfully", knowledgeBases)
	c.JSON(http.StatusOK, response)
}

func (h *KnowledgeBaseHandler) AddKnowledgeProgram(c *gin.Context) {
	var request AddKnowledgeProgramRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.AddKnowledgeProgram(request.ID, request.KnowledgeProgram); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Knowledge Program added successfully", request)
	c.JSON(http.StatusOK, response)
}

func (h *KnowledgeBaseHandler) RemoveKnowledgeProgram(c *gin.Context) {
	var request RemoveKnowledgeProgramRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.RemoveKnowledgeProgram(request.ID, request.ProgramName); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Knowledge Program removed successfully", request)
	c.JSON(http.StatusOK, response)
}

func (h *KnowledgeBaseHandler) ListKnowledgePrograms(c *gin.Context) {
	var request RequestWithID
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	programs, err := h.service.ListKnowledgePrograms(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Knowledge Programs fetched successfully", programs)
	c.JSON(http.StatusOK, response)
}

func (h *KnowledgeBaseHandler) UpdateKnowledgeProgram(c *gin.Context) {

	var request UpdareKnowledgeProgramRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewResponseError(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.service.UpdateKnowledgeProgram(request.ID, request.OldName, request.KnowledgeProgram); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewResponseError(http.StatusInternalServerError, err.Error()))
		return
	}

	response := errors.NewResponseData(http.StatusOK, "Knowledge Program update successfully", request)
	c.JSON(http.StatusOK, response)
}
