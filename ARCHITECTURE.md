# Go MarkLogic Client Library - Architecture Guide

## Overview

This document provides a comprehensive guide to the architecture and design patterns of the Go MarkLogic client library. It explains how the library is structured, how to use it effectively, and how to extend it with new functionality.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Components](#core-components)
3. [Service Pattern](#service-pattern)
4. [Handle Interface and Format Abstraction](#handle-interface-and-format-abstraction)
5. [Authentication](#authentication)
6. [Service-to-REST Endpoint Mapping](#service-to-rest-endpoint-mapping)
7. [Common Utilities and Patterns](#common-utilities-and-patterns)
8. [Connection Lifecycle](#connection-lifecycle)
9. [Format Handling](#format-handling)
10. [Advanced Topics](#advanced-topics)
11. [Best Practices](#best-practices)
12. [Adding a New Service](#adding-a-new-service)

## Architecture Overview

The Go MarkLogic client library uses a **service-oriented architecture with three key layers**:

```
┌─────────────────────────────────────────────────────┐
│         Root Package (clientAPI.go)                 │
│    Public Client interface aggregates services      │
└──────────────┬────────────────────────────────────┬─┘
               │                                    │
┌──────────────▼──────┐              ┌─────────────▼─────────────┐
│   Service Layer     │              │   Management APIs         │
│  (search/, docs/)   │              │  (admin/, management/)    │
└──────────────┬──────┘              └─────────────┬─────────────┘
               │                                    │
┌──────────────▼────────────────────────────────────▼─┐
│         Client Layer (clients/)                     │
│  Low-level HTTP handling with pluggable auth       │
└───────────────┬─────────────────────────────────────┘
                │
┌───────────────▼──────────────────────────────────┐
│    Handle Layer (handle/)                        │
│  Format abstraction (JSON/XML/Mixed)            │
│  Serialization/Deserialization bridge           │
└─────────────────────────────────────────────────┘
```

### Design Principles

1. **Separation of Concerns**: Each layer handles a specific responsibility
2. **Format Agnostic**: Services operate through handles that abstract serialization
3. **Composability**: Handles and utilities can be combined for complex operations
4. **Type Safety**: Strong typing with interfaces for extensibility
5. **Testability**: Package-level functions enable testing without mocking

## Core Components

### 1. Root Client (clientAPI.go)

The root package exports the public `Client` type and aggregates all services:

```go
// Client is the main interface for MarkLogic operations
type Client clients.Client

// NewClient creates and connects to MarkLogic
func NewClient(host string, port int64, username string, password string, authType int) (*Client, error)

// Service accessors (thin wrappers)
func (c *Client) Search() *search.Service
func (c *Client) Documents() *documents.Service
func (c *Client) Semantics() *semantics.Service
// ... etc
```

**Key Pattern**: The Client type is a type alias that wraps the internal `clients.Client` type, providing a clean public API while maintaining internal organization.

### 2. Client Layer (clients/)

Handles low-level HTTP communication and authentication:

- **`client.go`**: Main REST client with pluggable authentication
- **`adminClient.go`**: Admin REST client (uses `/admin/v1` endpoints)
- **`managementClient.go`**: Management client (uses `/manage/v2` endpoints)

Each client supports three authentication methods:
- **BasicAuth (0)**: HTTP Basic Authentication
- **DigestAuth (1)**: HTTP Digest Authentication with challenge/response
- **None (2)**: No authentication

### 3. Service Layer

Each service package follows a consistent pattern:

```
service/
├── service.go      # Public Service type with method signatures
├── [service].go    # Implementation functions
└── handles.go      # Request/Response handle types (if needed)
```

**Example**: `search/`
```
search/
├── service.go      # Service struct and method wrappers
├── search.go       # Search implementation + ResponseHandle
├── suggest.go      # Suggestions implementation + SuggestionsResponseHandle
└── query.go        # Query types and builders
```

### 4. Handle Layer (handle/)

Provides unified serialization across formats through the `Handle` interface:

```go
type Handle interface {
    io.ReadWriter
    GetFormat() int                    // JSON, XML, MIXED, etc.
    Serialize(interface{})              // Object → Buffer (encode)
    Deserialize([]byte)                 // Buffer → Object (decode)
    Serialized() string                 // Get raw serialized string
    Deserialized() interface{}          // Get parsed object
}

type ResponseHandle interface {
    Handle
    AcceptResponse(*http.Response) error  // Parse HTTP response
}
```

**Key Pattern**: Handles maintain both serialized (raw) and deserialized (typed) representations, enabling:
- Programmatic access to typed structs
- Access to raw response for debugging/passthrough
- Format negotiation with REST API

## Service Pattern

Every service follows a consistent, tested pattern that separates concerns and enables clean testing:

### Pattern Structure

```go
// Package service provides operations for X
package service

import "github.com/ryanjdew/go-marklogic-go/clients"

// Service struct holds client reference
type Service struct {
    client *clients.Client
}

// NewService factory constructor
func NewService(c *clients.Client) *Service {
    return &Service{client: c}
}

// Public method - thin wrapper delegating to package function
func (s *Service) DoOperation(param string, response handle.ResponseHandle) error {
    return doOperation(s.client, param, response)
}

// Private package function - actual implementation
func doOperation(c *clients.Client, param string, response handle.ResponseHandle) error {
    req, err := util.BuildRequestFromHandle(c, "GET", "/v1/endpoint", nil)
    if err != nil {
        return err
    }
    return util.Execute(c, req, response)
}
```

### Why This Pattern?

1. **Testability**: Package functions can be tested directly without creating Service instances
2. **Reusability**: Package functions can be called from multiple places
3. **Clean API**: Service methods are simple entry points
4. **Consistency**: All services follow the same pattern
5. **Flexibility**: Easy to add convenience methods that compose package functions

### Example: Search Service

```go
// Public method delegates to package function
func (s *Service) Search(text string, start int64, pageLength int64, 
                        transaction *util.Transaction, 
                        response handle.ResponseHandle) error {
    return Search(s.client, text, start, pageLength, transaction, response)
}

// Package function implements the logic
func Search(c *clients.Client, text string, start int64, pageLength int64,
           transaction *util.Transaction,
           response handle.ResponseHandle) error {
    params := "?q=" + url.QueryEscape(text)
    params += "&start=" + strconv.FormatInt(start, 10)
    params += "&pageLength=" + strconv.FormatInt(pageLength, 10)
    params = util.AddTransactionParam(params, transaction)
    
    req, err := http.NewRequest("GET", c.Base()+"/search"+params, nil)
    if err != nil {
        return err
    }
    return util.Execute(c, req, response)
}
```

## Handle Interface and Format Abstraction

The `Handle` interface is central to the library's design, providing format abstraction between client code and MarkLogic REST API.

### Format Constants

Handles map internal format enums to MarkLogic REST API MIME types:

```go
const (
    JSON       = iota     // application/json
    XML                   // application/xml
    MIXED                 // multipart/mixed (bulk operations)
    TEXTPLAIN             // text/plain
    TEXT_URI_LIST         // text/uri-list
    TEXTHTML              // text/html
    UNKNOWN               // application/octet-stream
)

// Conversion function used by util.Execute and util.BuildRequestFromHandle
func FormatEnumToMimeType(formatEnum int) string {
    switch formatEnum {
    case JSON:
        return "application/json"
    case XML:
        return "application/xml"
    // ...
    }
}
```

### Dual Representation Pattern

The key innovation of the Handle interface is maintaining both representations:

```go
// Example: ResponseHandle from search package
type ResponseHandle struct {
    Format int
    buffer bytes.Buffer        // Raw serialized response
    info   *Response           // Parsed typed response
}

// Serialized access (debugging, passthrough)
func (rh *ResponseHandle) Serialized() string {
    return rh.buffer.String()
}

// Deserialized access (programmatic usage)
func (rh *ResponseHandle) Deserialized() interface{} {
    return rh.info  // Returns *search.Response
}
```

### Handle Lifecycle

1. **Creation**: Create handle with desired format
```go
respHandle := search.ResponseHandle{Format: handle.JSON}
```

2. **Usage**: Pass to service method
```go
err := client.Search().Search("query", 1, 10, nil, &respHandle)
```

3. **Access**: Get results as needed
```go
// Typed access
results := respHandle.Get().(*search.Response)

// Raw access
jsonString := respHandle.Serialized()
```

## Authentication

The library supports three authentication strategies implemented in the clients layer:

### Basic Authentication

Username and password encoded in HTTP headers with each request:

```go
client, err := marklogic.NewClient(host, port, user, pass, marklogic.BasicAuth)
// Subsequent requests include: Authorization: Basic <base64(user:pass)>
```

**Pros**: Simple, stateless
**Cons**: Credentials sent with every request (use HTTPS only)

### Digest Authentication

Challenge-response authentication with cached credentials:

```go
client, err := marklogic.NewClient(host, port, user, pass, marklogic.DigestAuth)
```

**Implementation Details**:
1. On first connection, client fetches digest challenge from server
2. Challenge is cached in thread-safe `sync.RWMutex`
3. Each request computes response hash based on cached challenge
4. Server challenges trigger cache update

**Pros**: Credentials not sent plaintext
**Cons**: Stateful, requires initial handshake

### No Authentication

For open MarkLogic instances:

```go
client, err := marklogic.NewClient(host, port, "", "", marklogic.None)
```

### Authentication Lifecycle

```go
// 1. Connection initialization
client, err := marklogic.NewClient("localhost", 8050, "admin", "password", marklogic.DigestAuth)

// 2. On first request, util.Execute calls clients.ApplyAuth
// 3. For DigestAuth: Challenge is fetched and cached
// 4. Subsequent requests: Response hash computed from cached challenge
// 5. If server challenges: Cache is updated, request retried
```

## Service-to-REST Endpoint Mapping

The library maps Go services to MarkLogic REST API endpoints (base path: `/LATEST`):

| Service | Package | REST Endpoints | Purpose |
|---------|---------|---|---|
| **Search** | `search/` | `/v1/search` (GET/POST/DELETE) | Full-text & structured search with faceting, suggestions |
| **Documents** | `documents/` | `/v1/documents` (GET/PUT/POST/DELETE/PATCH) | CRUD for single/bulk documents + metadata |
| **Semantics** | `semantics/` | `/v1/graphs`, `/v1/graphs/sparql`, `/v1/graphs/things` | Triple store, SPARQL queries, semantic operations |
| **Config** | `config/` | `/v1/config/query`, `/v1/config/transforms`, `/v1/config/properties` | Query options, transforms, extensions management |
| **Resources** | `resources/` | `/v1/resources/{name}` (GET/PUT/POST/DELETE) | User-defined resource extensions |
| **Values** | `values/` | `/v1/values/{name}` (GET/POST) | Lexicon enumeration, distinct values, aggregation |
| **Data Movement** | `datamovement/` | `/v1/documents` (bulk batching) | Optimized batch write operations |
| **Rows** | `rows-management/` | `/v1/rows` (GET/POST) | Optic DSL for structured row queries |
| **Data Services** | `dataservices/` | `/v1/invoke` | Server-side module evaluation |
| **Alerting** | `alert/` | `/v1/alert/rules`, `/v1/alert/match` | Alert rule matching |
| **Temporal** | `temporal/` | `/v1/documents` + temporal params | Document versioning and time-based queries |
| **Indexes** | `indexes/` | `/v1/config/indexes` | Range and field index management |
| **Metadata** | `metadata/` | `/v1/metadata/{uri}` | Metadata extraction and validation |
| **Transactions** | `transactions/` | `/v1/transactions` | Multi-statement transaction management |
| **Admin** | `admin/` | `/admin/v1` | Server initialization and configuration |
| **Management** | `management/` | `/manage/v2` | Cluster and resource management |

## Common Utilities and Patterns

The `util/` package provides critical utilities for building and executing requests:

### HTTP Request/Response Utilities

**`BuildRequestFromHandle(c, method, uri, handle)`**
- Creates HTTP request with `Content-Type` header set from handle's format
- Uses `FormatEnumToMimeType()` to convert format enum to MIME type
- Returns `*http.Request` ready for util.Execute

**`Execute(c, req, responseHandle)`**
- Executes request with automatic authentication via `clients.ApplyAuth`
- Sets `Accept` header from responseHandle format
- Deserializes response via handle's `AcceptResponse()`
- Treats status ≥400 as errors
- Returns error if deserialization fails

### URL Parameter Building

The library uses string concatenation for parameter building (not `net/url.QueryBuilder`) because MarkLogic requires specific parameter ordering and non-standard encoding:

**`AddDatabaseParam(params, client)`**
- Appends `?database=X` if client has non-empty database
- Example: `"?database=Documents"`

**`AddTransactionParam(params, transaction)`**
- Appends `?txid=X` for multi-statement transactions
- Auto-begins transaction if not started
- Example: `"?txid=12345678"`

**`RepeatingParameters(params, label, []values)`**
- Builds repeated query parameters for multi-valued options
- Example: `&category=content&category=metadata`

**`MappedParameters(params, prefix, map[k]v)`**
- Builds mapped parameters with prefix
- Example: `?pref:key1=val1&pref:key2=val2`

### Types

**`Transaction`**
- Holds transaction ID (`txid`)
- `Begin()` creates server-side transaction
- `Commit()` commits all operations
- `Rollback()` aborts all operations
- See util/transaction.go for implementation

**`Transform`**
- Wraps server-side transform name and parameters
- `ToParameters()` serializes to query string
- Example: `"trans:extract-metadata?parm1=value1&parm2=value2"`

## Connection Lifecycle

The complete lifecycle of a connection and request:

```go
// 1. INITIALIZATION
client, err := marklogic.NewClient(
    "localhost",               // host
    8050,                      // REST API port
    "admin",                   // username
    "admin",                   // password
    marklogic.DigestAuth,      // auth type
)
if err != nil {
    log.Fatal(err)
}

// Internally:
// - Creates clients.Client with Connection struct
// - For DigestAuth: Fetches and caches digest challenge
// - For BasicAuth: Stores base64-encoded credentials
// - Root Client type wraps internal client

// 2. SERVICE ACCESS
searchService := client.Search()  // Gets search.Service with embedded client

// 3. OPERATION EXECUTION
respHandle := search.ResponseHandle{Format: handle.JSON}
err := searchService.Search("query", 1, 10, nil, &respHandle)

// Internally:
// a) Search() method calls Search(s.client, ...)
// b) Search() package function builds request URL
// c) util.BuildRequestFromHandle creates HTTP request with Content-Type
// d) util.Execute:
//    i. Calls clients.ApplyAuth to add auth header
//    ii. Sets Accept header from responseHandle format
//    iii. Sends HTTP request
//    iv. Calls respHandle.AcceptResponse(httpResp)
//    v. Returns any errors

// 4. RESULT ACCESS
results := respHandle.Get()           // *search.Response
jsonStr := respHandle.Serialized()    // Raw JSON string
```

## Format Handling

Format handling connects client code to MarkLogic REST API through Content-Type/Accept negotiation:

### Request Flow (Serialization)

```go
// 1. Create handle with desired request format
query := search.Query{...}
qh := search.QueryHandle{Format: handle.XML}

// 2. Serialize object to handle
qh.Serialize(query)
// Internally: Marshals query struct to XML string

// 3. Build request - util.BuildRequestFromHandle adds Content-Type
req, err := util.BuildRequestFromHandle(client, "POST", "/v1/search", &qh)
// Internally: 
//   - Gets format from qh.GetFormat() → handle.XML
//   - Converts to MIME type: "application/xml"
//   - Sets req.Header["Content-Type"] = "application/xml"
//   - Sets request body from qh.buffer (serialized content)

// 4. Execute request
util.Execute(client, req, respHandle)
```

### Response Flow (Deserialization)

```go
// 1. Create response handle with desired format
respHandle := search.ResponseHandle{Format: handle.JSON}

// 2. Execute request - util.Execute handles deserialization
util.Execute(client, req, &respHandle)
// Internally:
//   - Sets Accept header: "application/json"
//   - Sends request
//   - Gets response from server
//   - Calls respHandle.AcceptResponse(httpResponse)
//   - AcceptResponse unmarshals JSON into *search.Response
//   - Also buffers raw response bytes

// 3. Access results both ways
typedResults := respHandle.Get()      // *search.Response (parsed)
rawJSON := respHandle.Serialized()    // string (raw JSON)
```

### Multipart/Mixed for Bulk Operations

For bulk read operations returning multiple documents with metadata:

```go
// Request with multiple URIs
uris := []string{"/doc1.xml", "/doc2.xml"}
respHandle := documents.ResponseHandle{Format: handle.MIXED}

client.Documents().Read(uris, []string{"content", "metadata"}, nil, nil, &respHandle)
// Server returns multipart/mixed with:
// --boundary
// Content-Disposition: inline; category=content
// Content-Type: application/xml
// <doc1>...</doc1>
// --boundary
// Content-Disposition: inline; category=metadata
// <metadata>...</metadata>
// --boundary
```

The ResponseHandle parses this into structured data with proper categorization.

## Advanced Topics

### Custom Handles

Create custom handles for specialized serialization needs:

```go
type CustomHandle struct {
    Format int
    buffer bytes.Buffer
    data   interface{}
}

func (ch *CustomHandle) GetFormat() int {
    return ch.Format
}

func (ch *CustomHandle) Serialize(v interface{}) {
    // Custom serialization logic
    data := json.Marshal(v)
    ch.buffer.Write(data)
}

func (ch *CustomHandle) Deserialize(b []byte) {
    // Custom deserialization logic
    json.Unmarshal(b, &ch.data)
}

// ... implement other required methods
```

### Multi-Statement Transactions

Execute multiple operations atomically:

```go
txn := &util.Transaction{}
txn.Begin()

// All operations use same transaction
err1 := client.Documents().Write([]*documents.DocumentDescription{...}, nil, txn, respHandle1)
err2 := client.Search().Search("query", 1, 10, txn, respHandle2)

if err1 == nil && err2 == nil {
    txn.Commit()
} else {
    txn.Rollback()
}
```

### Custom REST Extensions

Call user-defined REST extensions:

```go
// Assumes server-side extension at /custom-extension
resp := resources.ResponseHandle{Format: handle.JSON}
err := client.Resources().Get("custom-extension", nil, &resp)

// Leverage RESTful interface
err = client.Resources().Post("custom-extension", &requestHandle, &resp)
```

### Server-Side Code Execution

Execute XQuery or JavaScript with external variables:

```go
// Data Services: Call named module with parameters
vars := map[string]string{"param1": "value1"}
resp := dataservices.ResponseHandle{Format: handle.JSON}
err := client.DataServices().Invoke("my-module", vars, nil, &resp)

// Direct evaluation with Eval service
code := `xdmp:database-name(xdmp:database())`
resp := eval.ResponseHandle{Format: handle.JSON}
err := client.Eval().Xquery(code, nil, &resp)
```

### Semantic Graph Operations

Work with RDF triples and SPARQL:

```go
// Execute SPARQL query
sparqlQuery := `SELECT ?s ?p ?o WHERE { ?s ?p ?o }`
resp := semantics.ResponseHandle{Format: handle.JSON}
err := client.Semantics().Sparql(sparqlQuery, &resp)

// Get semantic information for IRIs
iris := []string{"http://example.com/resource/1"}
resp := semantics.ResponseHandle{Format: handle.JSON}
err := client.Semantics().Things(iris, &resp)
```

### Optic Queries

Structured queries on relational data:

```go
// Build Optic query plan (DSL)
plan := `op:col("MySchema", "MyTable") 
    |> op:where(op:eq(op:col("status"), "active")) 
    |> op:limit(10)`

resp := rowsManagement.ResponseHandle{Format: handle.JSON}
err := client.RowsManagement().Rows(planHandle, &resp)
```

## Best Practices

### 1. Connection Management

```go
// Create once, reuse across application lifetime
client, err := marklogic.NewClient(host, port, user, pass, auth)
if err != nil {
    log.Fatal(err)
}
// client persists for application lifetime
// No explicit Close() needed
```

### 2. Error Handling

```go
// Always check errors at each step
respHandle := search.ResponseHandle{Format: handle.JSON}
if err := client.Search().Search(query, 1, 10, nil, &respHandle); err != nil {
    // Handle network errors, auth failures, 4xx/5xx responses
    log.Printf("Search failed: %v", err)
    return
}

// Verify response deserialization succeeded
results := respHandle.Get()
if results == nil {
    log.Printf("Response deserialization failed")
    return
}
```

### 3. Format Selection

```go
// Choose format based on use case:

// For programmatic processing: Use structured handle
resp := search.ResponseHandle{Format: handle.JSON}
results := resp.Get().(*search.Response)  // Type-safe

// For debugging/logging: Get raw response
rawResponse := resp.Serialized()  // Full JSON/XML string

// For format negotiation: Let handle manage both
// Handles maintain both representations internally
```

### 4. Concurrent Requests

```go
// Services are goroutine-safe; client connection can be shared
go func() {
    respHandle := search.ResponseHandle{Format: handle.JSON}
    client.Search().Search("query1", 1, 10, nil, &respHandle)
}()

go func() {
    respHandle := documents.ResponseHandle{Format: handle.JSON}
    client.Documents().Read([]string{"/doc"}, nil, nil, nil, &respHandle)
}()

// Each goroutine has independent response handle and results
```

### 5. Transaction Usage

```go
// Transactions are optional - use when atomic operations needed
txn := &util.Transaction{}

// Begin must be called before operations
if !txn.Begin() {
    return fmt.Errorf("failed to begin transaction")
}

// Use same txn instance for all related operations
op1 := client.Documents().Write(docs1, nil, txn, resp1)
op2 := client.Search().Delete(criteria, txn, resp2)

// Commit only if all operations succeeded
if op1 == nil && op2 == nil {
    txn.Commit()
} else {
    txn.Rollback()
}
```

### 6. Parameter Building

```go
// Use util parameter functions for consistency
params := ""

// Add database targeting
params = util.AddDatabaseParam(params, client)

// Add transaction
params = util.AddTransactionParam(params, txn)

// Add repeated parameters
params = util.RepeatingParameters(params, "category", 
    []string{"content", "metadata"})

// Resulting URL: "?database=mydb&txid=12345&category=content&category=metadata"
```

### 7. Handle Patterns

```go
// Pattern 1: Reusable handle instances
var jsonHandle = handle.JSONHandle{Format: handle.JSON}

// Pattern 2: Package-specific handles for typed access
var searchResp = search.ResponseHandle{Format: handle.JSON}
results := searchResp.Get() // *search.Response

// Pattern 3: Format flexibility with same handle
handle1 := CustomHandle{Format: handle.JSON}
handle2 := CustomHandle{Format: handle.XML}
// Same structure, different formats
```

### 8. Metadata Operations

```go
// Set metadata on documents
metadata := &documents.Metadata{
    Collections: []string{"collection1", "collection2"},
    Properties: map[string]string{"status": "active"},
}

doc := &documents.DocumentDescription{
    URI:      "/doc",
    Content:  &contentBuffer,
    Metadata: metadata,
}

client.Documents().Write([]*documents.DocumentDescription{doc}, nil, nil, nil, resp)
```

### 9. Error Recovery

```go
// Retry logic for transient failures
var resp search.ResponseHandle

for attempt := 0; attempt < 3; attempt++ {
    err := client.Search().Search(query, 1, 10, nil, &resp)
    if err == nil {
        break
    }
    if attempt < 2 {
        time.Sleep(time.Second * time.Duration(math.Pow(2, float64(attempt))))
    }
}
```

## Adding a New Service

When adding a new service to the library:

### Step 1: Create Package Structure

```
newservice/
├── service.go       # Public Service type
├── newservice.go    # Implementation functions
└── handles.go       # Custom handles (if needed)
```

### Step 2: Implement Service Struct

```go
// newservice/service.go
package newservice

import "github.com/ryanjdew/go-marklogic-go/clients"

// Service provides operations for [description]
type Service struct {
    client *clients.Client
}

// NewService creates a new newservice.Service
func NewService(c *clients.Client) *Service {
    return &Service{client: c}
}

// OperationName description
func (s *Service) OperationName(param string, response handle.ResponseHandle) error {
    return operationName(s.client, param, response)
}
```

### Step 3: Implement Package Functions

```go
// newservice/newservice.go
package newservice

import (
    "github.com/ryanjdew/go-marklogic-go/clients"
    "github.com/ryanjdew/go-marklogic-go/util"
    "github.com/ryanjdew/go-marklogic-go/handle"
)

// operationName performs the actual operation
func operationName(c *clients.Client, param string, response handle.ResponseHandle) error {
    // Build request
    req, err := util.BuildRequestFromHandle(c, "GET", "/v1/endpoint", nil)
    if err != nil {
        return err
    }
    
    // Execute with response handle
    return util.Execute(c, req, response)
}
```

### Step 4: Create Response Handle (if needed)

```go
// newservice/handles.go
package newservice

import (
    "bytes"
    "encoding/json"
    "net/http"
    handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// ResponseHandle provides a handle for service responses
type ResponseHandle struct {
    Format int
    buffer bytes.Buffer
    info   *Response
}

func (rh *ResponseHandle) GetFormat() int { return rh.Format }

func (rh *ResponseHandle) Serialize(v interface{}) {
    data, _ := json.Marshal(v)
    rh.buffer.Write(data)
}

func (rh *ResponseHandle) Deserialize(b []byte) {
    json.Unmarshal(b, &rh.info)
}

// ... implement other required Handle methods
```

### Step 5: Register Service in Root Client

```go
// clientAPI.go
package gomarklogicgo

// NewService method
func (c *Client) NewService() *newservice.Service {
    return newservice.NewService(convertToSubClient(c))
}
```

### Step 6: Write Tests

Create `newservice_test.go` with:
- Unit tests for package functions
- Integration tests for REST API calls
- Handle tests for serialization/deserialization

### Step 7: Update Documentation

- Add package documentation comment
- Document public Service methods
- Add examples to README.md
- Update this architecture guide

## Conclusion

The Go MarkLogic client library's architecture balances simplicity with power. By understanding these patterns and design decisions, you can:

1. Use the library effectively and idiomatically
2. Extend it with new services following established patterns
3. Debug issues through clear separation of concerns
4. Write maintainable, testable code that integrates with MarkLogic

For specific examples and API reference, see `API_REFERENCE.md`.
