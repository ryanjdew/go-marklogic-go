package documents

import (
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

func read(c *clients.Client, uris []string, categories []string, transform *util.Transform, response handle.Handle) error {
	params := buildParameters(uris, categories, nil, transform)
	req, err := http.NewRequest("GET", c.Base()+"/documents"+params, nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

func write(c *clients.Client, uri string, document handle.Handle, collections []string, transform *util.Transform, response handle.Handle) error {
	params := buildParameters([]string{uri}, nil, collections, transform)
	req, err := util.BuildRequestFromHandle(c, "PUT", "/documents"+params, document)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

func buildParameters(uris []string, categories []string, collections []string, transform *util.Transform) string {
	params := "?"
	params = util.RepeatingParameters(params, "uri", uris)
	params = util.RepeatingParameters(params, "category", categories)
	params = util.RepeatingParameters(params, "collection", collections)
	if transform != nil {
		params = params + transform.ToParameters()
	}
	return params
}
