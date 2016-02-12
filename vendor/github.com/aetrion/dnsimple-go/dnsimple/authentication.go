package dnsimple

import (
	"encoding/base64"
	"fmt"
)

const (
	httpHeaderDomainToken   = "X-DNSimple-Domain-Token"
	httpHeaderApiToken      = "X-DNSimple-Token"
	httpHeaderAuthorization = "Authorization"
)

// Provides credentials that can be used for authenticating with DNSimple.
//
// See https://developer.dnsimple.com/v2/#authentication
type Credentials interface {
	// Get the HTTP header key and value to use for authentication.
	HttpHeader() (string, string)
}

// Domain token authentication

type domainTokenCredentials struct {
	domainToken string
}

// NewDomainTokenCredentials construct Credentials using the DNSimple Domain Token method.
func NewDomainTokenCredentials(domainToken string) Credentials {
	return &domainTokenCredentials{domainToken: domainToken}
}

func (c *domainTokenCredentials) HttpHeader() (string, string) {
	return httpHeaderDomainToken, c.domainToken
}

// HTTP basic authentication

type httpBasicCredentials struct {
	email    string
	password string
}

// NewHttpBasicCredentials construct Credentials using HTTP Basic Auth.
func NewHttpBasicCredentials(email, password string) Credentials {
	return &httpBasicCredentials{email, password}
}

func (c *httpBasicCredentials) HttpHeader() (string, string) {
	return httpHeaderAuthorization, "Basic " + basicAuth(c.email, c.password)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// OAuth token authentication

type oauthTokenCredentials struct {
	oauthToken string
}

// NewOauthTokenCredentials construct Credentials using the OAuth access token.
func NewOauthTokenCredentials(oauthToken string) Credentials {
	return &oauthTokenCredentials{oauthToken: oauthToken}
}

func (c *oauthTokenCredentials) HttpHeader() (string, string) {
	return httpHeaderAuthorization, fmt.Sprintf("Bearer %v", c.oauthToken)
}
