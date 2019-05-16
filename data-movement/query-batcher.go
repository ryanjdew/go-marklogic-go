package dataMovement

import (
	"sync"

	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// QueryBatcher reads documents in bulk
type QueryBatcher struct {
	client      *clients.Client
	clients     []*clients.Client
	batchSize   uint16
	query       handle.Handle
	threadCount uint8
	timestamp   string
	listeners   []chan<- *QueryBatch
	waitGroup   *sync.WaitGroup
}

// BatchSize is the number documents we'll retrieve in a single batch
func (qbr *QueryBatcher) BatchSize() uint16 {
	return qbr.batchSize
}

// ThreadCount is the number threads we'll create
func (qbr *QueryBatcher) ThreadCount() uint8 {
	return qbr.threadCount
}

// Timestamp for the read operation
func (qbr *QueryBatcher) Timestamp() string {
	return qbr.timestamp
}

// WithThreadCount set the thread count
func (qbr *QueryBatcher) WithThreadCount(threadCount uint8) *QueryBatcher {
	qbr.threadCount = threadCount
	return qbr
}

// WithBatchSize set the batch size
func (qbr *QueryBatcher) WithBatchSize(batchSize uint16) *QueryBatcher {
	qbr.batchSize = batchSize
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

// Run the WriteBatcher
func (qbr *QueryBatcher) Run() *QueryBatcher {
	// TODO Use forestInfo for creating Clients for gathering URIs
	//forestInfoHandle := &handle.MapHandle{Format: handle.JSON}
	//resources.NewService(queryBatcher.client).Get("internal/forestinfo", map[string]string{}, forestInfoHandle)
	//forestInfo := *forestInfoHandle.Get()

	threadCount := int(qbr.ThreadCount())
	qbr.waitGroup = &sync.WaitGroup{}
	for i := 0; i < threadCount; i++ {
		qbr.waitGroup.Add(1)
		go runReadThread(qbr)
	}
	return qbr
}

// Wait on the WriteBatcher to finish
func (qbr *QueryBatcher) Wait() *QueryBatcher {
	qbr.waitGroup.Wait()
	return qbr
}

func runReadThread(queryBatcher *QueryBatcher) {
	listeners := queryBatcher.listeners
	batchSize := int(queryBatcher.BatchSize())
	wg := queryBatcher.waitGroup
	defer wg.Done()
	queryBatch := &QueryBatch{
		URIs: make([]string, 0, batchSize),
	}
	// TODO Call search.Values endpoint when implemented
	// provide queryBatch back to listeners
	for _, listener := range listeners {
		listener <- queryBatch
	}
}

// QueryBatch batch of URIs matching a query and relevant meta information
type QueryBatch struct {
	documentsService *documents.Service
	URIs             []string
	timestamp        string
}

// DocumentsService used to with forest for the documents
func (qb *QueryBatch) DocumentsService() *documents.Service {
	return qb.documentsService
}
