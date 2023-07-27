package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Data map[string]interface{} 


type Document struct {
	ID		 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Filename string				`json:"filename,omitempty" bson:"filename,omitempty"`
	FileData []Data				`json:"data,omitempty" bson:"data,omitempty"`
}
