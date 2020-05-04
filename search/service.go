// Package search interacts with the MarkLogic search API
package search

import (
	clients "github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	"github.com/cchatfield/go-marklogic-go/util"
)

// Service is used for the search service
type Service struct {
	client *clients.Client
}

// NewService returns a new search.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Search with text value
func (s *Service) Search(text string, start int64, pageLength int64, transaction *util.Transaction, response handle.ResponseHandle) error {
	return Search(s.client, text, start, pageLength, transaction, response)
}

// StructuredSearch searches with a structured query
func (s *Service) StructuredSearch(query handle.Handle, start int64, pageLength int64, transaction *util.Transaction, response handle.ResponseHandle) error {
	return StructuredSearch(s.client, query, start, pageLength, transaction, response)
}

// StructuredSuggestions suggests query text based off of a structured query
func (s *Service) StructuredSuggestions(query handle.Handle, partialQ string, limit int64, options string, transaction *util.Transaction, response handle.ResponseHandle) error {
	return StructuredSuggestions(s.client, query, partialQ, limit, options, response)
}

// Delete documents that match specified collection, directory, etc.
func (s *Service) Delete(parameters map[string]string, transaction *util.Transaction, response handle.ResponseHandle) error {
	return Delete(s.client, parameters, transaction, response)
}
