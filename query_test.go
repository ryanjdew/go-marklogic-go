package goMarklogicGo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestXMLQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, body)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	want := "<query xmlns=\"http://marklogic.com/appservices/search\"><term-query xmlns=\"http://marklogic.com/appservices/search\"><text xmlns=\"http://marklogic.com/appservices/search\">data</text></term-query></query>"
	query := NewQuery(XML)
	query.Queries = []interface{}{
		&TermQuery{
			Terms: []string{"data"},
		},
	}
	result := query.Encode().String()
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}

func TestJSONQueryEncode(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, body)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	want := `{"query":{"queries":[{"term-query":{"text":["data"]}}]}}`
	query := NewQuery(JSON)
	query.Queries = []interface{}{
		&TermQuery{
			Terms: []string{"data"},
		},
	}
	result := strings.TrimSpace(query.Encode().String())
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}

func TestJSONQueryDecode(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, body)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	want := `{"query":{"queries":[{"term-query":{"text":["data"]}}]}}`
	query := NewQuery(JSON)
	err := query.Decode(want)
	result := strings.TrimSpace(query.Encode().String())
	if err != nil {
		t.Errorf("Failed with error = %+v", err)
	} else if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}
