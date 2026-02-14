package datamovement

import (
	"context"
	"io"
	"testing"
)

func TestQueryBatcherIterator_NoForests_ReturnsEOF(t *testing.T) {
	qbr := &QueryBatcher{}
	it := qbr.Iterator(context.Background())
	defer it.Close()
	_, err := it.Next(context.Background())
	if err != io.EOF {
		t.Fatalf("expected io.EOF for empty forests, got: %v", err)
	}
}
