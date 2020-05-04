// Package rowsManagement executes an Optic API plan
package rowsManagement

import (
	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
)

// Service is used for the rowsManagement service
type Service struct {
	client *clients.Client
}

// NewService returns a new rowsManagement.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Rows find rules that match a given query
func (s *Service) Rows(opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	return rows(s.client, opticPlan, params, response)
}
