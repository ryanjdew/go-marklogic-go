package datamovement

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// WriteBatcher writes documents in bulk
type WriteBatcher struct {
	client                 *clients.Client
	documentsServiceByHost map[string]*documents.Service
	clientsByHost          map[string]*clients.Client
	batchSize              uint16
	threadCount            uint8
	writeChannel           <-chan *documents.DocumentDescription
	timestamp              string
	listeners              []chan<- *WriteBatch
	waitGroup              *sync.WaitGroup
	forestInfo             []util.ForestInfo
	transform              *util.Transform
	transaction            *util.Transaction
}

// WriteBatchIterator provides a pull-style iterator for WriteBatch results
type WriteBatchIterator interface {
	Next(ctx context.Context) (*WriteBatch, error)
	Close() error
}

// BatchSize is the number documents we'll create in a single batch
func (wbr *WriteBatcher) BatchSize() uint16 {
	return wbr.batchSize
}

// ThreadCount is the number threads we'll create
func (wbr *WriteBatcher) ThreadCount() uint8 {
	return wbr.threadCount
}

// WriteChannel is the channel that documents to write will be pulled from
func (wbr *WriteBatcher) WriteChannel() <-chan *documents.DocumentDescription {
	return wbr.writeChannel
}

// Timestamp for the write operation
func (wbr *WriteBatcher) Timestamp() string {
	return wbr.timestamp
}

// WithThreadCount set the thread count
func (wbr *WriteBatcher) WithThreadCount(threadCount uint8) *WriteBatcher {
	wbr.threadCount = threadCount
	return wbr
}

// WithBatchSize set the batch size
func (wbr *WriteBatcher) WithBatchSize(batchSize uint16) *WriteBatcher {
	wbr.batchSize = batchSize
	return wbr
}

// WithTransform transform to apply to the documents
func (wbr *WriteBatcher) WithTransform(transform *util.Transform) *WriteBatcher {
	wbr.transform = transform
	return wbr
}

// WithTransaction perform writes in given transaction
func (wbr *WriteBatcher) WithTransaction(transaction *util.Transaction) *WriteBatcher {
	wbr.transaction = transaction
	return wbr
}

// WithWriteChannel add a channel od documents to write
func (wbr *WriteBatcher) WithWriteChannel(writeChannel <-chan *documents.DocumentDescription) *WriteBatcher {
	wbr.writeChannel = writeChannel
	return wbr
}

// WithListener add a listener channel
func (wbr *WriteBatcher) WithListener(listener chan<- *WriteBatch) *WriteBatcher {
	wbr.listeners = append(wbr.listeners, listener)
	return wbr
}

// RemoveListener remove the a listener
func (wbr *WriteBatcher) RemoveListener(listener chan<- *WriteBatch) *WriteBatcher {
	for i, compareListener := range wbr.listeners {
		if compareListener == listener {
			copy(wbr.listeners[i:], wbr.listeners[i+1:])
			break
		}
	}
	return wbr
}

// Run the WriteBatcher
func (wbr *WriteBatcher) Run() *WriteBatcher {
	hosts := make([]string, len(wbr.documentsServiceByHost))
	for host := range wbr.documentsServiceByHost {
		hosts = append(hosts, host)
	}
	threadCount := int(wbr.ThreadCount())
	wbr.waitGroup = &sync.WaitGroup{}
	forestLength := len(wbr.forestInfo)
	hostLength := len(hosts)
	distributeByForest := forestLength > 0 && threadCount >= forestLength
	// debug logs removed
	roundRobinCounter := 0
	var roundRobinLength int
	if distributeByForest {
		roundRobinLength = forestLength
	} else {
		roundRobinLength = hostLength
	}
	for range threadCount {
		wbr.waitGroup.Add(1)
		var selectedHost string
		if distributeByForest {
			selectedHost = wbr.forestInfo[roundRobinCounter].PreferredHost()
		} else {
			selectedHost = hosts[roundRobinCounter]
		}
		documentsService := wbr.documentsServiceByHost[selectedHost]
		go runWriteThread(wbr, wbr.WriteChannel(), documentsService)
		roundRobinCounter = (roundRobinCounter + 1) % roundRobinLength
	}
	return wbr
}

// Wait on the WriteBatcher to finish
func (wbr *WriteBatcher) Wait() *WriteBatcher {
	wbr.waitGroup.Wait()
	return wbr
}

func runWriteThread(writeBatcher *WriteBatcher, writeChannel <-chan *documents.DocumentDescription, documentsService *documents.Service) {
	listeners := writeBatcher.listeners
	batchSizeInt := int(writeBatcher.BatchSize())
	wg := writeBatcher.waitGroup
	defer wg.Done()
	var writeBatch *WriteBatch
	for {
		if writeBatch == nil {
			writeBatch = &WriteBatch{
				documentsService:     documentsService,
				documentDescriptions: make([]*documents.DocumentDescription, 0, batchSizeInt),
			}
		}
		writeDoc, ok := <-writeChannel
		if writeDoc != nil {
			writeBatch.documentDescriptions = append(writeBatch.documentDescriptions, writeDoc)
			if len(writeBatch.documentDescriptions) >= batchSizeInt {
				submitBatch(writeBatch, writeBatcher.transform, writeBatcher.transaction, listeners)
				writeBatch = nil
			}
		} else if !ok && len(writeChannel) == 0 {
			if len(writeBatch.documentDescriptions) > 0 {
				submitBatch(writeBatch, writeBatcher.transform, writeBatcher.transaction, listeners)
				writeBatch = nil
			}
			return
		} else {
			time.Sleep(time.Nanosecond)
		}
	}
}

