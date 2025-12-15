package datamovement

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// TestQueryBatcherIterator_Integration_ReturnsURIs tests the QueryBatcher iterator
// against a mock HTTP server returning URI lists for /internal/uris.
func TestQueryBatcherIterator_Integration_ReturnsURIs(t *testing.T) {
	// Create a client that will respond to /internal/uris with two URIs
	// Create a client backed by a httptest server that will respond to /internal/uris
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/internal/uris" {
			w.Header().Set("Content-Type", "text/uri-list")
			// precise CRLF-separated URIs
			w.Write([]byte("/doc1.json\r\n/doc2.json\r\n"))
			return
		}
		w.WriteHeader(404)
	}))
	defer server.Close()
	client, _ := clients.NewClient(&clients.Connection{Host: "localhost", Port: 8000, Username: "admin", Password: "admin", AuthenticationType: clients.BasicAuth})
	client.SetBase(server.URL)

	// Create a single forest and ensure the client is mapped for that host
	// Use host name "localhost" for simplicity
	b := &QueryBatcher{
		mutex:         &sync.Mutex{},
		client:        client,
		clientsByHost: map[string]*clients.Client{"localhost": client},
		forestInfo:    []util.ForestInfo{{ID: "1", Name: "f1", Host: "localhost"}},
		batchSize:     10,
	}

	it := b.Iterator(context.Background())
	defer it.Close()

	batch, err := it.Next(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(batch.URIs) != 2 {
		t.Fatalf("expected 2 URIs, got %d", len(batch.URIs))
	}
	// Next should return EOF
	_, err = it.Next(context.Background())
	if err != io.EOF {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}
