package admin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// InitializeProperties represents a request for /admin/v1/init API
type InitializeProperties struct {
	XMLName    xml.Name `xml:"http://marklogic.com/manage init" json:"-"`
	LicenseKey string   `xml:"license-key" json:"license-key"`
	Licensee   string   `xml:"licensee" json:"licensee"`
}

// InitHandle returns the raw string results of JSON or XML
type InitHandle struct {
	*bytes.Buffer
	Format               int
	initializeProperties InitializeProperties
	timestamp            string
}

// RestartResponseHandle is a handle that places the results into
// a Response struct
type RestartResponseHandle struct {
	*bytes.Buffer
	Format    int
	response  RestartResponse
	timestamp string
}

// GetFormat returns int that represents XML or JSON
func (rh *RestartResponseHandle) GetFormat() int {
	return rh.Format
}

func (rh *RestartResponseHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *RestartResponseHandle) Deserialize(bytes []byte) {
	rh.resetBuffer()
	rh.Write(bytes)
	rh.response = RestartResponse{}
	if rh.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, &rh.response)
	} else {
		xml.Unmarshal(bytes, &rh.response)
	}
}

// AcceptResponse handles an *http.Response
func (rh *RestartResponseHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(rh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *RestartResponseHandle) Serialize(response interface{}) {
	rh.response = response.(RestartResponse)
	rh.resetBuffer()
	if rh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(rh.Buffer)
		enc.Encode(&rh.response)
	} else {
		enc := xml.NewEncoder(rh.Buffer)
		enc.Encode(&rh.response)
	}
}

// Get returns string of XML or JSON
func (rh *RestartResponseHandle) Get() *RestartResponse {
	return &rh.response
}

// Serialized returns string of XML or JSON
func (rh *RestartResponseHandle) Serialized() string {
	rh.Serialize(rh.response)
	return rh.String()
}

// SetTimestamp sets the timestamp
func (rh *RestartResponseHandle) SetTimestamp(timestamp string) {
	rh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (rh *RestartResponseHandle) Timestamp() string {
	return rh.timestamp
}

// RestartResponse represents a response from /admin/v1/init API
type RestartResponse struct {
	XMLName     xml.Name           `xml:"http://marklogic.com/manage restart" json:"-"`
	LastStartup LastStartupElement `xml:"last-startup" json:"last-startup,omitempty"`
	Link        LinkElement        `xml:"link" json:"link,omitempty"`
	Message     string             `xml:"message" json:"link,omitempty"`
	timestamp   string
}

// LastStartupElement represents MarkLogic last startup information
type LastStartupElement struct {
	XMLName xml.Name `xml:"last-startup" json:"-"`
	Value   string   `xml:",chardata"`
	HostID  string   `xml:"host-id,attr"`
}

// LinkElement represents link information
type LinkElement struct {
	XMLName xml.Name `xml:"link" json:"-"`
	KindRef string   `xml:"kindref"`
	URIRef  string   `xml:"uriref"`
}

// GetFormat returns int that represents XML or JSON
func (rh *InitHandle) GetFormat() int {
	return rh.Format
}

func (rh *InitHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *InitHandle) Deserialize(bytes []byte) {
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
func (rh *InitHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(rh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *InitHandle) Serialize(response interface{}) {
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
func (rh *InitHandle) Get() *InitializeProperties {
	return &rh.initializeProperties
}

// Serialized returns string of XML or JSON
func (rh *InitHandle) Serialized() string {
	rh.Serialize(rh.initializeProperties)
	return rh.String()
}

// SetTimestamp sets the timestamp
func (rh *InitHandle) SetTimestamp(timestamp string) {
	rh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (rh *InitHandle) Timestamp() string {
	return rh.timestamp
}

// Initialize MarkLogic instance
func initialize(ac *clients.AdminClient, license handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(ac, "POST", "/init", license)
	if err != nil {
		return err
	}
	return util.Execute(ac, req, response)
}
