package temporal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// TemporalAxis represents a temporal axis configuration
type TemporalAxis struct {
	AxisName  string `json:"axis-name" xml:"axis-name"`
	ScalarURI string `json:"scalar-uri,omitempty" xml:"scalar-uri,omitempty"`
}

// TemporalCollection represents a collection enabled for temporal operations
type TemporalCollection struct {
	CollectionURI       string `json:"collection-uri" xml:"collection-uri"`
	CollectionSystemURI string `json:"collection-system-uri" xml:"collection-system-uri"`
	SystemAxis          string `json:"system-axis" xml:"system-axis"`
	ValidAxis           string `json:"valid-axis,omitempty" xml:"valid-axis,omitempty"`
}

// TemporalDocument represents metadata about a temporal document
type TemporalDocument struct {
	URI         string `json:"uri" xml:"uri"`
	ValidStart  string `json:"valid-start" xml:"valid-start"`
	ValidEnd    string `json:"valid-end" xml:"valid-end"`
	SystemStart string `json:"system-start" xml:"system-start"`
	SystemEnd   string `json:"system-end" xml:"system-end"`
	VersionID   string `json:"version-id" xml:"version-id"`
}

// SystemTime represents the current system time for temporal operations
type SystemTime struct {
	Timestamp string `json:"timestamp" xml:"timestamp"`
}

// TemporalHandle implements handle.ResponseHandle for temporal responses
type TemporalHandle struct {
	*handle.RawHandle
}

// NewTemporalHandle creates a new TemporalHandle with the specified format
func NewTemporalHandle(format int) *TemporalHandle {
	return &TemporalHandle{
		RawHandle: &handle.RawHandle{
			Format: format,
		},
	}
}

// AcceptResponse handles HTTP response deserialization
func (h *TemporalHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(h, resp)
}

func createAxis(c *clients.Client, axisName string, requestBody handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "POST", "/temporal/axes/"+url.QueryEscape(axisName), requestBody)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func getAxis(c *clients.Client, axisName string, response handle.ResponseHandle) error {
	req, err := http.NewRequest("GET", c.Base()+"/temporal/axes/"+url.QueryEscape(axisName), nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func listAxes(c *clients.Client, response handle.ResponseHandle) error {
	req, err := http.NewRequest("GET", c.Base()+"/temporal/axes", nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func deleteAxis(c *clients.Client, axisName string, response handle.ResponseHandle) error {
	req, err := http.NewRequest("DELETE", c.Base()+"/temporal/axes/"+url.QueryEscape(axisName), nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func enableCollectionTemporal(c *clients.Client, collection string, temporalConfig handle.Handle, response handle.ResponseHandle) error {
	params := "?"
	params = params + "collection=" + url.QueryEscape(collection)
	req, err := util.BuildRequestFromHandle(c, "POST", "/temporal/collections"+params, temporalConfig)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func disableCollectionTemporal(c *clients.Client, collection string, response handle.ResponseHandle) error {
	params := "?"
	params = params + "collection=" + url.QueryEscape(collection)
	req, err := http.NewRequest("DELETE", c.Base()+"/temporal/collections"+params, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func getTemporalDocument(c *clients.Client, uri string, timestamp string, response handle.ResponseHandle) error {
	params := "?"
	params = params + "uri=" + url.QueryEscape(uri)
	if timestamp != "" {
		params = params + "&timestamp=" + url.QueryEscape(timestamp)
	}
	req, err := http.NewRequest("GET", c.Base()+"/temporal/documents"+params, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", handle.FormatEnumToMimeType(response.GetFormat()))
	return util.Execute(c, req, response)
}

func advanceSystemTime(c *clients.Client, timestamp string, response handle.ResponseHandle) error {
	payload := map[string]string{"timestamp": timestamp}
	jsonPayload, _ := json.Marshal(payload)
	reqHandle := &handle.RawHandle{
		Format: handle.JSON,
		Buffer: bytes.NewBuffer(jsonPayload),
	}

	req, err := util.BuildRequestFromHandle(c, "POST", "/temporal/system-time", reqHandle)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func getSystemTime(c *clients.Client, response handle.ResponseHandle) error {
	req, err := http.NewRequest("GET", c.Base()+"/temporal/system-time", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", handle.FormatEnumToMimeType(response.GetFormat()))
	return util.Execute(c, req, response)
}
