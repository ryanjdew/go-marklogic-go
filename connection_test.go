package go_marklogic_go

import (
	"testing"
)

func TestConnection(t *testing.T) {
	expectedBase := "http://localhost:8050/v1"
	// Using BASIC_AUTH, so it doesn't start authenitcation
	client, err := NewClient("localhost", 8050, "admin", "admin", BASIC_AUTH)
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base, expectedBase)
	}
}
