package transactions

import (
	"encoding/json"
	"net/http"
	"strings"

	clients "github.com/ryanjdew/go-marklogic-go/clients"
	handle "github.com/ryanjdew/go-marklogic-go/handle"
	util "github.com/ryanjdew/go-marklogic-go/util"
)

// begin initiates a new multi-statement transaction
func begin(c *clients.Client, response handle.ResponseHandle) error {
	req, err := http.NewRequest("POST", c.Base()+"/transactions", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", handle.FormatEnumToMimeType(response.GetFormat()))
	return util.Execute(c, req, response)
}

// commit commits a transaction by ID
func commit(c *clients.Client, txid string, response handle.ResponseHandle) error {
	req, err := http.NewRequest("POST", c.Base()+"/transactions/"+txid, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return util.Execute(c, req, response)
}

// rollback rolls back a transaction by ID
func rollback(c *clients.Client, txid string, response handle.ResponseHandle) error {
	req, err := http.NewRequest("DELETE", c.Base()+"/transactions/"+txid, nil)
	if err != nil {
		return err
	}
	return util.Execute(c, req, response)
}

// status retrieves transaction status by ID
func status(c *clients.Client, txid string, response handle.ResponseHandle) error {
	req, err := http.NewRequest("GET", c.Base()+"/transactions/"+txid, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", handle.FormatEnumToMimeType(response.GetFormat()))
	return util.Execute(c, req, response)
}

// TransactionInfo represents transaction metadata from the API
type TransactionInfo struct {
	TxID      string `json:"txid"`
	Status    string `json:"status"`
	StartTime string `json:"start-time"`
}

// TransactionHandle deserializes transaction responses
type TransactionHandle struct {
	Format int
	info   TransactionInfo
	buffer []byte
}

func (h *TransactionHandle) GetFormat() int { return h.Format }

func (h *TransactionHandle) Serialize(v any) {
	h.info = v.(TransactionInfo)
}

func (h *TransactionHandle) Deserialize(b []byte) {
	h.buffer = b
	h.info = TransactionInfo{}
	if h.GetFormat() == handle.JSON {
		json.Unmarshal(b, &h.info)
	}
}

func (h *TransactionHandle) Deserialized() any { return h.info }

func (h *TransactionHandle) Serialized() string { return string(h.buffer) }

func (h *TransactionHandle) SetTimestamp(ts string) {}

func (h *TransactionHandle) Timestamp() string { return "" }

func (h *TransactionHandle) Read(p []byte) (int, error) {
	copy(p, h.buffer)
	return len(h.buffer), nil
}

func (h *TransactionHandle) Write(p []byte) (int, error) {
	h.buffer = append(h.buffer, p...)
	return len(p), nil
}

func (h *TransactionHandle) AcceptResponse(resp *http.Response) error {
	return handle.CommonHandleAcceptResponse(h, resp)
}
