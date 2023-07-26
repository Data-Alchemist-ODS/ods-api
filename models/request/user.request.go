package request

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserCreateRequest struct {
	Name 		string `json:"name"`
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

type UserLoginRequest struct {
	ID			primitive.ObjectID  `bson:"_id,omitempty"`
	Email 		string 				`json:"email"`
	Password 	string 				`json:"password"`
}