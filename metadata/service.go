// Package metadata manages metadata extraction and validation in MarkLogic
package metadata

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the metadata service
type Service struct {
	client *clients.Client
}

// NewService returns a new metadata.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// ExtractMetadata extracts metadata from documents based on configured extractors
func (s *Service) ExtractMetadata(uris []string, options map[string]string, response handle.ResponseHandle) error {
	return extractMetadata(s.client, uris, options, response)
}

// ExtractMetadataFromQuery extracts metadata from documents matching a query
func (s *Service) ExtractMetadataFromQuery(query handle.Handle, options map[string]string, response handle.ResponseHandle) error {
	return extractMetadataFromQuery(s.client, query, options, response)
}

// ValidateDocuments validates documents against configured rules
func (s *Service) ValidateDocuments(uris []string, validationRules handle.Handle, response handle.ResponseHandle) error {
	return validateDocuments(s.client, uris, validationRules, response)
}

// ValidateQuery validates documents matching a query against rules
func (s *Service) ValidateQuery(query handle.Handle, validationRules handle.Handle, response handle.ResponseHandle) error {
	return validateQuery(s.client, query, validationRules, response)
}

// GetValidationRules retrieves configured validation rules
func (s *Service) GetValidationRules(response handle.ResponseHandle) error {
	return getValidationRules(s.client, response)
}

// SetValidationRules sets validation rules for the database
func (s *Service) SetValidationRules(rules handle.Handle, response handle.ResponseHandle) error {
	return setValidationRules(s.client, rules, response)
}

// ExtractMetadataFromURI extracts metadata from a single document
func (s *Service) ExtractMetadataFromURI(uri string, options map[string]string, response handle.ResponseHandle) error {
	return extractMetadataFromURI(s.client, uri, options, response)
}

// ValidateURI validates a single document against rules
func (s *Service) ValidateURI(uri string, validationRules handle.Handle, response handle.ResponseHandle) error {
	return validateURI(s.client, uri, validationRules, response)
}
