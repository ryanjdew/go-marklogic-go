package eval

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

var exampleXQueryResponse = `<?xml version="1.0" encoding="UTF-8"?>
<result>
  <value>42</value>
</result>`

var exampleJavaScriptResponse = `<?xml version="1.0" encoding="UTF-8"?>
<result>
  <output>Hello from JavaScript</output>
</result>`

// captureRequestHandler captures request details for inspection
type captureRequestHandler struct {
	requestBody    string
	requestHeaders http.Header
	requestURL     *url.URL
	responseBody   string
}

func (h *captureRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.requestURL = r.URL
	h.requestHeaders = r.Header.Clone()
	body, _ := io.ReadAll(r.Body)
	h.requestBody = string(body)
	fmt.Fprintln(w, h.responseBody)
}

func newCaptureClient(resp string, database string) (*clients.Client, *httptest.Server, *captureRequestHandler) {
	handler := &captureRequestHandler{responseBody: resp}
	server := httptest.NewServer(handler)
	conn := &clients.Connection{
		Host:               "localhost",
		Port:               8000,
		Username:           "admin",
		Password:           "admin",
		AuthenticationType: clients.BasicAuth,
		Database:           database,
	}
	client, _ := clients.NewClient(conn)
	client.SetBase(server.URL)
	return client, server, handler
}

func TestEvalXQuery(t *testing.T) {
	client, server, handler := newCaptureClient(exampleXQueryResponse, "")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `declare variable $num as xs:integer external; $num * 2`
	err := service.EvalXQuery(code, nil, nil, respHandle)
	if err != nil {
		t.Errorf("EvalXQuery returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty response")
	}

	// Verify request format
	if handler.requestHeaders.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type: application/x-www-form-urlencoded, got %s", handler.requestHeaders.Get("Content-Type"))
	}

	// Verify form-encoded body contains xquery parameter
	if !strings.Contains(handler.requestBody, "xquery=") {
		t.Errorf("Expected form body to contain 'xquery=' parameter, got: %s", handler.requestBody)
	}

	if !strings.Contains(handler.requestBody, url.QueryEscape(code)) {
		t.Errorf("Expected form body to contain encoded code")
	}
}

func TestEvalXQueryWithExternalVariables(t *testing.T) {
	client, server, handler := newCaptureClient(exampleXQueryResponse, "")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `declare variable $num as xs:integer external; $num * 2`
	params := map[string]string{
		"num": "42",
	}

	err := service.EvalXQuery(code, params, nil, respHandle)
	if err != nil {
		t.Errorf("EvalXQuery returned error: %v", err)
	}

	// Verify form-encoded body contains both xquery and vars parameters
	if !strings.Contains(handler.requestBody, "xquery=") {
		t.Errorf("Expected form body to contain 'xquery=' parameter")
	}

	if !strings.Contains(handler.requestBody, "vars=") {
		t.Errorf("Expected form body to contain 'vars=' parameter for external variables")
	}

	// Verify vars JSON contains the parameter
	if !strings.Contains(handler.requestBody, "num") {
		t.Errorf("Expected form body to contain variable name 'num'")
	}
}

func TestEvalJavaScript(t *testing.T) {
	client, server, handler := newCaptureClient(exampleJavaScriptResponse, "")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `'Hello from JavaScript'`
	err := service.EvalJavaScript(code, nil, nil, respHandle)
	if err != nil {
		t.Errorf("EvalJavaScript returned error: %v", err)
	}

	if respHandle.Serialized() == "" {
		t.Errorf("Expected non-empty response")
	}

	// Verify request format
	if handler.requestHeaders.Get("Content-Type") != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type: application/x-www-form-urlencoded, got %s", handler.requestHeaders.Get("Content-Type"))
	}

	// Verify form-encoded body contains javascript parameter
	if !strings.Contains(handler.requestBody, "javascript=") {
		t.Errorf("Expected form body to contain 'javascript=' parameter, got: %s", handler.requestBody)
	}

	if !strings.Contains(handler.requestBody, url.QueryEscape(code)) {
		t.Errorf("Expected form body to contain encoded code")
	}
}

