// Package datamovement provides way to read and write in bulk
package datamovement

import (
	"sync"

	"github.com/ryanjdew/go-marklogic-go/util"

	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
)

// Service is used for the datamovement service
type Service struct {
	client        *clients.Client
	clientsByHost map[string]*clients.Client
	forestInfo    []util.ForestInfo
}

// NewService returns a new datamovement.Service
func NewService(client *clients.Client) *Service {
	forestInfo := util.GetForestInfo(client)
	clientsByHost := util.GetClientsByHost(client, forestInfo)

	return &Service{
		client:        client,
		clientsByHost: clientsByHost,
		forestInfo:    forestInfo,
	}
}

// WriteBatcher for writing documents in bulk
func (s *Service) WriteBatcher() *WriteBatcher {
	documentsServiceByHost := make(map[string]*documents.Service)
	for host, client := range s.clientsByHost {
		documentsServiceByHost[host] = documents.NewService(client)
	}
	return &WriteBatcher{
		documentsServiceByHost: documentsServiceByHost,
		client:                 s.client,
		clientsByHost:          s.clientsByHost,
		threadCount:            uint8(len(s.forestInfo) * 2),
		batchSize:              250,
		forestInfo:             s.forestInfo,
	}
}

// QueryBatcher for reading documents in bulk
func (s *Service) QueryBatcher() *QueryBatcher {
	return &QueryBatcher{
		mutex:         &sync.Mutex{},
		client:        s.client,
		clientsByHost: s.clientsByHost,
		batchSize:     1000,
		forestInfo:    s.forestInfo,
	}
}
