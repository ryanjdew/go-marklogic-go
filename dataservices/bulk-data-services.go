package dataservices

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/util"
)

// WorkPhase describes the state of the Bulk Data Service
type WorkPhase int

const (
	// INITIALIZING BulkDataService has not run called the Data Service endpoint yet
	INITIALIZING WorkPhase = iota
	// RUNNING BulkDataService is calling the Data Service endpoint
	RUNNING
	// INTERRUPTING request to interupt the BulkDataService has been made
	INTERRUPTING
	// INTERRUPTED BulkDataService has been interupted
	INTERRUPTED
	// COMPLETED BulkDataService has finished
	COMPLETED
)

// BulkDataService reads documents in bulk
type BulkDataService struct {
	endpoint             string
	workPhase            WorkPhase
	workUnits            []any
	batchSize            uint16
	threadCount          uint8
	mutex                *sync.Mutex
	client               *clients.Client
	clientsByHost        map[string]*clients.Client
	statusChangeListener []chan<- *WorkPhase
	forestInfo           []util.ForestInfo
	inputChannel         <-chan *handle.Handle
	outputListeners      []chan<- []byte
	waitGroup            *sync.WaitGroup
	workIsForestBased    bool
	endpointState        []byte
}

// WithOutputListener adds a listener to the output from BulkDataServices
func (bds *BulkDataService) WithOutputListener(listener chan []byte) *BulkDataService {
	bds.outputListeners = append(bds.outputListeners, listener)
	return bds
}

// WithInputChannel adds a channel to feed input to BulkDataServices
func (bds *BulkDataService) WithInputChannel(inputChannel <-chan *handle.Handle) *BulkDataService {
	bds.inputChannel = inputChannel
	return bds
}

// WithWorkUnits sets up work units to send to BulkDataServices
func (bds *BulkDataService) WithWorkUnits(workUnits ...interface{}) *BulkDataService {
	bds.workUnits = workUnits
	bds.workIsForestBased = false
	return bds
}

// WithForestBasedWorkUnits sets work units based off of forests in the database
func (bds *BulkDataService) WithForestBasedWorkUnits() *BulkDataService {
	forestBasedWorkUnits := make([]any, 0, len(bds.forestInfo))
	for _, forestInfo := range bds.forestInfo {
		workUnit := map[string]string{}
		workUnit["forestId"] = forestInfo.ID
		forestBasedWorkUnits = append(forestBasedWorkUnits, workUnit)
	}
	bds.workUnits = forestBasedWorkUnits
	bds.workIsForestBased = true
	return bds
}

// WithEndpointState sets an intial Endpoint State for the Data Service
func (bds *BulkDataService) WithEndpointState(endpointState []byte) *BulkDataService {
	bds.endpointState = endpointState
	return bds
}

// BatchSize is the number documents we'll create in a single batch
func (bds *BulkDataService) BatchSize() uint16 {
	return bds.batchSize
}

// WithBatchSize sets the number documents we'll create in a single batch
func (bds *BulkDataService) WithBatchSize(batchSize uint16) *BulkDataService {
	bds.batchSize = batchSize
	return bds
}

// WithThreadCount set the thread count
func (bds *BulkDataService) WithThreadCount(threadCount uint8) *BulkDataService {
	bds.threadCount = threadCount
	return bds
}

// Cancel interupts the service
func (bds *BulkDataService) Cancel() *BulkDataService {
	bds.workPhase = INTERRUPTING
	return bds.Wait()
}

// ThreadCount is the number threads we'll create
func (bds *BulkDataService) ThreadCount() uint8 {
	return bds.threadCount
}

// DataServiceBatch submits documents in bulk
type DataServiceBatch struct {
	endpoint      string
	endpointState []byte
	input         []*handle.Handle
}

// Input contains the data to send to the Data Service
func (dsb *DataServiceBatch) Input() []*handle.Handle {
	return dsb.input
}

// Run the BulkDataService
func (bds *BulkDataService) Run() *BulkDataService {
	hosts := make([]string, 0, len(bds.clientsByHost))
	for host := range bds.clientsByHost {
		hosts = append(hosts, host)
	}
	forestLength := len(bds.forestInfo)
	hostLength := len(hosts)
	threadCount := int(bds.ThreadCount())
	roundRobinCounter := 0
	var roundRobinLength int
	if bds.workIsForestBased {
		roundRobinLength = forestLength
		threadCount = forestLength
	} else {
		roundRobinLength = hostLength
	}
	workUnitRoundRobinCounter := 0
	workUnitRoundRobinLength := len(bds.workUnits)
	bds.workPhase = RUNNING
	for i := 0; i < threadCount; i++ {
		bds.waitGroup.Add(1)
		var selectedHost string
		if bds.workIsForestBased {
			selectedHost = bds.forestInfo[roundRobinCounter].PreferredHost()
		} else {
			selectedHost = hosts[roundRobinCounter]
		}
		client := bds.clientsByHost[selectedHost]
		var currentWorkUnit *interface{} = nil
		if workUnitRoundRobinLength > 0 {
			currentWorkUnit = &bds.workUnits[workUnitRoundRobinCounter]
			workUnitRoundRobinCounter = (workUnitRoundRobinCounter + 1) % workUnitRoundRobinLength
		}
		if bds.inputChannel != nil {
			go runInputThread(bds, currentWorkUnit, bds.inputChannel, client)
		} else {
			go runProcessThread(bds, currentWorkUnit, client)
		}
		roundRobinCounter = (roundRobinCounter + 1) % roundRobinLength
	}
	return bds
}

