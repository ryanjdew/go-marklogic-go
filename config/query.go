package config

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// ListQueryOptions shows all the installed REST query options
func listQueryOptions(c *clients.Client, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/query", nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// DeleteAllQueryOptions removes all the installed REST query options
func deleteAllQueryOptions(c *clients.Client, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "DELETE", "/config/query", nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// SetQueryOptions shows all the installed REST extensions
func setQueryOptions(c *clients.Client, optionsName string, options handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "PUT", "/config/query/"+optionsName, options)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// GetQueryOptions returns the named REST query options
func getQueryOptions(c *clients.Client, name string, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/query/"+name, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// DeleteQueryOptions removes the named REST query options
func deleteQueryOptions(c *clients.Client, name string, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "DELETE", "/config/query/"+name, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
