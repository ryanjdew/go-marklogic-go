package documents

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// Service is used for the documents service
type Service struct {
	client *clients.Client
}

// NewService returns a new search.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Read documents
func (s *Service) Read(uris []string, categories []string, transform *util.Transform, response handle.ResponseHandle) error {
	return read(s.client, uris, categories, transform, response)
}

// Write documents according to the DocumentDescription slice passed
func (s *Service) Write(documents []DocumentDescription, transform *util.Transform, response handle.ResponseHandle) error {
	return write(s.client, documents, transform, response)
}
