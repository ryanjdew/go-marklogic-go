package dataservices

import (
	"context"
	"io"
	"testing"
	"time"

	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

func TestBulkDataServiceIterator_NoClients_ReturnsEOF(t *testing.T) {
	bds := &BulkDataService{}
	it := bds.Iterator(context.Background())
	defer it.Close()
	_, err := it.Next(context.Background())
	if err != io.EOF {
		t.Fatalf("expected io.EOF for no clients, got: %v", err)
	}
}

func TestBulkDataServiceIterator_CancelledContext_ReturnsDeadlineExceeded(t *testing.T) {
	inputChannel := make(chan *handle.Handle)
	bds := &BulkDataService{
		clientsByHost: map[string]*clients.Client{"localhost": {}},
		inputChannel:  inputChannel,
		endpoint:      "/v1/mock",
		batchSize:     10,
		threadCount:   1,
	}
	it := bds.Iterator(context.Background())
	defer it.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := it.Next(ctx)
	if err == nil {
		t.Fatalf("expected deadline exceeded error, got nil")
	}
}

func TestBulkDataServiceIterator_CloseReturnsEOF(t *testing.T) {
	inputChannel := make(chan *handle.Handle)
	bds := &BulkDataService{
		clientsByHost: map[string]*clients.Client{"localhost": {}},
		inputChannel:  inputChannel,
		endpoint:      "/v1/mock",
		batchSize:     10,
		threadCount:   1,
	}
	it := bds.Iterator(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, err := it.Next(context.Background())
		resultCh <- err
	}()
	time.Sleep(20 * time.Millisecond)
	it.Close()
	err := <-resultCh
	if err != io.EOF {
		t.Fatalf("expected io.EOF after Close, got: %v", err)
	}
}

func TestBulkDataServiceIterator_Integration_MultipartResponses(t *testing.T) {
	// Build multipart response body with two parts
	boundary := "myBoundary"
	part1 := "part1-value"
	part2 := "part2-value"
	body := strings.Join([]string{
		fmt.Sprintf("--%s", boundary),
		"Content-Disposition: form-data; name=\"part1\"",
		"",
		part1,
		fmt.Sprintf("--%s", boundary),
		"Content-Disposition: form-data; name=\"part2\"",
		"",
		part2,
		fmt.Sprintf("--%s--", boundary),
		"",
	}, "\r\n")

	// Setup test server returning multipart/mixed
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "multipart/mixed; boundary="+boundary)
		w.Write([]byte(body))
	}))
	defer ts.Close()

	client, _ := clients.NewClient(&clients.Connection{Host: "localhost", Port: 8000, Username: "admin", Password: "admin", AuthenticationType: clients.BasicAuth})
	client.SetBase(ts.URL)

	bds := &BulkDataService{
		clientsByHost: map[string]*clients.Client{"localhost": client},
		endpoint:      "/v1/invoke",
		batchSize:     1,
		threadCount:   1,
	}

	it := bds.Iterator(context.Background())
	defer it.Close()
	// read two parts
	v, err := it.Next(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(v) != part1 && string(v) != part2 {
		t.Fatalf("unexpected first part: %v", string(v))
	}
	// second
	v2, err := it.Next(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(v2) != part1 && string(v2) != part2 {
		t.Fatalf("unexpected second part: %v", string(v2))
	}
	// then EOF
	_, err = it.Next(context.Background())
	if err != io.EOF {
		t.Fatalf("expected EOF after parts, got: %v", err)
	}
}
