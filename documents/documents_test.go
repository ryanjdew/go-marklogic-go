package documents

import (
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/cchatfield/go-marklogic-go/util"
)

func TestBuildParameters(t *testing.T) {
	want := "?uri=%2Ftest%2F1.json&uri=%2Ftest%2F2.json&category=metadata&" +
		"category=content&collection=docs&perm:user=read&prop:date=2015-05-01&" +
		"transform=my-transform&trans:param1=val1"
	transform := &util.Transform{
		Name: "my-transform",
		Parameters: map[string]string{
			"param1": "val1",
		},
	}
	result := buildParameters(
		[]string{"/test/1.json", "/test/2.json"},
		[]string{"metadata", "content"},
		[]string{"docs"},
		map[string]string{
			"user": "read",
		},
		map[string]string{
			"date": "2015-05-01",
		},
		transform,
	)
	if want != result {
		t.Errorf("Build Parameters Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}

func TestMetadataSerialize(t *testing.T) {
	want := `{"collections":["col1"],"permissions":[{"role-name":"rest-writer","capabilities":["update"]},{"role-name":"rest-reader","capabilities":["read"]}],"properties":{"custom-prop":"some property value"},"metadataValues":{"MetaData1":"metadata"},"quality":1}`
	metadataHandle := &MetadataHandle{
		metadata: Metadata{
			Collections:    []string{"col1"},
			Permissions:    []Permission{Permission{RoleName: "rest-writer", Capability: []string{"update"}}, Permission{RoleName: "rest-reader", Capability: []string{"read"}}},
			Properties:     map[string]string{"custom-prop": "some property value"},
			Quality:        1,
			MetadataValues: map[string]string{"MetaData1": "metadata"},
		},
	}
	result := strings.Trim(metadataHandle.Serialized(), "\n")
	if want != result {
		t.Errorf("Build Parameters Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}
