package rowsManagement

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// rows executes an Optic plan and returns results
func rows(c *clients.Client, opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := util.BuildRequestFromHandle(c, "POST", "/rows"+paramsStr, opticPlan)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// explain returns query plan and execution information for an Optic plan
func explain(c *clients.Client, opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := util.BuildRequestFromHandle(c, "POST", "/rows/explain"+paramsStr, opticPlan)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// sample returns a sample of rows from an Optic plan execution
func sample(c *clients.Client, opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := util.BuildRequestFromHandle(c, "POST", "/rows/sample"+paramsStr, opticPlan)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// plan returns the optimized query plan for an Optic operation
func plan(c *clients.Client, opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := util.BuildRequestFromHandle(c, "POST", "/rows/plan"+paramsStr, opticPlan)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
