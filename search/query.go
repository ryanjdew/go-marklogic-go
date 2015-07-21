package search

import (
	"bytes"
	"encoding/xml"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

var mapperFunction = func(str string) interface{} {
	return stringToQueryStruct(str)
}

// QueryHandle is a handle that places the results into
// a Query struct
type QueryHandle struct {
	*bytes.Buffer
	Format int
	query  Query
}

// GetFormat returns int that represents XML or JSON
func (qh *QueryHandle) GetFormat() int {
	return qh.Format
}

func (qh *QueryHandle) resetBuffer() {
	if qh.Buffer == nil {
		qh.Buffer = new(bytes.Buffer)
	}
	qh.Reset()
}

// Deserialize returns Query struct that represents XML or JSON
func (qh *QueryHandle) Deserialize(bytes []byte) {
	qh.resetBuffer()
	qh.Write(bytes)
	qh.query = Query{}
	if qh.GetFormat() == handle.JSON {
		unwrapped, _ := unwrapJSON(bytes, mapperFunction)
		qh.query = unwrapped.(Query)
	} else {
		xml.Unmarshal(bytes, &qh.query)
	}
}

// Serialize returns []byte of XML or JSON that represents the Query struct
func (qh *QueryHandle) Serialize(query interface{}) {
	var readBytes []byte
	qh.query = query.(Query)
	qh.resetBuffer()
	if qh.GetFormat() == handle.JSON {
		readBytes, _ = wrapJSON(qh.query)
		qh.Write(readBytes)
	} else {
		enc := xml.NewEncoder(qh)
		enc.Encode(qh.query)
	}
}

// Get returns string of XML or JSON
func (qh *QueryHandle) Get() *Query {
	return &qh.query
}

// Serialized returns string of XML or JSON
func (qh *QueryHandle) Serialized() string {
	qh.Serialize(qh.query)
	return qh.String()
}

// Query represents http://docs.marklogic.com/guide/search-dev/structured-query#id_85307
type Query struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search query" json:"-"`
	Format  int           `xml:"-" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// OrQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_64259
type OrQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search or-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// AndQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83674
type AndQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search and-query" json:"-"`
	Ordered bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty" json:"ordered,omitempty"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// TermQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_56027
type TermQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search term-query" json:"-"`
	Terms   []string `xml:"http://marklogic.com/appservices/search text" json:"text"`
	Weight  float64  `xml:"http://marklogic.com/appservices/search weight,omitempty" json:"weight,omitempty"`
}

// AndNotQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type AndNotQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search and-not-query" json:"-"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query" json:"positive-query,omitempty"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query" json:"negative-query,omitempty"`
}

// PositiveQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type PositiveQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search positive-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// NegativeQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_65108
type NegativeQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search negative-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// NotQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_39488
type NotQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search not-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// NotInQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_90794
type NotInQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search not-in-query" json:"-"`
	PositiveQuery PositiveQuery `xml:"http://marklogic.com/appservices/search positive-query" json:"positive-query,omitempty"`
	NegativeQuery NegativeQuery `xml:"http://marklogic.com/appservices/search negative-query" json:"negative-query,omitempty"`
}

// NearQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_48512
type NearQuery struct {
	XMLName        xml.Name      `xml:"http://marklogic.com/appservices/search near-query" json:"-"`
	Queries        []interface{} `xml:",any" json:"queries"`
	Ordered        bool          `xml:"http://marklogic.com/appservices/search ordered,omitempty" json:"ordered,omitempty"`
	Distance       int64         `xml:"http://marklogic.com/appservices/search distance,omitempty" json:"distance,omitempty"`
	DistanceWeight float64       `xml:"http://marklogic.com/appservices/search distance-weight,omitempty" json:"distance-weight,omitempty"`
}

// BoostQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type BoostQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search boost-query" json:"-"`
	MatchingQuery PositiveQuery `xml:"http://marklogic.com/appservices/search macthing-query" json:"macthing-query,omitempty"`
	BoostingQuery NegativeQuery `xml:"http://marklogic.com/appservices/search boosting-query" json:"boosting-query,omitempty"`
}

// MatchingQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type MatchingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search matching-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// BoostingQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_25949
type BoostingQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search boosting-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// PropertiesQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_67222
type PropertiesQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search properties-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// DirectoryQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_94821
type DirectoryQuery struct {
	XMLName  xml.Name `xml:"http://marklogic.com/appservices/search directory-query" json:"-"`
	URIs     []string `xml:"http://marklogic.com/appservices/search uri" json:"uri,omitempty"`
	Infinite bool     `xml:"http://marklogic.com/appservices/search infinite,omitempty" json:"infinite,omitempty"`
}

// CollectionQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_76890
type CollectionQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search collection-query" json:"-"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri" json:"uri,omitempty"`
}

// ContainerQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87231
type ContainerQuery struct {
	XMLName       xml.Name      `xml:"http://marklogic.com/appservices/search container-query" json:"-"`
	Element       QueryElement  `xml:"http://marklogic.com/appservices/search element,omitempty" json:"element,omitempty"`
	JSONKey       string        `xml:"http://marklogic.com/appservices/search json-key,omitempty" json:"json-key,omitempty"`
	FragmentScope string        `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty" json:"fragment-scope,omitempty"`
	Queries       []interface{} `xml:",any" json:"queries"`
}

// QueryElement represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87231
type QueryElement struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search element" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// QueryAttribute represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type QueryAttribute struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search attribute" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// DocumentQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_27172
type DocumentQuery struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search document-query" json:"-"`
	URIs    []string `xml:"http://marklogic.com/appservices/search uri" json:"uri,omitempty"`
}

// DocumentFragmentQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_30556
type DocumentFragmentQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search document-fragment-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
}

// LocksQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_53441
type LocksQuery struct {
	XMLName xml.Name      `xml:"http://marklogic.com/appservices/search locks-query" json:"-"`
	Queries []interface{} `xml:",any" json:"queries"`
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

// FieldReference represents http://docs.marklogic.com/guide/search-dev/structured-query#id_83393
type FieldReference struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search field" json:"-"`
	Name      string   `xml:"name,attr"`
	Collation string   `xml:"collation,attr"`
}

// ValueQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_39758
type ValueQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search value-query" json:"-"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty" json:"element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty" json:"attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty" json:"json-key,omitempty"`
	Field         FieldReference `xml:"http://marklogic.com/appservices/search field,omitempty" json:"field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty" json:"fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty" json:"text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty" json:"term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty" json:"weight,omitempty"`
}

// WordQuery represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18990
type WordQuery struct {
	XMLName       xml.Name       `xml:"http://marklogic.com/appservices/search word-query" json:"-"`
	Element       QueryElement   `xml:"http://marklogic.com/appservices/search element,omitempty" json:"element,omitempty"`
	Attribute     QueryAttribute `xml:"http://marklogic.com/appservices/search attribute,omitempty" json:"attribute,omitempty"`
	JSONKey       string         `xml:"http://marklogic.com/appservices/search json-key,omitempty" json:"json-key,omitempty"`
	Field         FieldReference `xml:"http://marklogic.com/appservices/search field,omitempty" json:"field,omitempty"`
	FragmentScope string         `xml:"http://marklogic.com/appservices/search fragment-scope,omitempty" json:"fragment-scope,omitempty"`
	Text          []string       `xml:"http://marklogic.com/appservices/search text,omitempty" json:"text,omitempty"`
	TermOptions   []string       `xml:"http://marklogic.com/appservices/search term-option,omitempty" json:"term-option,omitempty"`
	Weight        float64        `xml:"http://marklogic.com/appservices/search weight,omitempty" json:"weight,omitempty"`
}

// QueryParent represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type QueryParent struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search parent" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
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

// Point represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Point struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search point" json:"-"`
	Latitude  float64  `xml:"http://marklogic.com/appservices/search latitude" json:"latitude,omitempty"`
	Longitude float64  `xml:"http://marklogic.com/appservices/search longitude" json:"longitude,omitempty"`
}

// Box represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Box struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search box" json:"-"`
	South   float64  `xml:"http://marklogic.com/appservices/search south" json:"south,omitempty"`
	West    float64  `xml:"http://marklogic.com/appservices/search west" json:"west,omitempty"`
	North   float64  `xml:"http://marklogic.com/appservices/search north" json:"north,omitempty"`
	East    float64  `xml:"http://marklogic.com/appservices/search east" json:"east,omitempty"`
}

// Circle represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Circle struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search circle" json:"-"`
	Radius  float64  `xml:"http://marklogic.com/appservices/search radius" json:"radius,omitempty"`
	Point   Point    `xml:"http://marklogic.com/appservices/search point" json:"point,omitempty"`
}

