// Package indexes manages flexible range indexes and field indexes in MarkLogic
package indexes

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the indexes service
type Service struct {
	client *clients.Client
}

// NewService returns a new indexes.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// ListIndexes retrieves all configured indexes
func (s *Service) ListIndexes(params map[string]string, response handle.ResponseHandle) error {
	return listIndexes(s.client, params, response)
}

// GetIndex retrieves a specific index configuration by name
func (s *Service) GetIndex(indexName string, params map[string]string, response handle.ResponseHandle) error {
	return getIndex(s.client, indexName, params, response)
}

// CreateIndex creates a new range or field index
func (s *Service) CreateIndex(requestBody handle.Handle, response handle.ResponseHandle) error {
	return createIndex(s.client, requestBody, response)
}

// UpdateIndex updates an existing index configuration
func (s *Service) UpdateIndex(indexName string, requestBody handle.Handle, response handle.ResponseHandle) error {
	return updateIndex(s.client, indexName, requestBody, response)
}

// DeleteIndex removes an index by name
func (s *Service) DeleteIndex(indexName string, params map[string]string, response handle.ResponseHandle) error {
	return deleteIndex(s.client, indexName, params, response)
}
