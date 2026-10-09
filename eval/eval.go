package eval

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// evalCode executes code (XQuery or JavaScript) on the server.
func evalCode(c *clients.Client, language string, code string, params map[string]string, transaction *util.Transaction, response handle.ResponseHandle) error {
	// Build query parameters
	pathParams := ""
	pathParams = util.AddDatabaseParam(pathParams, c)
	pathParams = util.AddTransactionParam(pathParams, transaction)

	// Build form-encoded body with code and parameters
	formData := url.Values{}

	// Add the code with appropriate parameter name based on language
	if language == "javascript" {
		formData.Set("javascript", code)
	} else {
		formData.Set("xquery", code)
	}

	// Add external variables if any parameters provided
	if len(params) > 0 {
		varsJSON, err := json.Marshal(params)
		if err != nil {
			return err
		}
		formData.Set("vars", string(varsJSON))
	}

	// Create request with form-encoded body
	req, err := http.NewRequest("POST", c.Base()+"/eval"+pathParams, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return err
	}

	// Set Content-Type to form-urlencoded
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return util.Execute(c, req, response)
}
