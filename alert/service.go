// Package alert works with the MarkLogic Alerting API
package alert

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the alert service
type Service struct {
	client *clients.Client
}

// NewService returns a new alert.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// MatchQuery find rules that match a given query
func (s *Service) MatchQuery(query handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	return matchQuery(s.client, query, params, response)
}

// MatchDocument finds rules that match a give document
func (s *Service) MatchDocument(document documents.DocumentDescription, params map[string]string, response handle.ResponseHandle) error {
	return matchDocument(s.client, document, params, response)
}

// ListRules lists rules
func (s *Service) ListRules(params map[string]string, response handle.ResponseHandle) error {
	return listRules(s.client, params, response)
}

// GetRule get a rule
func (s *Service) GetRule(ruleName string, params map[string]string, response handle.ResponseHandle) error {
	return getRule(s.client, ruleName, params, response)
}

// AddRule add a rule
func (s *Service) AddRule(ruleName string, rule handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	return addRule(s.client, ruleName, rule, params, response)
}

// RemoveRule remove a rule
func (s *Service) RemoveRule(ruleName string, params map[string]string, response handle.ResponseHandle) error {
	return removeRule(s.client, ruleName, params, response)
}
