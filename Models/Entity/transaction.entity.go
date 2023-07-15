package entity

import (
	//local modules
	"time"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID            primitive.ObjectID `gorm:"column:id;primaryKey" json:"transactionID" bson:"_id,omitempty"`
	PartitionType string             `json:"partition_type"`
	ShardingKey   string             `json:"sharding_key"`
	Database      string             `json:"database"`
	CreatedAt     time.Time          `json:"CreatedAt"`
	UpdatedAt     time.Time          `json:"UpdatedAt"`
	Data          string         	 `json:"data"`
}
