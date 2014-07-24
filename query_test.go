package go_marklogic_go

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuery(t *testing.T) {
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
