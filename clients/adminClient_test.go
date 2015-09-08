package clients

import "testing"

func TestAdminConnection(t *testing.T) {
	expectedBase := "http://marklogic1:28001/admin/v1"
	client, err := NewAdminClient(&Connection{Host: "marklogic1", Port: 28001, Username: "admin", Password: "admin", AuthenticationType: BasicAuth})
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}

func TestAdminConnectionDefaultPort(t *testing.T) {
	expectedBase := "http://localhost:8001/admin/v1"
	client, err := NewAdminClient(&Connection{Host: "localhost", Username: "admin", Password: "admin", AuthenticationType: BasicAuth})
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}
