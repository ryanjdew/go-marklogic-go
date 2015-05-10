package goMarklogicGo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
)

func testClient(resp string) (*Client, *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	})
	server := httptest.NewServer(handler)
	client, _ := NewClient("localhost", 8000, "admin", "admin", BasicAuth)
	client.setBase(server.URL)
	return client, server
}

func testManagementClient(resp string) (*ManagementClient, *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	})
	server := httptest.NewServer(handler)
	client, _ := NewManagementClient("localhost", "admin", "admin", BasicAuth)
	client.setBase(server.URL)
	return client, server
}

func normalizeSpace(str string) string {
	reSpace := regexp.MustCompile("\\s+")
	normalizedSpace := string(reSpace.ReplaceAllString(str, ` `))
	reBrackets := regexp.MustCompile("(\\}|\\{|\"|,|:|\\])\\s+(,|\\[|\\}|\")")
	adjustBrackets := string(reBrackets.ReplaceAllString(normalizedSpace, `$1$2`))
	reProp := regexp.MustCompile("(\\}|,|:)\\s+([^\\s])")
	adjustProperties := string(reProp.ReplaceAllString(adjustBrackets, `$1$2`))
	return strings.TrimSpace(adjustProperties)
}
