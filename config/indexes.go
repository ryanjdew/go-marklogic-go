package config

import (
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// IndexesReport shows the status of indexes in query options
func indexesReport(c *clients.Client, response handle.ResponseHandle) error {
	req, err := http.NewRequest("GET", c.Base()+"/config/indexes", nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}
