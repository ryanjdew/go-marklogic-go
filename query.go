package goMarklogicGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
)

// Format options
const (
	XML = iota
	JSON
)

// QueryHandle is a handle that places the results into
// a Query struct
type QueryHandle struct {
	Format int
	bytes  []byte
	query  *Query
}

// GetFormat returns int that represents XML or JSON
func (qh *QueryHandle) GetFormat() int {
	return qh.Format
}

// Encode returns Query struct that represents XML or JSON
func (qh *QueryHandle) Encode(bytes []byte) {
	qh.bytes = bytes
	qh.query = &Query{}
	if qh.GetFormat() == JSON {
		json.Unmarshal(bytes, &qh.query)
	} else {
		xml.Unmarshal(bytes, &qh.query)
	}
}

// Decode returns []byte of XML or JSON that represents the Query struct
func (qh *QueryHandle) Decode(query interface{}) {
	qh.query = query.(*Query)
	buf := new(bytes.Buffer)
	if qh.GetFormat() == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(qh.query)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(qh.query)
	}
	qh.bytes = buf.Bytes()
}

// Get returns string of XML or JSON
func (qh *QueryHandle) Get() *Query {
	return qh.query
}

// Serialized returns string of XML or JSON
func (qh *QueryHandle) Serialized() string {
	buf := new(bytes.Buffer)
	if qh.GetFormat() == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(qh.query)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(qh.query)
	}
	qh.bytes = buf.Bytes()
	return string(qh.bytes)
}

// Query represents http://docs.marklogic.com/guide/search-dev/structured-query#id_85307
type Query struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search query" json:"-"`
	Format  int           `xml:"-" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeQuery Query

//MarshalJSON for Query struct in a special way to add wraping {"query":...}
func (q Query) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeQuery(q))
}

// OrQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_64259
type OrQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search or-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeOrQuery OrQuery

//MarshalJSON for OrQuery struct in a special way to add wraping {"or-query":...}
func (q OrQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeOrQuery(q))
}

// AndQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83674
type AndQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search and-query" json:"-"`
	Ordered bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty" json:"ordered,omitempty"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeAndQuery AndQuery

//MarshalJSON for AndQuery struct in a special way to add wraping {"and-query":...}
func (q AndQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeAndQuery(q))
}

// TermQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_56027
type TermQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search term-query" json:"-"`
	Terms   []string `xml:"http://marklogic.com/appservices/search text" json:"text"`
	Weight  float64  `xml:"http://marklogic.com/appservices/search weight,omitempty" json:"weight,omitempty"`
}

type fakeTermQuery TermQuery

//MarshalJSON for TermQuery struct in a special way to add wraping {"term-query":...}
func (q TermQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeTermQuery(q))
}

// AndNotQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type AndNotQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search and-not-query" json:"-"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query" json:"positive-query,omitempty"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query" json:"negative-query,omitempty"`
}

type fakeAndNotQuery AndNotQuery

//MarshalJSON for AndNotQuery struct in a special way to add wraping {"and-not-query":...}
func (q AndNotQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeAndNotQuery(q))
}

// PositiveQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type PositiveQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search positive-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakePositiveQuery PositiveQuery

//MarshalJSON for PositiveQuery struct in a special way to add wraping {"positive-query":...}
func (q PositiveQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakePositiveQuery(q))
}

// NegativeQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type NegativeQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search negative-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeNegativeQuery NegativeQuery

//MarshalJSON for NegativeQuery struct in a special way to add wraping {"negative-query":...}
func (q NegativeQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeNegativeQuery(q))
}

// NotQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_39488
type NotQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search not-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeNotQuery NotQuery

//MarshalJSON for NotQuery struct in a special way to add wraping {"not-query":...}
func (q NotQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeNotQuery(q))
}

// NotInQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_90794
type NotInQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search not-in-query" json:"-"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query" json:"positive-query,omitempty"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query" json:"negative-query,omitempty"`
}

type fakeNotInQuery NotInQuery

//MarshalJSON for NotInQuery struct in a special way to add wraping {"not-in-query":...}
func (q NotInQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeNotInQuery(q))
}

// NearQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_48512
type NearQuery struct {
	XMLName        xml.Name      `xml:"http://marklogic.com/appservices/search near-query" json:"-"`
	Queries        []interface{} `xml:",any" json:"queries"`
	Ordered        bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty" json:"ordered,omitempty"`
	Distance       int64         `xml:"http://marklogic.com/appservices/search distance,omitempty" json:"distance,omitempty"`
	DistanceWeight float64       `xml:"http://marklogic.com/appservices/search distance-weight,omitempty" json:"distance-weight,omitempty"`
}

type fakeNearQuery NearQuery

//MarshalJSON for NearQuery struct in a special way to add wraping {"near-query":...}
func (q NearQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeNearQuery(q))
}

// BoostQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type BoostQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search boost-query" json:"-"`
	MatchingQuery PositiveQuery `xml:"http://marklogic.com/appservices/search macthing-query" json:"macthing-query,omitempty"`
	BoostingQuery NegativeQuery `xml:"http://marklogic.com/appservices/search boosting-query" json:"boosting-query,omitempty"`
}

type fakeBoostQuery BoostQuery

//MarshalJSON for BoostQuery struct in a special way to add wraping {"boost-query":...}
func (q BoostQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeBoostQuery(q))
}

// MatchingQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type MatchingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search matching-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeMatchingQuery MatchingQuery

//MarshalJSON for MatchingQuery struct in a special way to add wraping {"matching-query":...}
func (q MatchingQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeMatchingQuery(q))
}

// BoostingQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type BoostingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search boosting-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeBoostingQuery BoostingQuery

//MarshalJSON for BoostingQuery struct in a special way to add wraping {"boosting-query":...}
func (q BoostingQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeBoostingQuery(q))
}

// PropertiesQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_67222
type PropertiesQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search properties-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakePropertiesQuery PropertiesQuery

//MarshalJSON for PropertiesQuery struct in a special way to add wraping {"properties-query":...}
func (q PropertiesQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakePropertiesQuery(q))
}

// DirectoryQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_94821
type DirectoryQuery struct {
	XMLName  xml.Name `xml:"http://marklogic.com/appservices/search directory-query" json:"-"`
	URIs     []string `xml:"http://marklogic.com/appservices/search uri" json:"uri,omitempty"`
	Infinite bool     `xml:"http://marklogic.com/appservices/search infinite,omitempty" json:"infinite,omitempty"`
}

type fakeDirectoryQuery DirectoryQuery

//MarshalJSON for DirectoryQuery struct in a special way to add wraping {"directory-query":...}
func (q DirectoryQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeDirectoryQuery(q))
}

// CollectionQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_76890
type CollectionQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search collection-query" json:"-"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri" json:"uri,omitempty"`
}

type fakeCollectionQuery CollectionQuery

//MarshalJSON for CollectionQuery struct in a special way to add wraping {"collection-query":...}
func (q CollectionQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeCollectionQuery(q))
}

// ContainerQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87231
type ContainerQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search container-query" json:"-"`
	Element       QueryElement  `xml:"http://marklogic.com/appservices/search element,omitempty" json:"element,omitempty"`
	JSONKey       string        `xml:"http://marklogic.com/appservices/search json-key,omitempty" json:"json-key,omitempty"`
	FragmentScope string        `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty" json:"fragment-scope,omitempty"`
	Queries       []interface{} `xml:",any" json:"queries"`
}

type fakeContainerQuery ContainerQuery

//MarshalJSON for ContainerQuery struct in a special way to add wraping {"container-query":...}
func (q ContainerQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeContainerQuery(q))
}

// QueryElement represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87231
type QueryElement struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search element" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type fakeQueryElement QueryElement

//MarshalJSON for QueryElement struct in a special way to add wraping {"element":...}
func (q QueryElement) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeQueryElement(q))
}

// QueryAttribute represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type QueryAttribute struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search attribute" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type fakeQueryAttribute QueryAttribute

//MarshalJSON for QueryAttribute struct in a special way to add wraping {"attribute":...}
func (q QueryAttribute) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeQueryAttribute(q))
}

// DocumentQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_27172
type DocumentQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search document-query" json:"-"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri" json:"uri,omitempty"`
}

type fakeDocumentQuery DocumentQuery

//MarshalJSON for DocumentQuery struct in a special way to add wraping {"document-query":...}
func (q DocumentQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeDocumentQuery(q))
}

// DocumentFragmentQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_30556
type DocumentFragmentQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search document-fragment-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeDocumentFragmentQuery DocumentFragmentQuery

//MarshalJSON for DocumentFragmentQuery struct in a special way to add wraping {"document-fragment-query":...}
func (q DocumentFragmentQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeDocumentFragmentQuery(q))
}

// LocksQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_53441
type LocksQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search locks-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

type fakeLocksQuery LocksQuery

//MarshalJSON for LocksQuery struct in a special way to add wraping {"locks-query":...}
func (q LocksQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeLocksQuery(q))
}

// RangeQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type RangeQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search range-query" json:"-"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty" json:"element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty" json:"attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty" json:"json-key,omitempty"`
	Field         FieldReference `xml:"http://marklogic.com/appservices/search field,omitempty" json:"field,omitempty"`
	PathIndex     string         `xml:"http://marklogic.com/appservices/search path-index,omitempty" json:"path-index,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty" json:"fragment-scope,omitempty"`
	Value         string         `xml:"http://marklogic.com/appservices/search value,omitempty" json:"value,omitempty"`
	RangeOperator string         `xml:"http://marklogic.com/appservices/search range-operator,omitempty" json:"range-operator,omitempty"`
	RangeOptions  []string       `xml:"http://marklogic.com/appservices/search range-option,omitempty" json:"range-option,omitempty"`
}

type fakeRangeQuery RangeQuery

//MarshalJSON for RangeQuery struct in a special way to add wraping {"range-query":...}
func (q RangeQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeRangeQuery(q))
}

// FieldReference represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type FieldReference struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search field" json:"-"`
	Name      string   `xml:"name,attr"`
	Collation string   `xml:"collation,attr"`
}

type fakeFieldReference FieldReference

//MarshalJSON for Field struct in a special way to add wraping {"field":...}
func (q FieldReference) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeFieldReference(q))
}

// ValueQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_39758
type ValueQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search value-query" json:"-"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty" json:"element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty" json:"attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty" json:"json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty" json:"field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty" json:"fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty" json:"text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty" json:"term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty" json:"weight,omitempty"`
}

type fakeValueQuery ValueQuery

//MarshalJSON for ValueQuery struct in a special way to add wraping {"value-query":...}
func (q ValueQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeValueQuery(q))
}

// WordQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18990
type WordQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search word-query" json:"-"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty" json:"element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty" json:"attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty" json:"json-key,omitempty"`
	Field         Field          `xml:"http://marklogic.com/appservices/search field,omitempty" json:"field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty" json:"fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty" json:"text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty" json:"term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty" json:"weight,omitempty"`
}

type fakeWordQuery WordQuery

//MarshalJSON for WordQuery struct in a special way to add wraping {"word-query":...}
func (q WordQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeWordQuery(q))
}

// QueryParent represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type QueryParent struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search parent" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type fakeQueryParent QueryParent

//MarshalJSON for QueryParent struct in a special way to add wraping {"parent":...}
func (q QueryParent) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeQueryParent(q))
}

// HeatMap represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type HeatMap struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search heatmap" json:"-"`
	North   float64  `xml:"n,attr"`
	East    float64  `xml:"e,attr"`
	South   float64  `xml:"s,attr"`
	West    float64  `xml:"w,attr"`
	Latdivs int64    `xml:"latdivs,attr"`
	Londivs int64    `xml:"londivs,attr"`
}

type fakeHeatMap HeatMap

//MarshalJSON for HeatMap struct in a special way to add wraping {"heatmap":...}
func (q HeatMap) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeHeatMap(q))
}

// Point represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Point struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search point" json:"-"`
	Latitude  float64  `xml:"http://marklogic.com/appservices/search latitude" json:"latitude,omitempty"`
	Longitude float64  `xml:"http://marklogic.com/appservices/search longitude" json:"longitude,omitempty"`
}

