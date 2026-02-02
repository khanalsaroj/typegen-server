package connection

type DatabaseConnectionsRequest struct {
	ConnectionId uint   `json:"connectionId"`
	Name         string `json:"name"`
	DbType       string `json:"dbType"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"databaseName"`
	SchemaName   string `json:"schemaName"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}
