// build integration

package integrationtests

import (
	"fmt"
	"testing"
)

func TestProcessDataService(t *testing.T) {
	clearDocs()
	writeExtensionModule("process.mjs", `
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
	writeExtensionModule("process.api", `
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
	bulkDataService := dataServices().BulkDataService("/ext/process.mjs").WithForestBasedWorkUnits().WithEndpointState([]byte("{}"))
	bulkDataService.Run()
	bulkDataService.Wait()
	uris := returnURIsWithReadBatcher(t, "processed")
	uriCountWant := int(testCount)
	uriCountResult := len(uris)
	if uriCountResult != uriCountWant {
		t.Errorf("URI Count = %d, Want = %d", uriCountResult, uriCountWant)
	}
}

func TestInputDataService(t *testing.T) {
	clearDocs()
	writeDocumentsWithDataServices()
	uris := returnURIsWithReadBatcher(t, "collection-1")
	uriCountWant := int(testCount)
	uriCountResult := len(uris)
	if uriCountResult != uriCountWant {
		t.Errorf("URI Count = %d, Want = %d", uriCountResult, uriCountWant)
	}
}

func writeDocumentsWithDataServices() {
	writeExtensionModule("/input.mjs", `
		declareUpdate();
		const input = external.input ? external.input.toArray().map(contentObj => contentObj.toObject()) : [];
		const outputArray = [];
		for (const content of input) {
			xdmp.documentInsert(content.uri, content.value, content.context);
			outputArray.push(content.uri);
		}
		Sequence.from(outputArray);
		`, "vnd.marklogic-js-module")
	writeExtensionModule("/input.api", `
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
	inputChanel := make(chan string, 25)
	bulkService := dataServices().BulkDataService("/ext/input.mjs").WithBatchSize(25).WithInputChannel(inputChanel)
	bulkService.Run()
	for i := 0; i < testCount; i++ {
		countStr := fmt.Sprint(i)
		contentObject := `
		{
			"uri": "/test-doc-` + countStr + `.json",
			"value": {
				"testId": ` + countStr + `
			},
			"context": {
				"collections": ["collection-1", "collection-2"]
			}
		}
		`
		inputChanel <- contentObject
	}
	close(inputChanel)
	bulkService.Wait()
}
