package integrationtests

import (
	marklogic "github.com/cchatfield/go-marklogic-go"
	dataMovementMod "github.com/cchatfield/go-marklogic-go/datamovement"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	searchMod "github.com/cchatfield/go-marklogic-go/search"
)

// TODO make this config
var _client, _ = marklogic.NewClient("localhost", 8000, "admin", "admin", marklogic.DigestAuth)
var _search = _client.Search()
var _dataMovement = _client.DataMovement()

func client() *marklogic.Client {
	return _client
}

func search() *searchMod.Service {
	return _search
}

func dataMovement() *dataMovementMod.Service {
	return _dataMovement
}

func clearDocs() error {
	return search().Delete(map[string]string{}, nil, &handle.RawHandle{})
}

func collectionCount(collection string) int64 {
	resp := &searchMod.ResponseHandle{}
	structuredQuery := &searchMod.QueryHandle{
		Query: searchMod.Query{
			Queries: []interface{}{searchMod.CollectionQuery{URIs: []string{collection}}},
		},
	}
	search().StructuredSearch(structuredQuery, 1, 0, nil, resp)
	return resp.Get().Total
}
