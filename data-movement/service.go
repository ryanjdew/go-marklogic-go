package dataMovement

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
)

// Service is used for the documents service
type Service struct {
	client  *clients.Client
	clients []*clients.Client
}

// NewService returns a new search.Service
func NewService(client *clients.Client) *Service {
	// TODO Find other hosts in cluster
	clients := make([]*clients.Client, 1)
	clients = append(clients, client)
	return &Service{
		client:  client,
		clients: clients,
	}
}

// WriteBatcher for writing documents in bulk
func (s *Service) WriteBatcher() *WriteBatcher {
	return &WriteBatcher{
		documentsService: documents.NewService(s.client),
		client:           s.client,
		clients:          s.clients,
		threadCount:      4,
		batchSize:        250,
	}
}

// QueryBatcher for reading documents in bulk
func (s *Service) QueryBatcher() *QueryBatcher {
	return &QueryBatcher{
		client:      s.client,
		clients:     s.clients,
		threadCount: 4,
		batchSize:   250,
	}
}
