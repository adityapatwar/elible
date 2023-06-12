package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type University struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Alias       string             `bson:"alias,omitempty" json:"alias,omitempty"`
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`
	Website     string             `bson:"website,omitempty" json:"website,omitempty"`
	Logo        string             `bson:"logo,omitempty" json:"logo,omitempty"`
	Image       string             `bson:"image,omitempty" json:"image,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	Contact     Contact            `bson:"contact,omitempty" json:"contact,omitempty"`
	SocialMedia []SocialMedia      `bson:"social_media,omitempty" json:"social_media,omitempty"`
}

type Contact struct {
	Email string `bson:"email,omitempty" json:"email,omitempty"`
	Phone string `bson:"phone,omitempty" json:"phone,omitempty"`
	Fax   string `bson:"fax,omitempty" json:"fax,omitempty"`
}

type SocialMedia struct {
	Platform string `bson:"platform,omitempty" json:"platform,omitempty"`
	Link     string `bson:"link,omitempty" json:"link,omitempty"`
}
