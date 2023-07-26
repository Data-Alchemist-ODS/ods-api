package entity

import (
	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID            primitive.ObjectID `gorm:"column:id;primaryKey" json:"transactionID" bson:"_id,omitempty"`
	PartitionType string             `json:"partition_type" validate:"required"`
	ShardingKey   string             `json:"sharding_key" validate:"required"`
	Database      string             `json:"database" validate:"required"`
	Data          string             `json:"data" validate:"required"`
}
