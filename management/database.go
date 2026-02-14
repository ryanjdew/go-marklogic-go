package management

import (
	"net/http"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// DatabaseService contains management operations for databases (provisioning workflows)
// This service is a thin wrapper around the management client.
type DatabaseService struct {
	mc *clients.ManagementClient
}

// NewDatabaseService returns a new DatabaseService
func NewDatabaseService(mc *clients.ManagementClient) *DatabaseService {
	return &DatabaseService{mc: mc}
}

// CreateDatabase provisions a database via the Management API.
// This is a stub intended for provisioning workflows. It will POST the
// provided payload to the management endpoint. The payload should be
// a serialized representation (XML/JSON) held in the provided handle.
func (s *DatabaseService) CreateDatabase(payload handle.Handle, response handle.ResponseHandle) error {
	// Management endpoints typically live under /manage/v2
	req, err := util.BuildRequestFromHandle(s.mc, "POST", "/manage/v2/databases", payload)
	if err != nil {
		return err
	}

	// Ensure we create with the expected content type via the payload handle
	return util.Execute(s.mc, req, response)
}

// DeleteDatabase removes a database by name. This is a stub that performs
// an HTTP DELETE against the management endpoints.
func (s *DatabaseService) DeleteDatabase(name string, response handle.ResponseHandle) error {
	req, err := http.NewRequest("DELETE", s.mc.Base()+"/manage/v2/databases/"+name+"?format=json", nil)
	if err != nil {
		return err
	}
	return util.Execute(s.mc, req, response)
}
