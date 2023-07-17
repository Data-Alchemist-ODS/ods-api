package request

type TransactionCreateRequest struct {
	PartitionType string `json:"partitionType"`
	ShardingKey string `json:"shardingKey"`
	Database string `json:"database"`
	FileData string `json:"data"`
}
