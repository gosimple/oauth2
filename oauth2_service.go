// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type oAuth2Service struct {
	AuthHeader http.Header
	*config
}

func Service(clientId, clientSecret string,
	authorizeURL, accessTokenURL string) *oAuth2Service {

	service := new(oAuth2Service)
	service.AuthHeader = make(http.Header)
	service.AuthHeader.Set("Accept", "application/json")
	service.AuthHeader.Set("Content-Type", "application/x-www-form-urlencoded")
	service.config = new(config)

	service.config.ClientId = clientId

	service.config.ClientSecret = clientSecret

	authURL, err := url.Parse(authorizeURL)
	if err != nil {
		panic("authorizeURL error: " + err.Error())
	}
	service.config.AuthorizeURL = *authURL

	tokenURL, err := url.Parse(accessTokenURL)
	if err != nil {
		panic("accessTokenURL error: " + err.Error())
	}
	service.config.AccessTokenURL = *tokenURL

	return service
}

// http://tools.ietf.org/html/rfc6749#section-4.1
// http://tools.ietf.org/html/rfc6749#section-4.2
func (service *oAuth2Service) GetAuthorizeURL(
	responseType, state string) string {
	params := url.Values{}

	params.Set("response_type", responseType)
	params.Set("client_id", service.config.ClientId)
	params.Set("redirect_uri", service.config.RedirectURL)
	params.Set("scope", service.config.Scope)
	params.Set("state", state)
	params.Set("access_type", service.config.AccessType)

	query := params.Encode()
	if service.config.AuthorizeURL.RawQuery == "" {
		service.config.AuthorizeURL.RawQuery = query
	} else {
		service.config.AuthorizeURL.RawQuery += "&" + query
	}
	return service.config.AuthorizeURL.String()
}

// http://tools.ietf.org/html/rfc6749#section-4.1.3
func (service *oAuth2Service) GetAccessToken(code string) (*Token, error) {
	params := url.Values{}

	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("scope", service.config.Scope)

	return service.GetToken(params)
}

// http://tools.ietf.org/html/rfc6749#section-4.3
func (service *oAuth2Service) GetAccessTokenPassword(
	username, password string) (*Token, error) {
	params := url.Values{}

	params.Set("grant_type", "password")
	params.Set("username", username)
	params.Set("password", password)
	params.Set("scope", service.config.Scope)

	return service.GetToken(params)
}

// http://tools.ietf.org/html/rfc6749#section-4.4
func (service *oAuth2Service) GetAccessTokenCredentials() (*Token, error) {
	params := url.Values{}

	params.Set("grant_type", "client_credentials")
	params.Set("scope", service.config.Scope)

	return service.GetToken(params)
}

// http://tools.ietf.org/html/rfc6749#section-6
func (service *oAuth2Service) RefreshAccessToken(
	refreshToken string) (*Token, error) {
	params := url.Values{}

	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)
	params.Set("scope", service.config.Scope)

	return service.GetToken(params)
}

// If you need more custom parameters to get access token or OAuth 2.0
// Extension Grants http://tools.ietf.org/html/rfc6749#section-4.5 you can
// build params yourself.
// "client_id", "client_secret", "redirect_uri" will be added by default.
//
//		service := oauth2.Service(clId, clSecret, authURL, tokenURL)
// 		params := url.Values{}
// 		params.Set("example_parameter1", "one")
//		params.Set("example_parameter2", "two")
//		myToken, err := service.GetToken(params)
func (service *oAuth2Service) GetToken(params url.Values) (*Token, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", service.AccessTokenURL.String(), nil)

	params.Set("client_id", service.config.ClientId)
	params.Set("client_secret", service.config.ClientSecret)
	params.Set("redirect_uri", service.config.RedirectURL)
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
		errorString := fmt.Sprintf("No access token found, response: %v", raw)
		return nil, errors.New(errorString)
	}

	return &token, nil
}
