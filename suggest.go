package goMarklogicGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
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
	if srh.GetFormat() == JSON {
		json.Unmarshal(bytes, &srh.suggestionsResponse)
	} else {
		xml.Unmarshal(bytes, &srh.suggestionsResponse)
	}
}

// Decode returns []byte of XML or JSON that represents the Query struct
func (srh *SuggestionsResponseHandle) Decode(suggestionsResponse interface{}) {
	srh.suggestionsResponse = suggestionsResponse.(*SuggestionsResponse)
	buf := new(bytes.Buffer)
	if srh.GetFormat() == JSON {
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
	if srh.GetFormat() == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(srh.suggestionsResponse)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(srh.suggestionsResponse)
	}
	srh.bytes = buf.Bytes()
	return string(srh.bytes)
}

// StructuredSuggestions searches with a structured query
func (c *Client) StructuredSuggestions(query Handle, partialQ string, limit int64, options string, response Handle) error {
	var reqType string
	if response.GetFormat() == JSON {
		reqType = "json"
	} else {
		reqType = "xml"
	}
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
	applyAuth(c, req)
	req.Header.Add("Content-Type", "application/"+reqType)
	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	response.Encode(contents)
	return err
}
