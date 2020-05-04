package resources

import (
	"net/http"

	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	util "github.com/cchatfield/go-marklogic-go/util"
)

func get(c *clients.Client, resourceName string, parameters map[string]string, response handle.ResponseHandle) error {
	params := util.MappedParameters("?", "rs", parameters)
	params = util.AddDatabaseParam(params, c)
	req, err := http.NewRequest("GET", c.Base()+"/resources/"+resourceName+params, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func post(c *clients.Client, resourceName string, parameters map[string]string, requestBody handle.Handle, response handle.ResponseHandle) error {
	params := util.MappedParameters("?", "rs", parameters)
	params = util.AddDatabaseParam(params, c)
	req, err := http.NewRequest("POST", c.Base()+"/resources/"+resourceName+params, requestBody)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func put(c *clients.Client, resourceName string, parameters map[string]string, requestBody handle.Handle, response handle.ResponseHandle) error {
	params := util.MappedParameters("?", "rs", parameters)
	params = util.AddDatabaseParam(params, c)
	req, err := http.NewRequest("PUT", c.Base()+"/resources/"+resourceName+params, requestBody)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

func delete(c *clients.Client, resourceName string, parameters map[string]string, response handle.ResponseHandle) error {
	params := util.MappedParameters("?", "rs", parameters)
	params = util.AddDatabaseParam(params, c)
	req, err := http.NewRequest("DELETE", c.Base()+"/resources/"+resourceName+params, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
