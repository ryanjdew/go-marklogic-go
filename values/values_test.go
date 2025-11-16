package values

import (
	"bytes"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

var exampleValuesResponse = `<?xml version="1.0" encoding="UTF-8"?>
<values-response xmlns="http://marklogic.com/appservices/search">
  <distinct-value frequency="42">active</distinct-value>
  <distinct-value frequency="38">archived</distinct-value>
  <distinct-value frequency="20">pending</distinct-value>
</values-response>`

func TestListValues(t *testing.T) {
	client, server := test.Client(exampleValuesResponse)
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	err := service.ListValues("status", nil, respHandle)
	if err != nil {
		t.Errorf("ListValues returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty response")
	}
}

func TestListValuesWithParams(t *testing.T) {
	client, server := test.Client(exampleValuesResponse)
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	params := map[string]string{
		"start":      "1",
		"pageLength": "10",
	}
	err := service.ListValues("status", params, respHandle)
	if err != nil {
		t.Errorf("ListValues returned error: %v", err)
	}
}

func TestQueryValues(t *testing.T) {
	queryBody := `<values-query xmlns="http://marklogic.com/appservices/search">
  <range name="status" type="xs:string">
    <range-operator>contains</range-operator>
    <value>active</value>
  </range>
</values-query>`

	client, server := test.Client(exampleValuesResponse)
	defer server.Close()

	service := NewService(client)

	queryHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	queryHandle.Write([]byte(queryBody))

	respHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	err := service.QueryValues("status", nil, queryHandle, respHandle)
	if err != nil {
		t.Errorf("QueryValues returned error: %v", err)
	}
}
