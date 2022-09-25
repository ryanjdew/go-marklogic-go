package util

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// SerializableStringMap is a map[string]string which can be converted to XML.
type SerializableStringMap map[string]string

// MarshalXML marshals map[string]string into XML.
func (s SerializableStringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{Space: "", Local: key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{Name: start.Name})

	for _, t := range tokens {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}
	if err := e.Flush(); err != nil {
		return err
	}

	return nil
}

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
	var req *http.Request
	var err error
	if reqHandle == nil || reqHandle.Serialized() == "" {
		req, err = http.NewRequest(method, c.Base()+uri, nil)
	} else {
		req, err = http.NewRequest(method, c.Base()+uri, reqHandle)
	}
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
	resp, err := c.Do(req)
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
	if resp != nil {
		io.ReadAll(resp.Body)
	}

	return nil
}

// PostForm submits a URL encoded form
func PostForm(c clients.RESTClient, endpoint string, atomicParams map[string][]string, unatomicParams map[string][]*handle.Handle, responseHandle handle.ResponseHandle, isDataService bool) error {
	var reader io.Reader
	var contentType string
	var contentLength int
	if len(unatomicParams) == 0 {
		// if there are not unatomic types, use urlencode
		data := url.Values{}
		for key, values := range atomicParams {
			for _, value := range values {
				data.Add(key, value)
			}
		}
		encodedData := data.Encode()
		reader = strings.NewReader(encodedData)
		contentType = "application/x-www-form-urlencoded; charset=UTF-8"
		contentLength = len(encodedData)
	} else {
		// if there are unatomic types use multipart form
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for key, values := range atomicParams {
			for _, value := range values {
				fw, _ := writer.CreateFormField(key)
				io.Copy(fw, strings.NewReader(value))
			}
		}
		for key, values := range unatomicParams {
			for _, value := range values {
				h := (*value)
				mimetype := handle.FormatEnumToMimeType(h.GetFormat())
				fw, _ := writer.CreatePart(textproto.MIMEHeader{"Content-Disposition": []string{"form-data; name=\"" + key + "\"; filename=\"\""}, "Content-Type": []string{mimetype}})
				io.Copy(fw, bytes.NewBufferString(h.Serialized()))
			}
		}
		contentType = writer.FormDataContentType()
		writer.Close()
		reader = bytes.NewBuffer(body.Bytes())
		contentLength = len(body.Bytes())
	}
	baseURL := c.Base()
	if isDataService {
		baseURL = strings.Replace(baseURL, "/LATEST", "", -1)
	}
	req, _ := http.NewRequest(http.MethodPost, baseURL+endpoint, reader) // URL-encoded payload
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Content-Length", strconv.Itoa(contentLength))

	clients.ApplyAuth(c, req)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if responseHandle != nil {
		return responseHandle.AcceptResponse(resp)
	}
	return nil
}
