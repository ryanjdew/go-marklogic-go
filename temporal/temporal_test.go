package temporal

import (
	"bytes"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

func TestCreateAxis(t *testing.T) {
	mockResponse := `{"axis-name":"date-axis","status":"created"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	requestHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"axis-name":"date-axis"}`),
	}

	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.CreateAxis("date-axis", requestHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetAxis(t *testing.T) {
	mockResponse := `{"axis-name":"date-axis","scalar-uri":"xs:dateTime"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.GetAxis("date-axis", responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestListAxes(t *testing.T) {
	mockResponse := `{"axes":[{"axis-name":"date-axis"},{"axis-name":"effective-date"}]}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.ListAxes(responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDeleteAxis(t *testing.T) {
	mockResponse := `{"status":"deleted"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.DeleteAxis("date-axis", responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestEnableCollectionTemporal(t *testing.T) {
	mockResponse := `{"status":"enabled"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)

	requestHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBufferString(`{"system-axis":"system-time","valid-axis":"date-axis"}`),
	}

	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.EnableCollectionTemporal("temporal-docs", requestHandle, responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDisableCollectionTemporal(t *testing.T) {
	mockResponse := `{"status":"disabled"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.DisableCollectionTemporal("temporal-docs", responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetTemporalDocument(t *testing.T) {
	mockResponse := `{"uri":"/temporal/doc1.json","valid-start":"2024-01-01T00:00:00Z","valid-end":"2024-12-31T23:59:59Z"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.GetTemporalDocument("/temporal/doc1.json", "2024-06-15T00:00:00Z", responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestAdvanceSystemTime(t *testing.T) {
	mockResponse := `{"timestamp":"2024-06-15T12:00:00Z"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.AdvanceSystemTime("2024-06-15T12:00:00Z", responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetSystemTime(t *testing.T) {
	mockResponse := `{"timestamp":"2024-06-15T12:00:00Z"}`

	client, server := test.Client(mockResponse)
	defer server.Close()

	service := NewService(client)
	responseHandle := &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: handle.JSON,
			Buffer: &bytes.Buffer{},
		},
	}

	err := service.GetSystemTime(responseHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
