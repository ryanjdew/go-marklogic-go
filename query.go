package goMarklogicGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Format options
const (
	XML = iota
	JSON
)

// Query represents http://docs.marklogic.com/guide/search-dev/structured-query#id_85307
type Query struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search query"`
	Format  int           `xml:"-"`
	Queries []interface{} `xml:",any"`
}

// OrQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_64259
type OrQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search or-query"`
	Queries []interface{} `xml:",any"`
}

// AndQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83674
type AndQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search and-query"`
	Ordered bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty"`
	Queries []interface{} `xml:",any"`
}

// TermQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_56027
type TermQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search term-query"`
	Terms   []string `xml:"http://marklogic.com/appservices/search text"`
	Weight  float64  `xml:"http://marklogic.com/appservices/search weight,omitempty"`
}

// AndNotQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type AndNotQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search and-not-query"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query"`
}

// PositiveQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type PositiveQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search positive-query"`
	Queries []interface{} `xml:",any"`
}

// NegativeQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type NegativeQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search negative-query"`
	Queries []interface{} `xml:",any"`
}

// NotQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_39488
type NotQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search not-query"`
	Queries []interface{} `xml:",any"`
}

// NotInQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_90794
type NotInQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search not-in-query"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query"`
}

// NearQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_48512
type NearQuery struct {
	XMLName        xml.Name      `xml:"http://marklogic.com/appservices/search near-query"`
	Queries        []interface{} `xml:",any"`
	Ordered        bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty"`
	Distance       int64         `xml:"http://marklogic.com/appservices/search distance,omitempty"`
	DistanceWeight float64       `xml:"http://marklogic.com/appservices/search distance-weight,omitempty"`
}

// BoostQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type BoostQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search boost-query"`
	MatchingQuery PositiveQuery `xml:"http://marklogic.com/appservices/search macthing-query"`
	BoostingQuery NegativeQuery `xml:"http://marklogic.com/appservices/search boosting-query"`
}

// MatchingQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type MatchingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search matching-query"`
	Queries []interface{} `xml:",any"`
}

// BoostingQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type BoostingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search boosting-query"`
	Queries []interface{} `xml:",any"`
}

// PropertiesQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_67222
type PropertiesQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search properties-query"`
	Queries []interface{} `xml:",any"`
}

// DirectoryQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_94821
type DirectoryQuery struct {
	XMLName  xml.Name `xml:"http://marklogic.com/appservices/search directory-query"`
	URIs     []string `xml:"http://marklogic.com/appservices/search uri"`
	Infinite bool     `xml:"http://marklogic.com/appservices/search infinite,omitempty"`
}

// CollectionQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_76890
type CollectionQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search collection-query"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri"`
}

// ContainerQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87231
type ContainerQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search container-query"`
	Element       QueryElement  `xml:"http://marklogic.com/appservices/search element,omitempty"`
	JSONKey       string        `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	FragmentScope string        `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Queries       []interface{} `xml:",any"`
}

// QueryElement represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87231
type QueryElement struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search element"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// QueryAttribute represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type QueryAttribute struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search attribute"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// DocumentQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_27172
type DocumentQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search document-query"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri"`
}

// DocumentFragmentQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_30556
type DocumentFragmentQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search document-fragment-query"`
	Queries []interface{} `xml:",any"`
}

// LocksQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_53441
type LocksQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search locks-query"`
	Queries []interface{} `xml:",any"`
}

// RangeQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type RangeQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search range-query"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty"`
	PathIndex     string         `xml:"http://marklogic.com/appservices/search path-index,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Value         string         `xml:"http://marklogic.com/appservices/search value,omitempty"`
	RangeOperator string         `xml:"http://marklogic.com/appservices/search range-operator,omitempty"`
	RangeOptions  []string       `xml:"http://marklogic.com/appservices/search range-option,omitempty"`
}

// Field represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type Field struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search field"`
	Name      string   `xml:"name,attr"`
	Collation string   `xml:"collation,attr"`
}

// ValueQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_39758
type ValueQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search value-query"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty"`
}

// WordQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18990
type WordQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search word-query"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty"`
}

// QueryParent represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type QueryParent struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search parent"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// HeatMap represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type HeatMap struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search heatmap"`
	North   float64  `xml:"n,attr"`
	East    float64  `xml:"e,attr"`
	South   float64  `xml:"s,attr"`
	West    float64  `xml:"w,attr"`
	Latdivs int64    `xml:"latdivs,attr"`
	Londivs int64    `xml:"londivs,attr"`
}

// Point represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Point struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search point"`
	Latitude  float64  `xml:"http://marklogic.com/appservices/search latitude"`
	Longitude float64  `xml:"http://marklogic.com/appservices/search longitude"`
}

// Box represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Box struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search box"`
	South   float64  `xml:"http://marklogic.com/appservices/search south"`
	West    float64  `xml:"http://marklogic.com/appservices/search west"`
	North   float64  `xml:"http://marklogic.com/appservices/search north"`
	East    float64  `xml:"http://marklogic.com/appservices/search east"`
}

// Circle represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Circle struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search circle"`
	Radius  float64  `xml:"http://marklogic.com/appservices/search radius"`
	Point   Point    `xml:"http://marklogic.com/appservices/search point"`
}

// Polygon represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Polygon struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search polygon"`
	Points  []*Point `xml:"http://marklogic.com/appservices/search point"`
}

// GeoElemQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
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

// Lat represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
type Lat struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lat"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// Lon represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
type Lon struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lon"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// GeoElemPairQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
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

// GeoAttrPairQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_67897
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

// GeoPathQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_58782
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

// NewQuery returns a Query struct pointer and accepts a format const as a parameter
func NewQuery(format int) *Query {
	return &Query{Format: format}
}

// Encode returns a buffer that contains the encoded Query struct string
func (q *Query) Encode() *bytes.Buffer {
	buf := new(bytes.Buffer)
	if q.Format == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(q)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(q)
	}
	fmt.Print(buf.String())
	return buf
}
