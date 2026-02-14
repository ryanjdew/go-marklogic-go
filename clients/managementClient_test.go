package clients

import "testing"

func TestManagementConnection(t *testing.T) {
	expectedBase := "http://localhost:28002/manage/v2"
	client, err := NewManagementClient(&Connection{Host: "localhost", Port: 28002, Username: "admin", Password: "admin", AuthenticationType: BasicAuth})
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}

func TestManagementConnectionDefaultPort(t *testing.T) {
	expectedBase := "http://localhost:8002/manage/v2"
	client, err := NewManagementClient(&Connection{Host: "localhost", Username: "admin", Password: "admin", AuthenticationType: BasicAuth})
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}
