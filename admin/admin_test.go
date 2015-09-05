package admin

import (
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/rlouapre/go-marklogic-go/test"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

func TestXMLInitializePropertiesSerialize(t *testing.T) {
	want := "<init xmlns=\"http://marklogic.com/manage\"><license-key xmlns=\"http://marklogic.com/manage\">1234-5678-90AB</license-key><licensee xmlns=\"http://marklogic.com/manage\">Your Licensee</licensee></init>"
	initializeProperties :=
		InitializeProperties{
			LicenseKey: "1234-5678-90AB",
			Licensee:   "Your Licensee",
		}
	qh := AdminHandle{}
	qh.Serialize(initializeProperties)
	result := qh.Serialized()
	if want != result {
		t.Errorf("InitializeProperties Results = %+v, Want = %+v", result, want)
	} else if !reflect.DeepEqual(&initializeProperties, qh.Get()) {
		t.Errorf("InitializeProperties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}

func TestJsonInitializePropertiesSerialize(t *testing.T) {
	want := `{"license-key":"1234-5678-90AB","licensee":"Your Licensee"}`
	initializeProperties :=
		InitializeProperties{
			LicenseKey: "1234-5678-90AB",
			Licensee:   "Your Licensee",
		}
	qh := AdminHandle{Format: handle.JSON}
	qh.Serialize(initializeProperties)
	result := strings.TrimSpace(qh.Serialized())
	if want != result {
		t.Errorf("InitializeProperties Results = %+v, Want = %+v", result, want)
	} else if !reflect.DeepEqual(&initializeProperties, qh.Get()) {
		t.Errorf("InitializeProperties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}
