package documents

import (
	"encoding/xml"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	testHelper "github.com/ryanjdew/go-marklogic-go/test"
)

func TestMetadataSerializeJSON(t *testing.T) {
	want := `{"collections":["col1"],"permissions":[{"role-name":"rest-writer","capabilities":["update"]},{"role-name":"rest-reader","capabilities":["read"]}],"properties":{"custom-prop":"some property value"},"metadataValues":{"MetaData1":"metadata"},"quality":1}`
	metadata := Metadata{
		Collections:    []string{"col1"},
		Permissions:    []Permission{Permission{RoleName: "rest-writer", Capability: []string{"update"}}, Permission{RoleName: "rest-reader", Capability: []string{"read"}}},
		Properties:     map[string]string{"custom-prop": "some property value"},
		Quality:        1,
		MetadataValues: map[string]string{"MetaData1": "metadata"},
	}
	testHelper.RoundTripSerialization(t, "Metadata", metadata, &MetadataHandle{Format: handle.JSON}, want)
}

func TestMetadataSerializeXML(t *testing.T) {
	want := "<metadata xmlns=\"http://marklogic.com/rest-api\"><collections xmlns=\"http://marklogic.com/rest-api\">col1</collections><permission xmlns=\"http://marklogic.com/rest-api\"><role-name xmlns=\"http://marklogic.com/rest-api\">rest-writer</role-name><capability xmlns=\"http://marklogic.com/rest-api\">update</capability></permission><permission xmlns=\"http://marklogic.com/rest-api\"><role-name xmlns=\"http://marklogic.com/rest-api\">rest-reader</role-name><capability xmlns=\"http://marklogic.com/rest-api\">read</capability></permission><properties xmlns=\"http://marklogic.com/xdmp/property\"><custom-prop>some property value</custom-prop></properties><metadata-values xmlns=\"http://marklogic.com/rest-api\"><MetaData1>metadata</MetaData1></metadata-values><quality xmlns=\"http://marklogic.com/rest-api\">1</quality></metadata>"
	metadata := Metadata{
		XMLName: xml.Name{
			Space: "http://marklogic.com/rest-api",
			Local: "metadata",
		},
		Collections: []string{"col1"},
		Permissions: []Permission{Permission{XMLName: xml.Name{
			Space: "http://marklogic.com/rest-api",
			Local: "permission",
		}, RoleName: "rest-writer", Capability: []string{"update"}}, Permission{XMLName: xml.Name{
			Space: "http://marklogic.com/rest-api",
			Local: "permission",
		}, RoleName: "rest-reader", Capability: []string{"read"}}},
		Properties:     map[string]string{"custom-prop": "some property value"},
		Quality:        1,
		MetadataValues: map[string]string{"MetaData1": "metadata"},
	}
	testHelper.RoundTripSerialization(t, "Metadata", metadata, &MetadataHandle{Format: handle.XML}, want)
}
