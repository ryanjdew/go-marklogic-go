package config

import (
	"io"
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// ListExtensions shows all the installed REST extensions
func listExtensions(c *clients.Client, directory string, response handle.Handle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/ext"+directory, nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

// DeleteExtensions shows all the installed REST extensions
func deleteExtensions(c *clients.Client, directory string, response handle.Handle) error {
	req, err := util.BuildRequestFromHandle(c, "DELETE", "/ext"+directory, nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

// createExtension shows all the installed REST extensions
func createExtension(c *clients.Client, assetName string, resource io.Reader, extensionType string, options map[string]string, response handle.Handle) error {
	params := mapToParams(options)
	req, err := http.NewRequest("PUT", c.Base()+"/ext"+assetName+params, resource)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/"+extensionType)
	return clients.Execute(c, req, response)
}

// ListResources shows all the installed REST service extensions
func listResources(c *clients.Client, response handle.Handle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/resources", nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

// GetResourceInfo shows all the installed REST extensions
func getResourceInfo(c *clients.Client, name string, response handle.Handle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/resources/"+name, nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

// CreateResource installs a REST service
func createResource(c *clients.Client, name string, resource io.Reader, extensionType string, options map[string]string, response handle.Handle) error {
	params := mapToParams(options)
	req, err := http.NewRequest("PUT", c.Base()+"/config/resources/"+name+params, resource)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/"+extensionType)
	return clients.Execute(c, req, response)
}

// DeleteResource removes a REST service
func deleteResource(c *clients.Client, name string, response handle.Handle) error {
	req, err := util.BuildRequestFromHandle(c, "DELETE", "/config/resources/"+name, nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

func mapToParams(options map[string]string) string {
	params := "?"
	for key, val := range options {
		separator := "&"
		if params == "?" {
			separator = ""
		}
		params = params + separator + key + "=" + val
	}
	return params
}
