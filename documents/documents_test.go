package documents

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ryanjdew/go-marklogic-go/util"
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
