// Package dataservices provides a way to call Data Services
package dataservices

import (
	"sync"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// Service is used for the datamservices service
type Service struct {
	client        *clients.Client
	clientsByHost map[string]*clients.Client
	forestInfo    []util.ForestInfo
}

// NewService returns a new dataservices.Service
func NewService(client *clients.Client) *Service {
	forestInfo := util.GetForestInfo(client)
	clientsByHost := util.GetClientsByHost(client, forestInfo)

	return &Service{
		client:        client,
		clientsByHost: clientsByHost,
		forestInfo:    forestInfo,
	}
}

// BulkDataService for bulk data service operations
func (s *Service) BulkDataService(endpoint string) *BulkDataService {
	return &BulkDataService{
		mutex:             &sync.Mutex{},
		waitGroup:         &sync.WaitGroup{},
		endpoint:          endpoint,
		client:            s.client,
		clientsByHost:     s.clientsByHost,
		batchSize:         1000,
		forestInfo:        s.forestInfo,
		threadCount:       uint8(len(s.forestInfo)),
		workIsForestBased: false,
	}
}

// CallDataService for bulk data service operations
func (s *Service) CallDataService(endpoint string, reqParams map[string][]string, responseHandle handle.ResponseHandle) error {
	return util.PostForm(s.client, endpoint, reqParams, responseHandle, true)
}
