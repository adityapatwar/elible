package handlers

import "elible/internal/app/models"

type RequestWithID struct {
	ID string `json:"id" binding:"required"`
}

type UpdateStudentRequest struct {
	ID      string         `json:"id" binding:"required"`
	Student models.Student `json:"student"`
}

type AddServiceRequest struct {
	ID      string             `json:"id" binding:"required"`
	Service models.TrackRecord `json:"service" binding:"required"`
}

type AddLobbyRequest struct {
	ID    string         `json:"id" binding:"required"`
	Lobby models.Student `json:"lobby" binding:"required"`
}

type UpdateKnowledgeBaseRequest struct {
	ID            string               `json:"id" binding:"required"`
	KnowledgeBase models.KnowledgeBase `json:"knowledgeBase"`
}

type AddKnowledgeProgramRequest struct {
	ID               string                  `json:"id" binding:"required"`
	KnowledgeProgram models.KnowledgeProgram `json:"KnowledgeProgram"`
}

type UpdareKnowledgeProgramRequest struct {
	ID               string                  `json:"id" binding:"required"`
	OldName          string                  `json:"oldName" binding:"required"`
	KnowledgeProgram models.KnowledgeProgram `json:"KnowledgeProgram"`
}

type RemoveKnowledgeProgramRequest struct {
	ID          string `json:"id" binding:"required"`
	ProgramName string `json:"program_name" binding:"required"`
}

type AddProgramRequest struct {
	KbYear       string              `json:"kbYear" binding:"required"`
	KpName       string              `json:"kpName" binding:"required"`
	StudyProgram models.StudyProgram `json:"study_program"`
}

type UpdateProgramRequest struct {
	ID           string              `json:"id" binding:"required"`
	StudyProgram models.StudyProgram `json:"study_program"`
}

type GetProgramRequest struct {
	KbYear   string `json:"kbYear" binding:"required"`
	KpName   string `json:"kpName" binding:"required"`
	Page     string `json:"page" binding:"required"`
	PageSize string `json:"pageSize" binding:"required"`
}

type UpdateUniversityRequest struct {
	ID         string            `json:"id" binding:"required"`
	University models.University `json:"university"`
}
