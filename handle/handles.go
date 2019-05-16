package goMarklogicGo

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// Format options
const (
	JSON = iota
	XML
	MIXED
	TEXTPLAIN
)

// FormatEnumToMimeType converts a format enum to a mime/type value for the REST API
func FormatEnumToMimeType(formatEnum int) string {
	var formatStr string
	if formatEnum == JSON {
		formatStr = "application/json"
	} else if formatEnum == MIXED {
		formatStr = "multipart/mixed"
	} else if formatEnum == TEXTPLAIN {
		formatStr = "text/plain"
	} else {
		formatStr = "application/xml"
	}
	return formatStr
}

// Handle interface
type Handle interface {
	io.ReadWriter
	GetFormat() int
	Deserialize([]byte)
	Serialize(interface{})
	Serialized() string
	SetTimestamp(string)
	Timestamp() string
}

// ResponseHandle interface
type ResponseHandle interface {
	Handle
	AcceptResponse(*http.Response) error
}

// RawHandle returns the raw string results of JSON or XML
type RawHandle struct {
	*bytes.Buffer
	timestamp string
	Format    int
}

// GetFormat returns int that represents XML or JSON
func (r *RawHandle) GetFormat() int {
	return r.Format
}

func (r *RawHandle) resetBuffer() {
	if r.Buffer == nil {
		r.Buffer = new(bytes.Buffer)
	}
	r.Reset()
}

// Deserialize returns the bytes that represent XML or JSON
func (r *RawHandle) Deserialize(bytes []byte) {
	r.resetBuffer()
	r.Write(bytes)
}

// AcceptResponse handles an *http.Response
func (r *RawHandle) AcceptResponse(resp *http.Response) error {
	return CommonHandleAcceptResponse(r, resp)
}

// Serialize returns the bytes that represent XML or JSON
func (r *RawHandle) Serialize(bytes interface{}) {
	r.Deserialize(bytes.([]byte))
}

// Get returns string of XML or JSON
func (r *RawHandle) Get() string {
	return r.String()
}

// Serialized returns string of XML or JSON
func (r *RawHandle) Serialized() string {
	return r.Get()
}

// SetTimestamp sets the timestamp
func (r *RawHandle) SetTimestamp(timestamp string) {
	r.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (r *RawHandle) Timestamp() string {
	return r.timestamp
}

// MapHandle returns the raw string results of JSON or XML
type MapHandle struct {
	*bytes.Buffer
	Format    int
	timestamp string
	mapItem   *map[string]interface{}
}

// GetFormat returns int that represents XML or JSON
func (m *MapHandle) GetFormat() int {
	return m.Format
}

func (m *MapHandle) resetBuffer() {
	if m.Buffer == nil {
		m.Buffer = new(bytes.Buffer)
	}
	m.Reset()
}

// Deserialize returns the bytes that represent XML or JSON
func (m *MapHandle) Deserialize(bytes []byte) {
	m.resetBuffer()
	m.Write(bytes)
}

// AcceptResponse handles an *http.Response
func (m *MapHandle) AcceptResponse(resp *http.Response) error {
	return CommonHandleAcceptResponse(m, resp)
}

// Serialize returns the bytes that represent XML or JSON
func (m *MapHandle) Serialize(mapItem interface{}) {
	m.mapItem = mapItem.(*map[string]interface{})
}

// Get returns string of XML or JSON
func (m *MapHandle) Get() *map[string]interface{} {
	return m.mapItem
}

// Serialized returns string of XML or JSON
func (m *MapHandle) Serialized() string {
	return m.String()
}

// SetTimestamp sets the timestamp
func (m *MapHandle) SetTimestamp(timestamp string) {
	m.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (m *MapHandle) Timestamp() string {
	return m.timestamp
}

// CommonHandleAcceptResponse handles an HTTP response
func CommonHandleAcceptResponse(genericHandle Handle, response *http.Response) error {
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	genericHandle.Deserialize(contents)
	genericHandle.SetTimestamp(response.Header.Get("ML-Effective-Timestamp"))
	return err
}
