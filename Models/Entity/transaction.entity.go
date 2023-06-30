package Entity

import (
	"time"

	"gorm.io/gorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	gorm.Model

	ID primitive.ObjectID `gorm:"primaryKey"` 
	PartitionType string `json:"partition_type"`
	ShardingKey string `json:"sharding_key"`
	Database string `json:"database"`
	Date time.Time  `json:"Datec"`
}
	