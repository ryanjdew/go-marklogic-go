package datamovement

import (
	"context"
	"io"
	"sync"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
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
	forestInfo    []util.ForestInfo
	transaction   *util.Transaction
}

// QueryBatchIterator provides a pull-style iterator for QueryBatch results
type QueryBatchIterator interface {
	Next(ctx context.Context) (*QueryBatch, error)
	Close() error
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

func runReadThread(queryBatcher *QueryBatcher, forest util.ForestInfo) {
	listeners := queryBatcher.listeners
	batchSize := int(queryBatcher.BatchSize())
	wg := queryBatcher.waitGroup
	forestClient := queryBatcher.clientsByHost[forest.PreferredHost()]
	after := ""
	defer wg.Done()
	for {
		urisHandle := &util.URIsHandle{}
		urisHandle.SetTimestamp(queryBatcher.timestamp)
		util.GetURIs(forestClient, queryBatcher.query, forest.Name, queryBatcher.transaction, 0, after, uint(batchSize), urisHandle)
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

// Iterator returns an iterator over QueryBatch results. It is a non-blocking
// call which starts internal goroutines to fetch batches and provides them
// through Next(). The iterator respects ctx cancellation and Close().
func (qbr *QueryBatcher) Iterator(ctx context.Context) QueryBatchIterator {
	ctx, cancel := context.WithCancel(ctx)
	results := make(chan *QueryBatch)
	wg := &sync.WaitGroup{}

	// If there are no forests, close results and return an iterator that
	// immediately yields EOF.
	if len(qbr.forestInfo) == 0 {
		close(results)
		return &queryBatchIter{results: results, cancel: cancel, wg: wg}
	}

	for _, forest := range qbr.forestInfo {
		wg.Add(1)
		go runReadThreadIterator(qbr, forest, results, wg, ctx)
	}

	// close results when all goroutines finished
	go func() {
		wg.Wait()
		close(results)
	}()
	return &queryBatchIter{results: results, cancel: cancel, wg: wg}
}

type queryBatchIter struct {
	results <-chan *QueryBatch
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
}

func (it *queryBatchIter) Next(ctx context.Context) (*QueryBatch, error) {
	select {
	case <-ctx.Done():
		it.cancel()
		return nil, ctx.Err()
	case batch, ok := <-it.results:
		if !ok {
			return nil, io.EOF
		}
		return batch, nil
	}
}

func (it *queryBatchIter) Close() error {
	it.cancel()
	// Wait for goroutines to finish (best-effort)
	it.wg.Wait()
	return nil
}

func runReadThreadIterator(queryBatcher *QueryBatcher, forest util.ForestInfo, results chan<- *QueryBatch, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	batchSize := int(queryBatcher.BatchSize())
	after := ""
	for {
		// short-circuit if context cancelled
		select {
		case <-ctx.Done():
			return
		default:
		}
		urisHandle := &util.URIsHandle{}
		urisHandle.SetTimestamp(queryBatcher.timestamp)
		client := queryBatcher.clientsByHost[forest.PreferredHost()]
		util.GetURIs(client, queryBatcher.query, forest.Name, queryBatcher.transaction, 0, after, uint(batchSize), urisHandle)
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
			client:    client,
			URIs:      uris,
			timestamp: queryBatcher.timestamp,
		}
		select {
		case <-ctx.Done():
			return
		case results <- queryBatch:
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
