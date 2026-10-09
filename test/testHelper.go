package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	testHelper "github.com/ryanjdew/go-marklogic-go/test/text"
)

// Client is exported for testing subpackages
func Client(resp string) (*clients.Client, *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	})
	server := httptest.NewServer(handler)
	client, _ := clients.NewClient(&clients.Connection{Host: "localhost", Port: 8000, Username: "admin", Password: "admin", AuthenticationType: clients.BasicAuth})
	client.SetBase(server.URL)
	return client, server
}

// ManagementClient is exported for testing subpackages
func ManagementClient(resp string) (*clients.ManagementClient, *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	})
	server := httptest.NewServer(handler)
	client, _ := clients.NewManagementClient(&clients.Connection{Host: "localhost", Username: "admin", Password: "admin", AuthenticationType: clients.BasicAuth})
	client.SetBase(server.URL)
	return client, server
}

// AdminClient is exported for testing subpackages
func AdminClient(resp string) (*clients.AdminClient, *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, resp)
	})
	server := httptest.NewServer(handler)
	client, _ := clients.NewAdminClient(&clients.Connection{Host: "localhost", Username: "admin", Password: "admin", AuthenticationType: clients.BasicAuth})
	client.SetBase(server.URL)
	return client, server
}

// RoundTripSerialization is a helper to simplify testing serialization/deserialization
func RoundTripSerialization(t *testing.T, msg string, deserializedInterface any, handle handle.Handle, serialization string) {
	handle.Serialize(deserializedInterface)
	normalizedSpace := testHelper.NormalizeSpace(serialization)
	roundTripSerialization := testHelper.NormalizeSpace(handle.Serialized())
	if roundTripSerialization != normalizedSpace {
		t.Errorf("Serialization of %s Results = %+v, Want = %+v", msg, spew.Sdump(roundTripSerialization), spew.Sdump(normalizedSpace))
	}
	handle.Deserialize([]byte(serialization))
	roundTripDeserialization := handle.Deserialized()
	var finalDeserialization any
	switch roundTripDeserialization.(type) {
	case *any:
		finalDeserialization = &roundTripDeserialization
	default:
		finalDeserialization = roundTripDeserialization
	}
	if !reflect.DeepEqual(deserializedInterface, finalDeserialization) {
		t.Errorf("Deserialization of %s Results = %+v, Want = %+v", msg, spew.Sdump(finalDeserialization), spew.Sdump(deserializedInterface))
	}
}
