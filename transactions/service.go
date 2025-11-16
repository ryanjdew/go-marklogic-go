// Package transactions provides multi-statement transaction management
package transactions

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service provides transaction management operations
type Service struct {
	client *clients.Client
}

// NewService returns a new Service for transaction operations
func NewService(c *clients.Client) *Service {
	return &Service{client: c}
}

// Begin starts a new multi-statement transaction.
// Returns the transaction ID (txid) as a string via the response handle.
func (s *Service) Begin(response handle.ResponseHandle) error {
	return begin(s.client, response)
}

// Commit commits an active transaction by its ID.
func (s *Service) Commit(txid string, response handle.ResponseHandle) error {
	return commit(s.client, txid, response)
}

// Rollback rolls back an active transaction by its ID.
func (s *Service) Rollback(txid string, response handle.ResponseHandle) error {
	return rollback(s.client, txid, response)
}

// Status retrieves the status of a transaction.
func (s *Service) Status(txid string, response handle.ResponseHandle) error {
	return status(s.client, txid, response)
}
