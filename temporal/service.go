// Package temporal manages temporal documents and axes in MarkLogic
package temporal

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the temporal service
type Service struct {
	client *clients.Client
}

// NewService returns a new temporal.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// CreateAxis creates a new temporal axis
func (s *Service) CreateAxis(axisName string, requestBody handle.Handle, response handle.ResponseHandle) error {
	return createAxis(s.client, axisName, requestBody, response)
}

// GetAxis retrieves axis configuration by name
func (s *Service) GetAxis(axisName string, response handle.ResponseHandle) error {
	return getAxis(s.client, axisName, response)
}

// ListAxes retrieves all configured temporal axes
func (s *Service) ListAxes(response handle.ResponseHandle) error {
	return listAxes(s.client, response)
}

// DeleteAxis removes a temporal axis
func (s *Service) DeleteAxis(axisName string, response handle.ResponseHandle) error {
	return deleteAxis(s.client, axisName, response)
}

// EnableCollectionTemporal enables temporal functionality on a collection
func (s *Service) EnableCollectionTemporal(collection string, temporalConfig handle.Handle, response handle.ResponseHandle) error {
	return enableCollectionTemporal(s.client, collection, temporalConfig, response)
}

// DisableCollectionTemporal disables temporal functionality on a collection
func (s *Service) DisableCollectionTemporal(collection string, response handle.ResponseHandle) error {
	return disableCollectionTemporal(s.client, collection, response)
}

// GetTemporalDocument retrieves a temporal document at a specific point in time
func (s *Service) GetTemporalDocument(uri string, timestamp string, response handle.ResponseHandle) error {
	return getTemporalDocument(s.client, uri, timestamp, response)
}

// AdvanceSystemTime advances the system time for temporal operations
func (s *Service) AdvanceSystemTime(timestamp string, response handle.ResponseHandle) error {
	return advanceSystemTime(s.client, timestamp, response)
}

// GetSystemTime retrieves the current system time for temporal operations
func (s *Service) GetSystemTime(response handle.ResponseHandle) error {
	return getSystemTime(s.client, response)
}
