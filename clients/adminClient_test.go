package clients

import "testing"

func TestAdminConnection(t *testing.T) {
	expectedBase := "http://localhost:8001/admin/v1"
	client, err := NewAdminClient("localhost", "admin", "admin", BasicAuth)
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}
