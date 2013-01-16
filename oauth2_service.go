// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package oauth2

import (
	"net/http"
	"net/url"
)

type oAuth2Service struct {
	*Auth
	*config
	*Request
}

func Service(clientId, clientSecret string,
	authorizeURL, accessTokenURL string) *oAuth2Service {

	service := &oAuth2Service{}
	service.Auth = new(Auth)
	service.Auth.Header = http.Header{}
	service.Auth.Params = url.Values{}
	service.config = new(config)
	service.Request = new(Request)
	service.Request.Header = http.Header{}
	service.Request.Params = url.Values{}

	service.Auth.Params.Set("client_id", clientId)
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

func (service *oAuth2Service) GetAuthorizeURL(responseType string) string {
	service.Auth.Params.Set("response_type", responseType)
	service.Auth.Params.Set("redirect_uri", service.config.RedirectURL)
	service.Auth.Params.Set("scope", service.config.Scope)
	service.Auth.Params.Set("access_type", service.config.AccessType)

	query := service.Auth.Params.Encode()
	if service.config.AuthorizeURL.RawQuery == "" {
		service.config.AuthorizeURL.RawQuery = query
	} else {
		service.config.AuthorizeURL.RawQuery += "&" + query
	}
	return service.config.AuthorizeURL.String()
}
