package util

import (
	"net/http"
	"net/url"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// RepeatingParameters is a utility function for putting slices to parameters
func RepeatingParameters(params string, valueLabel string, values []string) string {
	for _, value := range values {
		separator := "&"
		if params == "?" {
			separator = ""
		}
		params = params + separator + valueLabel + "=" + url.QueryEscape(value)
	}
	return params
}

// MappedParameters is a utility function for putting map[string]string to parameters
func MappedParameters(params string, prefix string, values map[string]string) string {
	if prefix != "" {
		prefix = prefix + ":"
	}
	for key, value := range values {
		separator := "&"
		if params == "?" {
			separator = ""
		}
		params = params + separator + prefix + key + "=" + url.QueryEscape(value)
	}
	return params
}

// BuildRequestFromHandle builds a *http.Request based off a handle.Handle
func BuildRequestFromHandle(c clients.RESTClient, method string, uri string, reqHandle handle.Handle) (*http.Request, error) {
	reqType := ""
	if reqHandle != nil {
		reqType = handle.FormatEnumToMimeType(reqHandle.GetFormat())
	}
	req, err := http.NewRequest(method, c.Base()+uri, reqHandle)
	if err == nil && reqType != "" {
		req.Header.Add("Content-Type", reqType)
	}
	return req, err
}
