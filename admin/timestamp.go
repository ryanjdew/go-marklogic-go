package admin

import (
	"bytes"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// TimestampResponseHandle is a handle that places the results into
// a Response struct
type TimestampResponseHandle struct {
	*bytes.Buffer
	Format    int
	timestamp string
}

// GetFormat returns int that represents XML or JSON
func (rh *TimestampResponseHandle) GetFormat() int {
	return rh.Format
}

func (rh *TimestampResponseHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *TimestampResponseHandle) Deserialize(bytes []byte) {
	rh.resetBuffer()
	rh.Write(bytes)
	if rh.GetFormat() == handle.TEXTPLAIN {
		rh.timestamp = string(bytes)
	}
}

// Deserialized returns deserialized timestamp string as interface{}
func (rh *TimestampResponseHandle) Deserialized() any {
	return rh.timestamp
}

// AcceptResponse handles an *http.Response
func (rh *TimestampResponseHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(rh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *TimestampResponseHandle) Serialize(response any) {
	rh.resetBuffer()
	if rh.GetFormat() == handle.TEXTPLAIN {
		rh.timestamp = response.(string)
		rh.WriteString(rh.timestamp)
	}
}

// Get returns string of timestamp
func (rh *TimestampResponseHandle) Get() *string {
	return &rh.timestamp
}

// Serialized returns string of XML or JSON
func (rh *TimestampResponseHandle) Serialized() string {
	rh.Serialize(rh.timestamp)
	log.Println("Serialized = " + spew.Sdump(rh.timestamp))
	return rh.String()
}

// SetTimestamp sets the timestamp
func (rh *TimestampResponseHandle) SetTimestamp(timestamp string) {
	rh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (rh *TimestampResponseHandle) Timestamp() string {
	return rh.timestamp
}

// Verify that MarkLogic Server is up and accepting requests.
// https://docs.marklogic.com/REST/GET/admin/v1/timestamp
func timestamp(ac *clients.AdminClient, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(ac, "GET", "/timestamp", nil)
	if err != nil {
		return err
	}
	return util.Execute(ac, req, response)
}
