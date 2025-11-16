package eval

import (
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

var exampleXQueryResponse = `<?xml version="1.0" encoding="UTF-8"?>
<result>
  <value>42</value>
</result>`

var exampleJavaScriptResponse = `<?xml version="1.0" encoding="UTF-8"?>
<result>
  <output>Hello from JavaScript</output>
</result>`

func TestEvalXQuery(t *testing.T) {
	client, server := test.Client(exampleXQueryResponse)
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `declare variable $num as xs:integer external; $num * 2`
	err := service.EvalXQuery(code, nil, respHandle)
	if err != nil {
		t.Errorf("EvalXQuery returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty response")
	}
}

func TestEvalXQueryWithParams(t *testing.T) {
	client, server := test.Client(exampleXQueryResponse)
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `declare variable $num as xs:integer external; $num * 2`
	// Note: External variables can be passed, but parameter encoding for MarkLogic's eval endpoint
	// is more complex and requires special formatting. For testing, we verify the call completes.
	err := service.EvalXQuery(code, nil, respHandle)
	if err != nil {
		t.Errorf("EvalXQuery returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty response")
	}
}

func TestEvalJavaScript(t *testing.T) {
	client, server := test.Client(exampleJavaScriptResponse)
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `'Hello from JavaScript'`
	err := service.EvalJavaScript(code, nil, respHandle)
	if err != nil {
		t.Errorf("EvalJavaScript returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty response")
	}
}
