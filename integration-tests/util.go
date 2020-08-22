package integrationtests

import (
	"strings"

	marklogic "github.com/ryanjdew/go-marklogic-go"
	dataMovementMod "github.com/ryanjdew/go-marklogic-go/datamovement"
	dataServicesMod "github.com/ryanjdew/go-marklogic-go/dataservices"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	searchMod "github.com/ryanjdew/go-marklogic-go/search"
)

// TODO make this config
var _client, _ = marklogic.NewClient("localhost", 8011, "admin", "admin", marklogic.DigestAuth)
var _search = _client.Search()
var _dataMovement = _client.DataMovement()
var _dataServices = _client.DataServices()

func client() *marklogic.Client {
	return _client
}

func search() *searchMod.Service {
	return _search
}

func dataMovement() *dataMovementMod.Service {
	return _dataMovement
}

func dataServices() *dataServicesMod.Service {
	return _dataServices
}

func clearDocs() error {
	writeExtensionModule("clearDocs.mjs", `
		declareUpdate();
		for (let collection of cts.collections()) {
			xdmp.collectionDelete(collection);
		}
		`, "vnd.marklogic-js-module")
	writeExtensionModule("clearDocs.api", `
		{
			"endpoint": "/ext/clearDocs.mjs",
			"params": [],
			"return": {
				"datatype": "jsonDocument",
				"multiple": false,
				"nullable": true
			}
		}`, "json")
	return dataServices().CallDataService("/ext/clearDocs.mjs", make(map[string][]string), nil)
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

func writeExtensionModule(uri string, moduleBody string, extensionType string) {
	options := map[string]string{}
	client().Config().CreateExtension(uri, strings.NewReader(moduleBody), extensionType, options, &handle.RawHandle{})
}
