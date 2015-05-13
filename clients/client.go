package clients

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	handle "github.com/ryanjdew/go-marklogic-go/handle"
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

// applyAuth adds the neccessary headers for authentication
func applyAuth(c RESTClient, req *http.Request) {
	pwd, _ := c.Userinfo().Password()
	if c.AuthType() == BasicAuth {
		req.SetBasicAuth(c.Userinfo().Username(), pwd)
	} else if c.AuthType() == DigestAuth {
		c.DigestHeaders().ApplyAuth(req)
	}
}

// Execute uses a client to run a request and places the results in the
// response Handle
func Execute(c RESTClient, req *http.Request, responseHandle handle.Handle) error {
	applyAuth(c, req)
	respType := handle.FormatEnumToString(responseHandle.GetFormat())
	req.Header.Add("Accept", "application/"+respType)
	resp, err := c.HTTPClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	responseHandle.Encode(contents)
	return err
}
