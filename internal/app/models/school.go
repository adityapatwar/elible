package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type School struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`
	Province    string             `bson:"province,omitempty" json:"province,omitempty"`
	City        string             `bson:"city,omitempty" json:"city,omitempty"`
	SchoolLogo  string             `bson:"school_logo,omitempty" json:"school_logo,omitempty"`
	SchoolImage string             `bson:"school_image,omitempty" json:"school_image,omitempty"`
	Phone       string             `bson:"phone,omitempty" json:"phone,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
