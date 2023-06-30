package Request 

type TransactionCreateRequest struct {
	PartitionType string `json:"partition_type"`
	ShardingKey string `json:"sharding_key"`
	Database string `json:"database"`
}