// Wait on the BulkDataService to finish
func (bds *BulkDataService) Wait() *BulkDataService {
	bds.waitGroup.Wait()
	if bds.workPhase == INTERRUPTING {
		bds.workPhase = INTERRUPTED
	} else {
		bds.workPhase = COMPLETED
	}
	return bds
}

func runInputThread(bds *BulkDataService, workUnit *interface{}, inputChannel <-chan *handle.Handle, client *clients.Client) {
	trackEndpointState := bds.endpointState != nil && len(bds.endpointState) > 0
	listeners := bds.outputListeners
	batchSizeInt := int(bds.BatchSize())
	wg := bds.waitGroup
	defer wg.Done()
	inputBatch := &DataServiceBatch{
		endpoint:      bds.endpoint,
		input:         make([]*handle.Handle, 0, batchSizeInt),
		endpointState: bds.endpointState,
	}
	for input := range inputChannel {
		if input != nil {
			if bds.workPhase == INTERRUPTING {
				return
			}
			inputBatch.input = append(inputBatch.input, input)
			if len(inputBatch.input) >= batchSizeInt {
				submitDataServiceBatch(inputBatch, workUnit, listeners, client)
				inputBatch.input = make([]*handle.Handle, 0, batchSizeInt)
				if trackEndpointState && len(inputBatch.endpointState) == 0 {
					return
				}
			}
		} else {
			time.Sleep(time.Nanosecond)
		}
	}
	bds.workPhase = COMPLETED
	if len(inputBatch.input) > 0 {
		submitDataServiceBatch(inputBatch, workUnit, listeners, client)
		inputBatch.input = make([]*handle.Handle, 0, batchSizeInt)
	}
}

func runProcessThread(bds *BulkDataService, workUnit *interface{}, client *clients.Client) {
	listeners := bds.outputListeners
	wg := bds.waitGroup
	defer wg.Done()
	batch := &DataServiceBatch{
		endpoint:      bds.endpoint,
		endpointState: bds.endpointState,
	}
	for {
		submitDataServiceBatch(batch, workUnit, listeners, client)
		if len(batch.endpointState) == 0 || bds.workPhase == INTERRUPTING {
			return
		}
	}
}

func submitDataServiceBatch(dataServiceBatch *DataServiceBatch, workUnit *interface{}, listeners []chan<- []byte, client *clients.Client) error {
	unatomicParams := map[string][]*handle.Handle{}
	if workUnit != nil {
		jsonBytes, err := json.Marshal(workUnit)
		if err != nil {
			log.Fatal(err)
		}
		var workUnitHandle handle.Handle = &handle.RawHandle{Format: handle.JSON}
		workUnitHandle.Deserialize(jsonBytes)
		unatomicParams["workUnit"] = []*handle.Handle{&workUnitHandle}
	}
	trackEndpointState := dataServiceBatch.endpointState != nil && len(dataServiceBatch.endpointState) > 0
	if trackEndpointState {
		var endpointStateHandle handle.Handle = &handle.RawHandle{Format: handle.JSON}
		endpointStateHandle.Deserialize(dataServiceBatch.endpointState)
		unatomicParams["endpointState"] = []*handle.Handle{&endpointStateHandle}
	}
	if len(dataServiceBatch.input) > 0 {
		unatomicParams["input"] = dataServiceBatch.input
	}
	respHandle := &handle.MultipartResponseHandle{}
	err := util.PostForm(client, dataServiceBatch.endpoint, make(map[string][]string), unatomicParams, respHandle, true)
	multipartOutput := respHandle.Get()
	if len(multipartOutput) == 0 {
		if trackEndpointState {
			dataServiceBatch.endpointState = []byte("")
		}
		return nil
	}
	for index, val := range multipartOutput {
		if trackEndpointState && index == 0 {
			dataServiceBatch.endpointState = val
		} else {
			for _, listener := range listeners {
				listener <- val
			}
		}
	}
	return err
}
