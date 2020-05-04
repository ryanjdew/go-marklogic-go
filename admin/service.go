// Package admin can initialize MarkLogic instances
package admin

import (
	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
)

// Service is used for the admin service
type Service struct {
	client *clients.AdminClient
}

// NewService returns a new admin.Service
func NewService(client *clients.AdminClient) *Service {
	return &Service{
		client: client,
	}
}

// Init MarkLogic instance
func (s *Service) Init(license handle.Handle, response handle.ResponseHandle) error {
	return initialize(s.client, license, response)
}

// InstanceAdmin install the admin username and password, and initialize the security database and objects.
func (s *Service) InstanceAdmin(username string, password string, realm string, response handle.ResponseHandle) error {
	return instanceAdmin(s.client, username, password, realm, response)
}
