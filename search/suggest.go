package search

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// SuggestionsResponse represents the search Suggestions from MarkLogic
type SuggestionsResponse struct {
	XMLName     xml.Name `xml:"http://marklogic.com/appservices/search suggestions" json:"-"`
	Suggestions []string `xml:"http://marklogic.com/appservices/search suggestion" json:"suggestions"`
}

// SuggestionsResponseHandle is a handle that places the results into
// a Query struct
type SuggestionsResponseHandle struct {
	Format              int
	bytes               []byte
	suggestionsResponse *SuggestionsResponse
}

// GetFormat returns int that represents XML or JSON
func (srh *SuggestionsResponseHandle) GetFormat() int {
	return srh.Format
}

// Encode returns Query struct that represents XML or JSON
func (srh *SuggestionsResponseHandle) Encode(bytes []byte) {
	srh.bytes = bytes
	srh.suggestionsResponse = &SuggestionsResponse{}
	if srh.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, &srh.suggestionsResponse)
	} else {
		xml.Unmarshal(bytes, &srh.suggestionsResponse)
	}
}

// Decode returns []byte of XML or JSON that represents the Query struct
func (srh *SuggestionsResponseHandle) Decode(suggestionsResponse interface{}) {
	srh.suggestionsResponse = suggestionsResponse.(*SuggestionsResponse)
	buf := new(bytes.Buffer)
	if srh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(srh.suggestionsResponse)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(srh.suggestionsResponse)
	}
	srh.bytes = buf.Bytes()
}

// Get returns string of XML or JSON
func (srh *SuggestionsResponseHandle) Get() *SuggestionsResponse {
	return srh.suggestionsResponse
}

// Serialized returns string of XML or JSON
func (srh *SuggestionsResponseHandle) Serialized() string {
	buf := new(bytes.Buffer)
	if srh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(srh.suggestionsResponse)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(srh.suggestionsResponse)
	}
	srh.bytes = buf.Bytes()
	return string(srh.bytes)
}

// StructuredSuggestions suggests query text based off of a structured query
func StructuredSuggestions(c *clients.Client, query handle.Handle, partialQ string, limit int64, options string, response handle.Handle) error {
	reqType := handle.FormatEnumToString(query.GetFormat())
	buf := new(bytes.Buffer)
	buf.Write([]byte(query.Serialized()))
	url := c.Base() + "/suggest?format=" + reqType + "&limit=" + strconv.FormatInt(limit, 10)
	if options != "" {
		url = url + "&options=" + options
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/"+reqType)
	return clients.Execute(c, req, response)
}
