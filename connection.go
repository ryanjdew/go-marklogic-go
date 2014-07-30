package goMarklogicGo

import (
	digestAuth "github.com/ryanjdew/http-digest-auth-client"
	"net/http"
	"net/url"
	"strconv"
)

// Authentication options
const (
	BasicAuth = iota
	DigestAuth
	None
)

// Client is used for connecting to the MarkLogic REST API.
type Client struct {
	Base          string
	Userinfo      *url.Userinfo
	AuthType      int
	HTTPClient    *http.Client
	DigestHeaders *digestAuth.DigestHeaders
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
			Base:          base,
			Userinfo:      url.UserPassword(username, password),
			AuthType:      authType,
			HTTPClient:    httpClient,
			DigestHeaders: digestHeaders,
		}
	}
	return client, err
}

// applyAuth adds the neccessary headers for authentication
func applyAuth(c *Client, req *http.Request) {
	pwd, _ := c.Userinfo.Password()
	if c.AuthType == BasicAuth {
		req.SetBasicAuth(c.Userinfo.Username(), pwd)
	} else if c.AuthType == DigestAuth {
		c.DigestHeaders.ApplyAuth(req)
	}
}
