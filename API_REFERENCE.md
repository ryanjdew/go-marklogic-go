# Go MarkLogic Client Library - API Reference

This document provides a comprehensive reference for all major services and APIs in the Go MarkLogic client library, with detailed examples for each.

## Table of Contents

1. [Connection & Setup](#connection--setup)
2. [Search Service](#search-service)
3. [Documents Service](#documents-service)
4. [Semantics Service](#semantics-service)
5. [Values Service](#values-service)
6. [Configuration Service](#configuration-service)
7. [Resources Service](#resources-service)
8. [Transactions](#transactions)
9. [Temporal Operations](#temporal-operations)
10. [Data Services](#data-services)
11. [Data Movement](#data-movement)
12. [Rows/Optic Queries](#rowsoptic-queries)
13. [Indexes](#indexes)
14. [Alerting](#alerting)
15. [Metadata Operations](#metadata-operations)
16. [Error Handling](#error-handling)

## Connection & Setup

### Basic Connection

```go
package main

import (
    "log"
    marklogic "github.com/ryanjdew/go-marklogic-go"
)

func main() {
    // Create a client with Digest authentication
    client, err := marklogic.NewClient(
        "localhost",              // Host
        8050,                     // Port (REST API server port)
        "admin",                  // Username
        "admin",                  // Password
        marklogic.DigestAuth,     // Auth type: DigestAuth, BasicAuth, or None
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Client is ready to use throughout application lifetime
}
```

### Custom Connection

```go
// For more control over connection parameters
conn := &marklogic.Connection{
    Host:               "localhost",
    Port:               8050,
    Username:           "admin",
    Password:           "admin",
    AuthenticationType: marklogic.DigestAuth,
}

client, err := marklogic.New(conn)
if err != nil {
    log.Fatal(err)
}
```

### Authentication Options

```go
// Option 1: Basic Authentication (credentials in every request)
client, _ := marklogic.NewClient(host, port, user, pass, marklogic.BasicAuth)

// Option 2: Digest Authentication (challenge/response)
client, _ := marklogic.NewClient(host, port, user, pass, marklogic.DigestAuth)

// Option 3: No Authentication
client, _ := marklogic.NewClient(host, port, "", "", marklogic.None)
```

## Search Service

The Search service provides full-text and structured query capabilities.

### Full-Text Search

```go
import (
    "fmt"
    "log"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    search "github.com/ryanjdew/go-marklogic-go/search"
)

// Basic text search
respHandle := search.ResponseHandle{Format: handle.JSON}
err := client.Search().Search(
    "Shakespeare",    // Query text
    1,               // Start (page number)
    10,              // Page length
    nil,             // Transaction (optional)
    &respHandle,     // Response handle
)
if err != nil {
    log.Fatal(err)
}

// Access results
results := respHandle.Get().(*search.Response)
fmt.Printf("Found %d results\n", results.Total)
for _, doc := range results.Results {
    fmt.Printf("- %s: %s\n", doc.URI, doc.Score)
}
```

### Structured Search

```go
// Build structured query
query := search.Query{
    Queries: []interface{}{
        search.TermQuery{
            Terms: []string{"active", "approved"},
        },
    },
    Operator: "or",
}

// Execute structured search
queryHandle := search.QueryHandle{Format: handle.JSON}
queryHandle.Serialize(query)

respHandle := search.ResponseHandle{Format: handle.JSON}
err := client.Search().StructuredSearch(
    &queryHandle,    // Structured query
    1,              // Start
    10,             // Page length
    nil,            // Transaction
    &respHandle,    // Response
)
if err != nil {
    log.Fatal(err)
}
```

### Search with Suggestions

```go
// Get suggestions for partial query
respHandle := search.SuggestionsResponseHandle{Format: handle.JSON}
err := client.Search().Suggest(
    "shake",    // Partial query
    5,         // Max suggestions
    "default", // Options name
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

suggestions := respHandle.Get().(*search.SuggestionsResponse)
for _, sugg := range suggestions.Suggestions {
    fmt.Println(sugg)
}
```

### Structured Suggestions

```go
// Suggestions within structured query context
query := search.Query{...}
queryHandle := search.QueryHandle{Format: handle.JSON}
queryHandle.Serialize(query)

respHandle := search.SuggestionsResponseHandle{Format: handle.JSON}
err := client.Search().StructuredSuggestions(
    &queryHandle,   // Base query
    "shake",        // Partial text
    5,             // Max suggestions
    "default",     // Options
    nil,           // Transaction
    &respHandle,
)
```

### Delete by Query

```go
import "github.com/ryanjdew/go-marklogic-go/util"

// Delete all documents matching a query
params := map[string]string{
    "q": "status:archived",
}

respHandle := search.ResponseHandle{Format: handle.JSON}
err := client.Search().Delete(
    params,           // Query parameters
    nil,             // Transaction
    &respHandle,
)
```

## Documents Service

CRUD operations for documents with metadata management.

### Reading Documents

```go
import (
    "log"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    documents "github.com/ryanjdew/go-marklogic-go/documents"
)

// Read single document
respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.Documents().Read(
    []string{"/path/to/document.json"}, // URIs
    []string{"content"},                // Categories (content, metadata, etc.)
    nil,                               // Transform (optional)
    nil,                               // Transaction
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

// Get document content
content := respHandle.Deserialized().(string)
fmt.Println(content)
```

### Reading Multiple Documents

```go
// Read multiple documents with metadata
uris := []string{"/doc1.json", "/doc2.json", "/doc3.json"}

respHandle := documents.ResponseHandle{Format: handle.MIXED}
err := client.Documents().Read(
    uris,
    []string{"content", "metadata"},  // Get both content and metadata
    nil,
    nil,
    &respHandle,
)
```

### Writing Documents

```go
// Create document descriptions
docs := []*documents.DocumentDescription{
    {
        URI:     "/users/alice.json",
        Content: bytes.NewBufferString(`{"name":"Alice","status":"active"}`),
        Metadata: &documents.Metadata{
            Collections: []string{"users", "active"},
            Properties: map[string]string{
                "created": "2024-01-15",
                "type":    "user",
            },
        },
        Format: handle.JSON,
    },
    {
        URI:     "/users/bob.json",
        Content: bytes.NewBufferString(`{"name":"Bob","status":"active"}`),
        Metadata: &documents.Metadata{
            Collections: []string{"users", "active"},
        },
        Format: handle.JSON,
    },
}

respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.Documents().Write(
    docs,           // Documents to write
    nil,           // Transform (optional)
    nil,           // Transaction
    &respHandle,
)
```

### Updating Documents

```go
// Update existing document
doc := &documents.DocumentDescription{
    URI:     "/users/alice.json",
    Content: bytes.NewBufferString(`{"name":"Alice","status":"inactive"}`),
    Format:  handle.JSON,
}

respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.Documents().Write([]*documents.DocumentDescription{doc}, nil, nil, &respHandle)
```

### Deleting Documents

```go
// Delete single document
respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.Documents().Delete(
    []string{"/users/alice.json"},  // URIs to delete
    nil,                            // Transaction
    &respHandle,
)

// Delete multiple documents
err = client.Documents().Delete(
    []string{"/doc1.json", "/doc2.json", "/doc3.json"},
    nil,
    &respHandle,
)
```

### Document Metadata

```go
// Set permissions
perms := &documents.Permissions{
    Permissions: map[string][]string{
        "admin": {"read", "update"},
        "user":  {"read"},
    },
}

// Set metadata
metadata := &documents.Metadata{
    Collections: []string{"sensitive", "archived"},
    Permissions: perms,
    Properties: map[string]string{
        "classification": "confidential",
        "reviewed":       "true",
    },
    Quality: 10,
}

doc := &documents.DocumentDescription{
    URI:      "/secure/document.json",
    Content:  buffer,
    Metadata: metadata,
    Format:   handle.JSON,
}

respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.Documents().Write([]*documents.DocumentDescription{doc}, nil, nil, &respHandle)
```

### Applying Transforms

```go
import "github.com/ryanjdew/go-marklogic-go/util"

// Apply server-side transform
transform := &util.Transform{
    Name: "extract-metadata",
    Params: map[string]string{
        "threshold": "0.8",
        "format":    "json",
    },
}

respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.Documents().Read(
    []string{"/documents/text.txt"},
    []string{"content"},
    transform,  // Transform applied server-side
    nil,
    &respHandle,
)
```

## Semantics Service

Work with RDF triples and semantic graph operations.

### SPARQL Queries

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    semantics "github.com/ryanjdew/go-marklogic-go/semantics"
)

// Execute SPARQL query
sparqlQuery := `
PREFIX foaf: <http://xmlns.com/foaf/0.1/>
SELECT ?name ?email 
WHERE {
    ?person foaf:name ?name .
    ?person foaf:mbox ?email .
}
LIMIT 10
`

respHandle := semantics.SparqlResponseHandle{Format: handle.JSON}
err := client.Semantics().Sparql(
    sparqlQuery,      // SPARQL query
    nil,             // Query options (optional)
    nil,             // Transaction
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

results := respHandle.Get()
fmt.Println(results)
```

### Graph CRUD Operations

```go
// Write triples to graph
triples := `
<http://example.com/alice> <http://xmlns.com/foaf/0.1/name> "Alice" .
<http://example.com/bob> <http://xmlns.com/foaf/0.1/name> "Bob" .
`

graphHandle := handle.Handle{Format: handle.XML}
graphHandle.Serialize(triples)

respHandle := semantics.GraphResponseHandle{Format: handle.JSON}
err := client.Semantics().PutGraph(
    "my-graph",      // Graph URI
    &graphHandle,    // Triple data
    &respHandle,
)

// Read graph
respHandle = semantics.GraphResponseHandle{Format: handle.XML}
err = client.Semantics().GetGraph("my-graph", &respHandle)
```

### Things (Entity Operations)

```go
// Get semantic information for IRIs
iris := []string{
    "http://example.com/alice",
    "http://example.com/bob",
}

respHandle := semantics.ThingsResponseHandle{Format: handle.JSON}
err := client.Semantics().Things(iris, &respHandle)
```

## Values Service

Enumerate distinct values from range indexes with aggregation and co-occurrence analysis.

### List Values

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    values "github.com/ryanjdew/go-marklogic-go/values"
)

// List distinct values from index
params := map[string]string{
    "skip":  "0",
    "limit": "100",
}

respHandle := values.ResponseHandle{Format: handle.JSON}
err := client.Values().ListValues(
    "category",      // Values name
    params,         // Query parameters
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

vals := respHandle.Get().(*values.ValuesResponse)
for _, val := range vals.Values {
    fmt.Printf("%s: %d\n", val.Value, val.Count)
}
```

### Query Values

```go
// Query values with constraints
queryHandle := handle.Handle{Format: handle.JSON}
queryHandle.Serialize(map[string]interface{}{
    "query": map[string]string{
        "value": "active*",  // Wildcard query
    },
})

respHandle := values.ResponseHandle{Format: handle.JSON}
err := client.Values().QueryValues(
    "status",       // Values name
    nil,           // Parameters
    &queryHandle,  // Query constraints
    &respHandle,
)
```

### Aggregate Values

```go
// Aggregate by values
aggregateHandle := handle.Handle{Format: handle.JSON}
aggregateHandle.Serialize(map[string]interface{}{
    "aggregate-results": true,
})

respHandle := values.ResponseHandle{Format: handle.JSON}
err := client.Values().AggregateValues(
    "category",          // Values name
    nil,                // Parameters
    &aggregateHandle,   // Aggregation config
    &respHandle,
)
```

### Co-occurrence Analysis

```go
// Analyze value co-occurrence
respHandle := values.ResponseHandle{Format: handle.JSON}
err := client.Values().CoOccurrenceValues(
    []string{"category", "status"},  // Values to correlate
    nil,                            // Parameters
    nil,                            // Request body
    &respHandle,
)
```

## Configuration Service

Manage query options, transforms, and server extensions.

### Query Options Management

```go
import (
    "bytes"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    config "github.com/ryanjdew/go-marklogic-go/config"
)

// Create/update query options
optionsXML := `
<search:options xmlns:search="http://marklogic.com/appservices/search">
    <search:term-option>wildcarded</search:term-option>
    <search:constraint name="collection">
        <search:collection-facet>true</search:collection-facet>
    </search:constraint>
</search:options>
`

optionsHandle := handle.Handle{Format: handle.XML}
optionsHandle.Serialize(optionsXML)

respHandle := config.ResponseHandle{Format: handle.JSON}
err := client.Config().SetQueryOptions(
    "my-options",       // Options name
    &optionsHandle,    // Options definition
    &respHandle,
)

// Retrieve query options
respHandle = config.ResponseHandle{Format: handle.XML}
err = client.Config().GetQueryOptions("my-options", &respHandle)
options := respHandle.Serialized()
```

### Transform Management

```go
// Create server-side transform
transformXML := `
<transform xmlns="http://marklogic.com/rest-api">
    <transform-source>
        xdmp:to-json(
            map:new((
                map:entry("original", $input),
                map:entry("modified", fn:current-dateTime())
            ))
        )
    </transform-source>
</transform>
`

transformHandle := handle.Handle{Format: handle.XML}
transformHandle.Serialize(transformXML)

respHandle := config.ResponseHandle{Format: handle.JSON}
err := client.Config().SetTransforms(
    "add-metadata",     // Transform name
    &transformHandle,  // Transform definition
    &respHandle,
)

// List all transforms
respHandle = config.ResponseHandle{Format: handle.JSON}
err = client.Config().ListTransforms(&respHandle)
transforms := respHandle.Get()
```

### Index Management

```go
// Create range index
indexXML := `
<range-index xmlns="http://marklogic.com/rest-api">
    <scalar-type>xs:string</scalar-type>
    <path-expression>//@status</path-expression>
    <collation>http://marklogic.com/collation/codepoint</collation>
</range-index>
`

indexHandle := handle.Handle{Format: handle.XML}
indexHandle.Serialize(indexXML)

respHandle := config.ResponseHandle{Format: handle.JSON}
err := client.Config().CreateIndex("my-index", &indexHandle, &respHandle)
```

## Resources Service

Call user-defined REST extensions.

### Custom GET Endpoint

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    resources "github.com/ryanjdew/go-marklogic-go/resources"
)

// Call custom GET extension
respHandle := resources.ResponseHandle{Format: handle.JSON}
err := client.Resources().Get(
    "my-extension",    // Extension name
    map[string]string{"param1": "value1"},  // Query parameters
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

result := respHandle.Serialized()
fmt.Println(result)
```

### Custom POST Endpoint

```go
// Call custom POST extension with data
requestHandle := handle.Handle{Format: handle.JSON}
requestHandle.Serialize(map[string]string{
    "operation": "process",
    "data":      "some-value",
})

respHandle := resources.ResponseHandle{Format: handle.JSON}
err := client.Resources().Post(
    "my-processor",    // Extension name
    map[string]string{"action": "analyze"},  // Query params
    &requestHandle,   // Request body
    &respHandle,
)
```

### Custom PUT Endpoint

```go
// Call custom PUT extension
requestHandle := handle.Handle{Format: handle.JSON}
requestHandle.Serialize(map[string]interface{}{
    "config": map[string]string{
        "threshold": "0.8",
        "mode":      "strict",
    },
})

respHandle := resources.ResponseHandle{Format: handle.JSON}
err := client.Resources().Put(
    "my-config",       // Extension name
    nil,              // Query parameters
    &requestHandle,   // Request body
    &respHandle,
)
```

## Transactions

Multi-statement transaction management for atomic operations.

### Basic Transaction

```go
import (
    "log"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    util "github.com/ryanjdew/go-marklogic-go/util"
    documents "github.com/ryanjdew/go-marklogic-go/documents"
    search "github.com/ryanjdew/go-marklogic-go/search"
)

// Create and begin transaction
txn := &util.Transaction{}
if !txn.Begin() {
    log.Fatal("Failed to begin transaction")
}

// Execute operations within transaction
doc := &documents.DocumentDescription{
    URI:     "/order/123",
    Content: buffer1,
    Format:  handle.JSON,
}

resp1 := documents.ResponseHandle{Format: handle.JSON}
err1 := client.Documents().Write(
    []*documents.DocumentDescription{doc},
    nil,
    txn,  // Use same transaction
    &resp1,
)

resp2 := search.ResponseHandle{Format: handle.JSON}
err2 := client.Search().Delete(
    map[string]string{"q": "status:pending"},
    txn,  // Use same transaction
    &resp2,
)

// Commit or rollback based on results
if err1 == nil && err2 == nil {
    txn.Commit()
} else {
    txn.Rollback()
}
```

### Transaction Status

```go
// Check transaction status
resp := util.TransactionStatusHandle{Format: handle.JSON}
err := client.Transactions().Status(txn.ID, &resp)
status := resp.Get().(*util.TransactionStatus)
fmt.Printf("Transaction state: %s\n", status.State)
```

## Temporal Operations

Work with temporal documents and version tracking.

### Create Temporal Axis

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    temporal "github.com/ryanjdew/go-marklogic-go/temporal"
)

// Create system time axis
axisXML := `
<temporal:axis xmlns:temporal="http://marklogic.com/xdmp/temporal">
    <temporal:axis-name>system-axis</temporal:axis-name>
    <temporal:axis-start>/date-start</temporal:axis-start>
    <temporal:axis-end>/date-end</temporal:axis-end>
</temporal:axis>
`

axisHandle := handle.Handle{Format: handle.XML}
axisHandle.Serialize(axisXML)

respHandle := temporal.ResponseHandle{Format: handle.JSON}
err := client.Temporal().CreateAxis(
    "system-axis",    // Axis name
    &axisHandle,     // Axis definition
    &respHandle,
)
```

### Enable Temporal on Collection

```go
// Enable temporal tracking for collection
configXML := `
<temporal:collection-temporal xmlns:temporal="http://marklogic.com/xdmp/temporal">
    <temporal:collection-name>versioned-documents</temporal:collection-name>
    <temporal:axis-name>system-axis</temporal:axis-name>
    <temporal:range-index-name>date-range</temporal:range-index-name>
</temporal:collection-temporal>
`

configHandle := handle.Handle{Format: handle.XML}
configHandle.Serialize(configXML)

respHandle := temporal.ResponseHandle{Format: handle.JSON}
err := client.Temporal().EnableCollectionTemporal(
    "versioned-documents",  // Collection
    &configHandle,         // Configuration
    &respHandle,
)
```

### Query Historical Documents

```go
// Get document at specific timestamp
respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.Temporal().GetTemporalDocument(
    "/versioned-doc.json",               // URI
    "2024-01-15T10:30:00Z",            // Timestamp
    &respHandle,
)

// Results contain document as it existed at that time
```

### Manage System Time

```go
// Advance system time for temporal testing
respHandle := temporal.ResponseHandle{Format: handle.JSON}
err := client.Temporal().AdvanceSystemTime(
    "2024-12-31T23:59:59Z",  // New system time
    &respHandle,
)

// Get current system time
respHandle = temporal.ResponseHandle{Format: handle.JSON}
err = client.Temporal().GetSystemTime(&respHandle)
time := respHandle.Serialized()
fmt.Println("System time:", time)
```

## Data Services

Execute server-side modules with parameters.

### Invoke Data Service

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    dataservices "github.com/ryanjdew/go-marklogic-go/dataservices"
)

// Call named Data Service module
params := map[string]string{
    "search-term": "active",
    "max-results": "100",
    "format":      "json",
}

respHandle := dataservices.ResponseHandle{Format: handle.JSON}
err := client.DataServices().Invoke(
    "my-data-service",  // Module name
    params,             // Parameters
    nil,               // Transaction
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

results := respHandle.Get()
fmt.Println(results)
```

## Data Movement

Bulk read/write operations optimized for performance.

### Bulk Write

```go
import (
    "bytes"
    datamovement "github.com/ryanjdew/go-marklogic-go/datamovement"
    documents "github.com/ryanjdew/go-marklogic-go/documents"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Create large batch of documents
docs := make([]*documents.DocumentDescription, 0)
for i := 0; i < 10000; i++ {
    doc := &documents.DocumentDescription{
        URI: fmt.Sprintf("/bulk-load/%d.json", i),
        Content: bytes.NewBufferString(
            fmt.Sprintf(`{"id":%d,"data":"value%d"}`, i, i)),
        Format: handle.JSON,
    }
    docs = append(docs, doc)
}

// Efficient bulk write
respHandle := documents.ResponseHandle{Format: handle.JSON}
err := client.DataMovement().WriteBatcher(
    docs,           // Documents
    nil,           // Transform
    nil,           // Transaction
    &respHandle,
)

if err != nil {
    log.Fatal(err)
}
```

### Bulk Read

```go
// Read large number of documents efficiently
uris := []string{}
for i := 0; i < 5000; i++ {
    uris = append(uris, fmt.Sprintf("/documents/%d.json", i))
}

respHandle := documents.ResponseHandle{Format: handle.MIXED}
err := client.DataMovement().ReadBatcher(
    uris,                          // URIs to read
    []string{"content"},           // Categories
    nil,                          // Transform
    nil,                          // Transaction
    &respHandle,
)
```

## Rows/Optic Queries

Structured queries on relational data using Optic API.

### Simple Optic Query

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    rowsManagement "github.com/ryanjdew/go-marklogic-go/rows-management"
)

// Build Optic query plan (DSL)
planDSL := `
op:col("MySchema", "Users")
  |> op:where(op:eq(op:col("status"), "active"))
  |> op:select(("id", "name", "email"))
  |> op:order-by(op:col("name"))
  |> op:limit(100)
`

planHandle := handle.Handle{Format: handle.JSON}
planHandle.Serialize(map[string]interface{}{
    "plan": planDSL,
})

respHandle := rowsManagement.ResponseHandle{Format: handle.JSON}
err := client.RowsManagement().Rows(
    &planHandle,   // Plan
    nil,          // Transaction
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

rows := respHandle.Get()
fmt.Println(rows)
```

### Optic with Joins

```go
// Complex Optic query with joins
planDSL := `
op:inner-join-doc(
    op:col("MySchema", "Orders"),
    op:col("DocURI"),
    "/orders/"
)
  |> op:where(op:gt(op:col("amount"), 1000))
  |> op:select(("order-id", "customer-name", "amount", "DocURI"))
`

planHandle := handle.Handle{Format: handle.JSON}
planHandle.Serialize(map[string]interface{}{
    "plan": planDSL,
})

respHandle := rowsManagement.ResponseHandle{Format: handle.JSON}
err := client.RowsManagement().Rows(&planHandle, nil, &respHandle)
```

## Indexes

Create and manage range and field indexes.

### Create Range Index

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    indexes "github.com/ryanjdew/go-marklogic-go/indexes"
)

// Create JSON property range index
indexDef := `
<range-index xmlns="http://marklogic.com/rest-api">
    <scalar-type>xs:string</scalar-type>
    <json-property>//@status</json-property>
    <collation>http://marklogic.com/collation/codepoint</collation>
</range-index>
`

indexHandle := handle.Handle{Format: handle.XML}
indexHandle.Serialize(indexDef)

respHandle := indexes.ResponseHandle{Format: handle.JSON}
err := client.Indexes().CreateIndex(
    "status-index",   // Index name
    &indexHandle,    // Index definition
    &respHandle,
)
```

### Create Field Index

```go
// Create field index
fieldIndexDef := `
<field-index xmlns="http://marklogic.com/rest-api">
    <field-name>search-field</field-name>
    <include-root>true</include-root>
    <tokenizer-override>
        <regex>\s+</regex>
        <function>*/my-tokenizer#1</function>
    </tokenizer-override>
</field-index>
`

indexHandle := handle.Handle{Format: handle.XML}
indexHandle.Serialize(fieldIndexDef)

respHandle := indexes.ResponseHandle{Format: handle.JSON}
err := client.Indexes().CreateIndex(
    "search-field-index",
    &indexHandle,
    &respHandle,
)
```

## Alerting

Rule-based alerting and notifications.

### Create Alert Rule

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    alert "github.com/ryanjdew/go-marklogic-go/alert"
)

// Define alert rule
ruleXML := `
<alert:alert xmlns:alert="http://marklogic.com/alert">
    <alert:alert-name>high-value-orders</alert:alert-name>
    <alert:query>
        <alert:search:query-text>amount:[1000 TO 999999]</alert:search:query-text>
    </alert:query>
    <alert:action>
        <alert:action-name>notify-admin</alert:action-name>
    </alert:action>
</alert:alert>
`

ruleHandle := handle.Handle{Format: handle.XML}
ruleHandle.Serialize(ruleXML)

respHandle := alert.ResponseHandle{Format: handle.JSON}
err := client.Alerting().CreateRule(
    &ruleHandle,   // Rule definition
    &respHandle,
)
```

### Match Document Against Rules

```go
// Check which alert rules match a document
docURI := "/orders/12345.json"

respHandle := alert.MatchResponseHandle{Format: handle.JSON}
err := client.Alerting().Match(
    docURI,        // Document URI
    nil,          // Query parameters
    &respHandle,
)

matches := respHandle.Get()
fmt.Printf("Matched %d alert rules\n", len(matches))
```

## Metadata Operations

Extract and validate metadata.

### Extract Metadata

```go
import (
    "fmt"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    metadata "github.com/ryanjdew/go-marklogic-go/metadata"
)

// Extract metadata from document
respHandle := metadata.ResponseHandle{Format: handle.JSON}
err := client.Metadata().ReadMetadata(
    "/document.json",           // Document URI
    nil,                       // Query parameters
    &respHandle,
)
if err != nil {
    log.Fatal(err)
}

meta := respHandle.Get()
fmt.Printf("Collections: %v\n", meta.Collections)
fmt.Printf("Quality: %d\n", meta.Quality)
fmt.Printf("Properties: %v\n", meta.Properties)
```

### Validate Metadata

```go
// Validate metadata against rules
validationConfig := `
{
    "rules": [
        {
            "property": "classification",
            "required": true,
            "values": ["public", "internal", "confidential"]
        },
        {
            "property": "owner",
            "required": true
        }
    ]
}
`

configHandle := handle.Handle{Format: handle.JSON}
configHandle.Serialize(validationConfig)

respHandle := metadata.ResponseHandle{Format: handle.JSON}
err := client.Metadata().ValidateMetadata(
    "/document.json",    // Document URI
    &configHandle,      // Validation rules
    &respHandle,
)
```

## Error Handling

Proper error handling patterns.

### Basic Error Handling

```go
import (
    "fmt"
    "log"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    search "github.com/ryanjdew/go-marklogic-go/search"
)

// Basic error checking
respHandle := search.ResponseHandle{Format: handle.JSON}
if err := client.Search().Search("query", 1, 10, nil, &respHandle); err != nil {
    // Handle error: network issues, auth failure, 4xx/5xx responses
    log.Printf("Search error: %v\n", err)
    return
}

// Verify response deserialization
results := respHandle.Get()
if results == nil {
    log.Fatal("Failed to parse search results")
}
```

### Retry Logic

```go
import (
    "math"
    "time"
)

// Exponential backoff retry
const maxRetries = 3

for attempt := 0; attempt < maxRetries; attempt++ {
    respHandle := search.ResponseHandle{Format: handle.JSON}
    err := client.Search().Search("query", 1, 10, nil, &respHandle)
    
    if err == nil {
        // Success - process results
        results := respHandle.Get()
        break
    }
    
    if attempt < maxRetries-1 {
        // Wait with exponential backoff: 1s, 2s, 4s
        backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
        time.Sleep(backoff)
    } else {
        log.Fatal("Search failed after retries:", err)
    }
}
```

### Transaction Rollback on Error

```go
// Transactional error handling
txn := &util.Transaction{}
txn.Begin()

// Perform operations
resp1 := documents.ResponseHandle{Format: handle.JSON}
err1 := client.Documents().Write(docs, nil, txn, &resp1)

resp2 := search.ResponseHandle{Format: handle.JSON}
err2 := client.Search().Delete(criteria, txn, &resp2)

// Automatic rollback on error
if err1 != nil || err2 != nil {
    txn.Rollback()
    return fmt.Errorf("transaction failed: %v, %v", err1, err2)
}

txn.Commit()
```

### Deferred Cleanup

```go
// Ensure resources are cleaned up
func processDocument(uri string) error {
    txn := &util.Transaction{}
    
    if !txn.Begin() {
        return fmt.Errorf("failed to begin transaction")
    }
    // Always rollback/commit before return
    defer func() {
        if txn.ID != "" {
            txn.Rollback()
        }
    }()
    
    // Process document...
    resp := documents.ResponseHandle{Format: handle.JSON}
    if err := client.Documents().Read([]string{uri}, nil, nil, txn, &resp); err != nil {
        return err
    }
    
    txn.Commit()
    return nil
}
```

---

## Complete Example: Building a Real Application

Here's a complete example combining multiple services:

```go
package main

import (
    "bytes"
    "fmt"
    "log"
    "time"
    
    marklogic "github.com/ryanjdew/go-marklogic-go"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
    documents "github.com/ryanjdew/go-marklogic-go/documents"
    search "github.com/ryanjdew/go-marklogic-go/search"
    util "github.com/ryanjdew/go-marklogic-go/util"
)

func main() {
    // Connect to MarkLogic
    client, err := marklogic.NewClient(
        "localhost", 8050, "admin", "admin", marklogic.DigestAuth,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Insert documents
    docs := make([]*documents.DocumentDescription, 0)
    for i := 1; i <= 5; i++ {
        doc := &documents.DocumentDescription{
            URI: fmt.Sprintf("/products/%d.json", i),
            Content: bytes.NewBufferString(
                fmt.Sprintf(`{
                    "id": %d,
                    "name": "Product %d",
                    "price": %d,
                    "status": "active"
                }`, i, i, i*100),
            ),
            Metadata: &documents.Metadata{
                Collections: []string{"products", "active"},
                Properties: map[string]string{
                    "created": time.Now().Format(time.RFC3339),
                },
            },
            Format: handle.JSON,
        }
        docs = append(docs, doc)
    }
    
    writeResp := documents.ResponseHandle{Format: handle.JSON}
    if err := client.Documents().Write(docs, nil, nil, &writeResp); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Documents inserted successfully")
    
    // Search documents
    searchResp := search.ResponseHandle{Format: handle.JSON}
    if err := client.Search().Search("active", 1, 10, nil, &searchResp); err != nil {
        log.Fatal(err)
    }
    
    results := searchResp.Get().(*search.Response)
    fmt.Printf("\nSearch Results: Found %d documents\n", results.Total)
    for _, doc := range results.Results {
        fmt.Printf("  - %s (relevance: %.2f)\n", doc.URI, doc.Score)
    }
    
    // Update within transaction
    txn := &util.Transaction{}
    txn.Begin()
    defer func() {
        if txn.ID != "" {
            txn.Rollback()
        }
    }()
    
    updateDocs := []*documents.DocumentDescription{
        {
            URI: "/products/1.json",
            Content: bytes.NewBufferString(`{
                "id": 1,
                "name": "Updated Product 1",
                "price": 150,
                "status": "inactive"
            }`),
            Format: handle.JSON,
        },
    }
    
    updateResp := documents.ResponseHandle{Format: handle.JSON}
    if err := client.Documents().Write(updateDocs, nil, txn, &updateResp); err != nil {
        log.Fatal(err)
    }
    
    // Delete within transaction
    deleteResp := search.ResponseHandle{Format: handle.JSON}
    if err := client.Search().Delete(
        map[string]string{"q": "price:[0 TO 50]"},
        txn,
        &deleteResp,
    ); err != nil {
        log.Fatal(err)
    }
    
    txn.Commit()
    fmt.Println("\nTransactional operations completed successfully")
}
```

For more information on architecture and design patterns, see [ARCHITECTURE.md](ARCHITECTURE.md).
