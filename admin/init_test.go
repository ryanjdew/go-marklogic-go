package admin

import (
	_ "encoding/xml"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	_ "github.com/cchatfield/go-marklogic-go/test"
)

func TestXMLInitializePropertiesSerialize(t *testing.T) {
	want := `<init xmlns="http://marklogic.com/manage"><license-key>1234-5678-90AB</license-key><licensee>Your Licensee</licensee></init>`
	initializeProperties :=
		InitializeProperties{
			LicenseKey: "1234-5678-90AB",
			Licensee:   "Your Licensee",
		}
	qh := InitHandle{Format: handle.XML}
	qh.Serialize(initializeProperties)
	result := qh.Serialized()
	if want != result {
		t.Errorf("Not equal - InitializeProperties Results = %+v, Want = %+v", result, want)
	} else if !reflect.DeepEqual(&initializeProperties, qh.Get()) {
		t.Errorf("Not deep equal - InitializeProperties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}

func TestJsonInitializePropertiesSerialize(t *testing.T) {
	want := `{"license-key":"1234-5678-90AB","licensee":"Your Licensee"}`
	initializeProperties :=
		InitializeProperties{
			LicenseKey: "1234-5678-90AB",
			Licensee:   "Your Licensee",
		}
	qh := InitHandle{Format: handle.JSON}
	qh.Serialize(initializeProperties)
	result := strings.TrimSpace(qh.Serialized())
	if want != result {
		t.Errorf("InitializeProperties Results = %+v, Want = %+v", result, want)
	} else if !reflect.DeepEqual(&initializeProperties, qh.Get()) {
		t.Errorf("InitializeProperties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}
