package datamovement

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	"github.com/cchatfield/go-marklogic-go/util"
)

// ForestInfoHandle is a handle that places the results into
// a ForestInfo struct
type ForestInfoHandle struct {
	*bytes.Buffer
	Format     int
	ForestInfo []ForestInfo
	timestamp  string
}

// GetFormat returns int that represents XML or JSON
func (fih *ForestInfoHandle) GetFormat() int {
	return fih.Format
}

func (fih *ForestInfoHandle) resetBuffer() {
	if fih.Buffer == nil {
		fih.Buffer = new(bytes.Buffer)
	}
	fih.Reset()
}

// Deserialize returns Query struct that represents XML or JSON
func (fih *ForestInfoHandle) Deserialize(bytes []byte) {
	fih.resetBuffer()
	fih.Write(bytes)
	fih.ForestInfo = []ForestInfo{}
	if fih.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, &fih.ForestInfo)
	} else {
		xml.Unmarshal(bytes, &fih.ForestInfo)
	}
}

// Serialize returns []byte of XML or JSON that represents the Query struct
func (fih *ForestInfoHandle) Serialize(forestInfo interface{}) {
	fih.ForestInfo = forestInfo.([]ForestInfo)
	fih.resetBuffer()
	if fih.GetFormat() == handle.JSON {
		enc := json.NewEncoder(fih)
		enc.Encode(fih.ForestInfo)
	} else {
		enc := xml.NewEncoder(fih)
		enc.Encode(fih.ForestInfo)
	}
}

// Read bytes
func (fih *ForestInfoHandle) Read(bytes []byte) (n int, err error) {
	if fih.Buffer == nil {
		fih.Serialize(fih.ForestInfo)
	}
	return fih.Buffer.Read(bytes)
}

// Get returns string of XML or JSON
func (fih *ForestInfoHandle) Get() *[]ForestInfo {
	return &fih.ForestInfo
}

// Serialized returns string of XML or JSON
func (fih *ForestInfoHandle) Serialized() string {
	fih.Serialize(fih.ForestInfo)
	return fih.String()
}

// SetTimestamp sets the timestamp
func (fih *ForestInfoHandle) SetTimestamp(timestamp string) {
	fih.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (fih *ForestInfoHandle) Timestamp() string {
	return fih.timestamp
}

// AcceptResponse handles an *http.Response
func (fih *ForestInfoHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(fih, resp)
}

// ForestInfo describes a forest with associated host info
type ForestInfo struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	UpdatesAllowed  string `json:"updatesAllowed,omitempty"`
	Database        string `json:"database,omitempty"`
	Host            string `json:"host,omitempty"`
	RequestHost     string `json:"requestHost,omitempty"`
	AlternateHost   string `json:"alternateHost,omitempty"`
	OpenReplicaHost string `json:"openReplicaHost,omitempty"`
}

// PreferredHost for the forest
func (fi *ForestInfo) PreferredHost() string {
	if fi.RequestHost != "" {
		return fi.RequestHost
	} else if fi.AlternateHost != "" {
		return fi.AlternateHost
	} else if fi.OpenReplicaHost != "" {
		return fi.OpenReplicaHost
	}
	return fi.Host
}

func getForestInfo(c *clients.Client) []ForestInfo {
	params := util.AddDatabaseParam("", c)
	req, _ := util.BuildRequestFromHandle(c, "GET", "/internal/forestinfo"+params, nil)
	forestInfoHandle := &ForestInfoHandle{}
	util.Execute(c, req, forestInfoHandle)
	return *forestInfoHandle.Get()
}

// URIsHandle for retrieving URIs from the internal/uris endpoint
type URIsHandle struct {
	*bytes.Buffer
	URIs      []string
	timestamp string
}

// GetFormat returns int that represents XML or JSON
func (uh *URIsHandle) GetFormat() int {
	return handle.TEXT_URI_LIST
}

func (uh *URIsHandle) resetBuffer() {
	if uh.Buffer == nil {
		uh.Buffer = new(bytes.Buffer)
	}
	uh.Reset()
}

// Deserialize returns Query struct that represents XML or JSON
func (uh *URIsHandle) Deserialize(bytes []byte) {
	uh.resetBuffer()
	uh.Write(bytes)
	uris := strings.Split(string(bytes), "\r\n")
	filteredURIs := make([]string, 0, len(uris))
	for _, uri := range uris {
		if uri != "" {
			filteredURIs = append(filteredURIs, uri)
		}
	}
	uh.URIs = filteredURIs
}

// Serialize returns []byte of XML or JSON that represents the Query struct
func (uh *URIsHandle) Serialize(uris interface{}) {
	uh.URIs = uris.([]string)
	uh.resetBuffer()
	uh.Write([]byte(strings.Join(uh.URIs, "\r\n")))
}

// Get returns string of URIs
func (uh *URIsHandle) Get() []string {
	return uh.URIs
}

// Serialized returns string of XML or JSON
func (uh *URIsHandle) Serialized() string {
	uh.Serialize(uh.URIs)
	return uh.String()
}

// SetTimestamp sets the timestamp
func (uh *URIsHandle) SetTimestamp(timestamp string) {
	uh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (uh *URIsHandle) Timestamp() string {
	return uh.timestamp
}

// AcceptResponse handles an *http.Response
func (uh *URIsHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(uh, resp)
}

func getURIs(
	c *clients.Client,
	query handle.Handle,
	forestName string,
	transaction *util.Transaction,
	start uint64,
	after string,
	pageLength uint,
	respHandle handle.ResponseHandle) error {
	params := "?"

	paramsMap := map[string]string{
		"forest-name": forestName,
		"after":       after,
		"pageLength":  strconv.FormatUint(uint64(pageLength), 10),
	}
	if start != 0 {
		paramsMap["start"] = strconv.FormatUint(start, 10)
	}
	params = util.MappedParameters(params, "", paramsMap)
	params = util.AddDatabaseParam(params, c)
	params = util.AddTransactionParam(params, transaction)

	req, err := util.BuildRequestFromHandle(c, "POST", "/internal/uris"+params, query)
	if err != nil {
		return err
	}
	return util.Execute(c, req, respHandle)
}