func TestEvalJavaScriptWithExternalVariables(t *testing.T) {
	client, server, handler := newCaptureClient(exampleJavaScriptResponse, "")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `word1 + " " + word2`
	params := map[string]string{
		"word1": "hello",
		"word2": "world",
	}

	err := service.EvalJavaScript(code, params, nil, respHandle)
	if err != nil {
		t.Errorf("EvalJavaScript returned error: %v", err)
	}

	// Verify form-encoded body contains both javascript and vars parameters
	if !strings.Contains(handler.requestBody, "javascript=") {
		t.Errorf("Expected form body to contain 'javascript=' parameter")
	}

	if !strings.Contains(handler.requestBody, "vars=") {
		t.Errorf("Expected form body to contain 'vars=' parameter for external variables")
	}

	// Verify vars JSON contains the parameters
	if !strings.Contains(handler.requestBody, "word1") || !strings.Contains(handler.requestBody, "word2") {
		t.Errorf("Expected form body to contain variable names 'word1' and 'word2'")
	}
}

func TestEvalWithDatabase(t *testing.T) {
	client, server, handler := newCaptureClient(exampleXQueryResponse, "mydb")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `1 + 1`
	err := service.EvalXQuery(code, nil, nil, respHandle)
	if err != nil {
		t.Errorf("EvalXQuery returned error: %v", err)
	}

	// Verify database parameter in URL query string
	if handler.requestURL.Query().Get("database") != "mydb" {
		t.Errorf("Expected database=mydb in query parameters, got: %s", handler.requestURL.RawQuery)
	}
}

func TestEvalXQueryMultilineCode(t *testing.T) {
	client, server, handler := newCaptureClient(exampleXQueryResponse, "")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `xquery version "1.0-ml";
declare variable $x as xs:integer external;
for $i in 1 to $x
return $i * 2`

	err := service.EvalXQuery(code, nil, nil, respHandle)
	if err != nil {
		t.Errorf("EvalXQuery returned error: %v", err)
	}

	// Verify multiline code is properly encoded in form body
	if !strings.Contains(handler.requestBody, "xquery=") {
		t.Errorf("Expected form body to contain xquery parameter with multiline code")
	}
}

func TestEvalJavaScriptReturnsMultipleValues(t *testing.T) {
	responseWithMultipleParts := `--boundary123
Content-Type: text/plain
X-Primitive: untypedAtomic

hello
--boundary123
Content-Type: text/plain
X-Primitive: untypedAtomic

world
--boundary123--`

	client, server, handler := newCaptureClient(responseWithMultipleParts, "")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	code := `Sequence.from([word1, word2])`
	params := map[string]string{
		"word1": "hello",
		"word2": "world",
	}

	err := service.EvalJavaScript(code, params, nil, respHandle)
	if err != nil {
		t.Errorf("EvalJavaScript returned error: %v", err)
	}

	// Verify form encoding includes both parameters
	if !strings.Contains(handler.requestBody, "javascript=") {
		t.Errorf("Expected form body to contain 'javascript=' parameter")
	}

	if !strings.Contains(handler.requestBody, "vars=") {
		t.Errorf("Expected form body to contain 'vars=' parameter")
	}
}

func TestEvalWithTransaction(t *testing.T) {
	client, server, handler := newCaptureClient(exampleXQueryResponse, "")
	defer server.Close()

	service := NewService(client)
	respHandle := &handle.RawHandle{Format: handle.XML}

	// Create a mock transaction with an ID
	transaction := &util.Transaction{ID: "test-txid-12345"}

	code := `xdmp:document-insert('/test.xml', <test/>)`
	err := service.EvalXQuery(code, nil, transaction, respHandle)
	if err != nil {
		t.Errorf("EvalXQuery with transaction returned error: %v", err)
	}

	// Verify transaction ID is included in query parameters
	if handler.requestURL.Query().Get("txid") != "test-txid-12345" {
		t.Errorf("Expected txid=test-txid-12345 in query parameters, got: %s", handler.requestURL.RawQuery)
	}

	// Verify form encoding still works
	if !strings.Contains(handler.requestBody, "xquery=") {
		t.Errorf("Expected form body to contain 'xquery=' parameter with transaction")
	}
}
