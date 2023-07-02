package Request

type TransactionCreateRequest struct {
    PartitionType  string   `json:"partitionType"`
    ShardingKey    string   `json:"shardingKey"`
    Database       string   `json:"database"`
    FileContentType string   `json:"fileContentType"`
    FileData       []byte   `json:"fileData"`
}
