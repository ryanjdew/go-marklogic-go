package search

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	test "github.com/ryanjdew/go-marklogic-go/test"
)

var exampleResponse = `
<search:response xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns="" xmlns:search="http://marklogic.com/appservices/search" snippet-format="snippet" total="1" start="1" page-length="10">
	<search:result index="1" uri="/resources/wikipedia/ru_id/341620.xml" path='fn:doc("/resources/wikipedia/ru_id/341620.xml")' score="178432" confidence="0.9790292" fitness="0.9790292" href="/v1/documents?uri=%2Fresources%2Fwikipedia%2Fru_id%2F341620.xml"
	mimetype="text/xml" format="xml">
		<search:snippet>
			<search:match path='fn:doc("/resources/wikipedia/ru_id/341620.xml")/resource/description'>Lieutenant Commander <search:highlight>Data</search:highlight> is a character in the fictional Star Trek universe portrayed by actor Brent Spiner....</search:match>
		</search:snippet>
		<search:metadata>
			<organization>Starfleet</organization>
			<gender>Male</gender>
			<rank>Lieutenant Commander</rank>
			<name>Data</name>
			<search:constraint-meta name="BirthPlace">Omicron Theta</search:constraint-meta>
			<search:constraint-meta name="Organization">Starfleet</search:constraint-meta>
			<search:constraint-meta name="Rank">Lieutenant Commander</search:constraint-meta>
			<search:constraint-meta name="Species">Android</search:constraint-meta>
			<search:constraint-meta name="Species">Artificial intelligence</search:constraint-meta>
			<search:constraint-meta name="Gender">Male</search:constraint-meta>
		</search:metadata>
	</search:result>
	<search:facet name="Organization" type="xs:string">
		<search:facet-value name="Starfleet" count="1">Starfleet</search:facet-value>
	</search:facet>
	<search:facet name="Species" type="xs:string">
		<search:facet-value name="Android" count="1">Android</search:facet-value>
		<search:facet-value name="Artificial intelligence" count="1">Artificial intelligence</search:facet-value>
	</search:facet>
	<search:qtext>data</search:qtext>
	<search:query>
		<cts:word-query xmlns:cts="http://marklogic.com/cts">
			<cts:text xml:lang="en">data</cts:text>
			<cts:option>punctuation-insensitive</cts:option>
		</cts:word-query>
	</search:query>
	<search:metrics>
		<search:query-resolution-time>PT0.005881S</search:query-resolution-time>
		<search:facet-resolution-time>PT0.002879S</search:facet-resolution-time>
		<search:snippet-resolution-time>PT0.002744S</search:snippet-resolution-time>
		<search:metadata-resolution-time>PT0.000247S</search:metadata-resolution-time>
		<search:total-time>PT0.013846S</search:total-time>
	</search:metrics>
</search:response>
`

func TestSearch(t *testing.T) {
	client, server := test.Client(exampleResponse)
	defer server.Close()
	want :=
		Response{
			Total:      1,
			Start:      1,
			PageLength: 10,
			Results: []Result{
				Result{
					URI:        "/resources/wikipedia/ru_id/341620.xml",
					Href:       "/v1/documents?uri=%2Fresources%2Fwikipedia%2Fru_id%2F341620.xml",
					MimeType:   "text/xml",
					Format:     "xml",
					Path:       "fn:doc(\"/resources/wikipedia/ru_id/341620.xml\")",
					Index:      1,
					Score:      178432,
					Confidence: 0.9790292,
					Fitness:    0.9790292,
					Snippets: []Snippet{
						Snippet{
							Matches: []Match{
								Match{
									Path: "fn:doc(\"/resources/wikipedia/ru_id/341620.xml\")/resource/description",
									Text: []Text{
										Text{
											Text:            "Lieutenant Commander ",
											HighlightedText: false,
										},
										Text{
											Text:            "Data",
											HighlightedText: true,
										},
										Text{
											Text:            " is a character in the fictional Star Trek universe portrayed by actor Brent Spiner....",
											HighlightedText: false,
										},
									},
								},
							},
						},
					},
				},
			},
			Facets: []Facet{
				Facet{
					Name: "Organization",
					Type: "xs:string",
					FacetValues: []FacetValue{
						FacetValue{
							Name:  "Starfleet",
							Label: "Starfleet",
							Count: 1,
						},
					},
				},
				Facet{
					Name: "Species",
					Type: "xs:string",
					FacetValues: []FacetValue{
						FacetValue{
							Name:  "Android",
							Label: "Android",
							Count: 1,
						},
						FacetValue{
							Name:  "Artificial intelligence",
							Label: "Artificial intelligence",
							Count: 1,
						},
					},
				},
			},
		}
	// Using Basic Auth for test so initial call isn't actually made
	respHandle := ResponseHandle{Format: handle.XML}
	err := Search(client, "data", 1, 10, &respHandle)
	resp := respHandle.Get()
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if resp == nil {
		t.Errorf("No response found")
	} else if !reflect.DeepEqual(want.Results[0], resp.Results[0]) {
		t.Errorf("Search Results = %+v, Want = %+v", spew.Sdump(resp.Results[0]), spew.Sdump(want.Results[0]))
	} else if !reflect.DeepEqual(resp.Facets, want.Facets) {
		t.Errorf("Search Facets = %+v, Want = %+v", spew.Sdump(resp.Facets), spew.Sdump(want.Facets))
	} else if !reflect.DeepEqual(*resp, want) {
		t.Errorf("Search Response = %+v, Want = %+v", spew.Sdump(*resp), spew.Sdump(want))
	}
	query :=
		Query{
			Queries: []interface{}{
				TermQuery{Terms: []string{"data"}},
			},
		}
	qh := QueryHandle{}
	qh.Decode(query)
	err = StructuredSearch(client, &qh, 1, 10, &respHandle)
	resp = respHandle.Get()
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if resp == nil {
		t.Errorf("No response found")
	} else if !reflect.DeepEqual(resp.Results[0], want.Results[0]) {
		t.Errorf("Search Results = %+v, Want = %+v", spew.Sdump(resp.Results[0]), spew.Sdump(want.Results[0]))
	} else if !reflect.DeepEqual(resp.Facets, want.Facets) {
		t.Errorf("Search Facets = %+v, Want = %+v", spew.Sdump(resp.Facets), spew.Sdump(want.Facets))
	} else if !reflect.DeepEqual(*resp, want) {
		t.Errorf("Search Response = %+v, Want = %+v", spew.Sdump(*resp), spew.Sdump(want))
	}
}
