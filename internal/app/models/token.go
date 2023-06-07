// internal/app/models/admin.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Token struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	AccessToken string             `bson:"access_token,omitempty"`
	AccessUUID  string             `bson:"access_uuid,omitempty"`
	AtExpires   int64              `bson:"at_expires,omitempty"`
}
