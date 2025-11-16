package indexes

import (
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// RangeIndex represents a range index configuration
type RangeIndex struct {
	ScalarType string `json:"scalar-type" xml:"scalar-type"`
	Collation  string `json:"collation,omitempty" xml:"collation,omitempty"`
	Namespace  string `json:"namespace,omitempty" xml:"namespace,omitempty"`
	LocalName  string `json:"localname" xml:"localname"`
	RangeValue string `json:"range-value-positions,omitempty" xml:"range-value-positions,omitempty"`
}

// FieldIndex represents a field index configuration
type FieldIndex struct {
	FieldName        string   `json:"field-name" xml:"field-name"`
	Collation        string   `json:"collation,omitempty" xml:"collation,omitempty"`
	Tokenizer        string   `json:"tokenizer,omitempty" xml:"tokenizer,omitempty"`
	ScalarTypes      []string `json:"scalar-types" xml:"scalar-types"`
	Namespaces       []string `json:"namespaces,omitempty" xml:"namespaces,omitempty"`
	IncludedElements []string `json:"included-elements,omitempty" xml:"included-elements,omitempty"`
	ExcludedElements []string `json:"excluded-elements,omitempty" xml:"excluded-elements,omitempty"`
}

// IndexHandle implements handle.ResponseHandle for index responses
type IndexHandle struct {
	*handle.RawHandle
}

// NewIndexHandle creates a new IndexHandle with the specified format
func NewIndexHandle(format int) *IndexHandle {
	return &IndexHandle{
		RawHandle: &handle.RawHandle{
			Format: format,
		},
	}
}

func listIndexes(c *clients.Client, params map[string]string, response handle.ResponseHandle) error {
	queryStr := "?"
	if params != nil {
		queryStr = util.MappedParameters(queryStr, "", params)
		if queryStr == "?" {
			queryStr = ""
		}
	}
	req, err := http.NewRequest("GET", c.Base()+"/config/indexes"+queryStr, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func getIndex(c *clients.Client, indexName string, params map[string]string, response handle.ResponseHandle) error {
	queryStr := "?"
	if params != nil {
		queryStr = util.MappedParameters(queryStr, "", params)
		if queryStr == "?" {
			queryStr = ""
		}
	}
	req, err := http.NewRequest("GET", c.Base()+"/config/indexes/"+indexName+queryStr, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func createIndex(c *clients.Client, requestBody handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "POST", "/config/indexes", requestBody)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func updateIndex(c *clients.Client, indexName string, requestBody handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "PUT", "/config/indexes/"+indexName, requestBody)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func deleteIndex(c *clients.Client, indexName string, params map[string]string, response handle.ResponseHandle) error {
	queryStr := "?"
	if params != nil {
		queryStr = util.MappedParameters(queryStr, "", params)
		if queryStr == "?" {
			queryStr = ""
		}
	}
	req, err := http.NewRequest("DELETE", c.Base()+"/config/indexes/"+indexName+queryStr, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
