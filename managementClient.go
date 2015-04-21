package goMarklogicGo

import (
	"net/http"
	"net/url"

	digestAuth "github.com/ryanjdew/http-digest-auth-client"
)

// ManagementClient is used for connecting to the MarkLogic Management API.
type ManagementClient struct {
	*BasicClient
}

// NewManagementClient creates the Client struct used for managing databases, etc.
func NewManagementClient(host string, username string, password string, authType int) (*ManagementClient, error) {
	base := "http://" + host + ":8002/manage/v2"
	httpClient := &http.Client{}
	var client *ManagementClient
	var digestHeaders *digestAuth.DigestHeaders
	var err error
	if authType == DigestAuth {
		digestHeaders = &digestAuth.DigestHeaders{}
		digestHeaders, err = digestHeaders.Auth(username, password, base+"?format=xml")
	}
	if err == nil {
		client = &ManagementClient{
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
