// Package semantics interacts with MarkLogic's semantic APIs
package semantics

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the semantics service
type Service struct {
	client *clients.Client
}

// NewService returns a new semantics.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}

// Things associated with the IRIs
func (s *Service) Things(iris []string, response handle.ResponseHandle) error {
	return things(s.client, iris, response)
}
