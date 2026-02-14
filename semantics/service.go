// Package semantics provides RDF triple store and semantic graph operations for
// working with linked data and semantic queries in MarkLogic. It supports SPARQL
// querying, graph CRUD operations, and semantic entity (Things) retrieval.
//
// Semantics Service enables:
//   - SPARQL query execution with JSON/XML results
//   - RDF graph management (create, read, update, delete)
//   - Semantic entity extraction and enrichment
//   - Named graph operations
//   - SPARQL Update statements
//   - Linked data and ontology support
//
// Example: Execute SPARQL query
//
//	query := `SELECT ?name WHERE { ?person foaf:name ?name }`
//	respHandle := semantics.SparqlResponseHandle{Format: handle.JSON}
//	client.Semantics().Sparql(query, nil, nil, &respHandle)
package semantics

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service provides methods for semantic and RDF graph operations.
// All operations are executed through the REST API with transparent
// format negotiation between JSON and XML.
type Service struct {
	client *clients.Client
}

// NewService creates and returns a new semantics.Service instance configured
// with the provided client for semantic and graph operations.
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Things retrieves semantic information and entity data for the specified IRIs
// (Internationalized Resource Identifiers). Returns semantic properties and
// related information for each IRI.
//
// Parameters:
//
//	iris: List of IRIs to retrieve semantic data for
//	response: ResponseHandle to populate with Things results
func (s *Service) Things(iris []string, response handle.ResponseHandle) error {
	return things(s.client, iris, response)
}
