package metadata

import (
	"net/http"
	"net/url"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// MetadataExtraction represents metadata extraction configuration
type MetadataExtraction struct {
	URIs       []string          `json:"uris" xml:"uris"`
	Extractors []string          `json:"extractors" xml:"extractors"`
	Options    map[string]string `json:"options,omitempty" xml:"options,omitempty"`
}

// ValidationRule represents a single validation rule
type ValidationRule struct {
	Name        string `json:"name" xml:"name"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
	RuleType    string `json:"rule-type" xml:"rule-type"`
	XPath       string `json:"xpath,omitempty" xml:"xpath,omitempty"`
	Action      string `json:"action" xml:"action"`
	Message     string `json:"message,omitempty" xml:"message,omitempty"`
}

// MetadataResult represents the result of metadata extraction
type MetadataResult struct {
	URI      string         `json:"uri" xml:"uri"`
	Metadata map[string]any `json:"metadata" xml:"metadata"`
	Valid    bool           `json:"valid" xml:"valid"`
	Errors   []string       `json:"errors,omitempty" xml:"errors,omitempty"`
}

// ValidationResult represents the result of document validation
type ValidationResult struct {
	URI    string   `json:"uri" xml:"uri"`
	Valid  bool     `json:"valid" xml:"valid"`
	Errors []string `json:"errors,omitempty" xml:"errors,omitempty"`
}

// MetadataHandle implements handle.ResponseHandle for metadata responses
type MetadataHandle struct {
	*handle.RawHandle
}

// NewMetadataHandle creates a new MetadataHandle with the specified format
func NewMetadataHandle(format int) *MetadataHandle {
	return &MetadataHandle{
		RawHandle: &handle.RawHandle{
			Format: format,
		},
	}
}

// AcceptResponse handles HTTP response deserialization
func (h *MetadataHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(h, resp)
}

func extractMetadata(c *clients.Client, uris []string, options map[string]string, response handle.ResponseHandle) error {
	params := "?"
	for _, uri := range uris {
		params = params + "&uri=" + url.QueryEscape(uri)
	}
	if options != nil {
		params = util.MappedParameters(params, "", options)
	}
	req, err := http.NewRequest("GET", c.Base()+"/metadata"+params, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", handle.FormatEnumToMimeType(response.GetFormat()))
	return util.Execute(c, req, response)
}

func extractMetadataFromQuery(c *clients.Client, query handle.Handle, options map[string]string, response handle.ResponseHandle) error {
	params := "?"
	if options != nil {
		params = util.MappedParameters(params, "", options)
	}
	req, err := util.BuildRequestFromHandle(c, "POST", "/metadata/query"+params, query)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func validateDocuments(c *clients.Client, uris []string, validationRules handle.Handle, response handle.ResponseHandle) error {
	params := "?"
	for _, uri := range uris {
		params = params + "&uri=" + url.QueryEscape(uri)
	}
	req, err := util.BuildRequestFromHandle(c, "POST", "/metadata/validate"+params, validationRules)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func validateQuery(c *clients.Client, query handle.Handle, validationRules handle.Handle, response handle.ResponseHandle) error {
	// First, create a combined payload with query and validation rules
	// For now, we'll use a sequential approach where validation rules are passed separately
	params := "?validate=true"
	req, err := util.BuildRequestFromHandle(c, "POST", "/metadata/query-validate"+params, query)
	if err != nil {
		return err
	}
	// Add validation rules as a custom header for this example
	req.Header.Set("X-Validation-Rules", validationRules.Serialized())
	return util.Execute(c, req, response)
}

func getValidationRules(c *clients.Client, response handle.ResponseHandle) error {
	req, err := http.NewRequest("GET", c.Base()+"/metadata/validation-rules", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", handle.FormatEnumToMimeType(response.GetFormat()))
	return util.Execute(c, req, response)
}

func setValidationRules(c *clients.Client, rules handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "PUT", "/metadata/validation-rules", rules)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func extractMetadataFromURI(c *clients.Client, uri string, options map[string]string, response handle.ResponseHandle) error {
	params := "?uri=" + url.QueryEscape(uri)
	if options != nil {
		params = util.MappedParameters(params, "", options)
	}
	req, err := http.NewRequest("GET", c.Base()+"/metadata/document"+params, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", handle.FormatEnumToMimeType(response.GetFormat()))
	return util.Execute(c, req, response)
}

func validateURI(c *clients.Client, uri string, validationRules handle.Handle, response handle.ResponseHandle) error {
	params := "?uri=" + url.QueryEscape(uri)
	req, err := util.BuildRequestFromHandle(c, "POST", "/metadata/validate-document"+params, validationRules)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
