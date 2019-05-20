package rowsManagement

import (
	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

func rows(c *clients.Client, opticPlan handle.Handle, params map[string]string, response handle.ResponseHandle) error {
	paramsStr := util.MappedParameters("?", "", params)
	paramsStr = util.AddDatabaseParam(paramsStr, c)
	req, err := util.BuildRequestFromHandle(c, "POST", "/rows"+paramsStr, opticPlan)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
