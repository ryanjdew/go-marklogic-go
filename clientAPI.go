package gomarklogicgo

import (
	"github.com/ryanjdew/go-marklogic-go/alert"
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/config"
	datamovement "github.com/ryanjdew/go-marklogic-go/datamovement"
	"github.com/ryanjdew/go-marklogic-go/dataservices"
	"github.com/ryanjdew/go-marklogic-go/documents"
	"github.com/ryanjdew/go-marklogic-go/eval"
	"github.com/ryanjdew/go-marklogic-go/indexes"
	"github.com/ryanjdew/go-marklogic-go/metadata"
	resources "github.com/ryanjdew/go-marklogic-go/resources"
	rowsManagement "github.com/ryanjdew/go-marklogic-go/rows-management"
	search "github.com/ryanjdew/go-marklogic-go/search"
	"github.com/ryanjdew/go-marklogic-go/semantics"
	temporal "github.com/ryanjdew/go-marklogic-go/temporal"
	transactions "github.com/ryanjdew/go-marklogic-go/transactions"
	"github.com/ryanjdew/go-marklogic-go/util"
	"github.com/ryanjdew/go-marklogic-go/values"
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

// Eval service
func (c *Client) Eval() *eval.Service {
	return eval.NewService(convertToSubClient(c))
}

// Indexes service
func (c *Client) Indexes() *indexes.Service {
	return indexes.NewService(convertToSubClient(c))
}

// Metadata service
func (c *Client) Metadata() *metadata.Service {
	return metadata.NewService(convertToSubClient(c))
}

// Resources service
func (c *Client) Resources() *resources.Service {
	return resources.NewService(convertToSubClient(c))
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

// Temporal service
func (c *Client) Temporal() *temporal.Service {
	return temporal.NewService(convertToSubClient(c))
}

// Transactions service
func (c *Client) Transactions() *transactions.Service {
	return transactions.NewService(convertToSubClient(c))
}

// Values service
func (c *Client) Values() *values.Service {
	return values.NewService(convertToSubClient(c))
}

// NewTransaction returns a new transaction struct
func (c *Client) NewTransaction() *util.Transaction {
	return util.NewTransaction(convertToSubClient(c))
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
