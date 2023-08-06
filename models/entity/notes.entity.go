package entity

import (
	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notes struct {
	ID			primitive.ObjectID `gorm:"column:id;primaryKey" json:"NotesID" bson:"_id,omitempty"`
	Date 		string 			   `json:"date"`
	Description string 			   `json:"description"`
	Story 		string 			   `json:"story"`
}