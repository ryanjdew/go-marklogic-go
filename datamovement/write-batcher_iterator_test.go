package datamovement

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/ryanjdew/go-marklogic-go/documents"
)

func TestWriteBatcherIterator_NoWrites_ReturnsEOF(t *testing.T) {
	wbr := &WriteBatcher{}
	it := wbr.Iterator(context.Background())
	defer it.Close()
	_, err := it.Next(context.Background())
	if err != io.EOF {
		t.Fatalf("expected io.EOF for no writes, got: %v", err)
	}
}

func TestWriteBatcherIterator_CancelledContext_ReturnsDeadlineExceeded(t *testing.T) {
	writeChannel := make(chan *documents.DocumentDescription)
	wbr := &WriteBatcher{
		writeChannel:           writeChannel,
		documentsServiceByHost: map[string]*documents.Service{"dummy": nil},
		batchSize:              10,
		threadCount:            1,
	}
	if len(wbr.documentsServiceByHost) == 0 {
		t.Fatalf("documentsServiceByHost map should be non-empty for iterator test")
	}
	it := wbr.Iterator(context.Background())
	defer it.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := it.Next(ctx)
	if err == nil {
		t.Fatalf("expected deadline exceeded error, got nil")
	}
}

func TestWriteBatcherIterator_CloseReturnsEOF(t *testing.T) {
	writeChannel := make(chan *documents.DocumentDescription)
	wbr := &WriteBatcher{
		writeChannel:           writeChannel,
		documentsServiceByHost: map[string]*documents.Service{"dummy": nil},
		batchSize:              10,
		threadCount:            1,
	}
	it := wbr.Iterator(context.Background())
	resultCh := make(chan error, 1)
	go func() {
		_, err := it.Next(context.Background())
		resultCh <- err
	}()
	// Give the goroutine a chance to block on Next
	time.Sleep(20 * time.Millisecond)
	it.Close()
	err := <-resultCh
	if err != io.EOF {
		t.Fatalf("expected io.EOF after Close, got: %v", err)
	}
}