func runWriteThreadIterator(writeBatcher *WriteBatcher, writeChannel <-chan *documents.DocumentDescription, documentsService *documents.Service, results chan<- *WriteBatch, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	batchSizeInt := int(writeBatcher.BatchSize())
	// ensure writeBatch is initialized before reads to avoid nil deref
	writeBatch := &WriteBatch{
		documentsService:     documentsService,
		documentDescriptions: make([]*documents.DocumentDescription, 0, batchSizeInt),
	}
	for {
		select {
		case <-ctx.Done():
			return
		case writeDoc, ok := <-writeChannel:
			if writeDoc != nil {
				writeBatch.documentDescriptions = append(writeBatch.documentDescriptions, writeDoc)
				if len(writeBatch.documentDescriptions) >= batchSizeInt {
					// submit and forward via results
					ch := make(chan *WriteBatch, 1)
					submitBatch(writeBatch, writeBatcher.transform, writeBatcher.transaction, []chan<- *WriteBatch{ch})
					select {
					case <-ctx.Done():
						return
					case res := <-ch:
						results <- res
					}
					// reset batch
					writeBatch = &WriteBatch{
						documentsService:     documentsService,
						documentDescriptions: make([]*documents.DocumentDescription, 0, batchSizeInt),
					}
				}
			} else if !ok && len(writeChannel) == 0 {
				if len(writeBatch.documentDescriptions) > 0 {
					ch := make(chan *WriteBatch, 1)
					submitBatch(writeBatch, writeBatcher.transform, writeBatcher.transaction, []chan<- *WriteBatch{ch})
					select {
					case <-ctx.Done():
						return
					case res := <-ch:
						results <- res
					}
					writeBatch = nil
				}
				return
			}
		}
	}
}

// Iterator returns an iterator which yields WriteBatch results after performing
// writes. If `WithWriteChannel` wasn't set, the iterator yields EOF immediately.
func (wbr *WriteBatcher) Iterator(ctx context.Context) WriteBatchIterator {
	if wbr.writeChannel == nil || len(wbr.documentsServiceByHost) == 0 {
		ch := make(chan *WriteBatch)
		close(ch)
		return &writeBatchIter{results: ch}
	}
	ctx, cancel := context.WithCancel(ctx)
	results := make(chan *WriteBatch)
	wg := &sync.WaitGroup{}
	threadCount := int(wbr.ThreadCount())
	hosts := make([]string, 0, len(wbr.documentsServiceByHost))
	for host := range wbr.documentsServiceByHost {
		hosts = append(hosts, host)
	}
	// debug logs removed
	forestLength := len(wbr.forestInfo)
	hostLength := len(hosts)
	distributeByForest := forestLength > 0 && threadCount >= forestLength
	roundRobinCounter := 0
	var roundRobinLength int
	if distributeByForest {
		roundRobinLength = forestLength
	} else {
		roundRobinLength = hostLength
	}
	for range threadCount {
		wg.Add(1)
		var selectedHost string
		if distributeByForest {
			selectedHost = wbr.forestInfo[roundRobinCounter].PreferredHost()
		} else {
			selectedHost = hosts[roundRobinCounter]
		}
		documentsService := wbr.documentsServiceByHost[selectedHost]
		go runWriteThreadIterator(wbr, wbr.WriteChannel(), documentsService, results, wg, ctx)
		roundRobinCounter = (roundRobinCounter + 1) % roundRobinLength
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	return &writeBatchIter{results: results, cancel: cancel, wg: wg}
}

type writeBatchIter struct {
	results <-chan *WriteBatch
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
}

func (it *writeBatchIter) Next(ctx context.Context) (*WriteBatch, error) {
	select {
	case <-ctx.Done():
		if it.cancel != nil {
			it.cancel()
		}
		return nil, ctx.Err()
	case batch, ok := <-it.results:
		if !ok {
			return nil, io.EOF
		}
		return batch, nil
	}
}

func (it *writeBatchIter) Close() error {
	if it.cancel != nil {
		it.cancel()
	}
	if it.wg != nil {
		it.wg.Wait()
	}
	return nil
}

func submitBatch(writeBatch *WriteBatch, transform *util.Transform, transaction *util.Transaction, listeners []chan<- *WriteBatch) {
	if len(writeBatch.DocumentDescriptions()) > 0 {
		responseHandle := &handle.RawHandle{}
		writeBatch.DocumentsService().WriteSet(writeBatch.DocumentDescriptions(), &documents.MetadataHandle{}, transform, transaction, responseHandle)
		writeBatch.WithResponse(responseHandle)

		// provide writeBatch back to listeners
		for _, listener := range listeners {
			listener <- writeBatch
		}
	}
}

// WriteBatch batch of DocumentDescriptions to write to MarkLogic and relevant meta information
type WriteBatch struct {
	documentsService     *documents.Service
	documentDescriptions []*documents.DocumentDescription
	timestamp            string
	response             handle.ResponseHandle
}

// DocumentsService used to write the documents
func (wb *WriteBatch) DocumentsService() *documents.Service {
	return wb.documentsService
}

// Client used to write the documents
func (wb *WriteBatch) Client() *clients.Client {
	return wb.documentsService.Client()
}

// DocumentDescriptions representations of the documents written
func (wb *WriteBatch) DocumentDescriptions() []*documents.DocumentDescription {
	return wb.documentDescriptions
}

// Timestamp for the write operation
func (wb *WriteBatch) Timestamp() string {
	return wb.timestamp
}

// WithResponse set write batch response
func (wb *WriteBatch) WithResponse(response handle.ResponseHandle) *WriteBatch {
	wb.response = response
	return wb
}

// Response return the write batch response
func (wb *WriteBatch) Response() handle.ResponseHandle {
	return wb.response
}
