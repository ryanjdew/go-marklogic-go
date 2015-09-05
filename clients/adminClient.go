package clients

import (
	"net/http"
	"net/url"

	digestAuth "github.com/ryanjdew/http-digest-auth-client"
)

// AdminClient is used for connecting to the MarkLogic Management API.
type AdminClient struct {
	*BasicClient
}

// NewAdminClient creates the Client struct used for managing admin features, etc.
func NewAdminClient(host string, username string, password string, authType int) (*AdminClient, error) {
	base := "http://" + host + ":8001/admin/v1"
	httpClient := &http.Client{}
	var client *AdminClient
	var digestHeaders *digestAuth.DigestHeaders
	var err error
	if authType == DigestAuth {
		digestHeaders = &digestAuth.DigestHeaders{}
		digestHeaders, err = digestHeaders.Auth(username, password, base+"?format=xml")
	}
	if err == nil {
		client = &AdminClient{
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
