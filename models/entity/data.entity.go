package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Documents []DataFields       `bson:"documents,omitempty"`
}

type DataFields struct {
	Fields map[string]interface{} `bson:"fields,omitempty"`
}

type DataResponse struct {
	ID     string            `json:"id"`
	Fields map[string]string `json:"fields"`
}