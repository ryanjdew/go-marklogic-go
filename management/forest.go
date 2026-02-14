package management

import (
	"net/http"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// ForestService contains management operations for forests (provisioning workflows)
// This service is a thin wrapper around the management client.
type ForestService struct {
	mc *clients.ManagementClient
}

// NewForestService returns a new ForestService
func NewForestService(mc *clients.ManagementClient) *ForestService {
	return &ForestService{mc: mc}
}

// CreateForest provisions a forest via the Management API.
// The payload handle should contain the serialized forest configuration.
func (s *ForestService) CreateForest(payload handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(s.mc, "POST", "/manage/v2/forests", payload)
	if err != nil {
		return err
	}
	return util.Execute(s.mc, req, response)
}

// DeleteForest deletes a forest by name via the Management API.
func (s *ForestService) DeleteForest(name string, response handle.ResponseHandle) error {
	req, err := http.NewRequest("DELETE", s.mc.Base()+"/manage/v2/forests/"+name+"?format=json", nil)
	if err != nil {
		return err
	}
	return util.Execute(s.mc, req, response)
}
