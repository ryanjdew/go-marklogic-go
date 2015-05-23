package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
)

// Client is exported for testing subpackages
func Client(resp string) (*clients.Client, *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	})
	server := httptest.NewServer(handler)
	client, _ := clients.NewClient("localhost", 8000, "admin", "admin", clients.BasicAuth)
	client.SetBase(server.URL)
	return client, server
}

// ManagementClient is exported for testing subpackages
func ManagementClient(resp string) (*clients.ManagementClient, *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	})
	server := httptest.NewServer(handler)
	client, _ := clients.NewManagementClient("localhost", "admin", "admin", clients.BasicAuth)
	client.SetBase(server.URL)
	return client, server
}
