package admin

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	test "github.com/cchatfield/go-marklogic-go/test"
)

var instanceAdminResponse = `<restart xmlns="http://marklogic.com/manage"><last-startup host-id="13544732455686476949">2013-04-01T10:35:19.09913-07:00</last-startup><link><kindref>timestamp</kindref><uriref>/admin/v1/timestamp</uriref></link><message>Check for new timestamp to verify host restart.</message></restart>`

func TestInstanceAdmin(t *testing.T) {
	client, server := test.AdminClient(instanceAdminResponse)
	defer server.Close()
	want :=
		RestartResponse{
			XMLName: xml.Name{Space: "http://marklogic.com/manage", Local: "restart"},
			LastStartup: LastStartupElement{
				XMLName: xml.Name{Space: "http://marklogic.com/manage", Local: "last-startup"},
				Value:   "2013-04-01T10:35:19.09913-07:00",
				HostID:  "13544732455686476949",
			},
			Link: LinkElement{
				XMLName: xml.Name{Space: "http://marklogic.com/manage", Local: "link"},
				KindRef: "timestamp",
				URIRef:  "/admin/v1/timestamp",
			},
			Message: "Check for new timestamp to verify host restart.",
		}
	// Using Basic Auth for test so initial call isn't actually made
	respHandle := RestartResponseHandle{Format: handle.XML}
	err := instanceAdmin(client, "admin", "password", "public", &respHandle)
	resp := respHandle.Get()
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if resp == nil {
		t.Errorf("No response found")
	} else if !reflect.DeepEqual(want.LastStartup, resp.LastStartup) {
		t.Errorf("InstanceAdmin LastStartup = %+v, Want = %+v", spew.Sdump(resp.LastStartup), spew.Sdump(want.LastStartup))
	} else if !reflect.DeepEqual(resp.Link, want.Link) {
		t.Errorf("InstanceAdmin Link = %+v, Want = %+v", spew.Sdump(resp.Link), spew.Sdump(want.Link))
	} else if !reflect.DeepEqual(*resp, want) {
		t.Errorf("InstanceAdmin Response = %+v, Want = %+v", spew.Sdump(*resp), spew.Sdump(want))
	}
}

var initResponse = `<restart xmlns="http://marklogic.com/manage"><last-startup host-id="13544732455686476949">2013-05-15T09:01:43.019261-07:00</last-startup><link><kindref>timestamp</kindref><uriref>/admin/v1/timestamp</uriref></link><message>Check for new timestamp to verify host restart.</message></restart>`

func TestInit(t *testing.T) {
	client, server := test.AdminClient(initResponse)
	defer server.Close()
	want :=
		RestartResponse{
			XMLName: xml.Name{Space: "http://marklogic.com/manage", Local: "restart"},
			LastStartup: LastStartupElement{
				XMLName: xml.Name{Space: "http://marklogic.com/manage", Local: "last-startup"},
				Value:   "2013-05-15T09:01:43.019261-07:00",
				HostID:  "13544732455686476949",
			},
			Link: LinkElement{
				XMLName: xml.Name{Space: "http://marklogic.com/manage", Local: "link"},
				KindRef: "timestamp",
				URIRef:  "/admin/v1/timestamp",
			},
			Message: "Check for new timestamp to verify host restart.",
		}
	ih := InitHandle{}
	license := InitializeProperties{
		LicenseKey: "1234-5678-90AB",
		Licensee:   "Your Licensee",
	}
	ih.Serialize(license)

	// Using Basic Auth for test so initial call isn't actually made
	respHandle := RestartResponseHandle{Format: handle.XML}
	err := initialize(client, &ih, &respHandle)
	resp := respHandle.Get()
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if resp == nil {
		t.Errorf("No response found")
	} else if !reflect.DeepEqual(want.LastStartup, resp.LastStartup) {
		t.Errorf("InstanceAdmin LastStartup = %+v, Want = %+v", spew.Sdump(resp.LastStartup), spew.Sdump(want.LastStartup))
	} else if !reflect.DeepEqual(resp.Link, want.Link) {
		t.Errorf("InstanceAdmin Link = %+v, Want = %+v", spew.Sdump(resp.Link), spew.Sdump(want.Link))
	} else if !reflect.DeepEqual(*resp, want) {
		t.Errorf("InstanceAdmin Response = %+v, Want = %+v", spew.Sdump(*resp), spew.Sdump(want))
	}
}
