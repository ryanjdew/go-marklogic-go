# Go MarkLogic Client Library - Copilot Instructions

## Architecture Overview

This is a Go REST client library for MarkLogic. The codebase uses a **service-oriented architecture with three key layers**:

1. **Client Layer** (`clients/`): Low-level HTTP handling with pluggable auth (Basic/Digest)
2. **Service Layer** (`search/`, `documents/`, `semantics/`, etc.): Domain-specific operations
3. **Handle Layer** (`handle/`): Serialization abstraction for JSON/XML/mixed formats

The root package (`clientAPI.go`) exposes a public Client interface that aggregates all services.

## Service Pattern

Each service follows this consistent pattern:

```go
// Package structure: pkg/service.go & pkg/pkg.go
type Service struct {
    client *clients.Client
}

func NewService(c *clients.Client) *Service {
    return &Service{client: c}
}

// Public methods delegate to package-level functions
func (s *Service) Search(text string, ..., response handle.ResponseHandle) error {
    return Search(s.client, text, ...)
}

// Implementation in pkg.go handles actual logic
func Search(c *clients.Client, ...) error { ... }
```

**Key pattern:** Service methods are thin wrappers; real logic lives in package-level functions. This enables testing without mocking and keeps concerns separated.

## Handle Interface - Format Abstraction

Handles provide unified serialization across formats (JSON/XML/mixed). They map to MarkLogic REST API `Content-Type` and `Accept` headers:

```go
// In handle/handles.go - constants for formats (map to MIME types)
const (
    JSON = iota              // application/json
    XML                      // application/xml
    MIXED                    // multipart/mixed (for bulk operations)
    TEXTPLAIN                // text/plain
    TEXT_URI_LIST            // text/uri-list
    TEXTHTML                 // text/html
    UNKNOWN                  // application/octet-stream
)

// FormatEnumToMimeType converts format enum to REST API MIME type
func FormatEnumToMimeType(formatEnum int) string { ... }

// All handles implement this interface
type Handle interface {
    io.ReadWriter
    GetFormat() int
    Serialize(interface{})        // object → buffer (encode to format)
    Deserialize([]byte)           // buffer → object (decode from format)
    Serialized() string           // get raw serialized string
    Deserialized() interface{}    // get parsed object
}

type ResponseHandle interface {
    Handle
    AcceptResponse(*http.Response) error  // parse HTTP response body
}
```

**Key pattern:** Handles wrap both serialized (raw string) and deserialized (typed object) representations. Example in `search/search.go`: `ResponseHandle` unmarshals into a typed `Response` struct (via JSON/XML unmarshaling) while also buffering raw response bytes. This enables both programmatic access to typed structs AND access to the raw response format for debugging or passthrough scenarios.

## Service-to-REST Endpoint Mapping

The library's services map to MarkLogic REST API endpoints (base path: `/LATEST`):

| Service | Package | REST Endpoints | Purpose |
|---------|---------|---|---|
| **Search** | `search/` | `/v1/search` (GET/POST/DELETE) | Full-text & structured search with faceting, suggestions |
| **Documents** | `documents/` | `/v1/documents` (GET/PUT/POST/DELETE/PATCH) | CRUD for single/bulk documents + metadata |
| **Semantics** | `semantics/` | `/v1/graphs` (GET/PUT/POST/DELETE), `/v1/graphs/sparql`, `/v1/graphs/things` | Triple store, SPARQL queries, semantic operations |
| **Config** | `config/` | `/v1/config/query`, `/v1/config/transforms`, `/v1/config/properties` | Query options, transforms, extensions management |
| **Resources** | `resources/` | `/v1/resources/{name}` (GET/PUT/POST/DELETE) | User-defined resource extensions |
| **DataMovement** | `datamovement/` | `/v1/documents` (bulk batching) | Optimized batch write operations |
| **RowsManagement** | `rows-management/` | `/v1/rows` (GET/POST) | Optic DSL for structured row queries |
| **DataServices** | `dataservices/` | `/v1/invoke` | Server-side module evaluation |
| **Alerting** | `alert/` | `/v1/alert/rules`, `/v1/alert/match` | Alert rule matching |

**Note:** Admin operations (`admin/`, `management/`) use `/admin/v1` and `/manage/v2` endpoints for server initialization and cluster configuration.

## Critical Utilities (util.go)

**HTTP Request/Response Utilities:**
- `BuildRequestFromHandle(c, method, uri, handle)` - Create HTTP request with `Content-Type` header set from handle's format via `FormatEnumToMimeType()`
- `Execute(c, req, responseHandle)` - Execute request with auth, set `Accept` header from responseHandle format, deserialize response via handle's `AcceptResponse()`, check status ≥400 as errors

**URL Parameter Builders:**
- `AddDatabaseParam(params, client)` - Appends `?database=X` if client has non-empty database
- `AddTransactionParam(params, transaction)` - Appends `?txid=X` for multi-statement transactions; auto-begins transaction if not started
- `RepeatingParameters(params, label, []values)` - Builds repeated query params: `&label=val1&label=val2`
- `MappedParameters(params, prefix, map[k]v)` - Builds mapped params: `?prefix:key1=val1&prefix:key2=val2`

