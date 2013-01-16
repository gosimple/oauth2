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

type Request struct {
	Header http.Header
	Params url.Values
}

//func (req *Request) Do(req *http.Request) (
//	resp *http.Response, err error) {
//	return service.Client().Do(req)
//}

//func (req *Request) Get(url string) (resp *http.Response, err error) {
//	return service.Client().Get(url)
//}

//func (req *Request) Head(url string) (
//	resp *http.Response, err error) {
//	return service.Client().Head(url)
//}
//func (req *Request) Post(url string, bodyType string, body io.Reader) (
//	resp *http.Response, err error) {
//	return service.Client().Post(url, bodyType, body)
//}

//func (req *Request) PostForm(url string, data url.Values) (
//	resp *http.Response, err error) {
//	return service.Client().PostForm(url, data)
//}
