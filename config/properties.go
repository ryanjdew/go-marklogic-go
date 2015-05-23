package config

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// GetProperties shows the REST API properties
func getProperties(c *clients.Client, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "GET", "/config/properties", nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// SetProperties sets the REST API properties
func setProperties(c *clients.Client, properties handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "PUT", "/config/properties", properties)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// ResetProperties resets the REST API properties to their default
func resetProperties(c *clients.Client, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "DELETE", "/config/properties", nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// SetPropertyValue sets a property of the REST API
func setPropertyValue(c *clients.Client, propertyName string, property handle.Handle, response handle.ResponseHandle) error {
	req, err := util.BuildRequestFromHandle(c, "PUT", "/config/properties/"+propertyName, property)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
