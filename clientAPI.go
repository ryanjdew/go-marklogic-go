package gomarklogicgo

import (
	"github.com/ryanjdew/go-marklogic-go/alert"
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/config"
	datamovement "github.com/ryanjdew/go-marklogic-go/datamovement"
	"github.com/ryanjdew/go-marklogic-go/dataservices"
	"github.com/ryanjdew/go-marklogic-go/documents"
	rowsManagement "github.com/ryanjdew/go-marklogic-go/rows-management"
	search "github.com/ryanjdew/go-marklogic-go/search"
	"github.com/ryanjdew/go-marklogic-go/semantics"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// Authentication options
const (
	BasicAuth  = clients.BasicAuth
	DigestAuth = clients.DigestAuth
	None       = clients.None
)

// Client is used for connecting to the MarkLogic REST API.
type Client clients.Client

// Connection is used for defining the connection to the MarkLogic REST API.
type Connection clients.Connection

// NewClient creates the Client struct used for searching, etc.
func NewClient(host string, port int64, username string, password string, authType int) (*Client, error) {
	return New(&Connection{Host: host, Port: port, Username: username, Password: password, AuthenticationType: authType})
}

// New creates the Client struct used for searching, etc.
func New(config *Connection) (*Client, error) {
	client, err := clients.NewClient(convertToSubConnection(config))
	return convertToClient(client), err
}

// Alerting service
func (c *Client) Alerting() *alert.Service {
	return alert.NewService(convertToSubClient(c))
}

// Config service
func (c *Client) Config() *config.Service {
	return config.NewService(convertToSubClient(c))
}

// DataMovement service
func (c *Client) DataMovement() *datamovement.Service {
	return datamovement.NewService(convertToSubClient(c))
}

// DataServices service
func (c *Client) DataServices() *dataservices.Service {
	return dataservices.NewService(convertToSubClient(c))
}

// Documents service
func (c *Client) Documents() *documents.Service {
	return documents.NewService(convertToSubClient(c))
}

// RowsManagement service
func (c *Client) RowsManagement() *rowsManagement.Service {
	return rowsManagement.NewService(convertToSubClient(c))
}

// Search service
func (c *Client) Search() *search.Service {
	return search.NewService(convertToSubClient(c))
}

// Semantics service
func (c *Client) Semantics() *semantics.Service {
	return semantics.NewService(convertToSubClient(c))
}

// NewTransaction returns a new transaction struct
func (c *Client) NewTransaction() *util.Transaction {
	return &util.Transaction{}
}

func convertToSubClient(c *Client) *clients.Client {
	converted := clients.Client(*c)
	return &converted
}

func convertToClient(c *clients.Client) *Client {
	converted := Client(*c)
	return &converted
}

func convertToSubConnection(c *Connection) *clients.Connection {
	converted := clients.Connection(*c)
	return &converted
}
