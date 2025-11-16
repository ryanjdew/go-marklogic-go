package rowsManagement

import (
	"bytes"
	"testing"

	h "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

var exampleRowsResponse = `<?xml version="1.0" encoding="UTF-8"?>
<rows xmlns="http://marklogic.com/optic">
  <row>
    <col name="id">1</col>
    <col name="name">Alice</col>
    <col name="department">Engineering</col>
  </row>
  <row>
    <col name="id">2</col>
    <col name="name">Bob</col>
    <col name="department">Sales</col>
  </row>
  <row>
    <col name="id">3</col>
    <col name="name">Carol</col>
    <col name="department">Engineering</col>
  </row>
</rows>`

var exampleExplainResponse = `<?xml version="1.0" encoding="UTF-8"?>
<plan xmlns="http://marklogic.com/optic">
  <estimates>
    <row-count>1000000</row-count>
    <estimated-cost>42.5</estimated-cost>
  </estimates>
  <plan-node type="sequence">
    <plan-node type="source-node" uri="employees" />
    <plan-node type="filter-node" />
    <plan-node type="project-node" />
  </plan-node>
</plan>`

var exampleSampleResponse = `<?xml version="1.0" encoding="UTF-8"?>
<rows xmlns="http://marklogic.com/optic">
  <row>
    <col name="id">1</col>
    <col name="name">Alice</col>
  </row>
  <row>
    <col name="id">5</col>
    <col name="name">Eve</col>
  </row>
</rows>`

var examplePlanResponse = `<?xml version="1.0" encoding="UTF-8"?>
<optimized-plan xmlns="http://marklogic.com/optic">
  <access-plan>
    <index-usage type="range-index">
      <index-name>department-index</index-name>
      <estimated-selectivity>0.15</estimated-selectivity>
    </index-usage>
  </access-plan>
  <steps>
    <step sequence="1">Access employees by range index</step>
    <step sequence="2">Filter results by expression</step>
    <step sequence="3">Project selected columns</step>
  </steps>
</optimized-plan>`

func TestRows(t *testing.T) {
	client, server := test.Client(exampleRowsResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{
		"_export": "EMPLOYEES",
		"_module": {
			"_import": "marklogic-optic",
			"fn": "op"
		},
		"_next": {
			"_export": "FROM_EMPLOYEES",
			"_invoke": {
				"_function": "from",
				"_args": ["EMPLOYEES", "employees"]
			}
		}
	}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	err := service.Rows(opticPlan, nil, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty response")
	}
}

func TestRowsWithParams(t *testing.T) {
	client, server := test.Client(exampleRowsResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	params := map[string]string{
		"pageLength": "100",
		"start":      "1",
	}
	err := service.Rows(opticPlan, params, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestExplain(t *testing.T) {
	client, server := test.Client(exampleExplainResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	err := service.Explain(opticPlan, nil, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty explain response")
	}
}

func TestExplainWithParams(t *testing.T) {
	client, server := test.Client(exampleExplainResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	params := map[string]string{
		"debug": "true",
	}
	err := service.Explain(opticPlan, params, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSample(t *testing.T) {
	client, server := test.Client(exampleSampleResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	err := service.Sample(opticPlan, nil, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty sample response")
	}
}

func TestSampleWithParams(t *testing.T) {
	client, server := test.Client(exampleSampleResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	params := map[string]string{
		"size": "10",
	}
	err := service.Sample(opticPlan, params, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPlan(t *testing.T) {
	client, server := test.Client(examplePlanResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	err := service.Plan(opticPlan, nil, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty plan response")
	}
}

func TestPlanWithParams(t *testing.T) {
	client, server := test.Client(examplePlanResponse)
	defer server.Close()

	service := NewService(client)

	opticPlan := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	opticPlan.Write([]byte(`{"_export": "EMPLOYEES"}`))

	respHandle := &h.RawHandle{Format: h.JSON, Buffer: &bytes.Buffer{}}
	params := map[string]string{
		"format": "json",
	}
	err := service.Plan(opticPlan, params, respHandle)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
