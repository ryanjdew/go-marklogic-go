package goMarklogicGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// DatabaseProperties represents the properties of a MarkLogic Database
type DatabaseProperties struct {
	DatabaseName                        string                       `xml:"http://marklogic.com/manage database-name" json:"database-name"`
	Forest                              []string                     `xml:"http://marklogic.com/manage forest" json:"forest"`
	SecurityDatabase                    string                       `xml:"http://marklogic.com/manage security-database" json:"security-database"`
	SchemaDatabase                      string                       `xml:"http://marklogic.com/manage schema-database" json:"schema-database"`
	Enabled                             bool                         `xml:"http://marklogic.com/manage enabled" json:"enabled"`
	RetiredForestCount                  int                          `xml:"http://marklogic.com/manage retired-forest-count" json:"retired-forest-count"`
	Language                            string                       `xml:"http://marklogic.com/manage language" json:"language"`
	StemmedSearches                     string                       `xml:"http://marklogic.com/manage stemmed-searches" json:"stemmed-searches"`
	WordSearches                        bool                         `xml:"http://marklogic.com/manage word-searches" json:"word-searches"`
	WordPositions                       bool                         `xml:"http://marklogic.com/manage word-positions" json:"word-positions"`
	FastPhraseSearches                  bool                         `xml:"http://marklogic.com/manage fast-phrase-searches" json:"fast-phrase-searches"`
	FastReverseSearches                 bool                         `xml:"http://marklogic.com/manage fast-reverse-searches" json:"fast-reverse-searches"`
	TripleIndex                         bool                         `xml:"http://marklogic.com/manage triple-index" json:"triple-index"`
	TriplePositions                     bool                         `xml:"http://marklogic.com/manage triple-positions" json:"triple-positions"`
	FastCaseSensitiveSearches           bool                         `xml:"http://marklogic.com/manage fast-case-sensitive-searches" json:"fast-case-sensitive-searches"`
	FastDiacriticSensitiveSearches      bool                         `xml:"http://marklogic.com/manage fast-diacritic-sensitive-searches" json:"fast-diacritic-sensitive-searches"`
	FastElementWordSearches             bool                         `xml:"http://marklogic.com/manage fast-element-word-searches" json:"fast-element-word-searches"`
	ElementWordPositions                bool                         `xml:"http://marklogic.com/manage element-word-positions" json:"element-word-positions"`
	FastElementPhraseSearches           bool                         `xml:"http://marklogic.com/manage fast-element-phrase-searches" json:"fast-element-phrase-searches"`
	ElementValuePositions               bool                         `xml:"http://marklogic.com/manage element-value-positions" json:"element-value-positions"`
	AttributeValuePositions             bool                         `xml:"http://marklogic.com/manage attribute-value-positions" json:"attribute-value-positions"`
	FieldValueSearches                  bool                         `xml:"http://marklogic.com/manage field-value-searches" json:"field-value-searches"`
	FieldValuePositions                 bool                         `xml:"http://marklogic.com/manage field-value-positions" json:"field-value-positions"`
	ThreeCharacterSearches              bool                         `xml:"http://marklogic.com/manage three-character-searches" json:"three-character-searches"`
	ThreeCharacterWordPositions         bool                         `xml:"http://marklogic.com/manage three-character-word-positions" json:"three-character-word-positions"`
	FastElementCharacterSearches        bool                         `xml:"http://marklogic.com/manage fast-element-character-searches" json:"fast-element-character-searches"`
	TrailingWildcardSearches            bool                         `xml:"http://marklogic.com/manage trailing-wildcard-searches" json:"trailing-wildcard-searches"`
	TrailingWildcardWordPositions       bool                         `xml:"http://marklogic.com/manage trailing-wildcard-word-positions" json:"trailing-wildcard-word-positions"`
	FastElementTrailingWildcardSearches bool                         `xml:"http://marklogic.com/manage fast-element-trailing-wildcard-searches" json:"fast-element-trailing-wildcard-searches"`
	TwoCharacterSearches                bool                         `xml:"http://marklogic.com/manage two-character-searches" json:"two-character-searches"`
	OneCharacterSearches                bool                         `xml:"http://marklogic.com/manage one-character-searches" json:"one-character-searches"`
	URILexicon                          bool                         `xml:"http://marklogic.com/manage uri-lexicon" json:"uri-lexicon"`
	CollectionLexicon                   bool                         `xml:"http://marklogic.com/manage collection-lexicon" json:"collection-lexicon"`
	ReindexerEnable                     bool                         `xml:"http://marklogic.com/manage reindexer-enable" json:"reindexer-enable"`
	ReindexerThrottle                   int                          `xml:"http://marklogic.com/manage reindexer-throttle" json:"reindexer-throttle"`
	ReindexerTimestamp                  int                          `xml:"http://marklogic.com/manage reindexer-timestamp" json:"reindexer-timestamp"`
	DirectoryCreation                   string                       `xml:"http://marklogic.com/manage directory-creation" json:"directory-creation"`
	MaintainLastModified                bool                         `xml:"http://marklogic.com/manage maintain-last-modified" json:"maintain-last-modified"`
	MaintainDirectoryLastModified       bool                         `xml:"http://marklogic.com/manage maintain-directory-last-modified" json:"maintain-directory-last-modified"`
	InheritPermissions                  bool                         `xml:"http://marklogic.com/manage inherit-permissions" json:"inherit-permissions"`
	InheritCollections                  bool                         `xml:"http://marklogic.com/manage inherit-collections" json:"inherit-collections"`
	InheritQuality                      bool                         `xml:"http://marklogic.com/manage inherit-quality" json:"inherit-quality"`
	InMemoryLimit                       int                          `xml:"http://marklogic.com/manage in-memory-limit" json:"in-memory-limit"`
	InMemoryListSize                    int                          `xml:"http://marklogic.com/manage in-memory-list-size" json:"in-memory-list-size"`
	InMemoryTreeSize                    int                          `xml:"http://marklogic.com/manage in-memory-tree-size" json:"in-memory-tree-size"`
	InMemoryRangeIndexSize              int                          `xml:"http://marklogic.com/manage in-memory-range-index-size" json:"in-memory-range-index-size"`
	InMemoryReverseIndexSize            int                          `xml:"http://marklogic.com/manage in-memory-reverse-index-size" json:"in-memory-reverse-index-size"`
	InMemoryTripleIndexSize             int                          `xml:"http://marklogic.com/manage in-memory-triple-index-size" json:"in-memory-triple-index-size"`
	LargeSizeThredhold                  int                          `xml:"http://marklogic.com/manage large-size-thredhold" json:"large-size-thredhold"`
	Locking                             string                       `xml:"http://marklogic.com/manage locking" json:"locking"`
	Journaling                          string                       `xml:"http://marklogic.com/manage journaling" json:"journaling"`
	JournalSize                         int                          `xml:"http://marklogic.com/manage journal-size" json:"journal-size"`
	JournalCount                        int                          `xml:"http://marklogic.com/manage journal-count" json:"journal-count"`
	PreallocateJournals                 bool                         `xml:"http://marklogic.com/manage preallocate-journals" json:"preallocate-journals"`
	PreloadMappedData                   bool                         `xml:"http://marklogic.com/manage preload-mapped-data" json:"preload-mapped-data"`
	PreloadReplicaMappedData            bool                         `xml:"http://marklogic.com/manage preload-replica-mapped-data" json:"preload-replica-mapped-data"`
	RangeIndexOptimize                  string                       `xml:"http://marklogic.com/manage range-index-optimize" json:"range-index-optimize"`
	PositionsListMaxSize                int                          `xml:"http://marklogic.com/manage positions-list-max-size" json:"positions-list-max-size"`
	FormatCompatibility                 string                       `xml:"http://marklogic.com/manage format-compatibility" json:"format-compatibility"`
	IndexDetection                      string                       `xml:"http://marklogic.com/manage index-detection" json:"index-detection"`
	ExpungeLocks                        string                       `xml:"http://marklogic.com/manage expunge-locks" json:"expunge-locks"`
	TfNormalization                     string                       `xml:"http://marklogic.com/manage tf-normalization" json:"tf-normalization"`
	MergePriority                       string                       `xml:"http://marklogic.com/manage merge-priority" json:"merge-priority"`
	MergeMaxSize                        int                          `xml:"http://marklogic.com/manage merge-max-size" json:"merge-max-size"`
	MergeMinSize                        int                          `xml:"http://marklogic.com/manage merge-min-size" json:"merge-min-size"`
	MergeMinRatio                       int                          `xml:"http://marklogic.com/manage merge-min-ratio" json:"merge-min-ratio"`
	MergeTimestamp                      int                          `xml:"http://marklogic.com/manage merge-timestamp" json:"merge-timestamp"`
	RetainUntilBackup                   bool                         `xml:"http://marklogic.com/manage retain-until-backup" json:"retain-until-backup"`
	ElementWordQueryThrough             []ElementWordQueryThrough    `xml:"http://marklogic.com/manage element-word-query-through" json:"element-word-query-through"`
	PhraseThrough                       []PhraseThrough              `xml:"http://marklogic.com/manage phrase-through" json:"phrase-through"`
	PhraseAround                        []PhraseAround               `xml:"http://marklogic.com/manage phrase-around" json:"phrase-around"`
	RangeElementIndex                   []RangeElementIndex          `xml:"http://marklogic.com/manage range-element-index" json:"range-element-index"`
	RangeElementAttributeIndex          []RangeElementAttributeIndex `xml:"http://marklogic.com/manage range-element-attribute-index" json:"range-element-attribute-index"`
	Field                               []Field                      `xml:"http://marklogic.com/manage field" json:"field"`
	RangeFieldIndex                     []RangeFieldIndex            `xml:"http://marklogic.com/manage range-field-index" json:"range-field-index"`
	DatabaseReplication                 string                       `xml:"http://marklogic.com/manage database-replication" json:"database-replication"`
	RebalancerEnable                    bool                         `xml:"http://marklogic.com/manage rebalancer-enable" json:"rebalancer-enable"`
	RebalancerThrottle                  int                          `xml:"http://marklogic.com/manage rebalancer-throttle" json:"rebalancer-throttle"`
	AssignmentPolicy                    AssignmentPolicy             `xml:"http://marklogic.com/manage assignment-policy" json:"assignment-policy"`
}

