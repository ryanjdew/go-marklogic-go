package semantics

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
)

// Service is used for the documents service
type Service struct {
	client *clients.Client
}

// NewService returns a new search.Service
func NewService(client *clients.Client) *Service {
	return &Service{
		client: client,
	}
}
