package go_marklogic_go

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
)

const (
	XML = iota
	JSON
)

// Query represents a structured query from the search API
type Query struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search query"`
	Format  int           `xml:"-"`
	Queries []interface{} `xml:",any"`
}

type OrQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search or-query"`
	Queries []interface{} `xml:",any"`
}

type AndQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search and-query"`
	Ordered bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty"`
	Queries []interface{} `xml:",any"`
}

type TermQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search term-query"`
	Terms   []string `xml:"http://marklogic.com/appservices/search text"`
	Weight  float64  `xml:"http://marklogic.com/appservices/search weight,omitempty"`
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

type NearQuery struct {
	XMLName        xml.Name      `xml:"http://marklogic.com/appservices/search near-query"`
	Queries        []interface{} `xml:",any"`
	Ordered        bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty"`
	Distance       int64         `xml:"http://marklogic.com/appservices/search distance,omitempty"`
	DistanceWeight float64       `xml:"http://marklogic.com/appservices/search distance-weight,omitempty"`
}

type BoostQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search boost-query"`
	MatchingQuery PositiveQuery `xml:"http://marklogic.com/appservices/search macthing-query"`
	BoostingQuery NegativeQuery `xml:"http://marklogic.com/appservices/search boosting-query"`
}

type MatchingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search matching-query"`
	Queries []interface{} `xml:",any"`
}

type BoostingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search boosting-query"`
	Queries []interface{} `xml:",any"`
}

type PropertiesQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search properties-query"`
	Queries []interface{} `xml:",any"`
}

type DirectoryQuery struct {
	XMLName  xml.Name `xml:"http://marklogic.com/appservices/search directory-query"`
	URIs     []string `xml:"http://marklogic.com/appservices/search uri"`
	Infinite bool     `xml:"http://marklogic.com/appservices/search infinite,omitempty"`
}

type CollectionQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search collection-query"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri"`
}

type ContainerQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search container-query"`
	Element       QueryElement  `xml:"http://marklogic.com/appservices/search element,omitempty"`
	JsonKey       string        `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	FragmentScope string        `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Queries       []interface{} `xml:",any"`
}

type QueryElement struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search element"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type QueryAttribute struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search attribute"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type DocumentQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search document-query"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri"`
}

type DocumentFragmentQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search document-fragment-query"`
	Queries []interface{} `xml:",any"`
}

type LocksQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search locks-query"`
	Queries []interface{} `xml:",any"`
}

type RangeQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search range-query"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty"`
	JsonKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty"`
	PathIndex     string         `xml:"http://marklogic.com/appservices/search path-index,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Value         string         `xml:"http://marklogic.com/appservices/search value,omitempty"`
	RangeOperator string         `xml:"http://marklogic.com/appservices/search range-operator,omitempty"`
	RangeOptions  []string       `xml:"http://marklogic.com/appservices/search range-option,omitempty"`
}

type Field struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search field"`
	Name      string   `xml:"name,attr"`
	Collation string   `xml:"collation,attr"`
}

type ValueQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search value-query"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty"`
	JsonKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty"`
}

type WordQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search word-query"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty"`
	JsonKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty"`
}

type QueryParent struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search parent"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type HeatMap struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search heatmap"`
	North   float64  `xml:"n,attr"`
	East    float64  `xml:"e,attr"`
	South   float64  `xml:"s,attr"`
	West    float64  `xml:"w,attr"`
	Latdivs int64    `xml:"latdivs,attr"`
	Londivs int64    `xml:"londivs,attr"`
}

type Point struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search point"`
	Latitude  float64  `xml:"http://marklogic.com/appservices/search latitude"`
	Longitude float64  `xml:"http://marklogic.com/appservices/search longitude"`
}

type Box struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search box"`
	South   float64  `xml:"http://marklogic.com/appservices/search south"`
	West    float64  `xml:"http://marklogic.com/appservices/search west"`
	North   float64  `xml:"http://marklogic.com/appservices/search north"`
	East    float64  `xml:"http://marklogic.com/appservices/search east"`
}

type Circle struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search circle"`
	Radius  float64  `xml:"http://marklogic.com/appservices/search radius"`
	Point   Point    `xml:"http://marklogic.com/appservices/search point"`
}

type Polygon struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search polygon"`
	Points  []*Point `xml:"http://marklogic.com/appservices/search point"`
}

