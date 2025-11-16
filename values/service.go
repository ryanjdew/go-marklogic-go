// Package values interacts with the MarkLogic values API for lexicon enumeration
package values

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the values service
type Service struct {
	client *clients.Client
}

// NewService returns a new values.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// ListValues retrieves lexicon values for a specified range index or field
// name: The name of the range index or field to query
// params: Optional parameters including start, pageLength, options, format
// response: The response handle for results
func (s *Service) ListValues(name string, params map[string]string, response handle.ResponseHandle) error {
	return listValues(s.client, name, params, response)
}

// QueryValues queries the values in a lexicon or range index with co-occurrence support
// name: The name of the range index or field to query
// params: Query parameters
// requestBody: Optional POST body for complex value queries
// response: The response handle for results
func (s *Service) QueryValues(name string, params map[string]string, requestBody handle.Handle, response handle.ResponseHandle) error {
	return queryValues(s.client, name, params, requestBody, response)
}
