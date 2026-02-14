package management

import (
	"bytes"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

func TestCreateForestStub(t *testing.T) {
	client, server := test.ManagementClient("{\"status\":\"ok\"}")
	defer server.Close()

	svc := NewForestService(client)
	payload := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	payload.Write([]byte(`{"forest-name":"ml-test-forest"}`))
	resp := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	if err := svc.CreateForest(payload, resp); err != nil {
		t.Fatalf("CreateForest failed: %v", err)
	}
}

func TestDeleteForestStub(t *testing.T) {
	client, server := test.ManagementClient("{\"status\":\"deleted\"}")
	defer server.Close()

	svc := NewForestService(client)
	resp := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	if err := svc.DeleteForest("ml-test-forest", resp); err != nil {
		t.Fatalf("DeleteForest failed: %v", err)
	}
}
