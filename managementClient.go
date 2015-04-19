package goMarklogicGo

import (
	"net/http"
	"net/url"

	digestAuth "github.com/ryanjdew/http-digest-auth-client"
)

// ServerProperties represents the properties of a MarkLogic AppServer
type ServerProperties struct {
	Name string
}

// DatabaseProperties represents the properties of a MarkLogic Database
type DatabaseProperties struct {
}

// NewManagementClient creates the Client struct used for managing databases, etc.
func NewManagementClient(host string, username string, password string, authType int) (*Client, error) {
	base := "http://" + host + ":8002" + "/manage/v2"
	httpClient := &http.Client{}
	var client *Client
	var digestHeaders *digestAuth.DigestHeaders
	var err error
	if authType == DigestAuth {
		digestHeaders = &digestAuth.DigestHeaders{}
		digestHeaders, err = digestHeaders.Auth(username, password, base+"?format=xml")
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
