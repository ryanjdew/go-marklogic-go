package alert

import (
	"encoding/xml"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/search"
	testHelper "github.com/ryanjdew/go-marklogic-go/test"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

func TestRulesResponseHandleDeserialize(t *testing.T) {
	want := `<rules xmlns="http://marklogic.com/rest-api"><rule><name>example</name><description>An example rule.</description><search xmlns="http://marklogic.com/appservices/search"><query xmlns="http://marklogic.com/appservices/search"></query><qtext xmlns="http://marklogic.com/appservices/search">xdmp</qtext><sparql xmlns="http://marklogic.com/appservices/search"></sparql></search><rule-metadata><author>me</author></rule-metadata></rule></rules>`
	initializeProperties :=
		&RulesResponse{
			XMLName: xml.Name{
				Space: "http://marklogic.com/rest-api",
				Local: "rules",
			},
			Rules: []Rule{
				{
					Name:        "example",
					Description: "An example rule.",
					Query: search.CombinedQuery{
						XMLName: xml.Name{
							Space: "http://marklogic.com/appservices/search",
							Local: "search",
						},
						QText: []string{"xdmp"},
					},
					RuleMetadata: util.SerializableStringMap{"author": "me"},
				},
			},
		}
	qh := &RulesResponseHandle{Format: handle.XML}
	testHelper.RoundTripSerialization(t, "RulesResponseHandle XML", initializeProperties, qh, want)
}
