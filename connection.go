package goMarklogicGo

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

// Client is used for connecting to the MarkLogic REST API.
type Client struct {
	*BasicClient
}

// NewClient creates the Client struct used for searching, etc.
func NewClient(host string, port int64, username string, password string, authType int) (*Client, error) {
	base := "http://" + host + ":" + strconv.FormatInt(port, 10) + "/v1"
	httpClient := &http.Client{}
	var client *Client
	var digestHeaders *digestAuth.DigestHeaders
	var err error
	if authType == DigestAuth {
		digestHeaders = &digestAuth.DigestHeaders{}
		digestHeaders, err = digestHeaders.Auth(username, password, base+"/config/resources?format=xml")
	}
	if err == nil {
		client = &Client{
			&BasicClient{
				base:          base,
				userinfo:      url.UserPassword(username, password),
				authType:      authType,
				httpClient:    httpClient,
				digestHeaders: digestHeaders,
			},
		}
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

func (bc *BasicClient) Base() string {
	return bc.base
}

func (bc *BasicClient) SetBase(base string) {
	bc.base = base
}

func (bc *BasicClient) Userinfo() *url.Userinfo {
	return bc.userinfo
}

func (bc *BasicClient) AuthType() int {
	return bc.authType
}

func (bc *BasicClient) HTTPClient() *http.Client {
	return bc.httpClient
}

func (bc *BasicClient) DigestHeaders() *digestAuth.DigestHeaders {
	return bc.digestHeaders
}

// applyAuth adds the neccessary headers for authentication
func applyAuth(c RESTClient, req *http.Request) {
	pwd, _ := c.Userinfo().Password()
	if c.AuthType() == BasicAuth {
		req.SetBasicAuth(c.Userinfo().Username(), pwd)
	} else if c.AuthType() == DigestAuth {
		c.DigestHeaders().ApplyAuth(req)
	}
}
