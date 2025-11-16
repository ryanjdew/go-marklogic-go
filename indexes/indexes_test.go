package indexes

import (
	"bytes"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

func TestListIndexes(t *testing.T) {
	mockResponse := `{"indexes":[]}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &IndexHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ListIndexes(nil, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetIndex(t *testing.T) {
	mockResponse := `{"name":"element-range-index-1","scalar-type":"xs:string"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &IndexHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.GetIndex("element-range-index-1", nil, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCreateIndex(t *testing.T) {
	mockResponse := `{"status":"created"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	requestHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"name":"new-index","scalar-type":"xs:string"}`),
	}

	responseHandle := &IndexHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.CreateIndex(requestHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestUpdateIndex(t *testing.T) {
	mockResponse := `{"status":"updated"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	requestHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"name":"element-range-index-1","scalar-type":"xs:int"}`),
	}

	responseHandle := &IndexHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.UpdateIndex("element-range-index-1", requestHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteIndex(t *testing.T) {
	mockResponse := `{"status":"deleted"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &IndexHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.DeleteIndex("element-range-index-1", nil, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
