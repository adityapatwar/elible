package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KnowledgeBase struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Year      string             `bson:"year,omitempty" json:"year,omitempty"`
	Programs  []KnowledgeProgram `bson:"programs,omitempty" json:"programs,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type KnowledgeProgram struct {
	Name          string               `bson:"name,omitempty" json:"name,omitempty"`
	DisplayName   string               `bson:"display_name,omitempty" json:"display_name,omitempty"`
	StudyPrograms []primitive.ObjectID `bson:"study_programs,omitempty" json:"study_programs,omitempty"`
}

type StudyProgram struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string             `bson:"name,omitempty" json:"name,omitempty"`
	ProgramDetails Program            `bson:"program_details,omitempty" json:"program_details,omitempty"`
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Program struct {
	University    primitive.ObjectID `bson:"university,omitempty" json:"university,omitempty"`
	Program       string             `bson:"program,omitempty" json:"program,omitempty"`
	ProgramType   string             `bson:"program_type,omitempty" json:"program_type,omitempty"`
	UKT           string             `bson:"ukt,omitempty" json:"ukt,omitempty"`
	SPI           string             `bson:"spi,omitempty" json:"spi,omitempty"`
	Capacity      string             `bson:"capacity,omitempty" json:"capacity,omitempty"`
	IsPacketC     bool               `bson:"is_packet_c,omitempty" json:"is_packet_c,omitempty"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
	Advantages    string             `bson:"advantages,omitempty" json:"advantages,omitempty"`
	Disadvantages string             `bson:"disadvantages,omitempty" json:"disadvantages,omitempty"`
	Articles      []Article          `bson:"articles,omitempty" json:"articles,omitempty"`
	Requirements  []string           `bson:"requirements,omitempty" json:"requirements,omitempty"`
	Registration  RegistrationDates  `bson:"registration,omitempty" json:"registration,omitempty"`
	Exam          ExamDates          `bson:"exam,omitempty" json:"exam,omitempty"`
	Announcement  time.Time          `bson:"announcement,omitempty" json:"announcement,omitempty"`
}

type RegistrationDates struct {
	Start time.Time `bson:"start,omitempty" json:"start,omitempty"`
	End   time.Time `bson:"end,omitempty" json:"end,omitempty"`
}

type ExamDates struct {
	Start time.Time `bson:"start,omitempty" json:"start,omitempty"`
	End   time.Time `bson:"end,omitempty" json:"end,omitempty"`
}

type Article struct {
	Title   string `bson:"title,omitempty" json:"title,omitempty"`
	Content string `bson:"content,omitempty" json:"content,omitempty"`
}

type StudyProgramWithUniversity struct {
	StudyProgram StudyProgram `bson:"study_program,omitempty" json:"study_program,omitempty"`
	University   University   `bson:"university,omitempty" json:"university,omitempty"`
}

type PagedStudyPrograms struct {
	CurrentPage  int
	TotalRecords int64
	TotalPages   int
	Records      []StudyProgramWithUniversity
}

type GetStudyProgramsFilter struct {
	KbYear      string `bson:"kbYear,omitempty" json:"kbYear,omitempty"`
	KpName      string `bson:"kpName,omitempty" json:"kpName,omitempty"`
	Page        int    `bson:"page,omitempty" json:"page,omitempty"`
	PageSize    int    `bson:"pageSize,omitempty" json:"pageSize,omitempty"`
	SearchQuery string `bson:"searchQuery,omitempty" json:"searchQuery,omitempty"`
	ProgramType string `bson:"programType,omitempty" programType:"school,omitempty"`
	Program     string `bson:"program,omitempty" json:"program,omitempty"`
}

type ImportResult struct {
	UniversityStats OperationStats `json:"university_stats"`
	ProgramStats    OperationStats `json:"program_stats"`
}

type OperationStats struct {
	CreatedCount int   `json:"created_count"`
	UpdatedCount int   `json:"updated_count"`
	FailedCount  int   `json:"failed_count"`
	FailedRows   []int `json:"failed_rows"`
}
