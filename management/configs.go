package management

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// DatabaseConfig is a minimal representation of a database creation payload
type DatabaseConfig struct {
	XMLName          xml.Name `xml:"http://marklogic.com/manage database" json:"-"`
	DatabaseName     string   `xml:"http://marklogic.com/manage database-name" json:"database-name"`
	Forests          []string `xml:"http://marklogic.com/manage forest" json:"forest"`
	SecurityDatabase string   `xml:"http://marklogic.com/manage security-database" json:"security-database"`
	SchemaDatabase   string   `xml:"http://marklogic.com/manage schema-database" json:"schema-database"`
	Enabled          bool     `xml:"http://marklogic.com/manage enabled" json:"enabled"`
}

// ForestConfig is a minimal representation of a forest creation payload
type ForestConfig struct {
	XMLName    xml.Name `xml:"http://marklogic.com/manage forest" json:"-"`
	ForestName string   `xml:"http://marklogic.com/manage forest-name" json:"forest-name"`
	DataDir    string   `xml:"http://marklogic.com/manage data-dir" json:"data-dir"`
	Host       string   `xml:"http://marklogic.com/manage host" json:"host"`
}

// DatabaseConfigHandle serializes/deserializes DatabaseConfig
type DatabaseConfigHandle struct {
	*bytes.Buffer
	Format    int
	config    DatabaseConfig
	timestamp string
}

func (h *DatabaseConfigHandle) GetFormat() int { return h.Format }

func (h *DatabaseConfigHandle) resetBuffer() {
	if h.Buffer == nil {
		h.Buffer = new(bytes.Buffer)
	}
	h.Reset()
}

func (h *DatabaseConfigHandle) Serialize(v interface{}) {
	h.config = v.(DatabaseConfig)
	h.resetBuffer()
	if h.GetFormat() == handle.JSON {
		enc := json.NewEncoder(h.Buffer)
		enc.Encode(h.config)
	} else {
		enc := xml.NewEncoder(h.Buffer)
		enc.Encode(h.config)
	}
}

func (h *DatabaseConfigHandle) Deserialize(b []byte) {
	h.resetBuffer()
	h.Write(b)
	h.config = DatabaseConfig{}
	if h.GetFormat() == handle.JSON {
		json.Unmarshal(b, &h.config)
	} else {
		xml.Unmarshal(b, &h.config)
	}
}

func (h *DatabaseConfigHandle) Deserialized() interface{} { return h.config }

func (h *DatabaseConfigHandle) Serialized() string {
	h.Serialize(h.config)
	return h.String()
}

// For request handles we don't need AcceptResponse, but implement to satisfy ResponseHandle when used
func (h *DatabaseConfigHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(h, resp)
}

// SetTimestamp sets the timestamp
func (h *DatabaseConfigHandle) SetTimestamp(timestamp string) {
	h.timestamp = timestamp
}

// Timestamp retrieves the timestamp
func (h *DatabaseConfigHandle) Timestamp() string {
	return h.timestamp
}

// ForestConfigHandle serializes/deserializes ForestConfig
type ForestConfigHandle struct {
	*bytes.Buffer
	Format    int
	config    ForestConfig
	timestamp string
}

func (h *ForestConfigHandle) GetFormat() int { return h.Format }

func (h *ForestConfigHandle) resetBuffer() {
	if h.Buffer == nil {
		h.Buffer = new(bytes.Buffer)
	}
	h.Reset()
}

func (h *ForestConfigHandle) Serialize(v interface{}) {
	h.config = v.(ForestConfig)
	h.resetBuffer()
	if h.GetFormat() == handle.JSON {
		enc := json.NewEncoder(h.Buffer)
		enc.Encode(h.config)
	} else {
		enc := xml.NewEncoder(h.Buffer)
		enc.Encode(h.config)
	}
}

func (h *ForestConfigHandle) Deserialize(b []byte) {
	h.resetBuffer()
	h.Write(b)
	h.config = ForestConfig{}
	if h.GetFormat() == handle.JSON {
		json.Unmarshal(b, &h.config)
	} else {
		xml.Unmarshal(b, &h.config)
	}
}

func (h *ForestConfigHandle) Deserialized() interface{} { return h.config }

func (h *ForestConfigHandle) Serialized() string {
	h.Serialize(h.config)
	return h.String()
}

func (h *ForestConfigHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(h, resp)
}

// SetTimestamp sets the timestamp
func (h *ForestConfigHandle) SetTimestamp(timestamp string) {
	h.timestamp = timestamp
}

// Timestamp retrieves the timestamp
func (h *ForestConfigHandle) Timestamp() string {
	return h.timestamp
}
