// +build integration

package integrationtests

import (
	"bytes"
	"os"
	"testing"

	"github.com/cchatfield/go-marklogic-go/documents"
	handle "github.com/cchatfield/go-marklogic-go/handle"
)

func TestMain(m *testing.M) {
	clearDocs()
	code := m.Run()
	clearDocs()
	os.Exit(code)
}

func TestWriteReadDocuments(t *testing.T) {
	client := client()
	docSet := []*documents.DocumentDescription{
		&documents.DocumentDescription{
			URI:     "/test.json",
			Format:  handle.JSON,
			Content: bytes.NewBufferString(`{ "test": "json"}`),
			Metadata: &documents.Metadata{
				Collections: []string{"collection-1", "collection-2"},
				MetadataValues: map[string]string{
					"metadata1": "val1",
					"metadata2": "val2",
				},
			},
		},
		&documents.DocumentDescription{
			URI:     "/test2.xml",
			Format:  handle.XML,
			Content: bytes.NewBufferString(`<root/>`),
			Metadata: &documents.Metadata{
				Collections: []string{"collection-2"},
				MetadataValues: map[string]string{
					"metadata1": "val1",
					"metadata2": "val2",
				},
			},
		},
	}
	response := &handle.RawHandle{}
	err := client.Documents().WriteSet(docSet, &documents.MetadataHandle{}, nil, nil, response)
	if err != nil {
		t.Errorf("Error writing documents: %+v", err)
	}
	collection1CountWant := int64(1)
	collection1CountResult := collectionCount("collection-1")
	if collection1CountResult != collection1CountWant {
		t.Errorf("Collection Count 'collection-1' = %d, Want = %d", collection1CountResult, collection1CountWant)
	}
	collection2CountWant := int64(2)
	collection2CountResult := collectionCount("collection-2")
	if collection2CountResult != collection2CountWant {
		t.Errorf("Collection Count 'collection-2' = %d, Want = %d", collection2CountResult, collection2CountWant)
	}
}
