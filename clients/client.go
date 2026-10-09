// Package clients provides clients for working with MarkLogic.
package clients

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	digestAuth "github.com/ryanjdew/http-digest-auth-client"
)

// Authentication options
const (
	BasicAuth = iota
	DigestAuth
	None
	DigestBasicAuth
	OAuthAuth
)

var digestLock *sync.RWMutex = &sync.RWMutex{}

// Connection contains the information needed for a proper MarkLogic connection
type Connection struct {
	Host               string
	Port               int64
	Username           string
	Password           string
	AuthenticationType int
	Database           string
	Protocol           string
	TLSConfig          *tls.Config
	HTTPClient         *http.Client
	BearerToken        string
}

// Client is used for connecting to the MarkLogic REST API.
type Client struct {
	*BasicClient
}

// ClientBuilder is a factory for MarkLogic clients
func ClientBuilder(connection *Connection, base string) (*BasicClient, error) {
	httpClient := connection.HTTPClient
	if httpClient == nil {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		if connection.TLSConfig != nil {
			transport.TLSClientConfig = connection.TLSConfig.Clone()
		}
		httpClient = &http.Client{Transport: transport}
	}
	var basicClient *BasicClient
	var digestHeaders *digestAuth.DigestHeaders
	var err error
	if connection.AuthenticationType == DigestAuth || connection.AuthenticationType == DigestBasicAuth {
		digestHeaders = &digestAuth.DigestHeaders{Client: httpClient}
		digestHeaders, err = digestHeaders.Auth(connection.Username, connection.Password, base+"/config/resources?format=xml")
	}
	if err == nil {
		basicClient =
			&BasicClient{
				base:           base,
				userinfo:       url.UserPassword(connection.Username, connection.Password),
				authType:       connection.AuthenticationType,
				httpClient:     httpClient,
				digestHeaders:  digestHeaders,
				database:       connection.Database,
				bearerToken:    connection.BearerToken,
				connectionInfo: connection,
			}
	}
	return basicClient, err
}

// NewClient creates the Client struct used for searching, etc.
func NewClient(connection *Connection /*host string, port int64, username string, password string, authType int, database string*/) (*Client, error) {
	var client *Client
	base, err := connectionBase(connection, "/LATEST")
	if err != nil {
		return nil, err
	}
	basicClient, err := ClientBuilder(connection, base)
	if err == nil {
		client = &Client{basicClient}
	}
	return client, err
}

// RESTClient is an inteface the different REST Clients (Client and ManagementClient)
type RESTClient interface {
	Base() string
	Userinfo() *url.Userinfo
	AuthType() int
	HTTPClient() *http.Client
	Do(*http.Request) (*http.Response, error)
	DigestHeaders() *digestAuth.DigestHeaders
}

// BasicClient is the basic parts that compose both
type BasicClient struct {
	base           string
	userinfo       *url.Userinfo
	authType       int
	httpClient     *http.Client
	digestHeaders  *digestAuth.DigestHeaders
	database       string
	bearerToken    string
	connectionInfo *Connection
}

func connectionBase(connection *Connection, path string) (string, error) {
	scheme := connection.Protocol
	if scheme == "" {
		scheme = "http"
	}
	scheme = strings.ToLower(scheme)
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("unsupported connection protocol %q: expected http or https", connection.Protocol)
	}
	return scheme + "://" + net.JoinHostPort(connection.Host, strconv.FormatInt(connection.Port, 10)) + path, nil
}

// Base provides the base of the REST calls that will be made
func (bc *BasicClient) Base() string {
	return bc.base
}

// SetBase is to only be used for testing purposes.
// It is exported for subpackage test access.
func (bc *BasicClient) SetBase(base string) {
	bc.base = base
}

// Userinfo returns the credentials for the RESTClient
func (bc *BasicClient) Userinfo() *url.Userinfo {
	return bc.userinfo
}

// AuthType returns the configured authentication type.
func (bc *BasicClient) AuthType() int {
	return bc.authType
}

// HTTPClient returns the *http.Client to use to make the REST calls
func (bc *BasicClient) HTTPClient() *http.Client {
	return bc.httpClient
}

// DigestHeaders returns the headers used for DigestAuth and DigestBasicAuth requests.
func (bc *BasicClient) DigestHeaders() *digestAuth.DigestHeaders {
	return bc.digestHeaders
}

// BearerToken returns the OAuth access token for the RESTClient.
func (bc *BasicClient) BearerToken() string {
	return bc.bearerToken
}

// Database returns the database the client is targeting
func (bc *BasicClient) Database() string {
	return bc.database
}

// ConnectionInfo returns the information on the connection
func (bc *BasicClient) ConnectionInfo() *Connection {
	return bc.connectionInfo
}

// Do makes request with HTTP Client
func (bc *BasicClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := bc.HTTPClient().Do(req)
	if bc.AuthType() == DigestAuth || bc.AuthType() == DigestBasicAuth {
		digestLock.Unlock()
	}
	return resp, err
}

// ApplyAuth adds the neccessary headers for authentication
func ApplyAuth(c RESTClient, req *http.Request) {
	pwd, _ := c.Userinfo().Password()
	if c.AuthType() == BasicAuth {
		req.SetBasicAuth(c.Userinfo().Username(), pwd)
	} else if c.AuthType() == DigestAuth || c.AuthType() == DigestBasicAuth {
		digestLock.Lock()
		c.DigestHeaders().ApplyAuth(req)
	} else if c.AuthType() == OAuthAuth {
		if tokenClient, ok := c.(interface{ BearerToken() string }); ok {
			req.Header.Set("Authorization", "Bearer "+tokenClient.BearerToken())
		}
	}
}
