package request

type TransactionCreateRequest struct {
	PartitionType string `json:"partitionType"`
	ShardingKey string `json:"shardingKey"`
	Database string `json:"database"`
	FileData string `json:"data"`
}

type Data struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	PartitionType string `json:"partition_type"`
	ShardingKey   string `json:"sharding_key"`
	Database      string `json:"database"`
	FileData      string `json:"file_data"`
}
