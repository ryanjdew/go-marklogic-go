package documents

import (
	"io"
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// DocumentDescription describes a document to write
type DocumentDescription struct {
	URI         string
	Content     io.ReadWriter
	Collections []string
	Permissions map[string]string
	Properties  map[string]string
	Quality     int
	VersionID   int
}

func toURIs(docs []DocumentDescription) []string {
	uris := []string{}
	for _, doc := range docs {
		uris = append(uris, doc.URI)
	}
	return uris
}

func read(c *clients.Client, uris []string, categories []string, transform *util.Transform, response handle.ResponseHandle) error {
	params := buildParameters(uris, categories, nil, nil, nil, transform)
	req, err := http.NewRequest("GET", c.Base()+"/documents"+params, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func write(c *clients.Client, documents []DocumentDescription, transform *util.Transform, response handle.ResponseHandle) error {
	channel := make(chan error)
	var errReturn error
	for _, doc := range documents {
		go func(doc DocumentDescription) {
			params := buildParameters([]string{doc.URI}, nil, doc.Collections, doc.Permissions, doc.Properties, transform)
			req, err := http.NewRequest("PUT", c.Base()+"/documents"+params, doc.Content)
			if err == nil {
				err = util.Execute(c, req, response)
			}
			channel <- err
		}(doc)
	}
	for _ = range documents {
		if errReturn == nil {
			errReturn = <-channel
		} else {
			<-channel
		}
	}
	return errReturn
}

func buildParameters(uris []string, categories []string, collections []string, permissions map[string]string, properties map[string]string, transform *util.Transform) string {
	params := "?"
	params = util.RepeatingParameters(params, "uri", uris)
	params = util.RepeatingParameters(params, "category", categories)
	params = util.RepeatingParameters(params, "collection", collections)
	params = util.MappedParameters(params, "perm", permissions)
	params = util.MappedParameters(params, "prop", properties)
	if transform != nil {
		params = params + transform.ToParameters()
	}
	return params
}
