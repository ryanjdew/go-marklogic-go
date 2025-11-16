# Go MarkLogic Client Library

A comprehensive Go client library for interacting with MarkLogic's REST APIs. This library provides type-safe, idiomatic Go interfaces to MarkLogic's document, search, semantic, and data management capabilities.

## Status

[![GoDoc](https://godoc.org/github.com/ryanjdew/go-marklogic-go?status.svg)](https://godoc.org/github.com/ryanjdew/go-marklogic-go) [![Build Status](https://cloud.drone.io/api/badges/ryanjdew/go-marklogic-go/status.svg)](https://cloud.drone.io/ryanjdew/go-marklogic-go)

## Features

- **Full-text & Structured Search** - Query text, use structured queries, get suggestions, and retrieve faceted results
- **Document CRUD** - Read, write, update, and delete documents with metadata management
- **Semantic/Graph Operations** - Work with RDF triples, SPARQL queries, and semantic data
- **Lexicon Value Enumeration** - Query distinct values from range indexes
- **Server-Side Code Execution** - Execute XQuery and JavaScript on the server with external variables
- **Query Suggestions** - Get suggestions for queries and implement autocomplete
- **Multi-Statement Transactions** - Atomically execute multiple document or search operations
- **User-Defined Resources** - Call custom REST extensions (GET/POST/PUT/DELETE)
- **Flexible Indexes** - Create and manage range and field indexes for optimized performance
- **Bulk Operations** - Efficiently batch read/write multiple documents
- **Query Management** - Install and manage query options, transforms, and extensions
- **Format Flexibility** - Seamless JSON/XML serialization with format negotiation
- **Authentication** - Basic Auth, Digest Auth, or no authentication support

## Installation

```bash
go get github.com/ryanjdew/go-marklogic-go
```

## Quick Start

### Connecting to MarkLogic

```go
package main

import (
	"fmt"
	"log"
	
	marklogic "github.com/ryanjdew/go-marklogic-go"
)

func main() {
	// Create a client with Digest authentication
	client, err := marklogic.NewClient(
		"localhost",           // host
		8050,                  // port (REST API port)
		"admin",               // username
		"admin",               // password
		marklogic.DigestAuth,  // authentication type
	)
	if err != nil {
		log.Fatal(err)
	}
	
	// Use client to interact with MarkLogic
	_ = client
}
```

### Authentication Options

The library supports three authentication methods:

```go
// Basic Authentication (username:password in HTTP header)
client, err := marklogic.NewClient(host, port, user, pass, marklogic.BasicAuth)

// Digest Authentication (challenge/response)
client, err := marklogic.NewClient(host, port, user, pass, marklogic.DigestAuth)

// No Authentication
client, err := marklogic.NewClient(host, port, "", "", marklogic.None)
```

## Usage Examples

### Searching

#### Full-Text Search

```go
import (
	"fmt"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	search "github.com/ryanjdew/go-marklogic-go/search"
)

// Perform a text search
respHandle := search.ResponseHandle{Format: handle.JSON}
err := client.Search().Search("Shakespeare", 1, 10, nil, &respHandle)
if err != nil {
	log.Fatal(err)
}

// Access results
results := respHandle.Deserialized().(*search.Response)
fmt.Printf("Total results: %d\n", results.Total)
for _, result := range results.Results {
	fmt.Printf("  - %s (score: %d)\n", result.URI, result.Score)
}

// Or get the raw JSON response
fmt.Println(respHandle.Serialized())
```

#### Structured Search

```go
// Build a structured query
query := search.Query{
	Queries: []any{
		search.TermQuery{
			Terms: []string{"Shakespeare", "tragedy"},
		},
	},
}

// Serialize to XML or JSON
qh := search.QueryHandle{Format: handle.XML}
qh.Serialize(query)

// Execute the search
respHandle := search.ResponseHandle{Format: handle.JSON}
err := client.Search().StructuredSearch(&qh, 1, 10, nil, &respHandle)
```

#### Search Suggestions

```go
// Get suggestions for partial query
sugRespHandle := search.SuggestionsResponseHandle{Format: handle.JSON}
err := client.Search().StructuredSuggestions(&qh, "shake", 10, "", nil, &sugRespHandle)

suggestions := sugRespHandle.Deserialized().(*search.Suggestions)
for _, suggestion := range suggestions.Suggestions {
	fmt.Println(suggestion)
}
```

#### Delete by Query

```go
// Delete all documents matching a query
params := map[string]string{
	"collection": "old-data",
}
err := client.Search().Delete(params, nil, nil)
```

### Document Operations

#### Reading Documents

```go
import (
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	documents "github.com/ryanjdew/go-marklogic-go/documents"
)

// Read single or multiple documents by URI
uris := []string{"/docs/doc1.xml", "/docs/doc2.json"}
respHandle := documents.ResponseHandle{Format: handle.JSON}

// Get content only
err := client.Documents().Read(uris, []string{"content"}, nil, nil, &respHandle)

// Get content + metadata
err := client.Documents().Read(
	uris,
	[]string{"content", "metadata", "permissions", "collections"},
	nil, nil, &respHandle,
)
```

#### Writing Documents

```go
// Create documents to write
docs := []*documents.DocumentDescription{
	{
		URI:    "/docs/doc1.xml",
		Content: bytes.NewBufferString("<root>Content</root>"),
		Format: handle.XML,
		Metadata: &documents.Metadata{
			Collections: []string{"my-collection"},
			Quality: 1,
		},
	},
	{
		URI:    "/docs/doc2.json",
		Content: bytes.NewBufferString(`{"key": "value"}`),
		Format: handle.JSON,
		Metadata: &documents.Metadata{
			Collections: []string{"my-collection"},
		},
	},
}

// Write documents
respHandle := documents.ResponseHandle{}
err := client.Documents().Write(docs, nil, nil, &respHandle)
```

#### Updating Documents

```go
// Use PATCH to update (or PUT to replace)
doc := &documents.DocumentDescription{
	URI:    "/docs/doc1.xml",
	Content: bytes.NewBufferString("<root>Updated</root>"),
	Format: handle.XML,
}

// PATCH for partial update
err := client.Documents().Write(
	[]*documents.DocumentDescription{doc},
	nil, nil, &respHandle,
)
```

### Semantic Operations

#### Query RDF Triples

```go
import (
	semantics "github.com/ryanjdew/go-marklogic-go/semantics"
)

// Get things (entities) from the triple index
iris := []string{
	"http://example.org/person/john",
	"http://example.org/person/jane",
}

respHandle := semantics.ResponseHandle{Format: handle.JSON}
err := client.Semantics().Things(iris, &respHandle)

things := respHandle.Deserialized()
fmt.Println(things)
```

#### SPARQL Queries

```go
// Execute a SPARQL query
sparqlQuery := `
SELECT ?subject ?predicate ?object
WHERE {
  ?subject ?predicate ?object
}
LIMIT 10
`

respHandle := semantics.ResponseHandle{Format: handle.JSON}
err := client.Semantics().Sparql(sparqlQuery, &respHandle)
```

### Lexicon Value Enumeration

#### List Values from a Lexicon

```go
import (
	"github.com/ryanjdew/go-marklogic-go/values"
)

// List all distinct values from a range index
respHandle := values.ResponseHandle{Format: handle.JSON}
err := client.Values().ListValues(
	"status",  // Lexicon name (e.g., range index name)
	nil,       // Optional parameters (pageLength, start, etc.)
	&respHandle,
)

values := respHandle.Deserialized()
fmt.Println(values)
```

#### Query Values with Constraints

```go
// Query values with structured constraints
queryBody := `<values-query xmlns="http://marklogic.com/appservices/search">
  <range name="status" type="xs:string">
    <range-operator>EQ</range-operator>
    <value>active</value>
  </range>
</values-query>`

queryHandle := handle.RawHandle{Format: handle.XML}
queryHandle.Write([]byte(queryBody))

respHandle := values.ResponseHandle{Format: handle.XML}
err := client.Values().QueryValues(
	"status",
	nil,
	&queryHandle,
	&respHandle,
)
```

### Server-Side Code Execution

#### Evaluate XQuery

```go
import (
	"github.com/ryanjdew/go-marklogic-go/eval"
)

// Execute XQuery code on the server
xqueryCode := `
declare variable $x as xs:integer external;
$x * 2 + 1
`

respHandle := eval.ResponseHandle{Format: handle.JSON}
err := client.Eval().EvalXQuery(xqueryCode, nil, &respHandle)

// Result is available as JSON
result := respHandle.Deserialized()
fmt.Println(result)
```

#### Evaluate XQuery with External Variables

```go
// Pass external variables to XQuery
params := map[string]string{
	"x": `{"type": "xs:integer", "value": 21}`,
	"y": `{"type": "xs:string", "value": "hello"}`,
}

respHandle := eval.ResponseHandle{Format: handle.JSON}
err := client.Eval().EvalXQuery(xqueryCode, params, &respHandle)
```

#### Evaluate JavaScript

```go
// Execute JavaScript code on the server
jsCode := `
var x = 42;
var result = { value: x, doubled: x * 2 };
result;
`

respHandle := eval.ResponseHandle{Format: handle.JSON}
err := client.Eval().EvalJavaScript(jsCode, nil, &respHandle)
```

### Query Suggestions

#### Simple Suggestions

```go
// Get suggestions based on a partial query string
respHandle := search.SuggestionsResponseHandle{Format: handle.JSON}
err := client.Search().Suggest(
	"mark",  // Partial query string
	10,      // Limit number of suggestions
	"",      // Optional query options name
	&respHandle,
)

suggestions := respHandle.Get()
for _, suggestion := range suggestions.Suggestions {
	fmt.Println(suggestion)
}
```

#### Structured Suggestions

```go
// Get suggestions based on a structured query
query := search.Query{
	Queries: []any{
		search.TermQuery{Terms: []string{"active"}},
	},
}

qh := search.QueryHandle{Format: handle.XML}
qh.Serialize(query)

respHandle := search.SuggestionsResponseHandle{Format: handle.JSON}
err := client.Search().StructuredSuggestions(
	&qh,
	"par",  // Partial term to suggest
	5,      // Limit
	"",     // Query options
	nil,    // Transaction
	&respHandle,
)
```

### Multi-Statement Transactions

#### Begin a Transaction

```go
import (
	"github.com/ryanjdew/go-marklogic-go/transactions"
)

// Start a new transaction
respHandle := transactions.TransactionHandle{Format: handle.JSON}
err := client.Transactions().Begin(&respHandle)

// Extract transaction ID
txInfo := respHandle.Deserialized().(transactions.TransactionInfo)
txid := txInfo.TxID
```

#### Execute Multiple Operations in a Transaction

```go
// Begin transaction
txResp := transactions.TransactionHandle{Format: handle.JSON}
client.Transactions().Begin(&txResp)
txid := txResp.Deserialized().(transactions.TransactionInfo).TxID

// Create transaction object for use in operations
txn := client.NewTransaction()
txn.ID = txid

// Now use txn in document or search operations
doc := &documents.DocumentDescription{
	URI:     "/txn-test/doc1.json",
	Content: bytes.NewBufferString(`{"status":"active"}`),
	Format:  handle.JSON,
}

err := client.Documents().Write([]*documents.DocumentDescription{doc}, nil, txn, &respHandle)

// Commit the transaction
commitResp := transactions.TransactionHandle{Format: handle.JSON}
client.Transactions().Commit(txid, &commitResp)
```

#### Rollback a Transaction

```go
// Rollback instead of commit to discard all changes
rollbackResp := transactions.TransactionHandle{Format: handle.JSON}
err := client.Transactions().Rollback(txid, &rollbackResp)
```

#### Check Transaction Status

```go
// Check the status of an active transaction
statusResp := transactions.TransactionHandle{Format: handle.JSON}
err := client.Transactions().Status(txid, &statusResp)

info := statusResp.Deserialized().(transactions.TransactionInfo)
fmt.Printf("Transaction %s is %s\n", info.TxID, info.Status)
```

### User-Defined Resource Extensions

#### Call a GET Resource Extension

```go
import (
	"github.com/ryanjdew/go-marklogic-go/resources"
)

// Call a custom GET extension with parameters
params := map[string]string{
	"format": "json",
	"limit":  "10",
}

respHandle := resources.ResourceExtensionHandle{Format: handle.JSON}
err := client.Resources().Get("my-resource", params, &respHandle)

resource := respHandle.Deserialized().(resources.ResourceExtension)
fmt.Printf("Resource: %s - %s\n", resource.Name, resource.Title)
```

#### Call a POST Resource Extension

```go
// POST data to a custom resource
payload := handle.RawHandle{Format: handle.JSON}
payload.Write([]byte(`{"action":"process","data":"input"}`))

resp := resources.ResourceExtensionHandle{Format: handle.JSON}
err := client.Resources().Post("my-resource", nil, &payload, &resp)
```

#### Call a PUT Resource Extension

```go
// PUT (update) data in a custom resource
payload := handle.RawHandle{Format: handle.JSON}
payload.Write([]byte(`{"title":"Updated Title","status":"active"}`))

resp := resources.ResourceExtensionHandle{Format: handle.JSON}
err := client.Resources().Put("my-resource", nil, &payload, &resp)
```

#### Call a DELETE Resource Extension

```go
// DELETE a resource
resp := resources.ResourceExtensionHandle{Format: handle.JSON}
err := client.Resources().Delete("my-resource", nil, &resp)
```

### Flexible Indexes

#### List All Configured Indexes

```go
import (
	"github.com/ryanjdew/go-marklogic-go/indexes"
)

// List all indexes in the system
respHandle := indexes.IndexHandle{Format: handle.JSON}
err := client.Indexes().ListIndexes(nil, &respHandle)
if err != nil {
	log.Fatal(err)
}
```

#### Create a Range Index

```go
// Create an element range index for efficient value comparisons
indexConfig := handle.RawHandle{Format: handle.JSON}
indexConfig.Write([]byte(`{
	"scalar-type": "xs:string",
	"namespace": "http://example.com/",
	"localname": "title",
	"collation": "http://marklogic.com/collation/codepoint"
}`))

respHandle := indexes.IndexHandle{Format: handle.JSON}
err := client.Indexes().CreateIndex(&indexConfig, &respHandle)
```

#### Create a Field Index

```go
// Create a field index for faster value retrieval across multiple elements
fieldIndex := handle.RawHandle{Format: handle.JSON}
fieldIndex.Write([]byte(`{
	"field-name": "text-search",
	"scalar-types": ["xs:string"],
	"tokenizer": "http://marklogic.com/xdmp/tokenizer/default",
	"included-elements": ["title", "body"]
}`))

respHandle := indexes.IndexHandle{Format: handle.JSON}
err := client.Indexes().CreateIndex(&fieldIndex, &respHandle)
```

#### Get Index Details

```go
// Retrieve configuration for a specific index
respHandle := indexes.IndexHandle{Format: handle.JSON}
err := client.Indexes().GetIndex("text-search", nil, &respHandle)
```

#### Update an Index

```go
// Update an existing index configuration
updatedIndex := handle.RawHandle{Format: handle.JSON}
updatedIndex.Write([]byte(`{
	"scalar-type": "xs:string",
	"collation": "http://marklogic.com/collation/en"
}`))

respHandle := indexes.IndexHandle{Format: handle.JSON}
err := client.Indexes().UpdateIndex("element-range-index-1", &updatedIndex, &respHandle)
```

#### Delete an Index

```go
// Remove an index from the system
respHandle := indexes.IndexHandle{Format: handle.JSON}
err := client.Indexes().DeleteIndex("text-search", nil, &respHandle)
```

### Query Configuration

#### Installing Query Options

```go
import (
	config "github.com/ryanjdew/go-marklogic-go/config"
)

// Create query options
opts := config.QueryOptions{
	ReturnQuery: true,
	PageLength: 20,
	// ... other options
}

// Install query options
err := client.Config().SetQueryOptions("my-options", &opts)
```

#### Getting Query Options

```go
// Retrieve installed query options
opts, err := client.Config().GetQueryOptions("my-options")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Page length: %d\n", opts.PageLength)
```

### Multi-Statement Transactions

```go
// Create a transaction
txn := client.NewTransaction()

// Use transaction ID in operations
respHandle := search.ResponseHandle{Format: handle.JSON}
err := client.Search().StructuredSearch(&qh, 1, 10, txn, &respHandle)

// Other operations with same txn...

// Commit or rollback happens through the transaction API
```

### Bulk Operations

```go
import (
	datamovement "github.com/ryanjdew/go-marklogic-go/datamovement"
)

// Use DataMovement service for optimized batch operations
batcher := client.DataMovement().WriteBatcher()

// Add documents to batch
for i := 0; i < 1000; i++ {
	doc := &documents.DocumentDescription{
		URI: fmt.Sprintf("/docs/doc%d.xml", i),
		Content: bytes.NewBufferString(fmt.Sprintf("<doc>%d</doc>", i)),
		Format: handle.XML,
	}
	batcher.Add(doc)
}

// Flush and wait for completion
err := batcher.Flush()
```

## Handle Pattern

The library uses a **Handle** abstraction for serialization/deserialization. This allows flexible format handling (JSON, XML, multipart/mixed) without format-specific code.

### Creating Handles

```go
// JSON format
jsonHandle := search.ResponseHandle{Format: handle.JSON}

// XML format
xmlHandle := search.ResponseHandle{Format: handle.XML}

// Raw handle for binary or passthrough data
rawHandle := handle.RawHandle{Format: handle.UNKNOWN}
```

### Using Handles

```go
// For requests: serialize your data
query := search.Query{...}
qh := search.QueryHandle{Format: handle.XML}
qh.Serialize(query)  // Encodes to XML

// Service automatically uses the format
respHandle := search.ResponseHandle{Format: handle.JSON}
err := client.Search().StructuredSearch(&qh, 1, 10, nil, &respHandle)

// For responses: access both serialized and deserialized forms
rawJSON := respHandle.Serialized()  // Get raw JSON string
results := respHandle.Deserialized() // Get typed Response struct
```

### Supported Formats

| Format | MIME Type | Use Case |
|--------|-----------|----------|
| `JSON` | `application/json` | Structured data in JSON format |
| `XML` | `application/xml` | Structured data in XML format |
| `MIXED` | `multipart/mixed` | Multiple documents with metadata |
| `TEXTPLAIN` | `text/plain` | Plain text content |
| `TEXT_URI_LIST` | `text/uri-list` | Lists of document URIs |
| `UNKNOWN` | `application/octet-stream` | Binary or unknown format |

## Error Handling

All service methods return an `error` as the last return value. Check errors immediately:

```go
respHandle := search.ResponseHandle{Format: handle.JSON}
if err := client.Search().Search("query", 1, 10, nil, &respHandle); err != nil {
	log.Printf("Search failed: %v", err)
	// Handle error (e.g., connection issues, HTTP 4xx/5xx)
}
```

HTTP errors (status ≥400) are converted to Go errors. Common scenarios:

- **Connection errors**: Host unreachable, authentication failed
- **404 Not Found**: Document URI doesn't exist, resource not found
- **400 Bad Request**: Invalid query syntax, malformed input
- **500+ Server Errors**: MarkLogic internal errors

## Development & Testing

### Setup

```bash
# Install dependencies
go mod tidy

# Start MarkLogic in Docker (requires Docker)
task setup

# Build
task build
```

### Running Tests

Tests require a running MarkLogic instance:

```bash
# Run all integration tests
task test

# Or manually with go test
go test ./... -tags=integration
```

### Building Your Own Service

To add a new service to the library:

1. Create a package (e.g., `myservice/`)
2. Implement service.go with factory:
   ```go
   type Service struct {
       client *clients.Client
   }
   
   func NewService(c *clients.Client) *Service {
       return &Service{client: c}
   }
   ```
3. Implement operations in `myservice.go` using `util.BuildRequestFromHandle()` and `util.Execute()`
4. Add factory to root Client in `clientAPI.go`

## Architecture

The library uses three layers:

- **Client Layer** (`clients/`): Low-level HTTP handling with authentication
- **Service Layer** (`search/`, `documents/`, etc.): Domain-specific operations
- **Handle Layer** (`handle/`): Format-agnostic serialization/deserialization

Service methods are thin wrappers delegating to package-level functions, enabling testability without mocking and clear separation of concerns.

## MarkLogic REST API Reference

This library maps to MarkLogic's REST API. For more details on specific operations:

- **Search**: [/v1/search](https://docs.marklogic.com/REST/GET/v1/search) - Full-text, structured queries, faceting
- **Documents**: [/v1/documents](https://docs.marklogic.com/REST/GET/v1/documents) - CRUD operations
- **Semantics**: [/v1/graphs](https://docs.marklogic.com/REST/GET/v1/graphs), [/v1/graphs/sparql](https://docs.marklogic.com/REST/POST/v1/graphs/sparql) - Triple store, SPARQL queries
- **Values**: [/v1/values](https://docs.marklogic.com/REST/GET/v1/values/{name}) - Lexicon value enumeration
- **Eval**: [/v1/eval](https://docs.marklogic.com/REST/POST/v1/eval) - Server-side XQuery/JavaScript execution
- **Suggest**: [/v1/suggest](https://docs.marklogic.com/REST/GET/v1/suggest) - Query suggestions and autocomplete
- **Indexes**: [/v1/config/indexes](https://docs.marklogic.com/REST/GET/v1/config/indexes) - Range and field index management
- **Configuration**: [/v1/config](https://docs.marklogic.com/REST/GET/v1/config/query) - Query options, transforms, extensions
- **Transactions**: [/v1/transactions](https://docs.marklogic.com/REST/POST/v1/transactions) - Multi-statement transactions
- **Resources**: [/v1/resources](https://docs.marklogic.com/REST/GET/v1/resources/{name}) - Custom REST extensions
- **Full API Reference**: [docs.marklogic.com/REST](https://docs.marklogic.com/REST)

## License

Apache License 2.0 - See LICENSE.txt
