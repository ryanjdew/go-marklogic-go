package clients

import (
	"fmt"
)

// AdminClient is used for connecting to the MarkLogic Management API.
type AdminClient struct {
	*BasicClient
}

// NewAdminClient creates the Client struct used for managing admin features, etc.
func NewAdminClient(connection *Connection) (*AdminClient, error) {
	var client *AdminClient
	if connection.Port <= 0 {
		connection.Port = 8001
	}
	base := fmt.Sprintf("http://%s:%v/admin/v1", connection.Host, connection.Port)
	basicClient, err := ClientBuilder(connection, base)
	if err == nil {
		client = &AdminClient{basicClient}
	}
	return client, err
}
