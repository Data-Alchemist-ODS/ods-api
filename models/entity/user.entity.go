package entity

import (
	//local modules
	"time"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID 			primitive.ObjectID `gorm:"column:id;primaryKey" json:"userID" bson:"_id,omitempty"`
	Name		string `json:"name"`
	Email		string `json:"email" gorm:"unique"`
	Password	string `json:"password"`  
	CreatedAt 	time.Time `json:"CreatedAt"`
	UpdatedAt 	time.Time `json:"UpdatedAt"`
}

