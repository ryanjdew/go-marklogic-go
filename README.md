# Go MarkLogic Client Library

A comprehensive Go client library for interacting with MarkLogic's REST APIs. This library provides type-safe, idiomatic Go interfaces to MarkLogic's document, search, semantic, and data management capabilities.

## Status

[![GoDoc](https://godoc.org/github.com/ryanjdew/go-marklogic-go?status.svg)](https://godoc.org/github.com/ryanjdew/go-marklogic-go) [![Build Status](https://cloud.drone.io/api/badges/ryanjdew/go-marklogic-go/status.svg)](https://cloud.drone.io/ryanjdew/go-marklogic-go)

## Features

- **Full-text & Structured Search** - Query text, use structured queries, get suggestions, and retrieve faceted results
- **Document CRUD** - Read, write, update, and delete documents with metadata management
- **Semantic/Graph Operations** - Work with RDF triples, SPARQL queries, and semantic data
- **Lexicon Value Enumeration** - Query distinct values from range indexes with aggregation and co-occurrence analysis
- **Server-Side Code Execution** - Execute XQuery and JavaScript on the server with external variables
- **Query Suggestions** - Get suggestions for queries and implement autocomplete
- **Multi-Statement Transactions** - Atomically execute multiple document or search operations
- **User-Defined Resources** - Call custom REST extensions (GET/POST/PUT/DELETE)
- **Flexible Indexes** - Create and manage range and field indexes for optimized performance
- **Temporal Document Operations** - Track document versions across time with temporal axes and system time management
- **Metadata Extraction & Validation** - Extract metadata from documents and validate against configurable rules
- **Bulk Operations** - Efficiently batch read/write multiple documents
- **Optic Queries** - Execute SQL-like queries on structured data with explain and query planning
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

#### Aggregate Values

```go
// Calculate aggregate operations on lexicon values (sum, count, min, max, etc.)
aggregateBody := `<values-query xmlns="http://marklogic.com/appservices/search">
  <range name="count" type="xs:int">
    <aggregate function="sum" />
  </range>
</values-query>`

queryHandle := handle.RawHandle{Format: handle.XML}
queryHandle.Write([]byte(aggregateBody))

respHandle := values.ResponseHandle{Format: handle.XML}
err := client.Values().AggregateValues(
	"count",
	nil,
	&queryHandle,
	&respHandle,
)

// Results contain sum, min, max, count, average
```

#### Co-occurrence Values

```go
// Find correlated values across multiple lexicons
queryHandle := handle.RawHandle{Format: handle.XML}
queryHandle.Write([]byte(`<values-query xmlns="http://marklogic.com/appservices/search"></values-query>`))

respHandle := values.ResponseHandle{Format: handle.XML}
err := client.Values().CoOccurrenceValues(
	[]string{"status", "priority"},  // Multiple lexicon names
	nil,
	&queryHandle,
	&respHandle,
)

// Results show value pairs and their co-occurrence frequency
```

#### Tuple Values

```go
// Generate tuples of related values from multiple lexicons
queryHandle := handle.RawHandle{Format: handle.XML}
queryHandle.Write([]byte(`<values-query xmlns="http://marklogic.com/appservices/search"></values-query>`))

respHandle := values.ResponseHandle{Format: handle.XML}
err := client.Values().TupleValues(
	[]string{"status", "priority", "date"},  // Multiple lexicon names
	nil,
	&queryHandle,
	&respHandle,
)

// Results contain tuples with positional values from each lexicon
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

### Temporal Document and Axis Operations

#### Create a Temporal Axis

```go
import (
	"github.com/ryanjdew/go-marklogic-go/temporal"
)

// Create an axis for tracking system time
axisConfig := handle.RawHandle{Format: handle.JSON}
axisConfig.Write([]byte(`{
	"axis-name": "system-time",
	"scalar-uri": "xs:dateTime"
}`))

respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().CreateAxis("system-time", &axisConfig, &respHandle)
```

#### List All Temporal Axes

```go
// Retrieve all configured temporal axes
respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().ListAxes(&respHandle)
```

#### Get Temporal Axis Configuration

```go
// Get details about a specific axis
respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().GetAxis("system-time", &respHandle)
```

#### Enable Temporal Support on a Collection

```go
// Enable temporal versioning on a document collection
temporalConfig := handle.RawHandle{Format: handle.JSON}
temporalConfig.Write([]byte(`{
	"system-axis": "system-time",
	"valid-axis": "valid-dates"
}`))

respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().EnableCollectionTemporal("temporal-docs", &temporalConfig, &respHandle)
```

#### Retrieve a Temporal Document at a Specific Time

```go
// Get a document as it existed at a specific point in time
respHandle := temporal.TemporalHandle{Format: handle.JSON}
timestamp := "2024-06-15T12:00:00Z"
err := client.Temporal().GetTemporalDocument("/documents/doc1.json", timestamp, &respHandle)
```

#### Advance System Time

```go
// Advance the system clock for temporal operations
newTime := "2024-12-31T23:59:59Z"
respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().AdvanceSystemTime(newTime, &respHandle)
```

#### Get Current System Time

```go
// Retrieve the current system time for temporal operations
respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().GetSystemTime(&respHandle)

sysTime := respHandle.Deserialized().(temporal.SystemTime)
fmt.Printf("Current system time: %s\n", sysTime.Timestamp)
```

#### Disable Temporal Support on a Collection

```go
// Remove temporal versioning from a collection
respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().DisableCollectionTemporal("temporal-docs", &respHandle)
```

#### Delete a Temporal Axis

```go
// Remove a temporal axis from the system
respHandle := temporal.TemporalHandle{Format: handle.JSON}
err := client.Temporal().DeleteAxis("system-time", &respHandle)
```

### Metadata Extraction and Validation

#### Extract Metadata from Specific Documents

```go
import (
	"github.com/ryanjdew/go-marklogic-go/metadata"
)

// Extract metadata from specific URIs
options := map[string]string{
	"extractors": "document-properties,element-values",
}

respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().ExtractMetadata([]string{"/doc1.json", "/doc2.json"}, options, &respHandle)

// Access the metadata results
results := respHandle.Deserialized()
fmt.Printf("Extraction results: %v\n", results)
```

#### Extract Metadata from Query Results

```go
// Extract metadata from documents matching a query
queryPayload := handle.RawHandle{Format: handle.JSON}
queryPayload.Write([]byte(`{
	"query": {
		"text-query": {
			"text": "important"
		}
	}
}`))

respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().ExtractMetadataFromQuery(&queryPayload, nil, &respHandle)
```

#### Extract Metadata from a Single Document

```go
// Extract metadata from one document
options := map[string]string{
	"format": "json",
}

respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().ExtractMetadataFromURI("/doc1.json", options, &respHandle)

result := respHandle.Deserialized().(metadata.MetadataResult)
fmt.Printf("Document: %s\nMetadata: %v\n", result.URI, result.Metadata)
```

#### Set Validation Rules

```go
// Configure validation rules for the database
rulesPayload := handle.RawHandle{Format: handle.JSON}
rulesPayload.Write([]byte(`{
	"rules": [
		{
			"name": "required-title",
			"rule-type": "xpath",
			"xpath": "title",
			"action": "fail",
			"message": "Document must contain a title element"
		},
		{
			"name": "valid-status",
			"rule-type": "enum",
			"xpath": "status",
			"values": ["draft", "published", "archived"],
			"action": "warn"
		}
	]
}`))

respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().SetValidationRules(&rulesPayload, &respHandle)
```

#### Get Current Validation Rules

```go
// Retrieve the configured validation rules
respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().GetValidationRules(&respHandle)
```

#### Validate Specific Documents

```go
// Validate documents against rules
rulesPayload := handle.RawHandle{Format: handle.JSON}
rulesPayload.Write([]byte(`{
	"rules": [{
		"name": "required-title",
		"rule-type": "xpath",
		"xpath": "title",
		"action": "fail"
	}]
}`))

respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().ValidateDocuments([]string{"/doc1.json"}, &rulesPayload, &respHandle)

// Check validation results
results := respHandle.Deserialized()
fmt.Printf("Validation results: %v\n", results)
```

#### Validate a Single Document

```go
// Validate one document
rulesPayload := handle.RawHandle{Format: handle.JSON}
rulesPayload.Write([]byte(`{
	"rules": [{
		"name": "required-title",
		"xpath": "title",
		"action": "fail"
	}]
}`))

respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().ValidateURI("/doc1.json", &rulesPayload, &respHandle)

result := respHandle.Deserialized().(metadata.ValidationResult)
fmt.Printf("Document: %s - Valid: %v\n", result.URI, result.Valid)
if !result.Valid {
	fmt.Printf("Errors: %v\n", result.Errors)
}
```

#### Validate Documents from Query Results

```go
// Validate documents matching a query against rules
queryPayload := handle.RawHandle{Format: handle.JSON}
queryPayload.Write([]byte(`{
	"query": {
		"text-query": {
			"text": "important"
		}
	}
}`))

rulesPayload := handle.RawHandle{Format: handle.JSON}
rulesPayload.Write([]byte(`{
	"rules": [{
		"name": "required-title",
		"xpath": "title",
		"action": "fail"
	}]
}`))

respHandle := metadata.MetadataHandle{Format: handle.JSON}
err := client.Metadata().ValidateQuery(&queryPayload, &rulesPayload, &respHandle)
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

### Optic Queries

The Rows Management service executes MarkLogic's Optic API for advanced SQL-like queries on structured data.

#### Execute Optic Plan

```go
import (
	rowsManagement "github.com/ryanjdew/go-marklogic-go/rows-management"
)

// Execute an Optic plan
opticPlan := handle.RawHandle{Format: handle.JSON}
opticPlan.Write([]byte(`{
	"_export": "EMPLOYEES",
	"_module": {"_import": "marklogic-optic", "fn": "op"},
	"_next": {
		"_export": "FROM_EMPLOYEES",
		"_invoke": {"_function": "from", "_args": ["EMPLOYEES", "employees"]}
	}
}`))

respHandle := handle.RawHandle{Format: handle.JSON}
err := client.RowsManagement().Rows(&opticPlan, nil, &respHandle)

// Results contain the rows in JSON format
results := respHandle.Deserialized()
```

#### Explain Query Plan

```go
// Get query plan and execution estimates
opticPlan := handle.RawHandle{Format: handle.JSON}
opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

respHandle := handle.RawHandle{Format: handle.JSON}
err := client.RowsManagement().Explain(&opticPlan, nil, &respHandle)

// Response contains query plan, estimated row counts, and cost analysis
```

#### Sample Query Results

```go
// Get a sample of results without computing full result set
opticPlan := handle.RawHandle{Format: handle.JSON}
opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

params := map[string]string{"size": "10"}
respHandle := handle.RawHandle{Format: handle.JSON}
err := client.RowsManagement().Sample(&opticPlan, params, &respHandle)

// Response contains sample rows
```

#### Get Optimized Plan

```go
// Get the optimized execution plan
opticPlan := handle.RawHandle{Format: handle.JSON}
opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

respHandle := handle.RawHandle{Format: handle.JSON}
err := client.RowsManagement().Plan(&opticPlan, nil, &respHandle)

// Response shows index usage and access path optimization
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
````
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
- **Values**: [/v1/values](https://docs.marklogic.com/REST/GET/v1/values/{name}) - Lexicon value enumeration with aggregation
- **Eval**: [/v1/eval](https://docs.marklogic.com/REST/POST/v1/eval) - Server-side XQuery/JavaScript execution
- **Suggest**: [/v1/suggest](https://docs.marklogic.com/REST/GET/v1/suggest) - Query suggestions and autocomplete
- **Indexes**: [/v1/config/indexes](https://docs.marklogic.com/REST/GET/v1/config/indexes) - Range and field index management
- **Temporal**: [/v1/temporal](https://docs.marklogic.com/REST/GET/v1/temporal/axes) - Temporal axes, collections, and document versioning
- **Metadata**: [/v1/metadata](https://docs.marklogic.com/REST/GET/v1/metadata) - Metadata extraction and document validation
- **Rows (Optic)**: [/v1/rows](https://docs.marklogic.com/REST/POST/v1/rows) - SQL-like queries with explain and sampling
- **Configuration**: [/v1/config](https://docs.marklogic.com/REST/GET/v1/config/query) - Query options, transforms, extensions
- **Transactions**: [/v1/transactions](https://docs.marklogic.com/REST/POST/v1/transactions) - Multi-statement transactions
- **Resources**: [/v1/resources](https://docs.marklogic.com/REST/GET/v1/resources/{name}) - Custom REST extensions
- **Full API Reference**: [docs.marklogic.com/REST](https://docs.marklogic.com/REST)

## License

Apache License 2.0 - See LICENSE.txt
