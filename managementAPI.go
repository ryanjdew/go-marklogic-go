package goMarklogicGo

import (
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	management "github.com/ryanjdew/go-marklogic-go/management"
)

// ManagementClient is used for connecting to the MarkLogic Management API.
type ManagementClient clients.ManagementClient

// NewManagementClient creates the Client struct used for managing databases, etc.
func NewManagementClient(host string, username string, password string, authType int) (*ManagementClient, error) {
	client, err := clients.NewManagementClient(&clients.Connection{Host: host, Username: username, Password: password, AuthenticationType: authType})
	return convertToManageClient(client), err
}

// SetDatabaseProperties sets the database properties
func (mc *ManagementClient) SetDatabaseProperties(databaseName string, propertiesHandle handle.ResponseHandle) error {
	return management.SetDatabaseProperties(convertToSubManageClient(mc), databaseName, propertiesHandle)
}

// GetDatabaseProperties sets the database properties
func (mc *ManagementClient) GetDatabaseProperties(databaseName string, propertiesHandle handle.ResponseHandle) error {
	return management.GetDatabaseProperties(convertToSubManageClient(mc), databaseName, propertiesHandle)
}

func convertToSubManageClient(mc *ManagementClient) *clients.ManagementClient {
	converted := clients.ManagementClient(*mc)
	return &converted
}

func convertToManageClient(mc *clients.ManagementClient) *ManagementClient {
	converted := ManagementClient(*mc)
	return &converted
}
