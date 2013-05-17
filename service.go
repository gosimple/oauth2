// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type OAuth2Service struct {
	// AuthHeader allows you to add custom headers that'll be added to each
	// access token request.
	AuthHeader http.Header
	*Config
}

// Service initializes basic values that can be used to get access token.
func Service(clientId, clientSecret string,
	authorizeURL, accessTokenURL string) *OAuth2Service {

	service := new(OAuth2Service)
	service.AuthHeader = make(http.Header)
	service.AuthHeader.Set("Accept", "application/json")
	service.AuthHeader.Set("Content-Type", "application/x-www-form-urlencoded")
	service.Config = new(Config)

	service.ClientId = clientId
	service.ClientSecret = clientSecret

	authURL, err := url.Parse(authorizeURL)
	if err != nil {
		panic("authorizeURL error: " + err.Error())
	}
	service.AuthorizeURL = *authURL

	tokenURL, err := url.Parse(accessTokenURL)
	if err != nil {
		panic("accessTokenURL error: " + err.Error())
	}
	service.AccessTokenURL = *tokenURL
	service.ResponseType = "code"

	return service
}

// GetAuthorizeURL
func (service *OAuth2Service) GetAuthorizeURL(state string) string {
	// http://tools.ietf.org/html/rfc6749#section-4.1
	// http://tools.ietf.org/html/rfc6749#section-4.2
	params := url.Values{}

	params.Set("response_type", service.ResponseType)
	params.Set("client_id", service.ClientId)
	(*MyUrlValues)(&params).CheckAndSet("redirect_uri", service.RedirectURL)
	(*MyUrlValues)(&params).CheckAndSet("scope", service.Scope)
	(*MyUrlValues)(&params).CheckAndSet("state", state)
	(*MyUrlValues)(&params).CheckAndSet("access_type", service.AccessType)

	query := params.Encode()
	if service.AuthorizeURL.RawQuery == "" {
		service.AuthorizeURL.RawQuery = query
	} else {
		service.AuthorizeURL.RawQuery += "&" + query
	}
	return service.AuthorizeURL.String()
}

// GetAccessToken
func (service *OAuth2Service) GetAccessToken(accessCode string) (
	*Token, error) {
	// http://tools.ietf.org/html/rfc6749#section-4.1.3
	params := url.Values{}

	params.Set("grant_type", "authorization_code")
	params.Set("code", accessCode)
	(*MyUrlValues)(&params).CheckAndSet("scope", service.Scope)

	return service.getToken(params)
}

// GetAccessTokenPassword
func (service *OAuth2Service) GetAccessTokenPassword(
	username, password string) (*Token, error) {
	// http://tools.ietf.org/html/rfc6749#section-4.3
	params := url.Values{}

	params.Set("grant_type", "password")
	params.Set("username", username)
	params.Set("password", password)
	(*MyUrlValues)(&params).CheckAndSet("scope", service.Scope)

	return service.getToken(params)
}

// GetAccessTokenCredentials
func (service *OAuth2Service) GetAccessTokenCredentials() (*Token, error) {
	// http://tools.ietf.org/html/rfc6749#section-4.4
	params := url.Values{}

	params.Set("grant_type", "client_credentials")
	(*MyUrlValues)(&params).CheckAndSet("scope", service.Scope)

	return service.getToken(params)
}

// RefreshAccessToken
func (service *OAuth2Service) RefreshAccessToken(refreshToken string) (
	*Token, error) {
	// http://tools.ietf.org/html/rfc6749#section-6
	params := url.Values{}

	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)
	(*MyUrlValues)(&params).CheckAndSet("scope", service.Scope)

	return service.getToken(params)
}

// If you need more custom parameters to get access token or OAuth 2.0
// Extension Grants http://tools.ietf.org/html/rfc6749#section-4.5 you can
// provide custom URL parameters.
// "client_id", "client_secret", "redirect_uri" (if exists) will be added
// by default.
//
//		service := oauth2.Service(clId, clSecret, authURL, tokenURL)
//		// get access code
//		code := "..."
// 		params := url.Values{}
// 		params.Set("example_parameter1", "one")
//		params.Set("example_parameter2", "two")
//		myToken, err := service.GetToken(code, params)
func (service *OAuth2Service) GetToken(accessCode string, params url.Values) (
	*Token, error) {
	params.Set("code", accessCode)
	return service.getToken(params)
}

