package config

import (
	"net/http"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// IndexesReport shows the status of indexes in query options
func indexesReport(c *clients.Client, response handle.Handle) error {
	req, err := http.NewRequest("GET", c.Base()+"/config/indexes", nil)
	if err != nil {
		return err
	}
	return clients.Execute(c, req, response)
}
