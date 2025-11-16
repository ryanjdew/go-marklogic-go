// Package config provides server configuration management including query options,
// transforms, indexes, and custom REST extensions. It enables programmatic configuration
// of server-side behaviors, query processing rules, and extension deployment.
//
// Config Service enables:
//   - Query options (search configuration with faceting, fields, constraints)
//   - Server-side transforms (data transformation during read/write)
//   - Range and field indexes for query optimization
//   - Custom REST extensions and resources deployment
//   - Namespace management
//   - Server properties and configuration retrieval
//
// Example: Create query options
//
//	opts := `<search:options>...</search:options>`
//	optHandle := handle.Handle{Format: handle.XML}
//	optHandle.Serialize(opts)
//	client.Config().SetQueryOptions("my-options", &optHandle, respHandle)
package config

import (
	"io"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service provides methods for managing MarkLogic server configuration.
// All configuration operations are persisted server-side and affect subsequent
// query and document operations.
type Service struct {
	client *clients.Client
}

// NewService creates and returns a new config.Service instance for server
// configuration management operations.
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// ListExtensions retrieves all installed REST extensions in the specified directory.
// Extensions include both built-in and custom REST service implementations.
//
// Parameters:
//
//	directory: Directory path to list extensions from (e.g., "/")
//	response: ResponseHandle to populate with extension list
func (s *Service) ListExtensions(directory string, response handle.ResponseHandle) error {
	return listExtensions(s.client, directory, response)
}

// DeleteExtensions removes all REST extensions in the specified directory path.
// Exercise caution as this permanently removes extension configurations.
//
// Parameters:
//
//	directory: Directory path for extensions to delete
func (s *Service) DeleteExtensions(directory string) error {
	return deleteExtensions(s.client, directory)
}

// CreateExtension installs or updates a REST service extension. Extensions can be
// written in XQuery or JavaScript and provide custom REST API endpoints.
//
// Parameters:
//
//	assetName: Name/path for the extension (e.g., "my-service")
//	resource: Reader containing extension source code
//	extensionType: Extension language ("xqy" for XQuery or "js" for JavaScript)
//	options: Configuration options for the extension
//	response: ResponseHandle for confirmation
func (s *Service) CreateExtension(assetName string, resource io.Reader, extensionType string, options map[string]string, response handle.ResponseHandle) error {
	return createExtension(s.client, assetName, resource, extensionType, options, response)

}

// ListResources lists all installed REST service resources and their metadata.
//
// Parameters:
//
//	response: ResponseHandle to populate with resource list
func (s *Service) ListResources(response handle.ResponseHandle) error {
	return listResources(s.client, response)
}

// GetResourceInfo retrieves metadata for a specific REST resource including its
// definition, associated code, and configuration.
//
// Parameters:
//
//	name: Resource name/identifier
//	response: ResponseHandle to populate with resource metadata
func (s *Service) GetResourceInfo(name string, response handle.ResponseHandle) error {
	return getResourceInfo(s.client, name, response)
}

// CreateResource deploys a new REST service resource to the server. Resources
// provide endpoints for custom operations beyond standard CRUD.
//
// Parameters:
//
//	name: Resource name (becomes part of REST endpoint)
//	resource: Reader containing resource source code
//	extensionType: Language ("xqy" or "js")
//	options: Configuration options
//	response: ResponseHandle for confirmation
func (s *Service) CreateResource(name string, resource io.Reader, extensionType string, options map[string]string, response handle.ResponseHandle) error {
	return createResource(s.client, name, resource, extensionType, options, response)
}

// DeleteResource removes a REST service
func (s *Service) DeleteResource(name string, response handle.ResponseHandle) error {
	return deleteResource(s.client, name, response)
}

// IndexesReport shows the status of indexes in query options
func (s *Service) IndexesReport(response handle.ResponseHandle) error {
	return indexesReport(s.client, response)
}

// ListNamespaces shows the namespaces used in queries
func (s *Service) ListNamespaces(response handle.ResponseHandle) error {
	return listNamespaces(s.client, response)
}

// SetNamespace shows the namespaces used in queries
func (s *Service) SetNamespace(namespace handle.Handle, response handle.ResponseHandle) error {
	return setNamespace(s.client, namespace, response)
}

// GetProperties shows the REST API properties
func (s *Service) GetProperties(response handle.ResponseHandle) error {
	return getProperties(s.client, response)
}

// SetProperties sets the REST API properties
func (s *Service) SetProperties(properties handle.Handle, response handle.ResponseHandle) error {
	return setProperties(s.client, properties, response)
}

// ResetProperties resets the REST API properties to their default
func (s *Service) ResetProperties(response handle.ResponseHandle) error {
	return resetProperties(s.client, response)
}

// SetPropertyValue sets a property of the REST API
func (s *Service) SetPropertyValue(propertyName string, property handle.Handle, response handle.ResponseHandle) error {
	return setPropertyValue(s.client, propertyName, property, response)
}

// ListQueryOptions shows all the installed REST query options
func (s *Service) ListQueryOptions(response handle.ResponseHandle) error {
	return listQueryOptions(s.client, response)
}

// DeleteAllQueryOptions removes all the installed REST query options
func (s *Service) DeleteAllQueryOptions(response handle.ResponseHandle) error {
	return deleteAllQueryOptions(s.client, response)
}

// SetQueryOptions shows all the installed REST extensions
func (s *Service) SetQueryOptions(optionsName string, options handle.Handle, response handle.ResponseHandle) error {
	return setQueryOptions(s.client, optionsName, options, response)
}

// GetQueryOptions returns the named REST query options
func (s *Service) GetQueryOptions(name string, response handle.ResponseHandle) error {
	return getQueryOptions(s.client, name, response)
}

// DeleteQueryOptions removes the named REST query options
func (s *Service) DeleteQueryOptions(name string, response handle.ResponseHandle) error {
	return deleteQueryOptions(s.client, name, response)
}

// ListTransforms shows all the installed REST service extensions
func (s *Service) ListTransforms(response handle.ResponseHandle) error {
	return listTransforms(s.client, response)
}

// GetTransformInfo shows all the installed REST extensions
func (s *Service) GetTransformInfo(name string, response handle.ResponseHandle) error {
	return getTransformInfo(s.client, name, response)
}

// CreateTransform installs a REST service
func (s *Service) CreateTransform(name string, resource io.Reader, extensionType string, options map[string]string, response handle.ResponseHandle) error {
	return createTransform(s.client, name, resource, extensionType, options, response)
}

// DeleteTransform removes a REST service
func (s *Service) DeleteTransform(name string, response handle.ResponseHandle) error {
	return deleteTransform(s.client, name, response)
}
