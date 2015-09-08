package admin

import (
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/rlouapre/go-marklogic-go/test"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

func TestXMLInstanceAdminRequestSerialize(t *testing.T) {
	want := `<instance-admin xmlns="http://marklogic.com/manage"><admin-username>admin</admin-username><admin-password>password</admin-password><realm>public</realm></instance-admin>`
	request :=
		InstanceAdminRequest{
			Username: "admin",
			Password: "password",
			Realm:    "public",
		}
	qh := InstanceAdminHandle{}
	qh.Serialize(request)
	result := qh.Serialized()
	if want != result {
		t.Errorf("InstanceAdminRequest Results = %+v, Want = %+v", result, want)
	} else if !reflect.DeepEqual(&request, qh.Get()) {
		t.Errorf("InstanceAdminRequest Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}

func TestJsonInstanceAdminRequestSerialize(t *testing.T) {
	want := `{"admin-username":"admin","admin-password":"password","realm":"public"}`
	request :=
		InstanceAdminRequest{
			Username: "admin",
			Password: "password",
			Realm:    "public",
		}
	qh := InstanceAdminHandle{Format: handle.JSON}
	qh.Serialize(request)
	result := strings.TrimSpace(qh.Serialized())
	if want != result {
		t.Errorf("InstanceAdminRequest Results = %+v, Want = %+v", result, want)
	} else if !reflect.DeepEqual(&request, qh.Get()) {
		t.Errorf("InstanceAdminRequest Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}
