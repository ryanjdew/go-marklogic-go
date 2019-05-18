package dataMovement

import (
	"sync"

	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
)

// Service is used for the documents service
type Service struct {
	client        *clients.Client
	clientsByHost map[string]*clients.Client
	forestInfo    []ForestInfo
}

// NewService returns a new search.Service
func NewService(client *clients.Client) *Service {
	forestInfo := getForestInfo(client)

	uniqueHosts := make(map[string]struct{})
	for _, forest := range forestInfo {
		uniqueHosts[forest.PreferredHost()] = struct{}{}
	}

	clientsByHost := make(map[string]*clients.Client, len(uniqueHosts))
	connectionInfo := client.BasicClient.ConnectionInfo()
	for host := range uniqueHosts {
		if host == connectionInfo.Host {
			clientsByHost[host] = client
		} else {
			clientsByHost[host], _ = clients.NewClient(&clients.Connection{
				Host:               host,
				Port:               connectionInfo.Port,
				Username:           connectionInfo.Username,
				Password:           connectionInfo.Password,
				AuthenticationType: connectionInfo.AuthenticationType,
				Database:           connectionInfo.Database,
			})
		}
	}
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
