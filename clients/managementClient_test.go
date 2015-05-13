package clients

import "testing"

func TestManagementConnection(t *testing.T) {
	expectedBase := "http://localhost:8002/manage/v2"
	client, err := NewManagementClient("localhost", "admin", "admin", BasicAuth)
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}
