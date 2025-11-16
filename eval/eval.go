package eval

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// evalCode executes code (XQuery or JavaScript) on the server.
func evalCode(c *clients.Client, language string, code string, params map[string]string, response handle.ResponseHandle) error {
	// Build base path with language parameter
	path := "/eval"
	if language != "" {
		path = path + "?rex:language=" + language
	}

	// Build additional parameters
	var parameters strings.Builder
	if language != "" {
		parameters.WriteString("&")
	} else {
		parameters.WriteString("?")
	}

	// Add mapped params if any
	if len(params) > 0 {
		for key, value := range params {
			parameters.WriteString("&")
			parameters.WriteString(key)
			parameters.WriteString("=")
			parameters.WriteString(value)
		}
	}

	// Add database param
	fullPath := path + parameters.String()
	fullPath = util.AddDatabaseParam(fullPath, c)
	fullPath = util.AddTransactionParam(fullPath, nil)

	// Create request with code body
	req, err := http.NewRequest("POST", c.Base()+fullPath, bytes.NewBufferString(code))
	if err != nil {
		return err
	}

	// Set Content-Type based on language
	if language == "javascript" {
		req.Header.Set("Content-Type", "application/javascript")
	} else {
		req.Header.Set("Content-Type", "text/xquery")
	}

	return util.Execute(c, req, response)
}
