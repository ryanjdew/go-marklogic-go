package integrationtests

import (
	"strings"

	"github.com/ryanjdew/go-marklogic-go/util"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	config "github.com/ryanjdew/go-marklogic-go/config"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	searchMod "github.com/ryanjdew/go-marklogic-go/search"
)

// TODO make this config
var _client, _ = clients.NewClient(&clients.Connection{Host: "localhost", Port: 8011, Username: "admin", Password: "admin", AuthenticationType: clients.DigestAuth})
var _search = searchMod.NewService(_client)
var _config = config.NewService(_client)

func Client() *clients.Client {
	return _client
}

func Search() *searchMod.Service {
	return _search
}

func ClearDocs() error {
	WriteExtensionModule("/clearDocs.mjs", `
		declareUpdate();
		for (let collection of cts.collections()) {
			xdmp.collectionDelete(collection);
		}
		`, "vnd.marklogic-js-module")
	WriteExtensionModule("/clearDocs.api", `
		{
			"endpoint": "/ext/clearDocs.mjs",
			"params": [],
			"return": {
				"datatype": "jsonDocument",
				"multiple": false,
				"nullable": true
			}
		}`, "json")
	return util.PostForm(Client(), "/ext/clearDocs.mjs", make(map[string][]string), make(map[string][]*handle.Handle), nil, true)
}

func CollectionCount(collection string) int64 {
	resp := &searchMod.ResponseHandle{}
	structuredQuery := &searchMod.QueryHandle{
		Query: searchMod.Query{
			Queries: []interface{}{searchMod.CollectionQuery{URIs: []string{collection}}},
		},
	}
	Search().StructuredSearch(structuredQuery, 1, 0, nil, resp)
	return resp.Get().Total
}

func WriteExtensionModule(uri string, moduleBody string, extensionType string) {
	options := map[string]string{}
	_config.CreateExtension(uri, strings.NewReader(moduleBody), extensionType, options, &handle.RawHandle{})
}
