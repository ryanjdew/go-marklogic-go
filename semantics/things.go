package semantics

import (
	"bytes"
	"encoding/xml"
	"net/http"

	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	"github.com/cchatfield/go-marklogic-go/util"
)

// ThingsHandle is a handle that places the results into
// a Response struct
type ThingsHandle struct {
	*bytes.Buffer
	things    Things
	timestamp string
}

// GetFormat returns int that represents XML or JSON
func (th *ThingsHandle) GetFormat() int {
	return handle.TEXTHTML
}

func (th *ThingsHandle) resetBuffer() {
	if th.Buffer == nil {
		th.Buffer = new(bytes.Buffer)
	}
	th.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (th *ThingsHandle) Deserialize(bytes []byte) {
	th.resetBuffer()
	th.Write(bytes)
	th.things = Things{}
	xml.Unmarshal(bytes, &th.things)
}

// AcceptResponse handles an *http.Response
func (th *ThingsHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(th, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (th *ThingsHandle) Serialize(response interface{}) {
	th.things = response.(Things)
	th.resetBuffer()
	enc := xml.NewEncoder(th.Buffer)
	enc.Encode(&th.things)
}

// Get returns string of XML or JSON
func (th *ThingsHandle) Get() *Things {
	return &th.things
}

// Serialized returns string of XML or JSON
func (th *ThingsHandle) Serialized() string {
	th.Serialize(th.things)
	return th.String()
}

// SetTimestamp sets the timestamp
func (th *ThingsHandle) SetTimestamp(timestamp string) {
	th.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (th *ThingsHandle) Timestamp() string {
	return th.timestamp
}

// Things represents a response from the things API
type Things struct {
	Subjects []string `xml:"http://www.w3.org/1999/xhtml a"`
}

func things(c *clients.Client, iris []string, response handle.ResponseHandle) error {
	params := util.AddDatabaseParam(util.RepeatingParameters("?", "iri", iris), c)
	req, err := http.NewRequest("GET", c.Base()+"/graphs/things"+params, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
