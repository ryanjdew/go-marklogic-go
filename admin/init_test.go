package admin

import (
	"encoding/xml"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	testHelper "github.com/ryanjdew/go-marklogic-go/test"
)

func TestXMLInitializePropertiesSerialize(t *testing.T) {
	want := `<init xmlns="http://marklogic.com/manage"><license-key>1234-5678-90AB</license-key><licensee>Your Licensee</licensee></init>`
	initializeProperties :=
		&InitializeProperties{
			XMLName: xml.Name{
				Space: "http://marklogic.com/manage",
				Local: "init",
			},
			LicenseKey: "1234-5678-90AB",
			Licensee:   "Your Licensee",
		}
	qh := &InitHandle{Format: handle.XML}
	testHelper.RoundTripSerialization(t, "InitHandle XML", initializeProperties, qh, want)
}

func TestJsonInitializePropertiesSerialize(t *testing.T) {
	want := `{"license-key":"1234-5678-90AB","licensee":"Your Licensee"}`
	initializeProperties :=
		&InitializeProperties{
			LicenseKey: "1234-5678-90AB",
			Licensee:   "Your Licensee",
		}
	qh := &InitHandle{Format: handle.JSON}
	testHelper.RoundTripSerialization(t, "InitHandle JSON", initializeProperties, qh, want)
}

func TestXMLRestartResponseSerialize(t *testing.T) {
	want := "<restart xmlns=\"http://marklogic.com/manage\"><last-startup host-id=\"hostId\">yesterday</last-startup><link><kindref>kindRef</kindref><uriref>uriRef</uriref></link><message>message</message></restart>"
	value :=
		&RestartResponse{
			XMLName: xml.Name{
				Space: "http://marklogic.com/manage",
				Local: "restart",
			},
			LastStartup: LastStartupElement{
				XMLName: xml.Name{
					Space: "http://marklogic.com/manage",
					Local: "last-startup",
				},
				HostID: "hostId",
				Value:  "yesterday",
			},
			Link: LinkElement{
				XMLName: xml.Name{
					Space: "http://marklogic.com/manage",
					Local: "link",
				},
				URIRef:  "uriRef",
				KindRef: "kindRef",
			},
			Message: "message",
		}
	qh := &RestartResponseHandle{Format: handle.XML}
	testHelper.RoundTripSerialization(t, "RestartResponseHandle XML", value, qh, want)
}

func TestJsonRestartResponseSerialize(t *testing.T) {
	want := "{\"last-startup\":{\"Value\":\"yesterday\",\"HostID\":\"hostId\"},\"link\":{\"KindRef\":\"kindRef\",\"URIRef\":\"uriRef\"},\"message\":\"message\"}"
	value :=
		&RestartResponse{
			LastStartup: LastStartupElement{
				HostID: "hostId",
				Value:  "yesterday",
			},
			Link: LinkElement{
				URIRef:  "uriRef",
				KindRef: "kindRef",
			},
			Message: "message",
		}
	qh := &RestartResponseHandle{Format: handle.JSON}
	testHelper.RoundTripSerialization(t, "InitHandle JSON", value, qh, want)
}
