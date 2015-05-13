package goMarklogicGo

import (
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	search "github.com/ryanjdew/go-marklogic-go/search"
)

// Authentication options
const (
	BasicAuth  = clients.BasicAuth
	DigestAuth = clients.DigestAuth
	None       = clients.None
)

// Client is used for connecting to the MarkLogic REST API.
type Client clients.Client

// NewClient creates the Client struct used for searching, etc.
func NewClient(host string, port int64, username string, password string, authType int) (*Client, error) {
	client, err := clients.NewClient(host, port, username, password, authType)
	return convertToClient(client), err
}

// Search with text value
func (c *Client) Search(text string, start int64, pageLength int64, response handle.Handle) error {
	return search.Search(convertToSubClient(c), text, start, pageLength, response)
}

// StructuredSearch searches with a structured query
func (c *Client) StructuredSearch(query handle.Handle, start int64, pageLength int64, response handle.Handle) error {
	return search.StructuredSearch(convertToSubClient(c), query, start, pageLength, response)
}

// StructuredSuggestions suggests query text based off of a structured query
func (c *Client) StructuredSuggestions(query handle.Handle, partialQ string, limit int64, options string, response handle.Handle) error {
	return search.StructuredSuggestions(convertToSubClient(c), query, partialQ, limit, options, response)
}
func convertToSubClient(c *Client) *clients.Client {
	converted := clients.Client(*c)
	return &converted
}

func convertToClient(c *clients.Client) *Client {
	converted := Client(*c)
	return &converted
}
