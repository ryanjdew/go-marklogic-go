package search

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// SuggestionsResponse represents the search Suggestions from MarkLogic
type SuggestionsResponse struct {
	XMLName     xml.Name `xml:"http://marklogic.com/appservices/search suggestions" json:"-"`
	Suggestions []string `xml:"http://marklogic.com/appservices/search suggestion" json:"suggestions"`
}

// SuggestionsResponseHandle is a handle that places the results into
// a Query struct
type SuggestionsResponseHandle struct {
	*bytes.Buffer
	Format              int
	suggestionsResponse *SuggestionsResponse
}

// GetFormat returns int that represents XML or JSON
func (srh *SuggestionsResponseHandle) GetFormat() int {
	return srh.Format
}

func (srh *SuggestionsResponseHandle) resetBuffer() {
	if srh.Buffer == nil {
		srh.Buffer = new(bytes.Buffer)
	}
	srh.Reset()
}

// Deserialize returns Query struct that represents XML or JSON
func (srh *SuggestionsResponseHandle) Deserialize(bytes []byte) {
	srh.resetBuffer()
	srh.Write(bytes)
	srh.suggestionsResponse = &SuggestionsResponse{}
	if srh.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, &srh.suggestionsResponse)
	} else {
		xml.Unmarshal(bytes, &srh.suggestionsResponse)
	}
}

// AcceptResponse handles an *http.Response
func (srh *SuggestionsResponseHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(srh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Query struct
func (srh *SuggestionsResponseHandle) Serialize(suggestionsResponse interface{}) {
	srh.suggestionsResponse = suggestionsResponse.(*SuggestionsResponse)
	srh.resetBuffer()
	if srh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(srh.Buffer)
		enc.Encode(srh.suggestionsResponse)
	} else {
		enc := xml.NewEncoder(srh.Buffer)
		enc.Encode(srh.suggestionsResponse)
	}
}

// Get returns string of XML or JSON
func (srh *SuggestionsResponseHandle) Get() *SuggestionsResponse {
	return srh.suggestionsResponse
}

// Serialized returns string of XML or JSON
func (srh *SuggestionsResponseHandle) Serialized() string {
	srh.Serialize(srh.suggestionsResponse)
	return srh.String()
}

// StructuredSuggestions suggests query text based off of a structured query
func StructuredSuggestions(c *clients.Client, query handle.Handle, partialQ string, limit int64, options string, response handle.ResponseHandle) error {
	uri := "/suggest?limit=" + strconv.FormatInt(limit, 10)
	if options != "" {
		uri = uri + "&options=" + options
	}
	req, err := util.BuildRequestFromHandle(c, "POST", uri, query)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
