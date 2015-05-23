package search

import (
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
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
func (s *Service) Search(text string, start int64, pageLength int64, response handle.ResponseHandle) error {
	return Search(s.client, text, start, pageLength, response)
}

// StructuredSearch searches with a structured query
func (s *Service) StructuredSearch(query handle.Handle, start int64, pageLength int64, response handle.ResponseHandle) error {
	return StructuredSearch(s.client, query, start, pageLength, response)
}

// StructuredSuggestions suggests query text based off of a structured query
func (s *Service) StructuredSuggestions(query handle.Handle, partialQ string, limit int64, options string, response handle.ResponseHandle) error {
	return StructuredSuggestions(s.client, query, partialQ, limit, options, response)
}
