package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name             string             `bson:"name,omitempty" json:"name,omitempty"`
	School           string             `bson:"school,omitempty" json:"school,omitempty"`
	Interest         string             `bson:"interest,omitempty" json:"interest,omitempty"`
	Gender           string             `bson:"gender,omitempty" json:"gender,omitempty"`
	Phone            string             `bson:"phone,omitempty" json:"phone,omitempty"`
	FinancialAbility string             `bson:"financial_ability,omitempty" json:"financial_ability,omitempty"`
	Progress         string             `bson:"progress,omitempty" json:"progress,omitempty"`
	DetailSiswaLink  string             `bson:"detailsiswa_link,omitempty" json:"detailsiswa_link,omitempty"`
	Image            string             `bson:"image,omitempty" json:"image,omitempty"`
	Category         string             `bson:"category,omitempty" json:"category,omitempty"`
	TrackRecords     []TrackRecord      `bson:"track_records,omitempty" json:"track_records,omitempty"`
	IsActive         bool               `bson:"is_active,omitempty" json:"is_active,omitempty"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt        time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
