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

func TestAggregateValues(t *testing.T) {
	aggregateResponse := `<?xml version="1.0" encoding="UTF-8"?>
<aggregate xmlns="http://marklogic.com/appservices/search">
  <sum>100</sum>
  <count>3</count>
  <min>20</min>
  <max>42</max>
</aggregate>`

	client, server := test.Client(aggregateResponse)
	defer server.Close()

	service := NewService(client)

	aggregateBody := `<values-query xmlns="http://marklogic.com/appservices/search">
  <range name="count" type="xs:int">
    <aggregate function="sum" />
  </range>
</values-query>`

	queryHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	queryHandle.Write([]byte(aggregateBody))

	respHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	err := service.AggregateValues("count", nil, queryHandle, respHandle)
	if err != nil {
		t.Errorf("AggregateValues returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty aggregate response")
	}
}

func TestCoOccurrenceValues(t *testing.T) {
	cooccurrenceResponse := `<?xml version="1.0" encoding="UTF-8"?>
<co-occurrence xmlns="http://marklogic.com/appservices/search">
  <pair count="5">
    <value name="status">active</value>
    <value name="priority">high</value>
  </pair>
  <pair count="3">
    <value name="status">active</value>
    <value name="priority">low</value>
  </pair>
</co-occurrence>`

	client, server := test.Client(cooccurrenceResponse)
	defer server.Close()

	service := NewService(client)

	queryHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	queryHandle.Write([]byte(`<values-query xmlns="http://marklogic.com/appservices/search"></values-query>`))

	respHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	err := service.CoOccurrenceValues([]string{"status", "priority"}, nil, queryHandle, respHandle)
	if err != nil {
		t.Errorf("CoOccurrenceValues returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty co-occurrence response")
	}
}

func TestTupleValues(t *testing.T) {
	tupleResponse := `<?xml version="1.0" encoding="UTF-8"?>
<tuples xmlns="http://marklogic.com/appservices/search">
  <tuple>
    <value position="1">active</value>
    <value position="2">high</value>
    <value position="3">2024-01-15</value>
  </tuple>
  <tuple>
    <value position="1">pending</value>
    <value position="2">medium</value>
    <value position="3">2024-01-16</value>
  </tuple>
</tuples>`

	client, server := test.Client(tupleResponse)
	defer server.Close()

	service := NewService(client)

	queryHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	queryHandle.Write([]byte(`<values-query xmlns="http://marklogic.com/appservices/search"></values-query>`))

	respHandle := &handle.RawHandle{Format: handle.XML, Buffer: &bytes.Buffer{}}
	err := service.TupleValues([]string{"status", "priority", "date"}, nil, queryHandle, respHandle)
	if err != nil {
		t.Errorf("TupleValues returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty tuple response")
	}
}
