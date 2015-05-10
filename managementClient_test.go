package goMarklogicGo

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestManagementConnection(t *testing.T) {
	expectedBase := "http://localhost:8002/manage/v2"
	client, err := NewManagementClient("localhost", "admin", "admin", BasicAuth)
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}

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
	client, server := testManagementClient(dbPropertiesWantResp)
	defer server.Close()
	propertiesHandle := DatabasePropertiesHandle{Format: JSON}
	client.GetDatabaseProperties("Documents", &propertiesHandle)
	want := normalizeSpace(dbPropertiesWantResp)
	result := normalizeSpace(propertiesHandle.Serialized())
	if !reflect.DeepEqual(want, result) {
		t.Errorf("DB Properties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}

var serverPropertiesWantResp = `{
  "server-name": "samplestack",
  "group-name": "Default",
  "server-type": "http",
  "enabled": true,
  "root": "/",
  "port": 8006,
  "webDAV": false,
  "execute": true,
  "display-last-login": false,
  "address": "0.0.0.0",
  "backlog": 512,
  "threads": 32,
  "request-timeout": 30,
  "keep-alive-timeout": 5,
  "session-timeout": 3600,
  "max-time-limit": 3600,
  "default-time-limit": 600,
  "max-inference-size": 500,
  "default-inference-size": 100,
  "static-expires": 3600,
  "pre-commit-trigger-depth": 1000,
  "pre-commit-trigger-limit": 10000,
  "collation": "http://marklogic.com/collation/",
  "authentication": "digest",
  "internal-security": true,
  "concurrent-request-limit": 0,
  "compute-content-length": true,
  "log-errors": false,
  "debug-allow": true,
  "profile-allow": true,
  "default-xquery-version": "1.0-ml",
  "multi-version-concurrency-control": "contemporaneous",
  "distribute-timestamps": "fast",
  "output-sgml-character-entities": "none",
  "output-encoding": "UTF-8",
  "output-method": "default",
  "output-byte-order-mark": "default",
  "output-cdata-section-namespace-uri": "",
  "output-cdata-section-localname": "",
  "output-doctype-public": "",
  "output-doctype-system": "",
  "output-escape-uri-attributes": "default",
  "output-include-content-type": "default",
  "output-indent": "default",
  "output-indent-untyped": "default",
  "output-media-type": "",
  "output-normalization-form": "none",
  "output-omit-xml-declaration": "default",
  "output-standalone": "omit",
  "output-undeclare-prefixes": "default",
  "output-version": "",
  "output-include-default-attributes": "default",
  "default-error-format": "json",
  "error-handler": "/MarkLogic/rest-api/error-handler.xqy",
  "url-rewriter": "/MarkLogic/rest-api/rewriter.xml",
  "rewrite-resolves-globally": true,
  "ssl-certificate-template": 0,
  "ssl-allow-sslv3": true,
  "ssl-allow-tls": true,
  "ssl-hostname": "",
  "ssl-ciphers": "ALL:!LOW:@STRENGTH",
  "ssl-require-client-certificate": true,
  "content-database": "samplestack",
  "modules-database": "samplestack-modules",
  "default-user": "nobody"
}`

func TestServerProperties(t *testing.T) {
	client, server := testManagementClient(serverPropertiesWantResp)
	defer server.Close()
	propertiesHandle := ServerPropertiesHandle{Format: JSON}
	client.GetServerProperties("Documents", "Default", &propertiesHandle)
	want := normalizeSpace(serverPropertiesWantResp)
	result := normalizeSpace(propertiesHandle.Serialized())
	if !reflect.DeepEqual(want, result) {
		t.Errorf("Server Properties Results = %+v, Want = %+v", spew.Sdump(result), spew.Sdump(want))
	}
}
