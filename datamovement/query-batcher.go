package datamovement

import (
	"sync"

	"github.com/cchatfield/go-marklogic-go/clients"
	handle "github.com/cchatfield/go-marklogic-go/handle"
	"github.com/cchatfield/go-marklogic-go/util"
)

// QueryBatcher reads documents in bulk
type QueryBatcher struct {
	mutex         *sync.Mutex
	client        *clients.Client
	clientsByHost map[string]*clients.Client
	batchSize     uint16
	query         handle.Handle
	timestamp     string
	listeners     []chan<- *QueryBatch
	waitGroup     *sync.WaitGroup
	forestInfo    []ForestInfo
	transaction   *util.Transaction
}

// BatchSize is the number documents we'll retrieve in a single batch
func (qbr *QueryBatcher) BatchSize() uint16 {
	return qbr.batchSize
}

// Timestamp for the read operation
func (qbr *QueryBatcher) Timestamp() string {
	return qbr.timestamp
}

// WithBatchSize set the batch size
func (qbr *QueryBatcher) WithBatchSize(batchSize uint16) *QueryBatcher {
	qbr.batchSize = batchSize
	return qbr
}

// WithQuery add a query
func (qbr *QueryBatcher) WithQuery(query handle.Handle) *QueryBatcher {
	qbr.query = query
	return qbr
}

// WithTransaction add a transaction
func (qbr *QueryBatcher) WithTransaction(transaction *util.Transaction) *QueryBatcher {
	qbr.transaction = transaction
	return qbr
}

// WithListener add a listener channel
func (qbr *QueryBatcher) WithListener(listener chan *QueryBatch) *QueryBatcher {
	qbr.listeners = append(qbr.listeners, listener)
	return qbr
}

// RemoveListener remove the a listener
func (qbr *QueryBatcher) RemoveListener(listener chan *QueryBatch) *QueryBatcher {
	for i, compareListener := range qbr.listeners {
		if compareListener == listener {
			copy(qbr.listeners[i:], qbr.listeners[i+1:])
			break
		}
	}
	return qbr
}

// Run the QueryBatcher
func (qbr *QueryBatcher) Run() *QueryBatcher {
	qbr.waitGroup = &sync.WaitGroup{}
	for _, forest := range qbr.forestInfo {
		qbr.waitGroup.Add(1)
		go runReadThread(qbr, forest)
	}
	return qbr
}

// Wait on the QueryBatcher to finish
func (qbr *QueryBatcher) Wait() *QueryBatcher {
	qbr.waitGroup.Wait()
	return qbr
}

func runReadThread(queryBatcher *QueryBatcher, forest ForestInfo) {
	listeners := queryBatcher.listeners
	batchSize := int(queryBatcher.BatchSize())
	wg := queryBatcher.waitGroup
	forestClient := queryBatcher.clientsByHost[forest.PreferredHost()]
	after := ""
	defer wg.Done()
	for {
		urisHandle := &URIsHandle{timestamp: queryBatcher.timestamp}
		getURIs(forestClient, queryBatcher.query, forest.Name, queryBatcher.transaction, 0, after, uint(batchSize), urisHandle)
		if queryBatcher.timestamp == "" {
			queryBatcher.mutex.Lock()
			if queryBatcher.timestamp == "" {
				queryBatcher.timestamp = urisHandle.Timestamp()
			}
			queryBatcher.mutex.Unlock()
		}
		uris := urisHandle.Get()
		if len(uris) > 0 {
			after = uris[len(uris)-1]
		}
		queryBatch := &QueryBatch{
			client:    forestClient,
			URIs:      uris,
			timestamp: queryBatcher.timestamp,
		}
		for _, listener := range listeners {
			listener <- queryBatch
		}
		if len(queryBatch.URIs) < batchSize {
			return
		}
	}
}

// QueryBatch batch of URIs matching a query and relevant meta information
type QueryBatch struct {
	client    *clients.Client
	URIs      []string
	timestamp string
}

// Client used to with forest for the documents
func (qb *QueryBatch) Client() *clients.Client {
	return qb.client
}

// Timestamp for transaction
func (qb *QueryBatch) Timestamp() string {
	return qb.timestamp
}
