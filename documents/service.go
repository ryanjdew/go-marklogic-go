// Package documents provides comprehensive CRUD (Create, Read, Update, Delete)
// operations for managing documents in MarkLogic. It supports single and bulk
// document operations, metadata management, format negotiation, and server-side
// transformations.
//
// Documents Service enables:
//   - Single and bulk document read/write/delete operations
//   - Metadata extraction and management (collections, permissions, properties)
//   - Format flexibility (JSON, XML, binary, text)
//   - Server-side transformations during read/write
//   - Multi-statement transactions for atomic updates
//   - Content negotiation for different document formats
//   - Efficient bulk operations using multipart/mixed encoding
//
// Example: Write JSON document with metadata
//
//	doc := &documents.DocumentDescription{
//	    URI: "/users/alice.json",
//	    Content: bytes.NewBufferString(`{"name":"Alice"}`),
//	    Metadata: &documents.Metadata{
//	        Collections: []string{"users", "active"},
//	    },
//	}
//	client.Documents().Write([]*documents.DocumentDescription{doc}, nil, nil, respHandle)
package documents

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// Service provides methods for document CRUD operations and metadata management.
// All read/write operations are executed through the HTTP REST API with
// transparent format negotiation based on handle types.
type Service struct {
	client *clients.Client
}

// NewService creates and returns a new documents.Service instance configured
// with the provided client. The Service coordinates all document operations
// through the REST API client.
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Client returns the underlying HTTP client associated with this Service.
// Can be used for advanced operations requiring direct client access.
func (s *Service) Client() *clients.Client {
	return s.client
}

// Read retrieves one or more documents from MarkLogic by URI.
// Supports reading content, metadata, or both. Metadata categories include:
// content, metadata, collections, permissions, properties, quality, metadata-values.
// Results can be single documents or multipart/mixed for multiple documents.
//
// Parameters:
//
//	uris: Document URIs to read (e.g., ["/doc1.json", "/doc2.json"])
//	categories: Metadata categories to include (nil for content only)
//	transform: Optional server-side transformation
//	transaction: Optional transaction for consistent reads
//	response: ResponseHandle for results; use MIXED format for multiple docs
func (s *Service) Read(uris []string, categories []string, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	return read(s.client, uris, categories, transform, transaction, response)
}

// Write creates or updates documents. Operations are executed concurrently per
// document for efficiency. Supports metadata (collections, permissions, properties)
// and server-side transformations. Returns successfully if all documents write succeed;
// returns error if any document fails.
//
// Parameters:
//
//	documents: DocumentDescription slice with URI, content, and metadata
//	transform: Optional server-side transformation
//	transaction: Optional transaction for atomic multi-document updates
//	response: ResponseHandle for results
func (s *Service) Write(documents []*DocumentDescription, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	return write(s.client, documents, transform, transaction, response)
}

// WriteSet writes multiple documents in a single multipart/mixed request.
// Provides efficient bulk writing with shared metadata. All documents use
// the same metadata instance and are sent as a single multipart message.
//
// Parameters:
//
//	documents: DocumentDescription slice to write
//	metadata: Metadata handle applied to all documents
//	transform: Optional server-side transformation
//	transaction: Optional transaction
//	response: ResponseHandle for results
func (s *Service) WriteSet(documents []*DocumentDescription, metadata handle.Handle, transform *util.Transform, transaction *util.Transaction, response handle.ResponseHandle) error {
	return writeSet(s.client, documents, metadata, transform, transaction, response)
}
