package util

import (
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
	"github.com/ryanjdew/go-marklogic-go/test/text"
)

func TestExecute(t *testing.T) {
	want := `{"success":true}`
	client, _ := test.Client(want)
	req, err := http.NewRequest("GET", client.Base(), nil)
	respHandle := handle.RawHandle{Format: handle.JSON}
	Execute(client, req, &respHandle)
	result := testHelper.NormalizeSpace(respHandle.Serialized())
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if want != result {
		t.Errorf("Result = %v, want %v", spew.Sdump(result), spew.Sdump(want))
	}
}
