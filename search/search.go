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

const (
	searchNamespace = "http://marklogic.com/appservices/search"
)

// ResponseHandle is a handle that places the results into
// a Response struct
type ResponseHandle struct {
	*bytes.Buffer
	Format   int
	response Response
}

// GetFormat returns int that represents XML or JSON
func (rh *ResponseHandle) GetFormat() int {
	return rh.Format
}

func (rh *ResponseHandle) resetBuffer() {
	if rh.Buffer == nil {
		rh.Buffer = new(bytes.Buffer)
	}
	rh.Reset()
}

// Deserialize returns Response struct that represents XML or JSON
func (rh *ResponseHandle) Deserialize(bytes []byte) {
	rh.resetBuffer()
	rh.Write(bytes)
	rh.response = Response{}
	if rh.GetFormat() == handle.JSON {
		json.Unmarshal(bytes, &rh.response)
	} else {
		xml.Unmarshal(bytes, &rh.response)
	}
}

// AcceptResponse handles an *http.Response
func (rh *ResponseHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(rh, resp)
}

// Serialize returns []byte of XML or JSON that represents the Response struct
func (rh *ResponseHandle) Serialize(response interface{}) {
	rh.response = response.(Response)
	rh.resetBuffer()
	if rh.GetFormat() == handle.JSON {
		enc := json.NewEncoder(rh.Buffer)
		enc.Encode(&rh.response)
	} else {
		enc := xml.NewEncoder(rh.Buffer)
		enc.Encode(&rh.response)
	}
}

// Get returns string of XML or JSON
func (rh *ResponseHandle) Get() *Response {
	return &rh.response
}

// Serialized returns string of XML or JSON
func (rh *ResponseHandle) Serialized() string {
	rh.Serialize(rh.response)
	return rh.String()
}

// Response represents a response from the search API
type Response struct {
	Total      int64    `xml:"total,attr" json:"total,omitempty"`
	Start      int64    `xml:"start,attr" json:"start,omitempty"`
	PageLength int64    `xml:"page-length,attr" json:"page-length,omitempty"`
	Results    []Result `xml:"http://marklogic.com/appservices/search result" json:"result,omitempty"`
	Facets     []Facet  `xml:"http://marklogic.com/appservices/search facet" json:"facet,omitempty"`
}

// Result is an individual document fragment found by the search
type Result struct {
	URI        string    `xml:"uri,attr" json:"uri,omitempty"`
	Href       string    `xml:"href,attr" json:"href,omitempty"`
	MimeType   string    `xml:"mimetype,attr" json:"mimetype,omitempty"`
	Format     string    `xml:"format,attr" json:"format,omitempty"`
	Path       string    `xml:"path,attr" json:"path,omitempty"`
	Index      int64     `xml:"index,attr" json:"index,omitempty"`
	Score      int64     `xml:"score,attr" json:"score,omitempty"`
	Confidence float64   `xml:"confidence,attr" json:"confidence,omitempty"`
	Fitness    float64   `xml:"fitness,attr" json:"fitness,omitempty"`
	Snippets   []Snippet `xml:"http://marklogic.com/appservices/search snippet" json:"snippet,omitempty"`
}

// Snippet represents a snippet
type Snippet struct {
	Location string  `xml:"uri,attr" json:"uri,omitempty"`
	Matches  []Match `xml:"http://marklogic.com/appservices/search match" json:"match,omitempty"`
}

// Match is a path in document that matches the query
type Match struct {
	Path string
	Text []Text
}

// Text in a match. HightlightedText tells whether the text matched the query
// or not.
type Text struct {
	Text            string
	HighlightedText bool
}

// Facet represents a facet and contains a slice of FacetValue
type Facet struct {
	Name        string       `xml:"name,attr" json:"name,omitempty"`
	Type        string       `xml:"type,attr" json:"type,omitempty"`
	FacetValues []FacetValue `xml:"http://marklogic.com/appservices/search facet-value" json:"facet-value,omitempty"`
}

// FacetValue is a value with the frequency that value occurs
type FacetValue struct {
	Name  string `xml:"name,attr" json:"name,omitempty"`
	Label string `xml:",chardata"`
	Count int64  `xml:"count,attr" json:"count,omitempty"`
}

// Search with text value
func Search(c *clients.Client, text string, start int64, pageLength int64, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/search?q="+text+"&start="+strconv.FormatInt(start, 10)+"&pageLength="+strconv.FormatInt(pageLength, 10), nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// StructuredSearch searches with a structured query
func StructuredSearch(c *clients.Client, query handle.Handle, start int64, pageLength int64, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "POST", "/search?start="+strconv.FormatInt(start, 10)+"&pageLength="+strconv.FormatInt(pageLength, 10), query)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

//UnmarshalXML for Match struct in a special way to handle highlighting matching text
func (m *Match) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for i := range start.Attr {
		attr := start.Attr[i]
		if attr.Name.Local == "path" {
			m.Path = attr.Value
			break
		}
	}
	for {
		if token, err := d.Token(); (err == nil) && (token != nil) {
			switch t := token.(type) {
			case xml.StartElement:
				var content string
				e := xml.StartElement(t)
				d.DecodeElement(&content, &e)
				text := Text{
					Text:            content,
					HighlightedText: e.Name.Space == searchNamespace && e.Name.Local == "highlight",
				}
				m.Text = append(m.Text, text)
			case xml.EndElement:
				e := xml.EndElement(t)
				if e.Name.Space == searchNamespace && e.Name.Local == "match" {
					return nil
				}
			case xml.CharData:
				b := xml.CharData(t)
				text := Text{
					Text:            string([]byte(b)),
					HighlightedText: false,
				}
				m.Text = append(m.Text, text)
			}
		} else {
			return err
		}
	}
}
