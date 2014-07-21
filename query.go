package go_marklogic_go

import (
	"encoding/xml"
	//"io"
)

// Query represents a structured query from the search API
type Query struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search query"`
	Queries []interface{} `xml:",any"`
}

type OrQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search or-query"`
	Queries []interface{} `xml:",any"`
}

type AndQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search and-query"`
	Queries []interface{} `xml:",any"`
	Ordered bool          `xml:"http://marklogic.com/appservices/search ordered"`
}

type TermQuery struct {
	XMLName xml.Name  `xml:"http://marklogic.com/appservices/search term-query"`
	Terms   []*string `xml:"http://marklogic.com/appservices/search text"`
	Weight  float64   `xml:"http://marklogic.com/appservices/search weight"`
}

type AndNotQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search and-not-query"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query"`
}

type PositiveQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search positive-query"`
	Queries []interface{} `xml:",any"`
}

type NegativeQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search negative-query"`
	Queries []interface{} `xml:",any"`
}

type NotQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search not-query"`
	Queries []interface{} `xml:",any"`
}

type NotInQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search not-in-query"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query"`
}
