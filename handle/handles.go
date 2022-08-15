package goMarklogicGo

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

// Format options
const (
	JSON = iota
	XML
	MIXED
	TEXTPLAIN
	TEXT_URI_LIST
	TEXTHTML
	UNKNOWN
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
	} else if formatEnum == TEXT_URI_LIST {
		formatStr = "text/uri-list"
	} else if formatEnum == XML {
		formatStr = "application/xml"
	} else if formatEnum == TEXTHTML {
		formatStr = "text/html"
	} else {
		formatStr = "application/octet-stream"
	}
	return formatStr
}

// Handle interface
type Handle interface {
	io.ReadWriter
	GetFormat() int
	Deserialize([]byte)
	Deserialized() interface{}
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

// Deserialized returns string of XML or JSON as interface
func (r *RawHandle) Deserialized() interface{} {
	return r.String()
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

// Deserialized returns map as interface
func (m *MapHandle) Deserialized() interface{} {
	return m.mapItem
}

// AcceptResponse handles an *http.Response
func (m *MapHandle) AcceptResponse(resp *http.Response) error {
	return CommonHandleAcceptResponse(m, resp)
}

// Serialize returns the bytes that represent XML or JSON
func (m *MapHandle) Serialize(mapItem interface{}) {
	m.mapItem = mapItem.(*map[string]interface{})
}

// Get returns the map[string]interface{} form
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

// MultipartResponseHandle is a handle that places the results into
// a Response struct
type MultipartResponseHandle struct {
	*bytes.Buffer
	Format    int
	response  [][]byte
	timestamp string
}

// GetFormat returns int that represents XML or JSON
func (rh *MultipartResponseHandle) GetFormat() int {
	return MIXED
}

func (rh *MultipartResponseHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *MultipartResponseHandle) Deserialize(bytes []byte) {
	rh.resetBuffer()
	rh.Write(bytes)
	rh.response = append(rh.response, bytes)
}

// Deserialized returns [][]byte as interface{}
func (rh *MultipartResponseHandle) Deserialized() interface{} {
	return rh.response
}

// AcceptResponse handles an *http.Response
func (rh *MultipartResponseHandle) AcceptResponse(resp *http.Response) error {
	defer resp.Body.Close()
	if resp.ContentLength == 0 {
		return nil
	}
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(resp.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			output, err := io.ReadAll(p)
			if err != nil {
				return err
			}
			rh.Deserialize(output)
			p.Close()
		}
	}
	return err
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *MultipartResponseHandle) Serialize(response interface{}) {
	rh.response = response.([][]byte)
	rh.resetBuffer()
}

// Get returns string of XML or JSON
func (rh *MultipartResponseHandle) Get() [][]byte {
	return rh.response
}

// Serialized returns string of XML or JSON
func (rh *MultipartResponseHandle) Serialized() string {
	rh.Serialize(rh.response)
	return rh.String()
}

// SetTimestamp sets the timestamp
func (rh *MultipartResponseHandle) SetTimestamp(timestamp string) {
	rh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (rh *MultipartResponseHandle) Timestamp() string {
	return rh.timestamp
}

// CommonHandleAcceptResponse handles an HTTP response
func CommonHandleAcceptResponse(genericHandle Handle, response *http.Response) error {
	var err error
	var contents []byte
	if response != nil {
		defer response.Body.Close()
		contents, err = io.ReadAll(response.Body)
		genericHandle.Deserialize(contents)
		genericHandle.SetTimestamp(response.Header.Get("ML-Effective-Timestamp"))
	}
	return err
}
