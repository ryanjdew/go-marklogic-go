package util

import (
	"fmt"
	"io/ioutil"
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
		if value != "" {
			separator := "&"
			if params == "?" {
				separator = ""
			}
			params = params + separator + prefix + key + "=" + url.QueryEscape(value)
		}
	}
	return params
}

// AddDatabaseParam is a utility function for adding the database parameter
func AddDatabaseParam(params string, client *clients.Client) string {
	if client.Database() != "" {
		separator := "&"
		if params == "" {
			separator = "?"
		}
		params = params + separator + "database=" + url.QueryEscape(client.Database())
	}
	return params
}

// AddTransactionParam is a utility function for adding a transaction parameter
func AddTransactionParam(params string, transaction *Transaction) string {
	if transaction != nil {
		if transaction.ID == "" {
			transaction.Begin()
		}
		separator := "&"
		if params == "" {
			separator = "?"
		}
		params = params + separator + "txid=" + url.QueryEscape(transaction.ID)
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
	} else {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, err
}

// Execute uses a client to run a request and places the results in the
// response Handle
func Execute(c clients.RESTClient, req *http.Request, responseHandle handle.ResponseHandle) error {
	clients.ApplyAuth(c, req)
	respHandleNotNil := responseHandle != nil
	var respType string
	if respHandleNotNil {
		respType = handle.FormatEnumToMimeType(responseHandle.GetFormat())
	}
	req.Header.Add("Accept", respType)
	resp, err := c.HTTPClient().Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP call returned status %v", resp.StatusCode)
	}
	if respHandleNotNil {
		return responseHandle.AcceptResponse(resp)
	}
	ioutil.ReadAll(resp.Body)
	
	return nil
}
