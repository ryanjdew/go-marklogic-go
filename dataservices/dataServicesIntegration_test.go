// build integration

package dataservices

import (
	"fmt"
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	integrationtests "github.com/ryanjdew/go-marklogic-go/integrationtests"
)

var dataServices = NewService(integrationtests.Client())
var testCount int64 = 1000

func TestProcessDataService(t *testing.T) {
	integrationtests.ClearDocs()
	integrationtests.WriteExtensionModule("process.mjs", `
		declareUpdate();
		let endpointState = external.endpointState ? fn.head(xdmp.fromJSON(external.endpointState)) : {};
		const workUnit = fn.head(xdmp.fromJSON(external.workUnit)) || {};
		const batchSize = workUnit.batchSize || 25;
		const queries = [];
	
		if (endpointState.lastProcessedUri) {
			queries.push(cts.rangeQuery(cts.uriReference(), ">", endpointState.lastProcessedUri));
		}
		const uris = cts.uris(null, ['limit=' + batchSize], cts.andQuery(queries), null, [workUnit.forestId]);
		if (fn.empty(uris)) {
			null;
		} else {
			const outputArray = [];
			for (const uri of uris) {
				xdmp.documentAddCollections(fn.string(uri), ['processed']);
				outputArray.push(uri);
			}
			endpointState = Object.assign({}, endpointState, { lastProcessedUri: outputArray[outputArray.length - 1]});
			Sequence.from([endpointState].concat(outputArray));
		}
		`, "vnd.marklogic-js-module")
	integrationtests.WriteExtensionModule("process.api", `
		{
			"endpoint": "/ext/process.mjs",
			"params": [
				{
					"name": "endpointState",
					"datatype": "jsonDocument",
					"multiple": false,
					"nullable": true
				},
				{
					"name": "workUnit",
					"datatype": "jsonDocument",
					"multiple": false,
					"nullable": false
				}
			],
			"return": {
				"datatype": "jsonDocument",
				"multiple": true,
				"nullable": true
			}
		}`, "json")
	writeDocumentsWithDataServices()
	bulkDataService := dataServices.BulkDataService("/ext/process.mjs").WithForestBasedWorkUnits().WithEndpointState([]byte("{}"))
	bulkDataService.Run()
	bulkDataService.Wait()
	uriCountWant := testCount
	uriCountResult := integrationtests.CollectionCount("processed")
	if uriCountResult != uriCountWant {
		t.Errorf("URI Count = %d, Want = %d", uriCountResult, uriCountWant)
	}
}

func TestInputDataService(t *testing.T) {
	integrationtests.ClearDocs()
	writeDocumentsWithDataServices()
	uriCountWant := testCount
	uriCountResult := integrationtests.CollectionCount("collection-1")
	if uriCountResult != uriCountWant {
		t.Errorf("URI Count = %d, Want = %d", uriCountResult, uriCountWant)
	}
}

func writeDocumentsWithDataServices() {
	integrationtests.WriteExtensionModule("/input.mjs", `
		declareUpdate();
		const input = external.input ? external.input.toArray().map(contentObj => contentObj.toObject()) : [];
		const outputArray = [];
		for (const content of input) {
			xdmp.documentInsert(content.uri, content.value, content.context);
			outputArray.push(content.uri);
		}
		Sequence.from(outputArray);
		`, "vnd.marklogic-js-module")
	integrationtests.WriteExtensionModule("/input.api", `
		{
			"endpoint": "/ext/input.mjs",
			"params": [
				{
					"name": "input",
					"datatype": "jsonDocument",
					"multiple": true,
					"nullable": true
				}
			],
			"return": {
				"datatype": "jsonDocument",
				"multiple": true,
				"nullable": true
			}
		}`, "json")
	inputChanel := make(chan *handle.Handle, 25)
	bulkService := dataServices.BulkDataService("/ext/input.mjs").WithBatchSize(25).WithInputChannel(inputChanel)
	bulkService.Run()
	for i := int64(0); i < testCount; i++ {
		countStr := fmt.Sprint(i)
		var contentObject handle.Handle = &handle.RawHandle{
			Format: handle.JSON,
		}
		contentObject.Deserialize([]byte(`{
			"uri": "/test-doc-` + countStr + `.json",
			"value": {
				"testId": ` + countStr + `
			},
			"context": {
				"collections": ["collection-1", "collection-2"]
			}
		}`))
		inputChanel <- &contentObject
	}
	close(inputChanel)
	bulkService.Wait()
}
