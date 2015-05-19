package search

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

var exampleResponseXML = `
<search:suggestions xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns="" xmlns:search="http://marklogic.com/appservices/search" snippet-format="snippet" total="1" start="1" page-length="10">
	<search:suggestion>data</search:suggestion>
  <search:suggestion>database</search:suggestion>
</search:suggestions>
`
var exampleResponseJSON = `
{
  "suggestions": [
    "data",
    "database"
  ]
}
`

func TestSuggestionXML(t *testing.T) {
	client, server := test.Client(exampleResponseXML)
	defer server.Close()
	want :=
		SuggestionsResponse{
			Suggestions: []string{"data", "database"},
		}
	// Using Basic Auth for test so initial call isn't actually made
	respHandle := SuggestionsResponseHandle{Format: handle.XML}
	query :=
		Query{
			Queries: []interface{}{
				TermQuery{Terms: []string{"data"}},
			},
		}
	qh := QueryHandle{}
	qh.Decode(query)
	err := StructuredSuggestions(client, &qh, "data", 10, "all", &respHandle)
	resp := respHandle.Get()
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if resp == nil {
		t.Errorf("No response found")
	} else if !reflect.DeepEqual(resp.Suggestions, want.Suggestions) {
		t.Errorf("Search Response = %+v, Want = %+v", spew.Sdump(*resp), spew.Sdump(want))
	}
}

func TestSuggestionJSON(t *testing.T) {
	client, server := test.Client(exampleResponseJSON)
	defer server.Close()
	want :=
		SuggestionsResponse{
			Suggestions: []string{"data", "database"},
		}
	// Using Basic Auth for test so initial call isn't actually made
	respHandle := SuggestionsResponseHandle{Format: handle.JSON}
	query :=
		Query{
			Queries: []interface{}{
				TermQuery{Terms: []string{"data"}},
			},
		}
	qh := QueryHandle{}
	qh.Decode(query)
	err := StructuredSuggestions(client, &qh, "data", 10, "all", &respHandle)
	resp := respHandle.Get()
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if resp == nil {
		t.Errorf("No response found")
	} else if !reflect.DeepEqual(*resp, want) {
		t.Errorf("Search Response = %+v, Want = %+v", spew.Sdump(*resp), spew.Sdump(want))
	}
}