**Types:**
- `Transaction` - Holds transaction ID (`txid`); `Begin()` creates server-side transaction
- `Transform` - Wraps server-side transform name and parameters; `ToParameters()` serializes to query string

**Pattern:** All request building uses string concatenation (not `net/url` Query builder) because parameters need specific ordering and non-standard encoding for MarkLogic (e.g., repeated params, map prefixes).

## Build & Test Workflow

Uses **Task** (Taskfile.yml) for setup automation:

```bash
# Full environment setup (Docker + MarkLogic instance)
task setup

# Build (no deps installed by build task)
task build

# Integration tests (requires MarkLogic running)
task test
```

Tests are tagged `//+build integration` - they require live MarkLogic server. Test files follow `*_test.go` pattern.

## When Adding a New Service

1. Create `newservice/` directory with `service.go` and `newservice.go`
2. Define Service struct with embedded `*clients.Client`
3. Service methods delegate to package functions in `newservice.go`
4. Implement in `newservice.go` using `util.BuildRequestFromHandle()` and `util.Execute()`
5. Add Factory method to root Client in `clientAPI.go`: `func (c *Client) NewService() *newservice.Service { ... }`

## Connection Lifecycle

```go
// Root package creates client
client, err := marklogic.NewClient(host, port, user, pass, authType)

// Gets wrapped client for each service
func (c *Client) Search() *search.Service {
    return search.NewService(convertToSubClient(c))
}
```

The wrapping is necessary because root package exports `Client` as a type alias, not embedded type.

## Format Handling Pattern

Handles connect to MarkLogic REST API via `Content-Type` (requests) and `Accept` (responses):

```go
// 1. SERIALIZE: Build request - struct → handle → HTTP Content-Type header
query := search.Query{
    Queries: []any{search.TermQuery{Terms: []string{"text"}}},
}
qh := search.QueryHandle{Format: handle.XML}  // Set desired format
qh.Serialize(query)  // Encodes query struct to XML string

// 2. EXECUTE: util.BuildRequestFromHandle adds Content-Type from handle
req, _ := util.BuildRequestFromHandle(client, "POST", "/search", &qh)
// Sets: req.Header["Content-Type"] = "application/xml"

// 3. RESPONSE: Deserialize from Accept header (multipart/mixed for bulk ops)
respHandle := search.ResponseHandle{Format: handle.JSON}
util.Execute(client, req, &respHandle)  // Sets Accept: application/json

// 4. ACCESS: Both serialized (raw) and deserialized (typed) forms available
fmt.Println(respHandle.Serialized())     // raw JSON response string
resp := respHandle.Deserialized()        // typed *search.Response struct
```

**MarkLogic REST mapping:** The format constant (`JSON`/`XML`/`MIXED`) directly translates to REST API `Content-Type`/`Accept` values via `FormatEnumToMimeType()`. Bulk operations use `MIXED` format for `multipart/mixed` responses containing multiple documents with metadata.

## Authentication

`clients/client.go` supports three auth types (configured at connection time):
- **BasicAuth (0)**: HTTP Basic Auth - encoded credentials sent with every request
- **DigestAuth (1)**: Digest Auth with challenge/response - requires pre-fetching digest headers from `/config/resources?format=xml` on connection; uses `sync.RWMutex` for thread-safe header updates
- **None (2)**: No authentication

**Pattern:** Authentication is applied via `clients.ApplyAuth(c, req)` in `util.Execute()` before sending each request. Digest auth is stateful—the library caches digest challenge/response headers to avoid re-authenticating per request while handling server challenges.

**Security Note:** Both BasicAuth and DigestAuth transmit credentials; use HTTPS in production. Credentials are stored in `Connection` struct during client initialization.

## Common Conventions

- **Error handling**: Functions return `error` as last return value, no panics expected
- **Nil checks**: ResponseHandle is often optional (checked with `!= nil`)
- **Goroutines in write**: `documents.write()` spawns goroutines per document with channel coordination
- **XML namespaces**: Some types use constant namespaces (e.g., `searchNamespace = "http://marklogic.com/appservices/search"`)
- **Parameters**: Query params are built as strings and appended to URLs, never using http.URL query builder
- **Multipart/mixed responses**: Set `Accept: multipart/mixed` for bulk read operations (multiple documents or document + metadata). Response parts include `Content-Disposition` headers with filename and category.
- **Metadata categories**: Document operations support `content`, `metadata`, `collections`, `permissions`, `properties`, `quality`, `metadata-values`; can specify multiple via repeated query params

## Code Quality Notes

- Imports use full import paths (not relative)
- Package-level functions never conflict with exported methods
- Tests import test-specific helpers from `test/` directory
- No external dependencies except `http-digest-auth-client` and `go-spew`
