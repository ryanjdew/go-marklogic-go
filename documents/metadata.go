package documents

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// MetadataHandle describes a document to write
type MetadataHandle struct {
	*bytes.Buffer
	VersionID int
	Format    int
	metadata  Metadata
	timestamp string
}

// Metadata describes a document to write
type Metadata struct {
	*bytes.Buffer
	XMLName        xml.Name          `xml:"http://marklogic.com/rest-api metadata" json:"-"`
	Collections    []string          `xml:"http://marklogic.com/rest-api collections" json:"collections,omitempty"`
	Permissions    []Permission      `xml:"http://marklogic.com/rest-api permissions" json:"permissions,omitempty"`
	Properties     map[string]string `xml:"http://marklogic.com/xdmp/property properties" json:"properties,omitempty"`
	MetadataValues map[string]string `xml:"http://marklogic.com/rest-api metadata-values" json:"metadataValues,omitempty"`
	Quality        int               `xml:"http://marklogic.com/rest-api quality" json:"quality,omitempty"`
}

// GetFormat returns int that represents XML or JSON
func (mh *MetadataHandle) GetFormat() int {
	return mh.Format
}

func (mh *MetadataHandle) resetBuffer() {
	if mh.Buffer == nil {
		mh.Buffer = new(bytes.Buffer)
	}
	mh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (mh *MetadataHandle) Deserialize(bytes []byte) {
	if mh.GetFormat() == handle.XML {
		xml.Unmarshal(bytes, mh)
	} else {
		json.Unmarshal(bytes, mh)
	}
}

// Serialize returns []byte of XML or JSON that represents the Metadata struct
func (mh *MetadataHandle) Serialize(metadata interface{}) {
	mh.metadata = metadata.(Metadata)
	mh.resetBuffer()
	if mh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(mh)
		enc.Encode(mh.metadata)
	} else {
		enc := xml.NewEncoder(mh)
		enc.Encode(mh.metadata)
	}
}

// Serialized returns string of XML or JSON
func (mh *MetadataHandle) Serialized() string {
	buffer := &bytes.Buffer{}
	var err error
	if mh.GetFormat() == handle.XML {
		enc := xml.NewEncoder(buffer)
		err = enc.Encode(mh.metadata)
	} else {
		enc := json.NewEncoder(buffer)
		err = enc.Encode(mh.metadata)
	}
	if err != nil {
		panic(err)
	}
	return string(buffer.Bytes())
}

// SetTimestamp sets the timestamp
func (mh *MetadataHandle) SetTimestamp(timestamp string) {
	mh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (mh *MetadataHandle) Timestamp() string {
	return mh.timestamp
}

// PermissionsMap associated with the MetadataHandle
func (m *Metadata) PermissionsMap() map[string]string {
	permissionsMap := make(map[string]string)
	for _, permission := range m.Permissions {
		currentValue, ok := permissionsMap[permission.RoleName]
		if ok {
			permissionsMap[permission.RoleName] = currentValue + "," + strings.Join(permission.Capability, ",")
		} else {
			permissionsMap[permission.RoleName] = strings.Join(permission.Capability, ",")
		}
	}
	return permissionsMap
}

// Collections describes collections a document belongs to
type Collections struct {
	Values []string `xml:"http://marklogic.com/rest-api collection" json:"-"`
}

// MetadataValue describes collections a document belongs to
type MetadataValue struct {
	XMLName xml.Name `xml:"http://marklogic.com/rest-api metadata-value" json:"-"`
	Key     string   `xml:"key,attr" json:"key,omitempty"`
	Value   string   `xml:",chardata" json:"value,omitempty"`
}

// Permission describes a permission assigned to a document
type Permission struct {
	XMLName    xml.Name `xml:"http://marklogic.com/rest-api permission" json:"-"`
	RoleName   string   `xml:"http://marklogic.com/rest-api role-name" json:"role-name,omitempty"`
	Capability []string `xml:"http://marklogic.com/rest-api capability" json:"capabilities,omitempty"`
}
