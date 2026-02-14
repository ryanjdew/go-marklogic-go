// Package search provides full-text and structured search capabilities for querying
// documents in MarkLogic. It supports complex search queries, result faceting,
// query suggestions, and multi-statement transactions for atomic search operations.
//
// Search Service enables:
//   - Full-text search using simple text queries
//   - Structured queries using the MarkLogic Query DSL (XML format)
//   - Query suggestions and autocomplete functionality
//   - Result pagination with start position and page length
//   - Faceted search results for result refinement
//   - Bulk deletion of documents matching search criteria
//   - Multi-statement transactions for consistent search operations
//
// Example: Full-text search with results
//
//	respHandle := search.ResponseHandle{Format: handle.JSON}
//	client.Search().Search("shakespeare", 1, 10, nil, &respHandle)
//	results := respHandle.Get().(*search.Response)
package search

import (
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// Service provides methods for executing search operations against MarkLogic.
// All operations return results through ResponseHandle which supports both
// typed (JSON/XML) and raw response access.
type Service struct {
	client *clients.Client
}

// NewService creates and returns a new search.Service instance configured
// with the provided client. The Service wraps low-level HTTP client operations
// and provides high-level search method interfaces.
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Search performs a full-text search using a simple text query string.
// Returns paginated results with scores and URIs. Results can be accessed
// both as typed objects (respHandle.Get()) or raw response text.
//
// Parameters:
//
//	text: Search query string (e.g., "shakespeare plays")
//	start: Result page number (1-based)
//	pageLength: Number of results per page
//	transaction: Optional transaction for consistent results; if nil, uses single query
//	response: ResponseHandle to populate with results
func (s *Service) Search(text string, start int64, pageLength int64, transaction *util.Transaction, response handle.ResponseHandle) error {
	return Search(s.client, text, start, pageLength, transaction, response)
}

// StructuredSearch executes a structured query against MarkLogic documents.
// Structured queries allow complex, type-safe query composition using the
// MarkLogic Query DSL. Supports nested queries, faceting, and advanced search options.
//
// Parameters:
//
//	query: Handle containing structured query definition (XML format typical)
//	start: Result page number (1-based)
//	pageLength: Number of results per page
//	transaction: Optional transaction for atomic operations
//	response: ResponseHandle to populate with results
func (s *Service) StructuredSearch(query handle.Handle, start int64, pageLength int64, transaction *util.Transaction, response handle.ResponseHandle) error {
	return StructuredSearch(s.client, query, start, pageLength, transaction, response)
}

// Suggest provides query suggestions based on a partial query string.
// Useful for implementing autocomplete and query refinement interfaces.
// Suggestions are based on server-side query options and field configuration.
//
// Parameters:
//
//	partialQ: Partial query text to suggest completions for
//	limit: Maximum number of suggestions to return
//	options: Query options name (e.g., "default")
//	response: SuggestionsResponseHandle to populate with suggestions
func (s *Service) Suggest(partialQ string, limit int64, options string, response handle.ResponseHandle) error {
	return Suggest(s.client, partialQ, limit, options, response)
}

// StructuredSuggestions provides suggestions within the context of a structured query.
// Enables advanced autocomplete with query refinement based on current search context.
//
// Parameters:
//
//	query: Base structured query for context
//	partialQ: Partial query text to suggest completions for
//	limit: Maximum number of suggestions
//	options: Query options name
//	transaction: Optional transaction
//	response: SuggestionsResponseHandle to populate with suggestions
func (s *Service) StructuredSuggestions(query handle.Handle, partialQ string, limit int64, options string, transaction *util.Transaction, response handle.ResponseHandle) error {
	return StructuredSuggestions(s.client, query, partialQ, limit, options, response)
}

// Delete removes documents matching the specified search criteria.
// Useful for bulk deletion operations. Can be used within transactions for
// consistent multi-step operations. Parameters are passed as query string pairs.
//
// Parameters:
//
//	parameters: Query parameters (e.g., map{"q": "status:archived"})
//	transaction: Optional transaction for atomic deletion
//	response: ResponseHandle to populate with deletion confirmation
func (s *Service) Delete(parameters map[string]string, transaction *util.Transaction, response handle.ResponseHandle) error {
	return Delete(s.client, parameters, transaction, response)
}
