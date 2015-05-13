package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

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

//NormalizeSpace is a function that takes whitespace out of the equation for comparing strings
func NormalizeSpace(str string) string {
	reSpace := regexp.MustCompile("\\s+")
	normalizedSpace := string(reSpace.ReplaceAllString(str, ` `))
	reBrackets := regexp.MustCompile("(\\}|\\{|\"|,|:|\\])\\s+(,|\\[|\\}|\")")
	adjustBrackets := string(reBrackets.ReplaceAllString(normalizedSpace, `$1$2`))
	reProp := regexp.MustCompile("(\\}|,|:)\\s+([^\\s])")
	adjustProperties := string(reProp.ReplaceAllString(adjustBrackets, `$1$2`))
	return strings.TrimSpace(adjustProperties)
}
