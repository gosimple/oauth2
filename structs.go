// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Spec: http://tools.ietf.org/html/rfc6749
*/

package oauth2

import (
	"net/url"
	"time"
)

// Config store configuration for OAuth 2.0 client.
type Config struct {
	ClientId       string
	ClientSecret   string
	Scope          string
	AuthorizeURL   url.URL
	AccessTokenURL url.URL
	RedirectURL    string
	ResponseType   string
	AccessType     string
}

// Token represents a successful Access Token Response.
//
// 'expires_in' is used to build ExpirationTime
type Token struct {
	// http://tools.ietf.org/html/rfc6749#section-5.1

	// The access token issued by the authorization server.
	AccessToken string `json:"access_token"`

	// The type of the token issued (bearer, mac)
	TokenType string `json:"token_type"`

	// The expiration time of the token, zero if unknown
	ExpirationTime time.Time

	// The refresh token, which can be used to obtain new
	// access tokens using the same authorization grant
	RefreshToken string `json:"refresh_token"`

	// The scope of the access token.
	Scope string `json:"scope"`

	// REQUIRED if the "state" parameter was present in the client
	// authorization request.  The exact value received from the
	// client.
	State string `json:"state"`
}

// Error represents a failed Access Token Response.
type TokenError struct {
	// http://tools.ietf.org/html/rfc6749#section-5.2

	// A single ASCII [USASCII] error code
	Error string `json:"error"`

	// A human-readable ASCII [USASCII] text providing
	// additional information, used to assist the client developer in
	// understanding the error that occurred.
	Description string `json:"error_description"`

	// A URI identifying a human-readable web page with
	// information about the error, used to provide the client
	// developer with additional information about the error.
	URI string `json:"error_uri"`

	// REQUIRED if a "state" parameter was present in the client
	// authorization request.  The exact value received from the
	// client.
	State string `json:"state"`
}
