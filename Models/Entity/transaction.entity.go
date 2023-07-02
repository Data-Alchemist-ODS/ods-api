package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	ID            primitive.ObjectID `gorm:"primaryKey" json:"transactionID"`
	PartitionType string             `json:"partition_type"`
	ShardingKey   string             `json:"sharding_key"`
	Database      string             `json:"database"`
	CreatedAt     time.Time          `json:"CreatedAt"`
	UpdatedAt     time.Time          `json:"UpdatedAt"`
	Data          [][]string         `json:"data"`
}
