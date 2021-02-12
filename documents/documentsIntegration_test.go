// +build integration

package documents

import (
	"bytes"
	"os"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/integrationtests"
)

func TestMain(m *testing.M) {
	integrationtests.ClearDocs()
	code := m.Run()
	integrationtests.ClearDocs()
	os.Exit(code)
}

func TestWriteReadDocuments(t *testing.T) {
	client := integrationtests.Client()
	docSet := []*DocumentDescription{
		&DocumentDescription{
			URI:     "/test.json",
			Format:  handle.JSON,
			Content: bytes.NewBufferString(`{ "test": "json"}`),
			Metadata: &Metadata{
				Collections: []string{"collection-1", "collection-2"},
				MetadataValues: map[string]string{
					"metadata1": "val1",
					"metadata2": "val2",
				},
			},
		},
		&DocumentDescription{
			URI:     "/test2.xml",
			Format:  handle.XML,
			Content: bytes.NewBufferString(`<root/>`),
			Metadata: &Metadata{
				Collections: []string{"collection-2"},
				MetadataValues: map[string]string{
					"metadata1": "val1",
					"metadata2": "val2",
				},
			},
		},
	}
	response := &handle.RawHandle{}
	err := NewService(client).WriteSet(docSet, &MetadataHandle{}, nil, nil, response)
	if err != nil {
		t.Errorf("Error writing documents: %+v", err)
	}
	collection1CountWant := int64(1)
	collection1CountResult := integrationtests.CollectionCount("collection-1")
	if collection1CountResult != collection1CountWant {
		t.Errorf("Collection Count 'collection-1' = %d, Want = %d", collection1CountResult, collection1CountWant)
	}
	collection2CountWant := int64(2)
	collection2CountResult := integrationtests.CollectionCount("collection-2")
	if collection2CountResult != collection2CountWant {
		t.Errorf("Collection Count 'collection-2' = %d, Want = %d", collection2CountResult, collection2CountWant)
	}
}
