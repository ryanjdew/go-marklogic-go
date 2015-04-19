package goMarklogicGo

import (
	//"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
)

// Response represents a response from the search API
type Response struct {
	Total      int64     `xml:"total,attr" json:"total,omitempty"`
	Start      int64     `xml:"start,attr" json:"start,omitempty"`
	PageLength int64     `xml:"page-length,attr" json:"page-length,omitempty"`
	Results    []*Result `xml:"http://marklogic.com/appservices/search result" json:"result,omitempty"`
	Facets     []*Facet  `xml:"http://marklogic.com/appservices/search facet" json:"facet,omitempty"`
}

// Result is an individual document fragment found by the search
type Result struct {
	URI        string     `xml:"uri,attr" json:"uri,omitempty"`
	Href       string     `xml:"href,attr" json:"href,omitempty"`
	MimeType   string     `xml:"mimetype,attr" json:"mimetype,omitempty"`
	Format     string     `xml:"format,attr" json:"format,omitempty"`
	Path       string     `xml:"path,attr" json:"path,omitempty"`
	Index      int64      `xml:"index,attr" json:"index,omitempty"`
	Score      int64      `xml:"score,attr" json:"score,omitempty"`
	Confidence float64    `xml:"confidence,attr" json:"confidence,omitempty"`
	Fitness    float64    `xml:"fitness,attr" json:"fitness,omitempty"`
	Snippets   []*Snippet `xml:"http://marklogic.com/appservices/search snippet" json:"snippet,omitempty"`
}

// Snippet represents a snippet
type Snippet struct {
	Location string   `xml:"uri,attr" json:"uri,omitempty"`
	Matches  []*Match `xml:"http://marklogic.com/appservices/search match" json:"match,omitempty"`
}

// Match is a path in document that matches the query
type Match struct {
	Path string
	Text []*Text
}

// Text in a match. HightlightedText tells whether the text matched the query
// or not.
type Text struct {
	Text            string
	HighlightedText bool
}

// Facet represents a facet and contains a slice of FacetValue
type Facet struct {
	Name        string        `xml:"name,attr" json:"name,omitempty"`
	Type        string        `xml:"type,attr" json:"type,omitempty"`
	FacetValues []*FacetValue `xml:"http://marklogic.com/appservices/search facet-value" json:"facet-value,omitempty"`
}

// FacetValue is a value with the frequency that value occurs
type FacetValue struct {
	Name  string `xml:"name,attr" json:"name,omitempty"`
	Label string `xml:",chardata"`
	Count int64  `xml:"count,attr" json:"count,omitempty"`
}

// Search with text value
func (c *Client) Search(text string, start int64, pageLength int64) (*Response, error) {
	req, _ := http.NewRequest("GET", c.Base+"/search?q="+text+"&format=xml&start="+strconv.FormatInt(start, 10)+"&pageLength="+strconv.FormatInt(pageLength, 10), nil)
	applyAuth(c, req)
	resp, _ := c.HTTPClient.Do(req)
	defer resp.Body.Close()
	return readResults(resp.Body)
}

// StructuredSearch searches with a structured query
func (c *Client) StructuredSearch(query *Query, start int64, pageLength int64) (*Response, error) {
	buf := query.Encode()
	req, _ := http.NewRequest("POST", c.Base+"/search?format=xml&start="+strconv.FormatInt(start, 10)+"&pageLength="+strconv.FormatInt(pageLength, 10), buf)
	applyAuth(c, req)
	req.Header.Add("Content-Type", "application/xml")
	resp, _ := c.HTTPClient.Do(req)
	defer resp.Body.Close()
	return readResults(resp.Body)
}

func readResults(reader io.Reader) (*Response, error) {
	results := &Response{}
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(results); err != nil {
		return nil, err
	}
	return results, nil
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
				text := &Text{
					Text:            content,
					HighlightedText: e.Name.Space == "http://marklogic.com/appservices/search" && e.Name.Local == "highlight",
				}
				m.Text = append(m.Text, text)
			case xml.EndElement:
				e := xml.EndElement(t)
				if e.Name.Space == "http://marklogic.com/appservices/search" && e.Name.Local == "match" {
					return nil
				}
			case xml.CharData:
				b := xml.CharData(t)
				text := &Text{
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
