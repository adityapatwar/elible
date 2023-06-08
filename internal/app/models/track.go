package models

import "time"

type TrackRecord struct {
	ServiceName string    `bson:"service_name,omitempty" json:"service_name,omitempty"`
	ServiceDate string    `bson:"service_date,omitempty" json:"service_date,omitempty"`
	ServiceCost string    `bson:"service_cost,omitempty" json:"service_cost,omitempty"`
	Status      string    `bson:"status,omitempty" json:"status,omitempty"`
	CreatedAt   time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type TrackLobby struct {
	Progress  string    `bson:"progress,omitempty" json:"progress,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type ServiceFilter struct {
	Name        *string `json:"name,omitempty"`
	ServiceName *string `json:"service_name,omitempty"`
	ServiceDate *string `json:"service_date,omitempty"`
	Status      *string `json:"status,omitempty"`
}
