package util

import (
	"bytes"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
)

// Transaction represents a transaction for a request
type Transaction struct {
	client    *clients.Client
	Name      string
	ID        string
	TimeLimit int
	Database  string
}

// TransactionStatusHandle is a handle that places the results into
// a TransactionStatus struct
type TransactionStatusHandle struct {
	*bytes.Buffer
	TransactionStatus TransactionStatus
	timestamp         string
}

// TransactionStatus represents a status of a transaction
type TransactionStatus struct {
	Host                 Host     `json:"host"`
	Server               Server   `json:"server"`
	Database             Database `json:"database"`
	TransactionID        string   `json:"transaction-id"`
	TransactionName      string   `json:"transaction-name"`
	TransactionMode      string   `json:"transaction-mode"`
	TransactionTimestamp string   `json:"transaction-timestamp"`
	TransactionState     string   `json:"transaction-state"`
	Canceled             string   `json:"canceled"`
	StartTime            string   `json:"start-time"`
	TimeLimit            string   `json:"time-limit"`
	MaxTimeLimit         string   `json:"max-time-limit"`
	User                 string   `json:"user"`
	Admin                string   `json:"admin"`
}

// Host represents a host in the MarkLogic cluster
type Host struct {
	ID   string `json:"host-id"`
	Name string `json:"host-name"`
}

// Server represents a server in the MarkLogic cluster
type Server struct {
	ID   string `json:"server-id"`
	Name string `json:"server-name"`
}

// Database represents a database in the MarkLogic cluster
type Database struct {
	ID   string `json:"database-id"`
	Name string `json:"database-name"`
}

// GetStatus returns the status of a Transaction
func (t *Transaction) GetStatus() TransactionStatus {
	transactionStatusHandle := &TransactionStatusHandle{}
	params := "?format=json"
	params = AddDatabaseParam(params, t.client)
	req, _ := BuildRequestFromHandle(t.client, "GET", "/v1/transactions/"+t.ID+params, nil)
	Execute(t.client, req, transactionStatusHandle)
	return *transactionStatusHandle.Get()
}

// Begin starts a Transaction
func (t *Transaction) Begin() bool {
	if t.Name == "" {
		t.Name = "go-client-txn-" + strconv.Itoa(rand.IntN(1000000))
	}
	params := "?name=" + t.Name
	if t.TimeLimit != 0 {
		params = params + "&timeLimit=" + strconv.Itoa(t.TimeLimit)
	}
	params = AddDatabaseParam(params, t.client)
	req, err := BuildRequestFromHandle(t.client, "POST", "/transactions"+params, nil)
	if err != nil {
		return false
	}
	clients.ApplyAuth(t.client, req)
	resp, err := t.client.Do(req)
	if err != nil {
		return false
	}
	if resp.Request.URL.Path != req.URL.Path {
		location := resp.Request.URL.Path
		t.ID = location[strings.LastIndex(location, "/")+1:]
		return true
	} else if resp.StatusCode >= 400 {
		return false
	}
	location := resp.Header.Get("Location")
	t.ID = location[strings.LastIndex(location, "/")+1:]
	return true
}

// Commit commits a Transaction
func (t *Transaction) Commit() bool {
	return actOnTransaction(t, "commit")
}

// Rollback rolls back a Transaction
func (t *Transaction) Rollback() bool {
	return actOnTransaction(t, "rollback")
}

func actOnTransaction(t *Transaction, action string) bool {
	params := "?result=" + action
	params = AddDatabaseParam(params, t.client)
	req, err := BuildRequestFromHandle(t.client, "POST", "/v1/transactions/"+t.ID+params, nil)
	if err != nil {
		return false
	}
	return Execute(t.client, req, nil) == nil
}

// GetFormat returns int that represents JSON
func (tsh *TransactionStatusHandle) GetFormat() int {
	return handle.JSON
}

func (tsh *TransactionStatusHandle) resetBuffer() {
	if tsh.Buffer == nil {
		tsh.Buffer = new(bytes.Buffer)
	}
	tsh.Reset()
}

// Deserialize returns Query struct that represents XML or JSON
func (tsh *TransactionStatusHandle) Deserialize(bytes []byte) {
	tsh.resetBuffer()
	tsh.Write(bytes)
	tsh.TransactionStatus = TransactionStatus{}
	json.Unmarshal(bytes, &tsh.TransactionStatus)
}

// Deserialized returns interface{}
func (tsh *TransactionStatusHandle) Deserialized() interface{} {
	return &tsh.TransactionStatus
}

// Serialize returns []byte of XML or JSON that represents the Query struct
func (tsh *TransactionStatusHandle) Serialize(transactionStatus interface{}) {
	tsh.TransactionStatus = transactionStatus.(TransactionStatus)
	tsh.resetBuffer()
	enc := json.NewEncoder(tsh)
	enc.Encode(tsh.TransactionStatus)
}

// Read bytes
func (tsh *TransactionStatusHandle) Read(bytes []byte) (n int, err error) {
	if tsh.Buffer == nil {
		tsh.Serialize(tsh.TransactionStatus)
	}
	return tsh.Buffer.Read(bytes)
}

// Get returns string of XML or JSON
func (tsh *TransactionStatusHandle) Get() *TransactionStatus {
	return &tsh.TransactionStatus
}

// Serialized returns string of XML or JSON
func (tsh *TransactionStatusHandle) Serialized() string {
	tsh.Serialize(tsh.TransactionStatus)
	return tsh.String()
}

// SetTimestamp sets the timestamp
func (tsh *TransactionStatusHandle) SetTimestamp(timestamp string) {
	tsh.timestamp = timestamp
}

// Timestamp retieves a timestamp
func (tsh *TransactionStatusHandle) Timestamp() string {
	return tsh.timestamp
}

// AcceptResponse handles an *http.Response
func (tsh *TransactionStatusHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(tsh, resp)
}
