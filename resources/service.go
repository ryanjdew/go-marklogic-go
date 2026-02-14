// Package resources provides access to custom REST extensions defined on the
// MarkLogic server. Extensions are user-defined REST endpoints that provide
// domain-specific functionality beyond standard document and search operations.
//
// Resources Service enables:
//   - GET requests against custom resource endpoints
//   - POST requests with data to custom endpoints
//   - PUT requests for resource updates
//   - DELETE requests for resource cleanup
//   - Query parameter passing to extensions
//   - Format negotiation (JSON/XML) with endpoints
//
// Example: Call custom GET resource
//
//	respHandle := resources.ResponseHandle{Format: handle.JSON}
//	client.Resources().Get("my-analysis",
//	    map[string]string{"format": "summary"}, &respHandle)
package resources

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service provides methods for calling custom REST resource extensions.
// All operations delegate to the underlying HTTP client for transparent
// communication with server-side resource implementations.
type Service struct {
	client *clients.Client
}

// NewService creates and returns a new resources.Service instance for invoking
// custom REST extensions and resources defined on the MarkLogic server.
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Get invokes a custom GET resource endpoint. Used for read-only operations
// against server-side resources. Query parameters are passed to the resource
// implementation on the server.
//
// Parameters:
//
//	resourceName: Name of the custom resource to call
//	parameters: Query parameters passed to resource (e.g., map{"id": "123"})
//	response: ResponseHandle to populate with results
func (s *Service) Get(resourceName string, parameters map[string]string, response handle.ResponseHandle) error {
	return get(s.client, resourceName, parameters, response)
}

// Post invokes a custom POST resource endpoint with request body data.
// Used for operations that modify server state or require input data.
// Request format is determined by requestBody handle format.
//
// Parameters:
//
//	resourceName: Name of the custom resource to call
//	parameters: Query parameters for resource
//	requestBody: Handle containing POST request body data
//	response: ResponseHandle to populate with results
func (s *Service) Post(resourceName string, parameters map[string]string, requestBody handle.Handle, response handle.ResponseHandle) error {
	return post(s.client, resourceName, parameters, requestBody, response)
}

// Put Call PUT against resource
func (s *Service) Put(resourceName string, parameters map[string]string, requestBody handle.Handle, response handle.ResponseHandle) error {
	return put(s.client, resourceName, parameters, requestBody, response)
}

// Delete Call DELETE against resource
func (s *Service) Delete(resourceName string, parameters map[string]string, response handle.ResponseHandle) error {
	return delete(s.client, resourceName, parameters, response)
}