// AssignmentPolicy struct reprenting an assignment policy in the database
type AssignmentPolicy struct {
	AssignmentPolicyName string `xml:"http://marklogic.com/manage namespace-uri" json:"assignment-policy-name"`
}

// ElementWordQueryThrough struct reprenting an element word query through in the database
type ElementWordQueryThrough struct {
	NamespaceURI string   `xml:"http://marklogic.com/manage namespace-uri" json:"namespace-uri"`
	Localname    []string `xml:"http://marklogic.com/manage localname" json:"localname"`
}

// PhraseThrough struct reprenting a phrase through in the database
type PhraseThrough struct {
	NamespaceURI string   `xml:"http://marklogic.com/manage namespace-uri" json:"namespace-uri"`
	Localname    []string `xml:"http://marklogic.com/manage localname" json:"localname"`
}

// PhraseAround struct reprenting a phrase around in the database
type PhraseAround struct {
	NamespaceURI string   `xml:"http://marklogic.com/manage namespace-uri" json:"namespace-uri"`
	Localname    []string `xml:"http://marklogic.com/manage localname" json:"localname"`
}

// RangeElementIndex struct reprenting an element range index in the database
type RangeElementIndex struct {
	ScalarType          string `xml:"http://marklogic.com/manage scalar-type" json:"scalar-type"`
	NamespaceURI        string `xml:"http://marklogic.com/manage namespace-uri" json:"namespace-uri"`
	Localname           string `xml:"http://marklogic.com/manage localname" json:"localname"`
	Collation           string `xml:"http://marklogic.com/manage collation" json:"collation"`
	RangeValuePositions bool   `xml:"http://marklogic.com/manage range-value-positions" json:"range-value-positions"`
	InvalidValues       string `xml:"http://marklogic.com/manage invalid-values" json:"invalid-values"`
}

