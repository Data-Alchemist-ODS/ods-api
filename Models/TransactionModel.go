package Models

import (
	"time"

	"gorm.io/gorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	gorm.Model

	ID primitive.ObjectID `gorm:"primaryKey"`
	PartitionType string
	ShardingKey string
	Database string
	Date time.Time
}
	