type GeoElemQuery struct {
	XMLName      xml.Name     `xml:"http://marklogic.com/appservices/search geo-elem-query"`
	Parent       QueryParent  `xml:"http://marklogic.com/appservices/search parent,omitempty"`
	Element      QueryElement `xml:"http://marklogic.com/appservices/search element"`
	GeoOptions   []string     `xml:"http://marklogic.com/appservices/search geo-option,omitempty"`
	FacetOptions []string     `xml:"http://marklogic.com/appservices/search facet-option,omitempty"`
	HeatMap      HeatMap      `xml:"http://marklogic.com/appservices/search heatmap,omitempty"`
	Points       []*Point     `xml:"http://marklogic.com/appservices/search point,omitempty"`
	Boxes        []*Box       `xml:"http://marklogic.com/appservices/search box,omitempty"`
	Circles      []*Circle    `xml:"http://marklogic.com/appservices/search circle,omitempty"`
	Polygons     []*Polygon   `xml:"http://marklogic.com/appservices/search polygon,omitempty"`
}

type Lat struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lat"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type Lon struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lon"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type GeoElemPairQuery struct {
	XMLName      xml.Name    `xml:"http://marklogic.com/appservices/search geo-elem-pair-query"`
	Parent       QueryParent `xml:"http://marklogic.com/appservices/search parent,omitempty"`
	Lat          Lat         `xml:"http://marklogic.com/appservices/search lat"`
	Lon          Lon         `xml:"http://marklogic.com/appservices/search lon"`
	GeoOptions   []string    `xml:"http://marklogic.com/appservices/search geo-option,omitempty"`
	FacetOptions []string    `xml:"http://marklogic.com/appservices/search facet-option,omitempty"`
	HeatMap      HeatMap     `xml:"http://marklogic.com/appservices/search heatmap,omitempty"`
	Points       []*Point    `xml:"http://marklogic.com/appservices/search point,omitempty"`
	Boxes        []*Box      `xml:"http://marklogic.com/appservices/search box,omitempty"`
	Circles      []*Circle   `xml:"http://marklogic.com/appservices/search circle,omitempty"`
	Polygons     []*Polygon  `xml:"http://marklogic.com/appservices/search polygon,omitempty"`
}

type GeoAttrPairQuery struct {
	XMLName      xml.Name    `xml:"http://marklogic.com/appservices/search geo-attr-pair-query"`
	Parent       QueryParent `xml:"http://marklogic.com/appservices/search parent"`
	Lat          Lat         `xml:"http://marklogic.com/appservices/search lat"`
	Lon          Lon         `xml:"http://marklogic.com/appservices/search lon"`
	GeoOptions   []string    `xml:"http://marklogic.com/appservices/search geo-option,omitempty"`
	FacetOptions []string    `xml:"http://marklogic.com/appservices/search facet-option,omitempty"`
	HeatMap      HeatMap     `xml:"http://marklogic.com/appservices/search heatmap,omitempty"`
	Points       []*Point    `xml:"http://marklogic.com/appservices/search point,omitempty"`
	Boxes        []*Box      `xml:"http://marklogic.com/appservices/search box,omitempty"`
	Circles      []*Circle   `xml:"http://marklogic.com/appservices/search circle,omitempty"`
	Polygons     []*Polygon  `xml:"http://marklogic.com/appservices/search polygon,omitempty"`
}

type GeoPathQuery struct {
	XMLName      xml.Name   `xml:"http://marklogic.com/appservices/search geo-path-query"`
	PathIndex    string     `xml:"http://marklogic.com/appservices/search path-index,omitempty"`
	GeoOptions   []string   `xml:"http://marklogic.com/appservices/search geo-option,omitempty"`
	FacetOptions []string   `xml:"http://marklogic.com/appservices/search facet-option,omitempty"`
	HeatMap      HeatMap    `xml:"http://marklogic.com/appservices/search heatmap,omitempty"`
	Points       []*Point   `xml:"http://marklogic.com/appservices/search point,omitempty"`
	Boxes        []*Box     `xml:"http://marklogic.com/appservices/search box,omitempty"`
	Circles      []*Circle  `xml:"http://marklogic.com/appservices/search circle,omitempty"`
	Polygons     []*Polygon `xml:"http://marklogic.com/appservices/search polygon,omitempty"`
}

func NewQuery(format int) *Query {
	return &Query{Format: format}
}

func (q *Query) Encode() *bytes.Buffer {
	buf := new(bytes.Buffer)
	if q.Format == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(q)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(q)
	}
	return buf
}
