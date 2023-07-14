package request 

import "gorm.io/gorm"

type CreateUserData struct {
	FileData string `json:"data"`
}

type Data struct {
	gorm.Model
	Fields map[string]string `gorm:"-"`
}