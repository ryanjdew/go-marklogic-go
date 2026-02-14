package metadata

import (
	"bytes"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

func TestExtractMetadata(t *testing.T) {
	mockResponse := `{"results":[{"uri":"/doc1.json","metadata":{"title":"Document 1"}}]}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ExtractMetadata([]string{"/doc1.json"}, nil, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestExtractMetadataFromQuery(t *testing.T) {
	mockResponse := `{"results":[{"uri":"/doc1.json","metadata":{"title":"Document 1"}}]}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	queryHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"query":{"match-all":{}}}`),
	}

	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ExtractMetadataFromQuery(queryHandle, nil, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestValidateDocuments(t *testing.T) {
	mockResponse := `{"results":[{"uri":"/doc1.json","valid":true}]}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	rulesHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"rules":[{"name":"required-title","xpath":"title"}]}`),
	}

	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ValidateDocuments([]string{"/doc1.json"}, rulesHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestValidateQuery(t *testing.T) {
	mockResponse := `{"results":[{"uri":"/doc1.json","valid":true}]}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	queryHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"query":{"match-all":{}}}`),
	}

	rulesHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"rules":[{"name":"required-title","xpath":"title"}]}`),
	}

	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ValidateQuery(queryHandle, rulesHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetValidationRules(t *testing.T) {
	mockResponse := `{"rules":[{"name":"required-title","xpath":"title"}]}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.GetValidationRules(responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSetValidationRules(t *testing.T) {
	mockResponse := `{"status":"updated"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	rulesHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"rules":[{"name":"required-title","xpath":"title"}]}`),
	}

	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.SetValidationRules(rulesHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestExtractMetadataFromURI(t *testing.T) {
	mockResponse := `{"uri":"/doc1.json","metadata":{"title":"Document 1","author":"John Doe"}}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ExtractMetadataFromURI("/doc1.json", nil, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestValidateURI(t *testing.T) {
	mockResponse := `{"uri":"/doc1.json","valid":true}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	rulesHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"rules":[{"name":"required-title","xpath":"title"}]}`),
	}

	responseHandle := &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ValidateURI("/doc1.json", rulesHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
