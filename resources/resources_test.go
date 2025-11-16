package resources

import (
	"bytes"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

var exampleResourceResponse = `{
  "name": "my-resource",
  "title": "My Custom Resource",
  "provider": "my-app",
  "uri": "/v1/resources/my-resource"
}`

func TestGetResource(t *testing.T) {
	client, server := test.Client(exampleResourceResponse)
	defer server.Close()

	svc := NewService(client)
	resp := &ResourceExtensionHandle{Format: handle.JSON}
	err := svc.Get("my-resource", nil, resp)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	resource := resp.Deserialized().(ResourceExtension)
	if resource.Name != "my-resource" {
		t.Errorf("Expected name my-resource, got %s", resource.Name)
	}
}

func TestPostResource(t *testing.T) {
	client, server := test.Client(`{"status":"created"}`)
	defer server.Close()

	svc := NewService(client)
	payload := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	payload.Write([]byte(`{"name":"new-resource"}`))

	resp := &ResourceExtensionHandle{Format: handle.JSON}
	err := svc.Post("new-resource", nil, payload, resp)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
}

func TestPutResource(t *testing.T) {
	client, server := test.Client(`{"status":"updated"}`)
	defer server.Close()

	svc := NewService(client)
	payload := &handle.RawHandle{Format: handle.JSON, Buffer: &bytes.Buffer{}}
	payload.Write([]byte(`{"title":"Updated Title"}`))

	resp := &handle.RawHandle{Format: handle.JSON}
	err := svc.Put("my-resource", nil, payload, resp)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
}

func TestDeleteResource(t *testing.T) {
	client, server := test.Client(`{"status":"deleted"}`)
	defer server.Close()

	svc := NewService(client)
	resp := &handle.RawHandle{Format: handle.JSON}
	err := svc.Delete("my-resource", nil, resp)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestResourceWithParameters(t *testing.T) {
	client, server := test.Client(exampleResourceResponse)
	defer server.Close()

	svc := NewService(client)
	params := map[string]string{
		"format": "json",
		"limit":  "10",
	}
	resp := &ResourceExtensionHandle{Format: handle.JSON}
	err := svc.Get("my-resource", params, resp)
	if err != nil {
		t.Fatalf("Get with params failed: %v", err)
	}
}
