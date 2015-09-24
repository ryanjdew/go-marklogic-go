package admin

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Service is used for the documents service
type Service struct {
	client *clients.AdminClient
}

// NewService returns a new search.Service
func NewService(client *clients.AdminClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Init(license handle.Handle, response handle.ResponseHandle) error {
	return initialize(s.client, license, response)
}

func (s *Service) InstanceAdmin(username string, password string, realm string, response handle.ResponseHandle) error {
	return instanceAdmin(s.client, username, password, realm, response)
}
