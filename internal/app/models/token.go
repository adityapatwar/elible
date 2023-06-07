package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Token struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AccessToken string             `bson:"accessToken,omitempty"`
	AccessUUID  string             `bson:"accessUUID,omitempty"`
	AtExpires   int64              `bson:"atExpires,omitempty"`
}
