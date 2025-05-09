//go:build integration

package datamovement

import (
	"bytes"
	"strconv"
	"sync"
	"testing"

	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	integrationtests "github.com/ryanjdew/go-marklogic-go/integrationtests"
	searchMod "github.com/ryanjdew/go-marklogic-go/search"
)

var dataMovement = NewService(integrationtests.Client())

var testCount int = 1000

func writeDocumentsWithWriteBatcher() {
	writeChannel := make(chan *documents.DocumentDescription, 100)
	writeBatcher := dataMovement.WriteBatcher().WithBatchSize(250).WithWriteChannel(writeChannel)
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
	integrationtests.ClearDocs()
	defer integrationtests.ClearDocs()
	writeDocumentsWithWriteBatcher()
	collection1CountWant := int64(testCount)
	collection1CountResult := integrationtests.CollectionCount("collection-1")
	if collection1CountResult != collection1CountWant {
		t.Errorf("Collection Count 'collection-1' = %d, Want = %d", collection1CountResult, collection1CountWant)
	}
}

func returnURIsWithReadBatcher(t *testing.T, collection string) []string {
	uris := make([]string, 0, testCount)
	listenerChannel := make(chan *QueryBatch, 100)
	query :=
		searchMod.Query{
			Queries: []any{
				searchMod.CollectionQuery{
					URIs: []string{collection},
				},
			},
		}
	qh := searchMod.QueryHandle{Format: handle.JSON}
	qh.Serialize(query)
	queryBatcher := dataMovement.QueryBatcher().WithQuery(&qh).WithBatchSize(100).WithListener(listenerChannel)
	queryBatcher.Run()
	timestamp := ""
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for queryBatch := range listenerChannel {
			if queryBatch != nil {
				if timestamp != "" && timestamp != queryBatch.Timestamp() {
					t.Errorf("Has timestamp = %s, Expected = %s", queryBatch.Timestamp(), timestamp)
				}
				uris = append(uris, queryBatch.URIs...)
			}
		}
		wg.Done()
	}()
	queryBatcher.Wait()
	close(listenerChannel)
	wg.Wait()
	return uris
}

func TestReadBatcher(t *testing.T) {
	integrationtests.ClearDocs()
	defer integrationtests.ClearDocs()
	writeDocumentsWithWriteBatcher()
	uris := returnURIsWithReadBatcher(t, "collection-1")
	uriCountWant := int(testCount)
	uriCountResult := len(uris)
	if uriCountResult != uriCountWant {
		t.Errorf("URI Count = %d, Want = %d", uriCountResult, uriCountWant)
	}
}
