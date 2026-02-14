package search

import (
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

func TestXMLQuerySerialize(t *testing.T) {
	want := "<query xmlns=\"http://marklogic.com/appservices/search\"><and-query xmlns=\"http://marklogic.com/appservices/search\"><ordered xmlns=\"http://marklogic.com/appservices/search\">true</ordered><term-query xmlns=\"http://marklogic.com/appservices/search\"><text xmlns=\"http://marklogic.com/appservices/search\">data</text></term-query></and-query></query>"
	query :=
		Query{
			Queries: []any{
				AndQuery{
					Ordered: true,
					Queries: []any{
						TermQuery{Terms: []string{"data"}},
					},
				},
			},
		}
	qh := QueryHandle{Format: handle.XML}
	qh.Serialize(query)
	result := qh.Serialized()
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	} else if !reflect.DeepEqual(&query, qh.Get()) {
		t.Errorf("Query Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}

func TestXMLQueryEncode(t *testing.T) {
	want := "<query xmlns=\"http://marklogic.com/appservices/search\"><term-query xmlns=\"http://marklogic.com/appservices/search\"><text xmlns=\"http://marklogic.com/appservices/search\">data</text></term-query></query>"
	qh := QueryHandle{Format: handle.XML}
	qh.Deserialize([]byte(want))
	result := qh.Serialized()
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}

func TestJSONQuerySerialize(t *testing.T) {
	want := `{"query":{"queries":[{"and-query":{"ordered":true,"queries":[{"term-query":{"text":["data"]}}]}}]}}`
	query :=
		Query{
			Queries: []any{
				AndQuery{
					Ordered: true,
					Queries: []any{
						TermQuery{Terms: []string{"data"}},
					},
				},
			},
		}
	qh := QueryHandle{Format: handle.JSON}
	qh.Serialize(query)
	result := strings.TrimSpace(qh.Serialized())
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}

func TestJSONQueryEncode(t *testing.T) {
	want := `{"query":{"queries":[{"and-query":{"ordered":true,"queries":[{"term-query":{"text":["data"]}}]}}]}}`
	qh := QueryHandle{Format: handle.JSON}
	qh.Deserialize([]byte(want))
	result := strings.TrimSpace((&qh).Serialized())
	if want != result {
		t.Errorf("Query Results = %+v, Want = %+v", result, want)
	}
}