type fakePoint Point

//MarshalJSON for Point struct in a special way to add wraping {"point":...}
func (q Point) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakePoint(q))
}

// Box represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Box struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search box" json:"-"`
	South   float64  `xml:"http://marklogic.com/appservices/search south" json:"south,omitempty"`
	West    float64  `xml:"http://marklogic.com/appservices/search west" json:"west,omitempty"`
	North   float64  `xml:"http://marklogic.com/appservices/search north" json:"north,omitempty"`
	East    float64  `xml:"http://marklogic.com/appservices/search east" json:"east,omitempty"`
}

type fakeBox Box

//MarshalJSON for Box struct in a special way to add wraping {"box":...}
func (q Box) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeBox(q))
}

// Circle represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Circle struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search circle" json:"-"`
	Radius  float64  `xml:"http://marklogic.com/appservices/search radius" json:"radius,omitempty"`
	Point   Point    `xml:"http://marklogic.com/appservices/search point" json:"point,omitempty"`
}

type fakeCircle Circle

//MarshalJSON for Circle struct in a special way to add wraping {"circle":...}
func (q Circle) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeCircle(q))
}

// Polygon represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Polygon struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search polygon" json:"-"`
	Points  []*Point `xml:"http://marklogic.com/appservices/search point" json:"point,omitempty"`
}

type fakePolygon Polygon

//MarshalJSON for Polygon struct in a special way to add wraping {"polygon":...}
func (q Polygon) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakePolygon(q))
}

// GeoElemQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type GeoElemQuery struct {
	XMLName      xml.Name     `xml:"http://marklogic.com/appservices/search geo-elem-query" json:"-"`
	Parent       QueryParent  `xml:"http://marklogic.com/appservices/search parent,omitempty" json:"parent,omitempty"`
	Element      QueryElement `xml:"http://marklogic.com/appservices/search element" json:"element,omitempty"`
	GeoOptions   []string     `xml:"http://marklogic.com/appservices/search geo-option,omitempty" json:"geo-option,omitempty"`
	FacetOptions []string     `xml:"http://marklogic.com/appservices/search facet-option,omitempty" json:"facet-option,omitempty"`
	HeatMap      HeatMap      `xml:"http://marklogic.com/appservices/search heatmap,omitempty" json:"heatmap,omitempty"`
	Points       []*Point     `xml:"http://marklogic.com/appservices/search point,omitempty" json:"point,omitempty"`
	Boxes        []*Box       `xml:"http://marklogic.com/appservices/search box,omitempty" json:"box,omitempty"`
	Circles      []*Circle    `xml:"http://marklogic.com/appservices/search circle,omitempty" json:"circle,omitempty"`
	Polygons     []*Polygon   `xml:"http://marklogic.com/appservices/search polygon,omitempty" json:"polygon,omitempty"`
}

type fakeGeoElemQuery GeoElemQuery

//MarshalJSON for GeoElemQuery struct in a special way to add wraping {"geo-elem-query":...}
func (q GeoElemQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeGeoElemQuery(q))
}

// Lat represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
type Lat struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lat" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type fakeLat Lat

//MarshalJSON for Lat struct in a special way to add wraping {"lat":...}
func (q Lat) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeLat(q))
}

// Lon represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
type Lon struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lon" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

type fakeLon Lon

//MarshalJSON for Lon struct in a special way to add wraping {"lon":...}
func (q Lon) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeLon(q))
}

// GeoElemPairQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
type GeoElemPairQuery struct {
	XMLName      xml.Name    `xml:"http://marklogic.com/appservices/search geo-elem-pair-query" json:"-"`
	Parent       QueryParent `xml:"http://marklogic.com/appservices/search parent,omitempty" json:"parent,omitempty"`
	Lat          Lat         `xml:"http://marklogic.com/appservices/search lat" json:"lat,omitempty"`
	Lon          Lon         `xml:"http://marklogic.com/appservices/search lon" json:"lon,omitempty"`
	GeoOptions   []string    `xml:"http://marklogic.com/appservices/search geo-option,omitempty" json:"geo-option,omitempty"`
	FacetOptions []string    `xml:"http://marklogic.com/appservices/search facet-option,omitempty" json:"facet-option,omitempty"`
	HeatMap      HeatMap     `xml:"http://marklogic.com/appservices/search heatmap,omitempty" json:"heatmap,omitempty"`
	Points       []*Point    `xml:"http://marklogic.com/appservices/search point,omitempty" json:"point,omitempty"`
	Boxes        []*Box      `xml:"http://marklogic.com/appservices/search box,omitempty" json:"box,omitempty"`
	Circles      []*Circle   `xml:"http://marklogic.com/appservices/search circle,omitempty" json:"circle,omitempty"`
	Polygons     []*Polygon  `xml:"http://marklogic.com/appservices/search polygon,omitempty" json:"polygon,omitempty"`
}