// RangeElementAttributeIndex struct reprenting an element attribute range index in the database
type RangeElementAttributeIndex struct {
	ScalarType          string `xml:"http://marklogic.com/manage scalar-type" json:"scalar-type"`
	ParentNamespaceURI  string `xml:"http://marklogic.com/manage parent-namespace-uri" json:"parent-namespace-uri"`
	ParentLocalname     string `xml:"http://marklogic.com/manage parent-localname" json:"parent-localname"`
	NamespaceURI        string `xml:"http://marklogic.com/manage namespace-uri" json:"namespace-uri"`
	Localname           string `xml:"http://marklogic.com/manage localname" json:"localname"`
	Collation           string `xml:"http://marklogic.com/manage collation" json:"collation"`
	RangeValuePositions bool   `xml:"http://marklogic.com/manage range-value-positions" json:"range-value-positions"`
	InvalidValues       string `xml:"http://marklogic.com/manage invalid-values" json:"invalid-values"`
}

// Field struct reprenting a field in the database
type Field struct {
	FieldName          string            `xml:"http://marklogic.com/manage field-name" json:"field-name"`
	IncludeRoot        bool              `xml:"http://marklogic.com/manage include-root" json:"include-root"`
	WordLexicon        []string          `xml:"http://marklogic.com/manage word-lexicon" json:"word-lexicon"`
	IncludedElement    []IncludedElement `xml:"http://marklogic.com/manage included-element" json:"included-element"`
	ExcludedElement    []IncludedElement `xml:"http://marklogic.com/manage included-element" json:"excluded-element"`
	TokenizerOverrides string            `xml:"http://marklogic.com/manage tokenizer-overrides" json:"tokenizer-overrides,omitempty"`
}

