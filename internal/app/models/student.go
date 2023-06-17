package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name             string             `bson:"name,omitempty" json:"name,omitempty"`
	Email            string             `bson:"email,omitempty" json:"email,omitempty"`
	School           string             `bson:"school,omitempty" json:"school,omitempty"`
	SchoolID         primitive.ObjectID `bson:"school_id,omitempty" json:"school_id,omitempty"`
	Interest         string             `bson:"interest,omitempty" json:"interest,omitempty"`
	Gender           string             `bson:"gender,omitempty" json:"gender,omitempty"`
	Phone            string             `bson:"phone,omitempty" json:"phone,omitempty"`
	FinancialAbility string             `bson:"financial_ability,omitempty" json:"financial_ability,omitempty"`
	Progress         string             `bson:"progress,omitempty" json:"progress,omitempty"`
	DetailSiswaLink  string             `bson:"detailsiswa_link,omitempty" json:"detailsiswa_link,omitempty"`
	Image            string             `bson:"image,omitempty" json:"image,omitempty"`
	Category         string             `bson:"category,omitempty" json:"category,omitempty"`
	Birthdate        string             `bson:"birthdate,omitempty" json:"birthdate,omitempty"`
	TrackRecords     []TrackRecord      `bson:"track_records,omitempty" json:"track_records,omitempty"`
	TrackLobby       []TrackLobby       `bson:"track_lobby,omitempty" json:"track_lobby,omitempty"`
	IsActive         bool               `bson:"is_active,omitempty" json:"is_active"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt        time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type StudentFilter struct {
	Name             *string `bson:"name,omitempty" json:"name,omitempty"`
	School           *string `bson:"school,omitempty" json:"school,omitempty"`
	Interest         *string `bson:"interest,omitempty" json:"interest,omitempty"`
	Gender           *string `bson:"gender,omitempty" json:"gender,omitempty"`
	Phone            *string `bson:"phone,omitempty" json:"phone,omitempty"`
	Birthdate        *string `bson:"birthdate,omitempty" json:"birthdate,omitempty"`
	FinancialAbility *string `bson:"financial_ability,omitempty" json:"financial_ability,omitempty"`
	Progress         *string `bson:"progress,omitempty" json:"progress,omitempty"`
	Category         *string `bson:"category,omitempty" json:"category,omitempty"`
	IsActive         *bool   `bson:"is_active,omitempty" json:"is_active,omitempty"`
	Page             *int    `bson:"page,omitempty" json:"page,omitempty"`
	PageSize         *int    `bson:"pageSize,omitempty" json:"pageSize,omitempty"`
}

type PagedStudents struct {
	CurrentPage  int
	TotalRecords int64
	TotalPages   int
	Records      []Student
}
