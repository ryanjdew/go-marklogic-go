package config

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// ListNamespaces shows the namespaces used in queries
func listNamespaces(c *clients.Client, response handle.Handle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/namespaces", nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}

// SetNamespace shows the namespaces used in queries
func setNamespace(c *clients.Client, namespace handle.Handle, response handle.Handle) error {
	req, err := util.BuildRequestFromHandle(c, "PUT", "/config/namespaces", namespace)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}
