package clients

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestConnection(t *testing.T) {
	expectedBase := "http://localhost:8000/LATEST"
	// Using BASIC_AUTH, so it doesn't start authenitcation
	// client, err := NewClient("localhost", 8000, "admin", "admin", BasicAuth)
	client, err := NewClient(&Connection{Host: "localhost", Port: 8000, Username: "admin", Password: "admin", AuthenticationType: BasicAuth})
	if err != nil {
		t.Errorf("Error = %v", err)
	} else if client.Base() != expectedBase {
		t.Errorf("Result = %v, want %v", client.Base(), expectedBase)
	}
}

func TestHTTPSOAuthConnection(t *testing.T) {
	const token = "test-access-token"
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer "+token {
			t.Errorf("Authorization = %q, want Bearer token", got)
		}
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	port := server.Listener.Addr().(*net.TCPAddr).Port
	roots := x509.NewCertPool()
	roots.AddCert(server.Certificate())
	client, err := NewClient(&Connection{
		Host:               "127.0.0.1",
		Port:               int64(port),
		Protocol:           "https",
		TLSConfig:          &tls.Config{RootCAs: roots},
		AuthenticationType: OAuthAuth,
		BearerToken:        token,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := client.Base(), "https://127.0.0.1:"+strconv.Itoa(port)+"/LATEST"; got != want {
		t.Fatalf("Base() = %q, want %q", got, want)
	}

	request, err := http.NewRequest(http.MethodGet, client.Base()+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	ApplyAuth(client, request)
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
}

func TestHTTPSDigestConnection(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.HasPrefix(request.Header.Get("Authorization"), "Digest ") {
			writer.Header().Set("WWW-Authenticate", `Digest realm="test", nonce="test-nonce", qop="auth,auth-int", algorithm=MD5`)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	port := server.Listener.Addr().(*net.TCPAddr).Port
	roots := x509.NewCertPool()
	roots.AddCert(server.Certificate())
	for _, authType := range []int{DigestAuth, DigestBasicAuth} {
		client, err := NewClient(&Connection{
			Host:               "127.0.0.1",
			Port:               int64(port),
			Protocol:           "https",
			TLSConfig:          &tls.Config{RootCAs: roots},
			Username:           "test-user",
			Password:           "test-password",
			AuthenticationType: authType,
		})
		if err != nil {
			t.Fatal(err)
		}
		if got, want := client.Base(), "https://127.0.0.1:"+strconv.Itoa(port)+"/LATEST"; got != want {
			t.Fatalf("Base() = %q, want %q", got, want)
		}

		request, err := http.NewRequest(http.MethodGet, client.Base()+"/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		ApplyAuth(client, request)
		response, err := client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
	}
}
