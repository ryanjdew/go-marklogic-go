// +build integration

package integrationTests

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

func TestWriteBatcher(t *testing.T) {
	clearDocs()
	defer clearDocs()
	testCount := 100
	writeChannel := make(chan *documents.DocumentDescription, 100)
	writeBatcher := dataMovement().WriteBatcher().WithBatchSize(20).WithWriteChannel(writeChannel)
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
	collection1CountWant := int64(testCount)
	collection1CountResult := collectionCount("collection-1")
	if collection1CountResult != collection1CountWant {
		t.Errorf("Collection Count 'collection-1' = %d, Want = %d", collection1CountResult, collection1CountWant)
	}
}

func TestReadBatcher(t *testing.T) {
	clearDocs()
	defer clearDocs()
	testCount := 100
	writeChannel := make(chan *documents.DocumentDescription, 100)
	writeBatcher := dataMovement().WriteBatcher().WithBatchSize(20).WithWriteChannel(writeChannel)
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
	spew.Dump("Closing channel...")
	close(writeChannel)
	spew.Dump("Channel closed")
	spew.Dump("Waiting on batcher to finish...")
	writeBatcher.Wait()
	spew.Dump("Batcher to finished")
	collection1CountWant := int64(testCount)
	collection1CountResult := collectionCount("collection-1")
	if collection1CountResult != collection1CountWant {
		t.Errorf("Collection Count 'collection-1' = %d, Want = %d", collection1CountResult, collection1CountWant)
	}
}
