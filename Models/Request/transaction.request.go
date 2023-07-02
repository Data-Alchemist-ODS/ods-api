package Request

type TransactionCreateRequest struct {
    PartitionType  string   `json:"partitionType"`
    ShardingKey    string   `json:"shardingKey"`
    Database       string   `json:"database"`
    FileData       []byte   `json:"fileData"`
}