// IncludedElement struct containing information about an element contained in a field
type IncludedElement struct {
	NamespaceURI          string `xml:"http://marklogic.com/manage namespace-uri" json:"namespace-uri"`
	Localname             string `xml:"http://marklogic.com/manage localname" json:"localname"`
	Weight                int    `xml:"http://marklogic.com/manage weight" json:"weight"`
	AttributeNamespaceURI string `xml:"http://marklogic.com/manage attribute-namespace-uri" json:"attribute-namespace-uri"`
	AttributeLocalname    string `xml:"http://marklogic.com/manage attribute-localname" json:"attribute-localname"`
	AttributeValue        string `xml:"http://marklogic.com/manage attribute-value" json:"attribute-value"`
}

// RangeFieldIndex struct reprenting a field range index in the database
type RangeFieldIndex struct {
	ScalarType          string `xml:"http://marklogic.com/manage scalar-type" json:"scalar-type"`
	FieldName           string `xml:"http://marklogic.com/manage field-name" json:"field-name"`
	Collation           string `xml:"http://marklogic.com/manage collation" json:"collation"`
	RangeValuePositions bool   `xml:"http://marklogic.com/manage range-value-positions" json:"range-value-positions"`
	InvalidValues       string `xml:"http://marklogic.com/manage invalid-values" json:"invalid-values"`
}

// SetDatabaseProperties sets the database properties
func (mc *ManagementClient) SetDatabaseProperties(databaseName string, propertiesHandle Handle) error {
	var reqType string
	if propertiesHandle.GetFormat() == JSON {
		reqType = "json"
	} else {
		reqType = "xml"
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(propertiesHandle.Serialized()))
	req, _ := http.NewRequest("PUT", mc.Base()+"/databases/"+databaseName+"/properties?format="+reqType, buf)
	applyAuth(mc, req)
	req.Header.Add("Content-Type", "application/"+reqType)
	resp, err := mc.HTTPClient().Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

// GetDatabaseProperties sets the database properties
func (mc *ManagementClient) GetDatabaseProperties(databaseName string, propertiesHandle Handle) error {
	var reqType string
	if propertiesHandle.GetFormat() == JSON {
		reqType = "json"
	} else {
		reqType = "xml"
	}
	req, _ := http.NewRequest("GET", mc.Base()+"/databases/"+databaseName+"/properties?format="+reqType, nil)
	applyAuth(mc, req)
	req.Header.Add("Content-Type", "application/"+reqType)
	resp, err := mc.HTTPClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	propertiesHandle.Encode([]byte(contents))
	return err
}

// DatabasePropertiesHandle is a handle that places the results into
// a DatabaseProperties struct
type DatabasePropertiesHandle struct {
	Format             int
	bytes              []byte
	databaseProperties *DatabaseProperties
}

// GetFormat returns int that represents XML or JSON
func (dh *DatabasePropertiesHandle) GetFormat() int {
	return dh.Format
}

// Encode returns Query struct that represents XML or JSON
func (dh *DatabasePropertiesHandle) Encode(bytes []byte) *DatabaseProperties {
	dh.bytes = bytes
	dh.databaseProperties = &DatabaseProperties{}
	if dh.GetFormat() == JSON {
		json.Unmarshal(bytes, &dh.databaseProperties)
	} else {
		xml.Unmarshal(bytes, &dh.databaseProperties)
	}
	return dh.databaseProperties
}

// Decode returns []byte of XML or JSON that represents the Query struct
func (dh *DatabasePropertiesHandle) Decode(databaseProperties *DatabaseProperties) []byte {
	dh.databaseProperties = databaseProperties
	buf := new(bytes.Buffer)
	if dh.GetFormat() == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(dh.databaseProperties)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(dh.databaseProperties)
	}
	dh.bytes = buf.Bytes()
	return dh.bytes
}

// Get returns string of XML or JSON
func (dh *DatabasePropertiesHandle) Get() *DatabaseProperties {
	return dh.databaseProperties
}

// Serialized returns string of XML or JSON
func (dh *DatabasePropertiesHandle) Serialized() string {
	buf := new(bytes.Buffer)
	if dh.GetFormat() == JSON {
		enc := json.NewEncoder(buf)
		enc.Encode(dh.databaseProperties)
	} else {
		enc := xml.NewEncoder(buf)
		enc.Encode(dh.databaseProperties)
	}
	dh.bytes = buf.Bytes()
	return string(dh.bytes)
}