// Polygon represents http://docs.marklogic.com/guide/search-dev/structured-query#id_87280
type Polygon struct {
	XMLName xml.Name `xml:"http://marklogic.com/appservices/search polygon" json:"-"`
	Points  []*Point `xml:"http://marklogic.com/appservices/search point" json:"point,omitempty"`
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

// Lat represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
type Lat struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lat" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
}

// Lon represents http://docs.marklogic.com/guide/search-dev/structured-query#id_18303
type Lon struct {
	XMLName   xml.Name `xml:"http://marklogic.com/appservices/search lon" json:"-"`
	Namespace string   `xml:"ns,attr"`
	Local     string   `xml:"name,attr"`
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

// UnmarshalXML Query converts to XML
func (q *Query) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML OrQuery converts to XML
func (q *OrQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

type fakeAndQuery AndQuery

// UnmarshalXML AndQuery converts to XML
func (q *AndQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fake := fakeAndQuery(*q)
	err := d.DecodeElement(&fake, &start)
	if err != nil {
		return err
	}
	q.Ordered = fake.Ordered
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML PositiveQuery converts to XML
func (q *PositiveQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML NegativeQuery converts to XML
func (q *NegativeQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML NotQuery converts to XML
func (q *NotQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

type fakeNearQuery NearQuery

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
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML MatchingQuery converts to XML
func (q *MatchingQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML BoostingQuery converts to XML
func (q *BoostingQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML PropertiesQuery converts to XML
func (q *PropertiesQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

type fakeContainerQuery ContainerQuery

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
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML DocumentFragmentQuery converts to XML
func (q *DocumentFragmentQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// UnmarshalXML LocksQuery converts to XML
func (q *LocksQuery) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	queries, err2 := SerializeXMLWithQueries(d, start)
	q.Queries = queries
	return err2
}

// SerializeXMLWithQueries Serializes text into Query struct
func SerializeXMLWithQueries(d *xml.Decoder, start xml.StartElement) ([]interface{}, error) {
	var queries []interface{}
	for {
		if token, err := d.Token(); (err == nil) && (token != nil) {
			switch t := token.(type) {
			case xml.StartElement:
				e := xml.StartElement(t)
				q := stringToQueryStruct(e.Name.Local)
				err = d.DecodeElement(q, &e)
				queries = append(queries, q)
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

func stringToQueryStruct(inputString string) interface{} {
	switch inputString {
	case "query":
		return &Query{}
	case "or-query":
		return &OrQuery{}
	case "and-query":
		return &AndQuery{}
	case "term-query":
		return &TermQuery{}
	case "and-not-query":
		return &AndNotQuery{}
	case "positive-query":
		return &PositiveQuery{}
	case "negative-query":
		return &NegativeQuery{}
	case "not-query":
		return &NotQuery{}
	case "not-in-query":
		return &NotInQuery{}
	case "near-query":
		return &NearQuery{}
	case "boost-query":
		return &BoostQuery{}
	case "matching-query":
		return &MatchingQuery{}
	case "boosting-query":
		return &BoostingQuery{}
	case "properties-query":
		return &PropertiesQuery{}
	case "directory-query":
		return &DirectoryQuery{}
	case "collection-query":
		return &CollectionQuery{}
	case "container-query":
		return &ContainerQuery{}
	case "element":
		return &QueryElement{}
	case "attribute":
		return &QueryAttribute{}
	case "document-query":
		return &DocumentQuery{}
	case "document-fragment-query":
		return &DocumentFragmentQuery{}
	case "locks-query":
		return &LocksQuery{}
	case "range-query":
		return &RangeQuery{}
	case "field":
		return &FieldReference{}
	case "value-query":
		return &ValueQuery{}
	case "word-query":
		return &WordQuery{}
	case "parent":
		return &QueryParent{}
	case "heatmap":
		return &HeatMap{}
	case "point":
		return &Point{}
	case "box":
		return &Box{}
	case "circle":
		return &Circle{}
	case "polygon":
		return &Polygon{}
	case "geo-elem-query":
		return &GeoElemQuery{}
	case "lat":
		return &Lat{}
	case "lon":
		return &Lon{}
	case "geo-elem-pair-query":
		return &GeoElemPairQuery{}
	case "geo-attr-pair-query":
		return &GeoAttrPairQuery{}
	case "geo-path-query":
		return &GeoPathQuery{}
	default:
		return nil
	}
}
