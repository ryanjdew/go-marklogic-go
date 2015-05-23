package config

import (
	"io"
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// ListTransforms shows all the installed REST service extensions
func listTransforms(c *clients.Client, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/transforms", nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// GetTransformInfo shows all the installed REST extensions
func getTransformInfo(c *clients.Client, name string, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/transforms/"+name, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// CreateTransform installs a REST service
func createTransform(c *clients.Client, name string, resource io.Reader, extensionType string, options map[string]string, response handle.ResponseHandle) error {
	params := mapToParams(options)
	req, err := http.NewRequest("PUT", c.Base()+"/config/transforms/"+name+params, resource)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/"+extensionType)
	return util.Execute(c, req, response)
}

// DeleteTransform removes a REST service
func deleteTransform(c *clients.Client, name string, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "DELETE", "/config/transforms/"+name, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
