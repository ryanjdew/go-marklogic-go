package management

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	test "github.com/ryanjdew/go-marklogic-go/test"
)

var dbPropertiesWantResp = `{
  "database-name": "samplestack",
  "forest": ["samplestack-1", "samplestack-2", "samplestack-3"],
  "security-database": "Security",
  "schema-database": "Schemas",
  "enabled": true,
  "retired-forest-count": 0,
  "language": "en",
  "stemmed-searches": "basic",
  "word-searches": false,
  "word-positions": false,
  "fast-phrase-searches": true,
  "fast-reverse-searches": false,
  "triple-index": true,
  "triple-positions": false,
  "fast-case-sensitive-searches": true,
  "fast-diacritic-sensitive-searches": true,
  "fast-element-word-searches": true,
  "element-word-positions": false,
  "fast-element-phrase-searches": true,
  "element-value-positions": false,
  "attribute-value-positions": false,
  "field-value-searches": false,
  "field-value-positions": false,
  "three-character-searches": true,
  "three-character-word-positions": false,
  "fast-element-character-searches": true,
  "trailing-wildcard-searches": false,
  "trailing-wildcard-word-positions": false,
  "fast-element-trailing-wildcard-searches": false,
  "two-character-searches": false,
  "one-character-searches": false,
  "uri-lexicon": true,
  "collection-lexicon": false,
  "reindexer-enable": true,
  "reindexer-throttle": 5,
  "reindexer-timestamp": 0,
  "directory-creation": "manual",
  "maintain-last-modified": true,
  "maintain-directory-last-modified": false,
  "inherit-permissions": false,
  "inherit-collections": false,
  "inherit-quality": false,
  "in-memory-limit": 262144,
  "in-memory-list-size": 512,
  "in-memory-tree-size": 128,
  "in-memory-range-index-size": 16,
  "in-memory-reverse-index-size": 16,
  "in-memory-triple-index-size": 64,
  "large-size-threshold": 1024,
  "locking": "fast",
  "journaling": "fast",
  "journal-size": 1365,
  "journal-count": 2,
  "preallocate-journals": false,
  "preload-mapped-data": false,
  "preload-replica-mapped-data": false,
  "range-index-optimize": "facet-time",
  "positions-list-max-size": 256,
  "format-compatibility": "automatic",
  "index-detection": "automatic",
  "expunge-locks": "none",
  "tf-normalization": "scaled-log",
  "merge-priority": "lower",
  "merge-max-size": 32768,
  "merge-min-size": 1024,
  "merge-min-ratio": 2,
  "merge-timestamp": 0,
  "retain-until-backup": false,
  "element-word-query-through": [{
    "namespace-uri": "http://schemas.microsoft.com/office/word/2003/wordml",
    "localname": "p"
  }, {
    "namespace-uri": "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
    "localname": "p"
  }],
  "phrase-through": [{
    "namespace-uri": "http://marklogic.com/entity",
    "localname": ["coordinate", "credit-card-number", "date", "email", "facility", "gpe", "id", "location", "money", "nationality", "organization", "percent", "person", "phone-number", "religion", "time", "url", "utm"]
  }, {
    "namespace-uri": "http://schemas.microsoft.com/office/word/2003/auxHint",
    "localname": "t"
  }, {
    "namespace-uri": "http://schemas.microsoft.com/office/word/2003/wordml",
    "localname": ["br", "cr", "fldChar", "fldData", "fldSimple", "hlink", "noBreakHyphen", "permEnd", "permStart", "pgNum", "proofErr", "r", "softHyphen", "sym", "t", "tab"]
  }, {
    "namespace-uri": "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
    "localname": ["bookmarkEnd", "bookmarkStart", "commentRangeEnd", "commentRangeStart", "customXml", "endnoteReference", "fldSimple", "footnoteReference", "hyperlink", "ins", "instrText", "proofErr", "r", "sdt", "sdtContent", "smartTag", "t"]
  }, {
    "namespace-uri": "http://www.w3.org/1999/xhtml",
    "localname": ["a", "abbr", "acronym", "b", "big", "br", "center", "cite", "code", "dfn", "em", "font", "i", "kbd", "q", "samp", "small", "span", "strong", "sub", "sup", "tt", "var"]
  }],
  "phrase-around": [{
    "namespace-uri": "http://schemas.microsoft.com/office/word/2003/wordml",
    "localname": ["delInstrText", "delText", "endnote", "footnote", "instrText", "pict", "rPr"]
  }, {
    "namespace-uri": "http://schemas.openxmlformats.org/wordprocessingml/2006/main",
    "localname": ["commentReference", "customXmlPr", "del", "pPr", "rPr", "sdtPr"]
  }],
  "range-element-index": [{
    "scalar-type": "string",
    "namespace-uri": "",
    "localname": "item",
    "collation": "http://marklogic.com/collation/",
    "range-value-positions": false,
    "invalid-values": "ignore"
  }, {
    "scalar-type": "string",
    "namespace-uri": "",
    "localname": "cat",
    "collation": "http://marklogic.com/collation/",
    "range-value-positions": false,
    "invalid-values": "ignore"
  }, {
    "scalar-type": "dateTime",
    "namespace-uri": "http://marklogic.com/xdmp/property",
    "localname": "last-modified",
		"collation": "",
    "range-value-positions": false,
    "invalid-values": "ignore"
  }],
  "field": [{
    "field-name": "",
    "include-root": true,
    "tokenizer-overrides": ""
  }],
  "database-replication": "",
  "rebalancer-enable": true,
  "rebalancer-throttle": 5,
  "assignment-policy": {
    "assignment-policy-name": "bucket"
  }
}`

func TestDatabaseProperties(t *testing.T) {
	client, server := test.ManagementClient(dbPropertiesWantResp)
	defer server.Close()
	propertiesHandle := DatabasePropertiesHandle{Format: handle.JSON}
	GetDatabaseProperties(client, "Documents", &propertiesHandle)
	want := test.NormalizeSpace(dbPropertiesWantResp)
	result := test.NormalizeSpace(propertiesHandle.Serialized())
	if !reflect.DeepEqual(want, result) {
		t.Errorf("DB Properties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}
