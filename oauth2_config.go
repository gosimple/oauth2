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
)

type config struct {
	ClientId       string
	ClientSecret   string
	Scope          string
	AuthorizeURL   url.URL
	AccessTokenURL url.URL
	RedirectURL    string
	ResponseType   string
	AccessType     string
}
