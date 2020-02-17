// +build integration

package integrationtests

import (
	"bytes"
	"strconv"
	"sync"
	"testing"

	dataMovementMod "github.com/ryanjdew/go-marklogic-go/datamovement"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

var testCount int = 1000

func writeDocumentsWithWriteBatcher() {
	writeChannel := make(chan *documents.DocumentDescription, 100)
	writeBatcher := dataMovement().WriteBatcher().WithBatchSize(250).WithWriteChannel(writeChannel)
	writeBatcher.Run()
	for i := 0; i < testCount; i++ {
		docDescription := &documents.DocumentDescription{
			URI:     "/test-" + strconv.Itoa(i) + ".json",
			Format:  handle.JSON,
			Content: bytes.NewBufferString(`{ "test": "json"}`),
			Metadata: &documents.Metadata{
				Collections: []string{"collection-1", "collection-2"},
				MetadataValues: map[string]string{
					"metadata1": "val1",
					"metadata2": "val2",
				},
			},
		}
		writeChannel <- docDescription
	}
	close(writeChannel)
	writeBatcher.Wait()
}

func TestWriteBatcher(t *testing.T) {
	clearDocs()
	defer clearDocs()
	writeDocumentsWithWriteBatcher()
	collection1CountWant := int64(testCount)
	collection1CountResult := collectionCount("collection-1")
	if collection1CountResult != collection1CountWant {
		t.Errorf("Collection Count 'collection-1' = %d, Want = %d", collection1CountResult, collection1CountWant)
	}
}

func returnURIsWithReadBatcher(t *testing.T) []string {
	uris := make([]string, 0, testCount)
	listenerChannel := make(chan *dataMovementMod.QueryBatch, 100)
	queryBatcher := dataMovement().QueryBatcher().WithBatchSize(100).WithListener(listenerChannel)
	queryBatcher.Run()
	timestamp := ""
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case queryBatch, ok := <-listenerChannel:
				if queryBatch != nil {
					if timestamp != "" && timestamp != queryBatch.Timestamp() {
						t.Errorf("Has timestamp = %s, Expected = %s", queryBatch.Timestamp(), timestamp)
					}
					uris = append(uris, queryBatch.URIs...)
				} else if !ok && len(listenerChannel) == 0 {
					wg.Done()
					return
				}
			}
		}
	}()
	queryBatcher.Wait()
	close(listenerChannel)
	wg.Wait()
	return uris
}

func TestReadBatcher(t *testing.T) {
	clearDocs()
	defer clearDocs()
	writeDocumentsWithWriteBatcher()
	uris := returnURIsWithReadBatcher(t)
	uriCountWant := int(testCount)
	uriCountResult := len(uris)
	if uriCountResult != uriCountWant {
		t.Errorf("URI Count = %d, Want = %d", uriCountResult, uriCountWant)
	}
}