type fakeGeoElemPairQuery GeoElemPairQuery

//MarshalJSON for GeoElemPairQuery struct in a special way to add wraping {"geo-elem-pair-query":...}
func (q GeoElemPairQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeGeoElemPairQuery(q))
}

// GeoAttrPairQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_67897
type GeoAttrPairQuery struct {
	XMLName      xml.Name    `xml:"http://marklogic.com/appservices/search geo-attr-pair-query" json:"-"`
	Parent       QueryParent `xml:"http://marklogic.com/appservices/search parent" json:"parent,omitempty"`
	Lat          Lat         `xml:"http://marklogic.com/appservices/search lat" json:"lat,omitempty"`
	Lon          Lon         `xml:"http://marklogic.com/appservices/search lon" json:"lon,omitempty"`
	GeoOptions   []string    `xml:"http://marklogic.com/appservices/search geo-option,omitempty" json:"geo-option,omitempty"`
	FacetOptions []string    `xml:"http://marklogic.com/appservices/search facet-option,omitempty" json:"facet-option,omitempty"`
	HeatMap      HeatMap     `xml:"http://marklogic.com/appservices/search heatmap,omitempty" json:"heatmap,omitempty"`
	Points       []*Point    `xml:"http://marklogic.com/appservices/search point,omitempty" json:"point,omitempty"`
	Boxes        []*Box      `xml:"http://marklogic.com/appservices/search box,omitempty" json:"box,omitempty"`
	Circles      []*Circle   `xml:"http://marklogic.com/appservices/search circle,omitempty" json:"circle,omitempty"`
	Polygons     []*Polygon  `xml:"http://marklogic.com/appservices/search polygon,omitempty" json:"polygon,omitempty"`
}

type fakeGeoAttrPairQuery GeoAttrPairQuery

//MarshalJSON for GeoAttrPairQuery struct in a special way to add wraping {"geo-attr-pair-query":...}
func (q GeoAttrPairQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeGeoAttrPairQuery(q))
}

// GeoPathQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_58782
type GeoPathQuery struct {
	XMLName      xml.Name   `xml:"http://marklogic.com/appservices/search geo-path-query" json:"-"`
	PathIndex    string     `xml:"http://marklogic.com/appservices/search path-index,omitempty" json:"path-index,omitempty"`
	GeoOptions   []string   `xml:"http://marklogic.com/appservices/search geo-option,omitempty" json:"geo-option,omitempty"`
	FacetOptions []string   `xml:"http://marklogic.com/appservices/search facet-option,omitempty" json:"facet-option,omitempty"`
	HeatMap      HeatMap    `xml:"http://marklogic.com/appservices/search heatmap,omitempty" json:"heatmap,omitempty"`
	Points       []*Point   `xml:"http://marklogic.com/appservices/search point,omitempty" json:"point,omitempty"`
	Boxes        []*Box     `xml:"http://marklogic.com/appservices/search box,omitempty" json:"box,omitempty"`
	Circles      []*Circle  `xml:"http://marklogic.com/appservices/search circle,omitempty" json:"circle,omitempty"`
	Polygons     []*Polygon `xml:"http://marklogic.com/appservices/search polygon,omitempty" json:"polygon,omitempty"`
}

type fakeGeoPathQuery GeoPathQuery

//MarshalJSON for GeoPathQuery struct in a special way to add wraping {"geo-path-query":...}
func (q GeoPathQuery) MarshalJSON() ([]byte, error) {
	return wrapJSON(fakeGeoPathQuery(q))
}

