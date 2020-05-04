// Package documents provides a way to read and write documents
package documents

import (
	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	"github.com/cchatfield/go-marklogic-go/util"
)

// Service is used for the documents service
type Service struct {
	client *clients.Client
}

// NewService returns a new documents.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Client associated with the documents service
func (s *Service) Client() *clients.Client {
	return s.client
}

// Read documents
func (s *Service) Read(uris []string, categories []string, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	return read(s.client, uris, categories, transform, transaction, response)
}

// Write documents according to the DocumentDescription slice passed
func (s *Service) Write(documents []*DocumentDescription, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	return write(s.client, documents, transform, transaction, response)
}

// WriteSet documents according to the DocumentDescription slice passed
func (s *Service) WriteSet(documents []*DocumentDescription, metadata handle.Handle, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	return writeSet(s.client, documents, metadata, transform, transaction, response)
}
