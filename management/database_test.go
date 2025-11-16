package management

import (
	"bytes"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

func TestCreateDatabaseStub(t *testing.T) {
	// Minimal test to ensure the service is callable; HTTP behavior is mocked by test.ManagementClient
	client, server := test.ManagementClient("{\"status\":\"ok\"}")
	defer server.Close()

	svc := NewDatabaseService(client)
	payload := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	payload.Write([]byte(`{"database-name":"ml-test-db"}`))
	resp := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	if err := svc.CreateDatabase(payload, resp); err != nil {
		t.Fatalf("CreateDatabase failed: %v", err)
	}
}

func TestDeleteDatabaseStub(t *testing.T) {
	client, server := test.ManagementClient("{\"status\":\"deleted\"}")
	defer server.Close()

	svc := NewDatabaseService(client)
	resp := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	if err := svc.DeleteDatabase("ml-test-db", resp); err != nil {
		t.Fatalf("DeleteDatabase failed: %v", err)
	}
}
