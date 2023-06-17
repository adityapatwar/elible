package models

type ImportResult struct {
	UniversityStats OperationStats `json:"university_stats"`
	ProgramStats    OperationStats `json:"program_stats"`
}

type ImportResultStudent struct {
	SchoolStats  OperationStats `json:"school_stats"`
	StudentStats OperationStats `json:"student_stats"`
}

type OperationStats struct {
	CreatedCount int   `json:"created_count"`
	UpdatedCount int   `json:"updated_count"`
	FailedCount  int   `json:"failed_count"`
	FailedRows   []int `json:"failed_rows"`
}
