package admin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	_ "strconv"

	clients "github.com/rlouapre/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// InitializeProperties represents the properties of a MarkLogic Database
type InitializeProperties struct {
	XMLName    xml.Name `xml:"http://marklogic.com/manage init" json:"-"`
	LicenseKey string   `xml:"http://marklogic.com/manage license-key" json:"license-key"`
	Licensee   string   `xml:"http://marklogic.com/manage licensee" json:"licensee"`
}

// RawHandle returns the raw string results of JSON or XML
type AdminHandle struct {
	*bytes.Buffer
	Format               int
	initializeProperties InitializeProperties
}

// GetFormat returns int that represents XML or JSON
func (rh *AdminHandle) GetFormat() int {
	return rh.Format
}

func (rh *AdminHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *AdminHandle) Deserialize(bytes []byte) {
	rh.resetBuffer()
	rh.Write(bytes)
	rh.initializeProperties = InitializeProperties{}
	if rh.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, &rh.initializeProperties)
	} else {
		xml.Unmarshal(bytes, &rh.initializeProperties)
	}
}

// AcceptResponse handles an *http.Response
func (rh *AdminHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(rh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *AdminHandle) Serialize(response interface{}) {
	rh.initializeProperties = response.(InitializeProperties)
	rh.resetBuffer()
	if rh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(rh.Buffer)
		enc.Encode(&rh.initializeProperties)
	} else {
		enc := xml.NewEncoder(rh.Buffer)
		enc.Encode(&rh.initializeProperties)
	}
}

// Get returns string of XML or JSON
func (rh *AdminHandle) Get() *InitializeProperties {
	return &rh.initializeProperties
}

// Serialized returns string of XML or JSON
func (rh *AdminHandle) Serialized() string {
	rh.Serialize(rh.initializeProperties)
	return rh.String()
}

// Response represents a response from the search API
type Response struct {
	XMLName     xml.Name           `xml:"http://marklogic.com/manage restart" json:"-"`
	LastStartup LastStartupElement `xml:"http://marklogic.com/manage last-startup" json:"last-startup,omitempty"`
	Link        LinkElement        `xml:"http://marklogic.com/manage link" json:"link,omitempty"`
	Message     string             `xml:"http://marklogic.com/manage message" json:"link,omitempty"`
}

type LastStartupElement struct {
	XMLName xml.Name `xml:"http://marklogic.com/manage last-startup" json:"-"`
	Value   string   `xml:",chardata"`
	HostId  string   `xml:"host-id,attr"`
}

type LinkElement struct {
	XMLName xml.Name `xml:"http://marklogic.com/manage link" json:"-"`
	KindRef string   `xml:"http://marklogic.com/manage kindref"`
	UriRef  string   `xml:"http://marklogic.com/manage uriref"`
}

// Initialize MarkLogic instance
func initialize(ac *clients.AdminClient, license handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(ac, "POST", "/init", license)
	if err != nil {
		return err
	}
	return util.Execute(ac, req, response)
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
