package go_marklogic_go

import (
	digestAuth "github.com/ryanjdew/http-digest-auth-client"
	"net/http"
	"net/url"
	"strconv"
)

const (
	BASIC_AUTH = iota
	DIGEST_AUTH
	NONE
)

// Client is used for connecting to the MarkLogic REST API.
type Client struct {
	Base          string
	Userinfo      *url.Userinfo
	AuthType      int
	HttpClient    *http.Client
	DigestHeaders *digestAuth.DigestHeaders
}

//
func NewClient(host string, port int64, username string, password string, authType int) (*Client, error) {
	base := "http://" + host + ":" + strconv.FormatInt(port, 10) + "/v1"
	var client *Client
	var digestHeaders *digestAuth.DigestHeaders
	var err error
	if authType == DIGEST_AUTH {
		digestHeaders = &digestAuth.DigestHeaders{}
		digestHeaders, err = digestHeaders.Auth(username, password, base+"/config/resources")
	}
	if err == nil {
		client = &Client{
			Base:          base,
			Userinfo:      url.UserPassword(username, password),
			AuthType:      authType,
			HttpClient:    &http.Client{},
			DigestHeaders: digestHeaders,
		}
	}
	return client, err
}

// This function adds the neccessary headers for authentication
func ApplyAuth(c *Client, req *http.Request) {
	pwd, _ := c.Userinfo.Password()
	if c.AuthType == BASIC_AUTH {
		req.SetBasicAuth(c.Userinfo.Username(), pwd)
	} else if c.AuthType == DIGEST_AUTH {
		c.DigestHeaders.ApplyAuth(req)
	}
}
