package config

import (
	"io"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the configuration service
type Service struct {
	client *clients.Client
}

// NewService returns a new search.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// ListExtensions shows all the installed REST extensions
func (s *Service) ListExtensions(directory string, response handle.ResponseHandle) error {
	return listExtensions(s.client, directory, response)
}

// DeleteExtensions removes all the installed REST extensions under the provided path
func (s *Service) DeleteExtensions(directory string) error {
	return deleteExtensions(s.client, directory)
}

// CreateExtension shows all the installed REST extensions
func (s *Service) CreateExtension(assetName string, resource io.Reader, extensionType string, options map[string]string, response handle.ResponseHandle) error {
	return createExtension(s.client, assetName, resource, extensionType, options, response)

}

// ListResources shows all the installed REST service extensions
func (s *Service) ListResources(response handle.ResponseHandle) error {
	return listResources(s.client, response)
}

// GetResourceInfo shows all the installed REST extensions
func (s *Service) GetResourceInfo(name string, response handle.ResponseHandle) error {
	return getResourceInfo(s.client, name, response)
}

// CreateResource installs a REST service
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
