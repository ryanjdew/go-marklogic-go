package goMarklogicGo

import (
	"bytes"
	"io"
)

// Format options
const (
	XML = iota
	JSON
	MIXED
)

// FormatEnumToMimeType converts a format enum to a mime/type value for the REST API
func FormatEnumToMimeType(formatEnum int) string {
	var formatStr string
	if formatEnum == JSON {
		formatStr = "application/json"
	} else if formatEnum == MIXED {
		formatStr = "multipart/mixed"
	} else {
		formatStr = "application/xml"
	}
	return formatStr
}

// Handle interface
type Handle interface {
	io.ReadWriter
	GetFormat() int
	Encode([]byte)
	Decode(interface{})
	Serialized() string
}

// RawHandle returns the raw string results of JSON or XML
type RawHandle struct {
	*bytes.Buffer
	Format int
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

// Encode returns the bytes that represent XML or JSON
func (r *RawHandle) Encode(bytes []byte) {
	r.resetBuffer()
	r.Write(bytes)
}

// Decode returns the bytes that represent XML or JSON
func (r *RawHandle) Decode(bytes interface{}) {
	r.Encode(bytes.([]byte))
}

// Get returns string of XML or JSON
func (r *RawHandle) Get() string {
	return r.String()
}

// Serialized returns string of XML or JSON
func (r *RawHandle) Serialized() string {
	return r.Get()
}

// MapHandle returns the raw string results of JSON or XML
type MapHandle struct {
	*bytes.Buffer
	Format  int
	mapItem *map[string]interface{}
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

// Encode returns the bytes that represent XML or JSON
func (m *MapHandle) Encode(bytes []byte) {
	m.resetBuffer()
	m.Write(bytes)
}

// Decode returns the bytes that represent XML or JSON
func (m *MapHandle) Decode(mapItem interface{}) {
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
