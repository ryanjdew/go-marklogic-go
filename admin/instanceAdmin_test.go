package admin

import (
	"encoding/xml"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	testHelper "github.com/ryanjdew/go-marklogic-go/test"
)

func TestXMLInstanceAdminRequestSerialize(t *testing.T) {
	want := `<instance-admin xmlns="http://marklogic.com/manage"><admin-username>admin</admin-username><admin-password>password</admin-password><realm>public</realm></instance-admin>`
	request :=
		&InstanceAdminRequest{
			XMLName: xml.Name{
				Space: "http://marklogic.com/manage",
				Local: "instance-admin",
			},
			Username: "admin",
			Password: "password",
			Realm:    "public",
		}
	qh := &InstanceAdminHandle{Format: handle.XML}
	testHelper.RoundTripSerialization(t, "InstanceAdminRequest", request, qh, want)
}

func TestJsonInstanceAdminRequestSerialize(t *testing.T) {
	want := `{"admin-username":"admin","admin-password":"password","realm":"public"}`
	request :=
		&InstanceAdminRequest{
			Username: "admin",
			Password: "password",
			Realm:    "public",
		}
	qh := &InstanceAdminHandle{Format: handle.JSON}
	testHelper.RoundTripSerialization(t, "InstanceAdminRequest", request, qh, want)
}
