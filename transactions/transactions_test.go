package transactions

import (
	"testing"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
	"github.com/ryanjdew/go-marklogic-go/test"
)

var exampleBeginResponse = `{
  "txid": "3251334055975682123",
  "status": "active",
  "start-time": "2014-03-20T14:23:56.639Z"
}`

func TestBeginTransaction(t *testing.T) {
	client, server := test.Client(exampleBeginResponse)
	defer server.Close()

	svc := NewService(client)
	resp := &TransactionHandle{Format: handle.JSON}
	err := svc.Begin(resp)
	if err != nil {
		t.Fatalf("Begin failed: %v", err)
	}

	info := resp.Deserialized().(TransactionInfo)
	if info.TxID != "3251334055975682123" {
		t.Errorf("Expected txid 3251334055975682123, got %s", info.TxID)
	}
}

func TestCommitTransaction(t *testing.T) {
	client, server := test.Client(`{"status":"committed"}`)
	defer server.Close()

	svc := NewService(client)
	resp := &TransactionHandle{Format: handle.JSON}
	err := svc.Commit("3251334055975682123", resp)
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
}

func TestRollbackTransaction(t *testing.T) {
	client, server := test.Client(`{"status":"rolled-back"}`)
	defer server.Close()

	svc := NewService(client)
	resp := &TransactionHandle{Format: handle.JSON}
	err := svc.Rollback("3251334055975682123", resp)
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
}

func TestStatusTransaction(t *testing.T) {
	client, server := test.Client(exampleBeginResponse)
	defer server.Close()

	svc := NewService(client)
	resp := &TransactionHandle{Format: handle.JSON}
	err := svc.Status("3251334055975682123", resp)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	info := resp.Deserialized().(TransactionInfo)
	if info.Status != "active" {
		t.Errorf("Expected status active, got %s", info.Status)
	}
}
