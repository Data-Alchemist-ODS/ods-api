package request

type TiDBConnectionRequest struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	ServerName string `json:"server_name"`
	Database   string `json:"database"`
}
