package Request 

type TransactionCreateRequest struct {
	Data []byte `json:"file_data" validate:"requiered"`
	PartitionType string `json:"partition_type" validate:"requiered"`
	ShardingKey string `json:"sharding_key" validate:"requiered"`
	Database string `json:"database" validate:"requiered"`
}