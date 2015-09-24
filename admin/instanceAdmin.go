package admin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	_ "strconv"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// InstanceAdminRequest represents the properties of a MarkLogic Database
type InstanceAdminRequest struct {
	XMLName  xml.Name `xml:"http://marklogic.com/manage instance-admin" json:"-"`
	Username string   `xml:"admin-username" json:"admin-username"`
	Password string   `xml:"admin-password" json:"admin-password"`
	Realm    string   `xml:"realm" json:"realm"`
}

// RawHandle returns the raw string results of JSON or XML
type InstanceAdminHandle struct {
	*bytes.Buffer
	Format  int
	request InstanceAdminRequest
}

// GetFormat returns int that represents XML or JSON
func (rh *InstanceAdminHandle) GetFormat() int {
	return rh.Format
}

func (rh *InstanceAdminHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *InstanceAdminHandle) Deserialize(bytes []byte) {
	rh.resetBuffer()
	rh.Write(bytes)
	rh.request = InstanceAdminRequest{}
	if rh.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, &rh.request)
	} else {
		xml.Unmarshal(bytes, &rh.request)
	}
}

// AcceptResponse handles an *http.Response
func (rh *InstanceAdminHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(rh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *InstanceAdminHandle) Serialize(response interface{}) {
	rh.request = response.(InstanceAdminRequest)
	rh.resetBuffer()
	if rh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(rh.Buffer)
		enc.Encode(&rh.request)
	} else {
		enc := xml.NewEncoder(rh.Buffer)
		enc.Encode(&rh.request)
	}
}

// Get returns string of XML or JSON
func (rh *InstanceAdminHandle) Get() *InstanceAdminRequest {
	return &rh.request
}

// Serialized returns string of XML or JSON
func (rh *InstanceAdminHandle) Serialized() string {
	rh.Serialize(rh.request)
	return rh.String()
}

// Install the admin username and password, and initialize the security database and objects.
func instanceAdmin(ac *clients.AdminClient, username string, password string, realm string, response handle.ResponseHandle) error {
	params := "?"
	params = util.RepeatingParameters(params, "admin-username", []string{username})
	params = util.RepeatingParameters(params, "admin-password", []string{password})
	params = util.RepeatingParameters(params, "realm", []string{realm})
	req, err := util.BuildRequestFromHandle(ac, "POST", "/instance-admin"+params, nil)
	if err != nil {
		return err
	}
	return util.Execute(ac, req, response)
}
