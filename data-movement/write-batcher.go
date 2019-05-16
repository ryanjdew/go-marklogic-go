package dataMovement

import (
	"sync"

	"github.com/ryanjdew/go-marklogic-go/clients"
	"github.com/ryanjdew/go-marklogic-go/documents"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// WriteBatcher writes documents in bulk
type WriteBatcher struct {
	client           *clients.Client
	documentsService *documents.Service
	clients          []*clients.Client
	batchSize        uint16
	threadCount      uint8
	writeChannel     <-chan *documents.DocumentDescription
	timestamp        string
	listeners        []chan<- *WriteBatch
	waitGroup        *sync.WaitGroup
}

// BatchSize is the number documents we'll create in a single batch
func (wbr *WriteBatcher) BatchSize() uint16 {
	return wbr.batchSize
}

// DocumentsService used to write the documents
func (wbr *WriteBatcher) DocumentsService() *documents.Service {
	return wbr.documentsService
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
	threadCount := int(wbr.ThreadCount())
	wbr.waitGroup = &sync.WaitGroup{}
	for i := 0; i < threadCount; i++ {
		wbr.waitGroup.Add(1)
		go runWriteThread(wbr, wbr.WriteChannel())
	}
	return wbr
}

// Wait on the WriteBatcher to finish
func (wbr *WriteBatcher) Wait() *WriteBatcher {
	wbr.waitGroup.Wait()
	return wbr
}

func runWriteThread(writeBatcher *WriteBatcher, writeChannel <-chan *documents.DocumentDescription) {
	listeners := writeBatcher.listeners
	batchSizeInt := int(writeBatcher.BatchSize())
	wg := writeBatcher.waitGroup
	defer wg.Done()
	var writeBatch *WriteBatch
	for {
		if writeBatch == nil {
			writeBatch = &WriteBatch{
				documentsService:     writeBatcher.DocumentsService(),
				documentDescriptions: make([]*documents.DocumentDescription, 0, batchSizeInt),
			}
		}
		select {
		case writeDoc, ok := <-writeChannel:
			if writeDoc != nil {
				writeBatch.documentDescriptions = append(writeBatch.documentDescriptions, writeDoc)
				if len(writeBatch.documentDescriptions) >= batchSizeInt {
					submitBatch(writeBatch, listeners)
					writeBatch = nil
				}
			} else if !ok && len(writeChannel) == 0 {
				if len(writeBatch.documentDescriptions) > 0 {
					submitBatch(writeBatch, listeners)
					writeBatch = nil
				}
				return
			}
		}
	}
}

func submitBatch(writeBatch *WriteBatch, listeners []chan<- *WriteBatch) {
	if len(writeBatch.DocumentDescriptions()) > 0 {
		responseHandle := &handle.RawHandle{}
		writeBatch.DocumentsService().WriteSet(writeBatch.DocumentDescriptions(), &documents.MetadataHandle{}, nil, responseHandle)
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
