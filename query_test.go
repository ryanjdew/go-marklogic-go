package goMarklogicGo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestXMLQueryDecode(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, body)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	want := "<query xmlns=\"http://marklogic.com/appservices/search\"><and-query xmlns=\"http://marklogic.com/appservices/search\"><ordered xmlns=\"http://marklogic.com/appservices/search\">true</ordered><term-query xmlns=\"http://marklogic.com/appservices/search\"><text xmlns=\"http://marklogic.com/appservices/search\">data</text></term-query></and-query></query>"
	query :=
		Query{
			Queries: []interface{}{
				AndQuery{
					Ordered: true,
					Queries: []interface{}{
						TermQuery{Terms: []string{"data"}},
					},
				},
			},
		}
	qh := QueryHandle{}
	qh.Decode(query)
	result := qh.Serialized()
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}

func TestXMLQueryEncode(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintln(w, body)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	want := "<query xmlns=\"http://marklogic.com/appservices/search\"><term-query xmlns=\"http://marklogic.com/appservices/search\"><text xmlns=\"http://marklogic.com/appservices/search\">data</text></term-query></query>"
	qh := QueryHandle{}
	qh.Encode([]byte(want))
	result := qh.Serialized()
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
	want := `{"query":{"queries":[{"and-query":{"ordered":true,"queries":[{"term-query":{"text":["data"]}}]}}]}}`
	query :=
		Query{
			Queries: []interface{}{
				AndQuery{
					Ordered: true,
					Queries: []interface{}{
						TermQuery{Terms: []string{"data"}},
					},
				},
			},
		}
	qh := QueryHandle{Format: JSON}
	qh.Decode(query)
	result := strings.TrimSpace(qh.Serialized())
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
	want := `{"query":{"queries":[{"and-query":{"ordered":true,"queries":[{"term-query":{"text":["data"]}}]}}]}}`
	qh := QueryHandle{Format: JSON}
	qh.Encode([]byte(want))
	result := strings.TrimSpace((&qh).Serialized())
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}
