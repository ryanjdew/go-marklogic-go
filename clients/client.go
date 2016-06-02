package clients

import (
	"net/http"
	"net/url"
	"strconv"

	digestAuth "github.com/ryanjdew/http-digest-auth-client"
)

// Authentication options
const (
	BasicAuth = iota
	DigestAuth
	None
)

// Connection contains the information needed for a proper MarkLogic connection
type Connection struct {
	Host               string
	Port               int64
	Username           string
	Password           string
	AuthenticationType int
}

// Client is used for connecting to the MarkLogic REST API.
type Client struct {
	*BasicClient
}

//ClientBuilder is a factory for MarkLogic clients
func ClientBuilder(connection *Connection, base string) (*BasicClient, error) {
	httpClient := &http.Client{}
	var basicClient *BasicClient
	var digestHeaders *digestAuth.DigestHeaders
	var err error
	if connection.AuthenticationType == DigestAuth {
		digestHeaders = &digestAuth.DigestHeaders{}
		digestHeaders, err = digestHeaders.Auth(connection.Username, connection.Password, base+"/config/resources?format=xml")
	}
	if err == nil {
		basicClient =
			&BasicClient{
				base:          base,
				userinfo:      url.UserPassword(connection.Username, connection.Password),
				authType:      connection.AuthenticationType,
				httpClient:    httpClient,
				digestHeaders: digestHeaders,
			}
	}
	return basicClient, err
}

// NewClient creates the Client struct used for searching, etc.
func NewClient(connection *Connection /*host string, port int64, username string, password string, authType int*/) (*Client, error) {
	var client *Client
	base := "http://" + connection.Host + ":" + strconv.FormatInt(connection.Port, 10) + "/v1"
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
	DigestHeaders() *digestAuth.DigestHeaders
}

// BasicClient is the basic parts that compose both
type BasicClient struct {
	base          string
	userinfo      *url.Userinfo
	authType      int
	httpClient    *http.Client
	digestHeaders *digestAuth.DigestHeaders
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

// AuthType returns the int that represents an authentication type (BasicAuth, DigestAuth)
func (bc *BasicClient) AuthType() int {
	return bc.authType
}

// HTTPClient returns the *http.Client to use to make the REST calls
func (bc *BasicClient) HTTPClient() *http.Client {
	return bc.httpClient
}

// DigestHeaders returns the digest headers that need updated with each DigestAuth call
func (bc *BasicClient) DigestHeaders() *digestAuth.DigestHeaders {
	return bc.digestHeaders
}

// ApplyAuth adds the neccessary headers for authentication
func ApplyAuth(c RESTClient, req *http.Request) {
	pwd, _ := c.Userinfo().Password()
	if c.AuthType() == BasicAuth {
		req.SetBasicAuth(c.Userinfo().Username(), pwd)
	} else if c.AuthType() == DigestAuth {
		c.DigestHeaders().ApplyAuth(req)
	}
}