// getToken makes request for token
func (service *OAuth2Service) getToken(params url.Values) (
	*Token, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", service.AccessTokenURL.String(), nil)
	if err != nil {
		return nil, err
	}

	params.Set("client_id", service.ClientId)
	params.Set("client_secret", service.ClientSecret)
	(*MyUrlValues)(&params).CheckAndSet("redirect_uri", service.RedirectURL)

	encParams := params.Encode()
	reader := strings.NewReader(encParams)
	req.Body = ioutil.NopCloser(reader)
	req.ContentLength = int64(len(encParams))

	req.Header = service.AuthHeader

	//for key, val := range req.Header {
	//	fmt.Println(key, val)
	//}

	//raw2, _ := ioutil.ReadAll(req.Body)
	//fmt.Println(string(raw2))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Get the response body
	raw, err := ioutil.ReadAll(resp.Body)
	//for key, val := range resp.Header {
	//	fmt.Println(key, val)
	//}
	//fmt.Println(string(raw))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Parse response body to get the localToken
	var localToken struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresInt64 int64  `json:"expires_in"`
		ExpiresIn    time.Duration
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		State        string `json:"state"`
	}

	content, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain", "text/html":
		if err != nil {
			return nil, err
		}
		vals, err := url.ParseQuery(string(raw))
		if err != nil {
			return nil, err
		}
		//for key, val := range vals {
		//	fmt.Println(key, val)
		//}
		localToken.AccessToken = vals.Get("access_token")
		localToken.TokenType = vals.Get("token_type")
		localToken.ExpiresIn, _ = time.ParseDuration(vals.Get("expires") + "s")
		localToken.RefreshToken = vals.Get("refresh_token")
		localToken.Scope = vals.Get("scope")
		localToken.State = vals.Get("state")
	default:
		if err := json.Unmarshal(raw, &localToken); err != nil {
			return nil, err
		}
		expiresIn := strconv.FormatInt(localToken.ExpiresInt64, 10)
		localToken.ExpiresIn, _ = time.ParseDuration(expiresIn + "s")
	}

	// Create return token
	token := Token{}
	token.AccessToken = localToken.AccessToken
	token.TokenType = localToken.TokenType
	if localToken.ExpiresIn == 0 {
		token.ExpirationTime = time.Time{}
	} else {
		token.ExpirationTime = time.Now().Add(localToken.ExpiresIn)
	}
	if len(localToken.RefreshToken) > 0 {
		token.RefreshToken = localToken.RefreshToken
	}
	token.Scope = localToken.Scope
	token.State = localToken.State

	if len(token.AccessToken) == 0 {
		tokenError := Error{}
		switch content {
		case "application/x-www-form-urlencoded", "text/plain", "text/html":
			vals, err := url.ParseQuery(string(raw))
			if err != nil {
				return nil, err
			}
			tokenError.Type = vals.Get("error")
			tokenError.Description = vals.Get("error_description")
			tokenError.URI = vals.Get("error_uri")
			tokenError.State = vals.Get("state")
		default:
			if err := json.Unmarshal(raw, &tokenError); err != nil {
				return nil, err
			}
		}
		return nil, fmt.Errorf("No access token found, response: %v", tokenError)
	}

	return &token, nil
}

// MyUrlValues is a wrapper to a url.Values
type MyUrlValues url.Values

// CheckAndSet sets the key to value if value is not empty. It replaces any
// existing values.
//
//		params := url.Values{}
//		params.Set("one", "")
//		params.Set("two", "")
//		(*MyUrlValues)(&params).CheckAndSet("three", "")
//		(*MyUrlValues)(&params).CheckAndSet("four", "4")
//		// params.Encode() => two=&one=&four=4
func (params *MyUrlValues) CheckAndSet(key, value string) {
	if len(value) > 0 {
		(*url.Values)(params).Set(key, value)
	}
}
