package clients

import (
	"fmt"
)

// ManagementClient is used for connecting to the MarkLogic Management API.
type ManagementClient struct {
	*BasicClient
}

// NewManagementClient creates the Client struct used for managing databases, etc.
func NewManagementClient(connection *Connection) (*ManagementClient, error) {
	var client *ManagementClient
	if connection.Port <= 0 {
		connection.Port = 8002
	}
	base := fmt.Sprintf("http://%s:%v/manage/v2", connection.Host, connection.Port)
	basicClient, err := ClientBuilder(connection, base)
	if err == nil {
		client = &ManagementClient{basicClient}
	}
	return client, err
}
