// Package rowsManagement executes an Optic API plan
package rowsManagement

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
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

// Rows executes an Optic plan and returns results
func (s *Service) Rows(opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	return rows(s.client, opticPlan, params, response)
}

// Explain returns query plan and execution information for an Optic plan
// useful for understanding performance characteristics
func (s *Service) Explain(opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	return explain(s.client, opticPlan, params, response)
}

// Sample returns a sample of rows from an Optic plan execution without full computation
func (s *Service) Sample(opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	return sample(s.client, opticPlan, params, response)
}

// Plan returns the optimized query plan for an Optic operation
func (s *Service) Plan(opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	return plan(s.client, opticPlan, params, response)
}
