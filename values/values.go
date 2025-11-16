package values

import (
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// listValues retrieves lexicon values for a specified range index or field
func listValues(c *clients.Client, name string, params map[string]string, response handle.ResponseHandle) error {
	parameters := util.MappedParameters("?", "", params)
	parameters = util.AddDatabaseParam(parameters, c)
	req, err := http.NewRequest("GET", c.Base()+"/values/"+name+parameters, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// queryValues queries the values in a lexicon or range index
func queryValues(c *clients.Client, name string, params map[string]string, requestBody handle.Handle, response handle.ResponseHandle) error {
	parameters := util.MappedParameters("?", "", params)
	parameters = util.AddDatabaseParam(parameters, c)
	req, err := util.BuildRequestFromHandle(c, "POST", "/values/"+name+parameters, requestBody)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