// UnmarshalJSON Query converts to JSON
func (q *Query) UnmarshalJSON(data []byte) error {
	fake := fakeQuery(*q)
	err := json.Unmarshal(data, &fake)
	if err != nil {
		return err
	}
	q.Format = fake.Format
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON OrQuery converts to JSON
func (q *OrQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON AndQuery converts to JSON
func (q *AndQuery) UnmarshalJSON(data []byte) error {
	fake := fakeAndQuery(*q)
	err := json.Unmarshal(data, &fake)
	if err != nil {
		return err
	}
	q.Ordered = fake.Ordered
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON PositiveQuery converts to JSON
func (q *PositiveQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON NegativeQuery converts to JSON
func (q *NegativeQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON NotQuery converts to JSON
func (q *NotQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON NearQuery converts to JSON
func (q *NearQuery) UnmarshalJSON(data []byte) error {
	fake := fakeNearQuery(*q)
	err := json.Unmarshal(data, &fake)
	if err != nil {
		return err
	}
	q.Ordered = fake.Ordered
	q.Distance = fake.Distance
	q.DistanceWeight = fake.DistanceWeight
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON MatchingQuery converts to JSON
func (q *MatchingQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON BoostingQuery converts to JSON
func (q *BoostingQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON PropertiesQuery converts to JSON
func (q *PropertiesQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON ContainerQuery converts to JSON
func (q *ContainerQuery) UnmarshalJSON(data []byte) error {
	fake := fakeContainerQuery(*q)
	err := json.Unmarshal(data, &fake)
	if err != nil {
		return err
	}
	q.Element = fake.Element
	q.JSONKey = fake.JSONKey
	q.FragmentScope = fake.FragmentScope
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON DocumentFragmentQuery converts to JSON
func (q *DocumentFragmentQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// UnmarshalJSON LocksQuery converts to JSON
func (q *LocksQuery) UnmarshalJSON(data []byte) error {
	queries, err2 := DecodeJSONWithQueries(data)
	q.Queries = queries
	return err2
}

// DecodeJSONWithQueries decodes text into Query struct
func DecodeJSONWithQueries(inputData []byte) ([]interface{}, error) {
	var queries []interface{}
	var rootQuery map[string]json.RawMessage
	err := json.Unmarshal(inputData, &rootQuery)
	var data map[string][]json.RawMessage
	for _, root := range rootQuery {
		err = json.Unmarshal(root, &data)
		for _, query := range data["queries"] {
			var queryStructs map[string]json.RawMessage
			err = json.Unmarshal(query, &queryStructs)
			for key, queryJSON := range queryStructs {
				var queryStruct interface{}
				switch key {
				case "query":
					q := Query{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "or-query":
					q := OrQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "and-query":
					q := AndQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "term-query":
					q := TermQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "and-not-query":
					q := AndNotQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "positive-query":
					q := PositiveQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "negative-query":
					q := NegativeQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "not-query":
					q := NotQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "not-in-query":
					q := NotInQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "near-query":
					q := NearQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "boost-query":
					q := BoostQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "matching-query":
					q := MatchingQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "boosting-query":
					q := BoostingQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "properties-query":
					q := PropertiesQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "directory-query":
					q := DirectoryQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "collection-query":
					q := CollectionQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "container-query":
					q := ContainerQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "element":
					q := QueryElement{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "attribute":
					q := QueryAttribute{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "document-query":
					q := DocumentQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "document-fragment-query":
					q := DocumentFragmentQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "locks-query":
					q := LocksQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "range-query":
					q := RangeQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "field":
					q := Field{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "value-query":
					q := ValueQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "word-query":
					q := WordQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "parent":
					q := QueryParent{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "heatmap":
					q := HeatMap{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "point":
					q := Point{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "box":
					q := Box{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "circle":
					q := Circle{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "polygon":
					q := Polygon{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "geo-elem-query":
					q := GeoElemQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "lat":
					q := Lat{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "lon":
					q := Lon{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "geo-elem-pair-query":
					q := GeoElemPairQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "geo-attr-pair-query":
					q := GeoAttrPairQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				case "geo-path-query":
					q := GeoPathQuery{}
					err = json.Unmarshal(queryJSON, &q)
					queryStruct = q
				default:
				}
				queries = append(queries, queryStruct)
			}
		}
	}
	return queries, err
}

// UnmarshalXML Query converts to XML
func (q *Query) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML OrQuery converts to XML
func (q *OrQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML AndQuery converts to XML
func (q *AndQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fake := fakeAndQuery(*q)
	err := d.DecodeElement(&fake, &start)
	if err != nil {
		return err
	}
	q.Ordered = fake.Ordered
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML PositiveQuery converts to XML
func (q *PositiveQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML NegativeQuery converts to XML
func (q *NegativeQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML NotQuery converts to XML
func (q *NotQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML NearQuery converts to XML
func (q *NearQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fake := fakeNearQuery(*q)
	err := d.DecodeElement(&fake, &start)
	if err != nil {
		return err
	}
	q.Ordered = fake.Ordered
	q.Distance = fake.Distance
	q.DistanceWeight = fake.DistanceWeight
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML MatchingQuery converts to XML
func (q *MatchingQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML BoostingQuery converts to XML
func (q *BoostingQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML PropertiesQuery converts to XML
func (q *PropertiesQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML ContainerQuery converts to XML
func (q *ContainerQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fake := fakeContainerQuery(*q)
	err := d.DecodeElement(&fake, &start)
	if err != nil {
		return err
	}
	q.Element = fake.Element
	q.JSONKey = fake.JSONKey
	q.FragmentScope = fake.FragmentScope
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML DocumentFragmentQuery converts to XML
func (q *DocumentFragmentQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML LocksQuery converts to XML
func (q *LocksQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := DecodeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// DecodeXMLWithQueries decodes text into Query struct
func DecodeXMLWithQueries(d *xml.Decoder, start xml.StartElement) ([]interface{}, error) {
	var queries []interface{}
	for {
		if token, err := d.Token(); (err == nil) && (token != nil) {
			switch t := token.(type) {
			case xml.StartElement:
				e := xml.StartElement(t)
				var queryStruct interface{}
				switch e.Name.Local {
				case "query":
					q := Query{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "or-query":
					q := OrQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "and-query":
					q := AndQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "term-query":
					q := TermQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "and-not-query":
					q := AndNotQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "positive-query":
					q := PositiveQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "negative-query":
					q := NegativeQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "not-query":
					q := NotQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "not-in-query":
					q := NotInQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "near-query":
					q := NearQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "boost-query":
					q := BoostQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "matching-query":
					q := MatchingQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "boosting-query":
					q := BoostingQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "properties-query":
					q := PropertiesQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "directory-query":
					q := DirectoryQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "collection-query":
					q := CollectionQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "container-query":
					q := ContainerQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "element":
					q := QueryElement{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "attribute":
					q := QueryAttribute{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "document-query":
					q := DocumentQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "document-fragment-query":
					q := DocumentFragmentQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "locks-query":
					q := LocksQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "range-query":
					q := RangeQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "field":
					q := Field{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "value-query":
					q := ValueQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "word-query":
					q := WordQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "parent":
					q := QueryParent{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "heatmap":
					q := HeatMap{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "point":
					q := Point{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "box":
					q := Box{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "circle":
					q := Circle{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "polygon":
					q := Polygon{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "geo-elem-query":
					q := GeoElemQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "lat":
					q := Lat{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "lon":
					q := Lon{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "geo-elem-pair-query":
					q := GeoElemPairQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "geo-attr-pair-query":
					q := GeoAttrPairQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				case "geo-path-query":
					q := GeoPathQuery{}
					err = d.DecodeElement(&q, &e)
					queryStruct = q
				default:
				}
				queries = append(queries, queryStruct)
			case xml.EndElement:
				e := xml.EndElement(t)
				if e.Name.Space == "http://marklogic.com/appservices/search" && e.Name.Local == "queries" {
					return queries, err
				}
			}
		} else {
			return queries, err
		}
	}
}
