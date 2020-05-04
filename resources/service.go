// Package resources allows for interacting with custom REST extensions
package resources

import (
	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
)

// Service is used for the resources service
type Service struct {
	client *clients.Client
}

// NewService returns a new resource.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Get Call GET against resource
func (s *Service) Get(resourceName string, parameters map[string]string, response handle.ResponseHandle) error {
	return get(s.client, resourceName, parameters, response)
}

// Post Call POST against resource
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
