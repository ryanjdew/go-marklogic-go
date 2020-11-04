package util

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"strings"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
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

// Deserialized returns *[]ForestInfo as interface{}
func (fih *ForestInfoHandle) Deserialized() interface{} {
	return &fih.ForestInfo
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

// GetForestInfo provides the forest information about the database
func GetForestInfo(c *clients.Client) []ForestInfo {
	params := AddDatabaseParam("", c)
	req, _ := BuildRequestFromHandle(c, "GET", "/internal/forestinfo"+params, nil)
	forestInfoHandle := &ForestInfoHandle{}
	Execute(c, req, forestInfoHandle)
	return *forestInfoHandle.Get()
}

// GetClientsByHost provides the forest information about the database
func GetClientsByHost(client *clients.Client, forestInfo []ForestInfo) map[string]*clients.Client {
	uniqueHosts := make(map[string]struct{})
	connectionInfo := client.BasicClient.ConnectionInfo()
	for _, forest := range forestInfo {
		uniqueHosts[forest.PreferredHost()] = struct{}{}
	}

	clientsByHost := make(map[string]*clients.Client, len(uniqueHosts))
	for host := range uniqueHosts {
		if host == connectionInfo.Host {
			clientsByHost[host] = client
		} else {
			clientsByHost[host], _ = clients.NewClient(&clients.Connection{
				Host:               host,
				Port:               connectionInfo.Port,
				Username:           connectionInfo.Username,
				Password:           connectionInfo.Password,
				AuthenticationType: connectionInfo.AuthenticationType,
				Database:           connectionInfo.Database,
			})
		}
	}
	return clientsByHost
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

// Deserialized returns string array of URIs as interface{}
func (uh *URIsHandle) Deserialized() interface{} {
	return uh.URIs
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

// GetURIs retrieves URIs from the internal API for the Data Movement SDK
func GetURIs(
	c *clients.Client,
	query handle.Handle,
	forestName string,
	transaction *Transaction,
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
	params = MappedParameters(params, "", paramsMap)
	params = AddDatabaseParam(params, c)
	params = AddTransactionParam(params, transaction)

	req, err := BuildRequestFromHandle(c, "POST", "/internal/uris"+params, query)
	if err != nil {
		return err
	}
	return Execute(c, req, respHandle)
}
