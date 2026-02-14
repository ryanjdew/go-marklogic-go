package datamovement

import (
	"context"
	"io"
	"testing"
	"time"

	"bytes"
	"strings"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	documents "github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	test "github.com/ryanjdew/go-marklogic-go/test"
)

// TestWriteBatcherIterator_Integration_WritesToServerAndYieldsResponse validates
// that write-batcher iterator uses documents.Service and the iterator yields
// a WriteBatch with a response when the server accepts writes.
func TestWriteBatcherIterator_Integration_WritesToServerAndYieldsResponse(t *testing.T) {
	// Respond with a plain OK body for writes
	client, server := test.Client("OK")
	defer server.Close()

	// Wire up the documents service for host `localhost`
	docsService := documents.NewService(client)

	writeChannel := make(chan *documents.DocumentDescription, 1)
	// create a simple document
	doc := &documents.DocumentDescription{URI: "/integration/doc.json", Content: bytes.NewBufferString("{" + "\"x\":1}"), Format: handle.JSON}
	writeChannel <- doc
	close(writeChannel)

	wbr := &WriteBatcher{
		writeChannel:           writeChannel,
		documentsServiceByHost: map[string]*documents.Service{"localhost": docsService},
		clientsByHost:          map[string]*clients.Client{"localhost": client},
		batchSize:              1,
		threadCount:            1,
	}

	it := wbr.Iterator(context.Background())
	defer it.Close()
	batch, err := it.Next(context.Background())
	if err == io.EOF {
		t.Fatalf("unexpected EOF from iterator: %v", err)
	}
	if err != nil {
		t.Fatalf("unexpected error from iterator: %v", err)
	}
	if batch == nil || len(batch.DocumentDescriptions()) == 0 {
		t.Fatalf("expected a batch with document descriptions, got %#v", batch)
	}
	// Verify response is present
	if batch.Response() == nil {
		t.Fatalf("expected response handle on batch")
	}
	if raw, ok := batch.Response().(*handle.RawHandle); !ok {
		t.Fatalf("expected RawHandle on response")
	} else {
		if raw.Get() == "" || !strings.Contains(raw.Get(), "OK") {
			t.Fatalf("expected 'OK' in raw response, got: %#v", raw.Get())
		}
	}
	// Acquire the response by letting the iterator drain
	// and ensure EOF afterwards
	time.Sleep(5 * time.Millisecond)
	_, err = it.Next(context.Background())
	if err != io.EOF {
		t.Fatalf("expected io.EOF after write, got: %v", err)
	}
}
